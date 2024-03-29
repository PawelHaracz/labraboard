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
	"labraboard/internal/aggregates"
	"labraboard/internal/entities"
	"labraboard/internal/helpers"
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
		//log.Fatalf("error running Init: %s", err)
	}

	serializer := helpers.NewSerializer[entities.IacTerraformPlanJson]()

	return &TofuIacService{
		iacFolderPath: iacFolderPath,
		tf:            tf,
		serializer:    serializer,
	}, nil
}

func (svc *TofuIacService) Plan(planId uuid.UUID, envs map[string]string, variables []string) (*Plan, error) {
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

	p, err := svc.tf.PlanJSON(context.Background(), jsonWriter, planConfig...)
	//p, err := svc.tf.PlanJSON(context.Background(), jsonWriter)
	//p, err := svc.tf.Plan(context.Background(), planConfig...)
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
	plans, err := svc.serializer.SerializeJsonl(r)
	if err != nil {
		return nil, errors.New("Cannot reade plan")
	}

	plan, err := aggregates.NewIacPlan(planId, aggregates.Tofu, jsonPlan, nil, nil)

	if err != nil {
		return nil, errors.New("Cannot create aggregate")
	}

	plan.AddChanges(plans...)

	return &Plan{
		//plan: result,
		Id:   planId,
		Type: Tofu,
		plan: plan,
	}, nil
}
