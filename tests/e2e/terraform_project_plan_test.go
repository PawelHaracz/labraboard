package e2e

import (
	json2 "encoding/json"
	"fmt"
	"github.com/google/uuid"
	tc "github.com/testcontainers/testcontainers-go/modules/compose"
	"io"
	"labraboard/internal/routers/api/dtos"
	"net/http"
	"testing"
)

// use this deprecated method becasue there is a bug with opentelemetry schema version between docker desktop and test containers
// conflicting Schema URL: https://opentelemetry.io/schemas/1.24.0 and https://opentelemetry.io/schemas/1.21.0
func TestTerraformProjectPlan(t *testing.T) {
	t.SkipNow()
	composeFilePaths := []string{"../../docker-compose.yaml"}
	compose := tc.NewLocalDockerCompose(composeFilePaths, uuid.New().String())
	compose.Cmd = []string{"up", "-d"}

	execError := compose.Invoke()
	//err = compose.Up(ctx)

	if execError.Error != nil {
		t.Errorf(execError.Error.Error())
	}
	defer func() {
		execError := compose.Down()
		if execError.Error != nil {
			t.Errorf(execError.Error.Error())
		}
	}()

	var baseUrl = "http://localhost:8080/api/v1"
	resp, err := http.Get(fmt.Sprintf("%s/project", baseUrl))
	if err != nil {
		t.Errorf(err.Error())
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	var listProjects []*dtos.GetProjectBaseDto
	err = json2.Unmarshal(body, &listProjects)
	if err != nil {
		t.Errorf(err.Error())
	}

	if len(listProjects) != 0 {
		t.Errorf("projects should be empty")
	}
	//resp, err := http.Post(fmt.Sprint("{0}/project", baseUrl), "application/json",
	//	strings.NewReader(`{"name": "test"}`))
}
