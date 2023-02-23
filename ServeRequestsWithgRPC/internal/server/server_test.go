package server

import (
	"context"
	"net"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	api "github.com/ppmasa8/proglog/api/v1"
	"github.com/ppmasa8/proglog/internal/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

func TestServer(t *testing.T) {
	for scenario, fn := range map[string]func(
		t *testing.T,
		client api.LogClient,
		config *Config
	){
		"produce/consume a message to/from the log succeeds": testProduceConsume,
		"produce.consume stream  succeeds": testProduceConsumeStream,
		"consume past log boundary fails": testConsumePastBoundary,
	} {
		t.Run(scenario, func(t *testing.T) {
			client, config, teardown := setupTest(t, nil)
			defer teardown()
			fn(t, client, config)
		})
	}
}

func setupTest(t *testing.T, fn func(*Config)) (
	client api.LogClient,
	cfg *Config,
	teardown func(),
) {
	t. Helper()

	l, err := new.Listen("tcp", ":0")
	require.NoError(t, err)
}