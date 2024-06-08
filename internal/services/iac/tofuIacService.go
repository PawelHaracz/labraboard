package iac

import (
	"bufio"
	"bytes"
	"context"
	json2 "encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hc-install/product"
	"github.com/hashicorp/hc-install/releases"
	"github.com/hashicorp/terraform-exec/tfexec"
	"labraboard/internal/entities"
	"labraboard/internal/helpers"
	"labraboard/internal/logger"
	"labraboard/internal/models"
	"os"
)

const PlanPath = "plan.tfplan"

type TofuIacService struct {
	iacFolderPath        string
	tf                   *tfexec.Terraform
	serializer           *helpers.Serializer[entities.IacTerraformOutputJson]
	diagnosticSerializer *helpers.Serializer[entities.IacTerraformDiagnosticJson]
}

func NewTofuIacService(iacFolderPath string, ctx context.Context) (*TofuIacService, error) {
	if iacFolderPath == "" {
		return nil, errors.New("iacFolderPath is empty")
	}

	installer := &releases.ExactVersion{
		Product: product.Terraform,
		Version: version.Must(version.NewVersion("1.7.5")),
	}

	execPath, err := installer.Install(ctx)
	if err != nil {
		return nil, errors.New("error installing Terraform")
	}

	tf, err := tfexec.NewTerraform(iacFolderPath, execPath)
	if err != nil {
		return nil, err
		//log.Fatalf("error running NewTerraform: %s", err)
	}

	var config = []tfexec.InitOption{
		tfexec.Upgrade(true),
	}

	err = tf.Init(ctx, config...)
	if err != nil {
		return nil, err
	}

	serializer := helpers.NewSerializer[entities.IacTerraformOutputJson]()
	diagnosticSerializer := helpers.NewSerializer[entities.IacTerraformDiagnosticJson]()

	return &TofuIacService{
		iacFolderPath:        iacFolderPath,
		tf:                   tf,
		serializer:           serializer,
		diagnosticSerializer: diagnosticSerializer,
	}, nil
}

func (svc *TofuIacService) Plan(envs map[string]string, variables []string, ctx context.Context) (*models.IacTerraformPlanJson, error) {
	log := logger.GetWitContext(ctx)
	var b bytes.Buffer

	var planPath = fmt.Sprintf("%s/%s", svc.tf.WorkingDir(), PlanPath)
	jsonWriter := bufio.NewWriter(&b)
	planConfig := []tfexec.PlanOption{
		tfexec.Lock(true),
		tfexec.Destroy(false),
		tfexec.Refresh(false),
		tfexec.Out(planPath),
	}

	if len(variables) > 0 {
		for _, v := range variables {
			planConfig = append(planConfig, tfexec.Var(v))
		}
	}
	if len(envs) > 0 {
		err := svc.tf.SetEnv(envs)
		if err != nil {
			return nil, err
		}
	}
	log.Info().Msg("Running plan")
	p, err := svc.tf.PlanJSON(context.Background(), jsonWriter, planConfig...)
	log.Info().Msg("finished running plan")
	if err != nil {
		if err = jsonWriter.Flush(); err != nil {
			return nil, errors.New("error running Flush")
		}
		r := bytes.NewReader(b.Bytes())
		iacDeserialized, err1 := svc.diagnosticSerializer.DeserializeJsonl(r)
		if err1 != nil {
			log.Error().Err(err1).Msg(string(b.Bytes()))
		}
		if len(iacDeserialized) != 0 {
			log.Error().Err(err).Msg(iacDeserialized[0].Message)
		} else {
			log.Error().Err(errors.New("deserialize doesn't contain any items")).Msg(string(b.Bytes()))
		}
		return nil, errors.Join(fmt.Errorf("%s: %v", "error running Plan", err), err, err1, errors.New(iacDeserialized[0].Message))

	}

	planJson, err := svc.tf.ShowPlanFile(context.Background(), planPath)
	if err != nil {
		return nil, fmt.Errorf("%s: %v", "error running ShowPlanFile", err)
	}
	if !p {
		log.Warn().Err(errors.New("plan is not finish well")).Msg("")
		//return nil, errors.Join(errors.New("plan is not finish well"), err)
	}

	jsonPlan, err := json2.Marshal(planJson)
	if err := jsonWriter.Flush(); err != nil {
		return nil, errors.New("error running Flush")
	}
	r := bytes.NewReader(b.Bytes())

	iacPlanDeserialized, err := svc.serializer.DeserializeJsonl(r)
	if err != nil {
		return nil, errors.New("cannot reade plan")
	}

	planContent, err := os.ReadFile(fmt.Sprintf("%s/%s", svc.iacFolderPath, PlanPath)) //read the content of file
	if err != nil {
		log.Error().Err(err).Msg(err.Error())
		return nil, err
	}

	planChanges := models.NewIacTerraformPlanJson(jsonPlan, iacPlanDeserialized, planContent)

	return planChanges, nil
}

func (svc *TofuIacService) Apply(planId uuid.UUID, envs map[string]string, ctx context.Context) ([]entities.IacTerraformOutputJson, error) {
	log := logger.GetWitContext(ctx)

	var b bytes.Buffer
	jsonWriter := bufio.NewWriter(&b)

	applyConfig := []tfexec.ApplyOption{
		tfexec.Lock(true),
		tfexec.Destroy(false),
		tfexec.Refresh(true),
		tfexec.DirOrPlan(fmt.Sprintf("%s/%s", svc.tf.WorkingDir(), PlanPath)),
	}

	if len(envs) > 0 {
		err := svc.tf.SetEnv(envs)
		if err != nil {
			return nil, err
		}
	}

	log.Info().Msgf("Apply plan %s", planId.String())

	err := svc.tf.ApplyJSON(ctx, jsonWriter, applyConfig...)
	if err != nil {
		if err = jsonWriter.Flush(); err != nil {
			return nil, errors.New("error running Flush")
		}
		r := bytes.NewReader(b.Bytes())
		iacDeserialized, err1 := svc.diagnosticSerializer.DeserializeJsonl(r)
		if err1 != nil {
			log.Error().Err(err1).Msg(string(b.Bytes()))
		}
		var msg = ""
		if len(iacDeserialized) != 0 {
			msg = iacDeserialized[len(iacDeserialized)-1].Message
			log.Error().Err(err).Msg(msg)
		} else {
			msg = "deserialize doesn't contain any items"
			log.Error().Err(errors.New(msg)).Msg(string(b.Bytes()))
		}
		return nil, errors.Join(err, err1, errors.New(msg))
	}
	if err = jsonWriter.Flush(); err != nil {
		return nil, errors.New("error running Flush")
	}
	r := bytes.NewReader(b.Bytes())

	iacDeserialized, err := svc.serializer.DeserializeJsonl(r)
	return iacDeserialized, nil
}
