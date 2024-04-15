package e2e

import (
	"context"
	"github.com/stretchr/testify/require"
	tc "github.com/testcontainers/testcontainers-go/modules/compose"
	"github.com/testcontainers/testcontainers-go/wait"
	"testing"
	"time"
)

func TestTerraformProjectPlan(t *testing.T) {
	identifier := tc.StackIdentifier("some_ident")
	compose, err := tc.NewDockerComposeWith(tc.WithStackFiles("../../docker-compose.yaml"), identifier)
	require.NoError(t, err, "NewDockerComposeAPIWith()")

	defer t.Cleanup(func() {
		require.NoError(t, compose.Down(context.Background(), tc.RemoveOrphans(true), tc.RemoveImagesLocal), "compose.Down()")
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer t.Cleanup(cancel)

	err = compose.
		WaitForService("api", wait.NewHTTPStrategy("/").WithPort("8080/tcp").WithStartupTimeout(10*time.Second)).
		Up(ctx, tc.Wait(true))

	require.NoError(t, err, "compose.Up()")
}
