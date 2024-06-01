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

type TofuIacService struct {
	iacFolderPath string
	tf            *tfexec.Terraform
	serializer    *helpers.Serializer[entities.IacTerraformPlanJson]
}

func NewTofuIacService(iacFolderPath string) (*TofuIacService, error) {
	if iacFolderPath == "" {
		return nil, errors.New("iacFolderPath is empty")
	}

	installer := &releases.ExactVersion{
		Product: product.Terraform,
		Version: version.Must(version.NewVersion("1.7.5")),
	}

	execPath, err := installer.Install(context.Background())
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

	err = tf.Init(context.Background(), config...)
	if err != nil {
		return nil, err
	}

	serializer := helpers.NewSerializer[entities.IacTerraformPlanJson]()

	return &TofuIacService{
		iacFolderPath: iacFolderPath,
		tf:            tf,
		serializer:    serializer,
	}, nil
}

func (svc *TofuIacService) Plan(envs map[string]string, variables []string, ctx context.Context) (*models.IacTerraformPlanJson, error) {
	log := logger.GetWitContext(ctx)
	var b bytes.Buffer
	var planPath = "plan.tfplan"
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
		return nil, fmt.Errorf("%s: %v", "error running Plan", err)

	}
	if !p {
		return nil, errors.New("plan is not finish well")
	}

	planJson, err := svc.tf.ShowPlanFile(context.Background(), planPath)
	if err != nil {
		return nil, fmt.Errorf("%s: %v", "error running ShowPlanFile", err)
	}

	jsonPlan, err := json2.Marshal(planJson)
	if err := jsonWriter.Flush(); err != nil {
		return nil, errors.New("error running Flush")
	}
	r := bytes.NewReader(b.Bytes())

	iacPlanDeserialized, err := svc.serializer.DeserializeJsonl(r)
	if err != nil {
		return nil, errors.New("Cannot reade plan")
	}

	planContent, err := os.ReadFile(fmt.Sprintf("%s/%s", svc.iacFolderPath, planPath)) //read the content of file
	if err != nil {
		log.Error().Err(err).Msg(err.Error())
		return nil, err
	}

	planChanges := models.NewIacTerraformPlanJson(jsonPlan, iacPlanDeserialized, planContent)

	return planChanges, nil
}

func (svc *TofuIacService) Apply(planId uuid.UUID, envs map[string]string, variables []string, planPath string, ctx context.Context) (interface{}, error) {
	log := logger.GetWitContext(ctx)
	if planPath == "" {
		var err = errors.New("plan path is empty")
		log.Error().Err(err)
		return nil, err
	}

	var b bytes.Buffer
	jsonWriter := bufio.NewWriter(&b)
	applyConfig := []tfexec.ApplyOption{
		tfexec.Lock(true),
		tfexec.Destroy(false),
		tfexec.Refresh(false),
		tfexec.DirOrPlan(planPath),
	}

	if len(variables) > 0 {
		for _, v := range variables {
			applyConfig = append(applyConfig, tfexec.Var(v))
		}
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
		log.Error().Err(err)
		return nil, err
	}
	if err = jsonWriter.Flush(); err != nil {
		return nil, errors.New("error running Flush")
	}
	r := bytes.NewReader(b.Bytes())
	iacDeserialized, err := svc.serializer.DeserializeJsonl(r)
	return iacDeserialized, nil //todo
}
