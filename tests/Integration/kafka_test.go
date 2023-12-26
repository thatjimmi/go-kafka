package integration

import (
	"context"
	"testing"

	tc "github.com/testcontainers/testcontainers-go/modules/compose"
)

func TestKafka(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	compose, err := tc.NewDockerCompose("../../docker-compose.yml")
	if err != nil {
		t.Fatalf("Error obtaining docker compose: %v", err)
	}
	defer func() {
		if err := compose.Down(ctx, tc.RemoveOrphans(true), tc.RemoveImagesLocal); err != nil {
			t.Fatalf("Error stopping docker compose: %v", err)
		}
	}()

	if err := compose.Up(ctx, tc.Wait(true)); err != nil {
		t.Fatalf("Error starting docker compose: %v", err)
	}
}
