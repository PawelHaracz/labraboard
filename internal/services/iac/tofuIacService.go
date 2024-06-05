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
	iacFolderPath string
	tf            *tfexec.Terraform
	serializer    *helpers.Serializer[entities.IacTerraformPlanJson]
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

	jsonWriter := bufio.NewWriter(&b)
	planConfig := []tfexec.PlanOption{
		tfexec.Lock(true),
		tfexec.Destroy(false),
		tfexec.Refresh(false),
		tfexec.Out(PlanPath),
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

	planJson, err := svc.tf.ShowPlanFile(context.Background(), PlanPath)
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

	planContent, err := os.ReadFile(fmt.Sprintf("%s/%s", svc.iacFolderPath, PlanPath)) //read the content of file
	if err != nil {
		log.Error().Err(err).Msg(err.Error())
		return nil, err
	}

	planChanges := models.NewIacTerraformPlanJson(jsonPlan, iacPlanDeserialized, planContent)

	return planChanges, nil
}

func (svc *TofuIacService) Apply(planId uuid.UUID, envs map[string]string, ctx context.Context) (interface{}, error) {

	log := logger.GetWitContext(ctx)
	//tfLogPath := filepath.Join("/tmp/13e92193-7c36-4d01-9589-16d8e3c96ff5/apply", "test.log")
	//if planPath == "" {
	//	err := errors.New("plan path is empty")
	//	log.Error().Err(err)
	//	return nil, err
	//}

	var b bytes.Buffer
	jsonWriter := bufio.NewWriter(&b)

	applyConfig := []tfexec.ApplyOption{
		tfexec.Lock(true),
		tfexec.Destroy(false),
		tfexec.Refresh(true),
		tfexec.DirOrPlan(fmt.Sprintf("%s/%s", svc.tf.WorkingDir(), PlanPath)),
	}

	//if len(variables) > 0 {
	//	for _, v := range variables {
	//		applyConfig = append(applyConfig, tfexec.Var(v))
	//	}
	//}
	if len(envs) > 0 {
		err := svc.tf.SetEnv(envs)
		if err != nil {
			return nil, err
		}
	}

	log.Info().Msgf("Apply plan %s", planId.String())
	//svc.tf.SetLog("")
	//err := svc.tf.SetLogPath(tfLogPath)
	//if err != nil {
	//	log.Error().Err(err)
	//	return nil, err
	//}

	//svc.tf.SetStderr(lpError)
	//svc.tf.SetStdout(lpDebug)
	//todo handle outputs: {"@level":"error","@message":"Error: Error creating Resource Group \"rg-starterterraform-staging-main\": resources.GroupsClient#CreateOrUpdate: Failure responding to request: StatusCode=400 -- Original Error: autorest/azure: Service returned an error. Status=400 Code=\"LocationNotAvailableForResourceGroup\" Message=\"The provided location 'polandcenter' is not available for resource group. List of available regions is 'eastasia,southeastasia,australiaeast,australiasoutheast,brazilsouth,canadacentral,canadaeast,switzerlandnorth,germanywestcentral,eastus2,eastus,centralus,northcentralus,francecentral,uksouth,ukwest,centralindia,southindia,jioindiawest,italynorth,japaneast,japanwest,koreacentral,koreasouth,mexicocentral,northeurope,norwayeast,polandcentral,qatarcentral,spaincentral,swedencentral,uaenorth,westcentralus,westeurope,westus2,westus,southcentralus,westus3,southafricanorth,australiacentral,australiacentral2,israelcentral,westindia'.\"","@module":"terraform.ui","@timestamp":"2024-06-05T23:04:28.550164+02:00","diagnostic":{"severity":"error","summary":"Error creating Resource Group \"rg-starterterraform-staging-main\": resources.GroupsClient#CreateOrUpdate: Failure responding to request: StatusCode=400 -- Original Error: autorest/azure: Service returned an error. Status=400 Code=\"LocationNotAvailableForResourceGroup\" Message=\"The provided location 'polandcenter' is not available for resource group. List of available regions is 'eastasia,southeastasia,australiaeast,australiasoutheast,brazilsouth,canadacentral,canadaeast,switzerlandnorth,germanywestcentral,eastus2,eastus,centralus,northcentralus,francecentral,uksouth,ukwest,centralindia,southindia,jioindiawest,italynorth,japaneast,japanwest,koreacentral,koreasouth,mexicocentral,northeurope,norwayeast,polandcentral,qatarcentral,spaincentral,swedencentral,uaenorth,westcentralus,westeurope,westus2,westus,southcentralus,westus3,southafricanorth,australiacentral,australiacentral2,israelcentral,westindia'.\"","detail":"","address":"azurerm_resource_group.main","range":{"filename":"main.tf","start":{"line":8,"column":42,"byte":277},"end":{"line":8,"column":43,"byte":278}},"snippet":{"context":"resource \"azurerm_resource_group\" \"main\"","code":"resource \"azurerm_resource_group\" \"main\" {","start_line":8,"highlight_start_offset":41,"highlight_end_offset":42,"values":[]}},"type":"diagnostic"}
	err := svc.tf.ApplyJSON(ctx, jsonWriter, applyConfig...)
	if err != nil {
		if err = jsonWriter.Flush(); err != nil {
			return nil, errors.New("error running Flush")
		}
		//r := bytes.NewReader(b.Bytes())
		//var iacDeserialized, err = svc.serializer.DeserializeJsonl(r)
		log.Error().Err(err).Msg(string(b.Bytes()))
		return nil, err
	}
	if err = jsonWriter.Flush(); err != nil {
		return nil, errors.New("error running Flush")
	}
	r := bytes.NewReader(b.Bytes())
	iacDeserialized, err := svc.serializer.DeserializeJsonl(r)
	return iacDeserialized, nil //todo
}
