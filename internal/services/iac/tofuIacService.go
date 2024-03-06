package iac

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hc-install/product"
	"github.com/hashicorp/hc-install/releases"
	"github.com/hashicorp/terraform-exec/tfexec"
	entities "labraboard/internal/entities"
	"log"
)

type TofuIacService struct {
	iacFolderPath string
	tf            *tfexec.Terraform
}

func NewTofuIacService(iacFolderPath string) (*TofuIacService, error) {
	if iacFolderPath == "" {
		return nil, errors.New("iacFolderPath is empty")
	}

	installer := &releases.ExactVersion{
		Product: product.Terraform,
		Version: version.Must(version.NewVersion("1.0.6")),
	}

	execPath, err := installer.Install(context.Background())
	if err != nil {
		log.Fatalf("error installing Terraform: %s", err)
		return nil, errors.New("error installing Terraform")
	}

	tf, err := tfexec.NewTerraform(iacFolderPath, execPath)
	if err != nil {
		log.Fatalf("error running NewTerraform: %s", err)
	}

	err = tf.Init(context.Background(), tfexec.Upgrade(true))
	if err != nil {
		log.Fatalf("error running Init: %s", err)
	}

	return &TofuIacService{
		iacFolderPath: iacFolderPath,
		tf:            tf,
	}, nil
}

func (svc *TofuIacService) Plan(planId uuid.UUID) (*Plan, error) {
	var b bytes.Buffer
	jsonWriter := bufio.NewWriter(&b)
	//planConfig := []tfexec.PlanOption{
	//	tfexec.Lock(true),
	//	tfexec.Destroy(false),
	//	tfexec.Refresh(false),
	//}
	//p, err := svc.tf.PlanJSON(context.Background(), jsonWriter, planConfig...)
	p, err := svc.tf.PlanJSON(context.Background(), jsonWriter)
	if err != nil {
		//log.Fatalf("error running Plan: %s", err)
		return nil, errors.New("error running Plan")

	}
	if !p {
		return nil, errors.New("plan is not finish well")
	}

	if err := jsonWriter.Flush(); err != nil {
		return nil, errors.New("error running Flush")
	}
	r := bytes.NewReader(b.Bytes())
	plans, err := entities.SerializeIacTerraformPlanJsons(r)
	if err != nil {
		return nil, errors.New("Cannot reade plan")
	}
	//todo convert plans to more better object

	return &Plan{
		//plan: result,
		Id:   planId,
		Type: Tofu,
		plan: plans,
	}, nil
}
