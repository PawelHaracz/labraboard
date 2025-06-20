package e2e

import (
	"context"
	json2 "encoding/json"
	"fmt"
	"io"
	"labraboard/internal/routers/api/dtos"
	"net/http"
	"testing"

	tc "github.com/testcontainers/testcontainers-go/modules/compose"
)

// use this deprecated method becasue there is a bug with opentelemetry schema version between docker desktop and test containers
// conflicting Schema URL: https://opentelemetry.io/schemas/1.24.0 and https://opentelemetry.io/schemas/1.21.0
func TestTerraformProjectPlan(t *testing.T) {
	composeFilePaths := []string{"../../docker-compose.yaml"}
	compose, err := tc.NewDockerCompose(composeFilePaths...)
	if err != nil {
		t.Errorf("failed to create compose: %v", err)
		return
	}

	err = compose.Up(context.Background(), tc.Wait(true))
	if err != nil {
		t.Errorf("compose up failed: %v", err)
	}
	defer func() {
		err := compose.Down(context.Background())
		if err != nil {
			t.Errorf("compose down failed: %v", err)
		}
	}()

	var baseUrl = "http://localhost:8080/api/v1"
	resp, err := http.Get(fmt.Sprintf("%s/project", baseUrl))
	if err != nil {
		t.Errorf("http request failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	var listProjects []*dtos.GetProjectBaseDto
	err = json2.Unmarshal(body, &listProjects)
	if err != nil {
		t.Errorf("unmarshal failed: %v", err)
	}

	if len(listProjects) != 0 {
		t.Errorf("projects should be empty")
	}
	//resp, err := http.Post(fmt.Sprint("{0}/project", baseUrl), "application/json",
	//	strings.NewReader(`{"name": "test"}`))
}
