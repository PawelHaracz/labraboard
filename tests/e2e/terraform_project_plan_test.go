package e2e

import (
	"context"
	tc "github.com/testcontainers/testcontainers-go/modules/compose"
	"testing"
)

func TestTerraformProjectPlan(t *testing.T) {
	ctx := context.Background()
	//compose, err := tc.NewDockerCompose("../../docker-compose.yaml")
	//require.NoError(t, err, "NewDockerComposeAPI()")
	//
	//t.Cleanup(func() {
	//	require.NoError(t, compose.Down(context.Background(), tc.RemoveOrphans(true), tc.RemoveImagesLocal), "compose.Down()")
	//})
	//
	//ctx, cancel := context.WithCancel(context.Background())
	//t.Cleanup(cancel)
	//require.NoError(t, compose.Up(ctx, tc.Wait(true)), "compose.Up()")
	//identifier := strings.ToLower(uuid.New().String())
	composeFilePaths := []string{"../../docker-compose.yaml"}
	compose, err := tc.NewDockerComposeWith(tc.WithStackFiles(composeFilePaths...))
	if err != nil {
		t.Errorf(err.Error())
	}
	err = compose.Up(ctx)
	//WithCommand([]string{"up", "-d"}).
	//	WithEnv(map[string]string{
	//		"key1": "value1",
	//		"key2": "value2",
	//	}).
	//	Invoke()
	//err := execError.Error
	if err != nil {
		t.Errorf(err.Error())
	}

	//req := testcontainers.ContainerRequest{
	//	Image:        "redis:latest",
	//	ExposedPorts: []string{"6379/tcp"},
	//	WaitingFor:   wait.ForLog("Ready to accept connections"),
	//}
	//redisC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
	//	ContainerRequest: req,
	//	Started:          true,
	//})
	//if err != nil {
	//	log.Fatalf("Could not start redis: %s", err)
	//}
	//defer func() {
	//	if err := redisC.Terminate(ctx); err != nil {
	//		log.Fatalf("Could not stop redis: %s", err)
	//	}
	//}()
	//
	//container, err := postgres.RunContainer(
	//	ctx,
	//	testcontainers.WithImage("docker.io/postgres:16-alpine"),
	//	postgres.WithDatabase("dbname"),
	//	postgres.WithUsername(user),
	//	postgres.WithPassword(password),
	//	testcontainers.WithWaitStrategy(
	//		wait.ForLog("database system is ready to accept connections").
	//			WithOccurrence(2).
	//			WithStartupTimeout(5*time.Second)),
	//)
	//if err != nil {
	//	t.Fatal(err)
	//}

	//dbName := "users"
	//dbUser := "user"
	//dbPassword := "password"
	//
	//postgresContainer, err := postgres.RunContainer(ctx,
	//	testcontainers.WithImage("docker.io/postgres:16-alpine"),
	//	postgres.WithInitScripts(filepath.Join("testdata", "init-user-db.sh")),
	//	postgres.WithConfigFile(filepath.Join("testdata", "my-postgres.conf")),
	//	postgres.WithDatabase(dbName),
	//	postgres.WithUsername(dbUser),
	//	postgres.WithPassword(dbPassword),
	//	testcontainers.WithWaitStrategy(
	//		wait.ForLog("database system is ready to accept connections").
	//			WithOccurrence(2).
	//			WithStartupTimeout(5*time.Second)),
	//)
	//if err != nil {
	//	log.Fatalf("failed to start container: %s", err)
	//}
	//
	//// Clean up the container
	//defer func() {
	//	if err := postgresContainer.Terminate(ctx); err != nil {
	//		log.Fatalf("failed to terminate container: %s", err)
	//	}
	//}()
}
