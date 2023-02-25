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
	clientOptions := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials())
	}
	cc, err := grpc.Dial(l.Addr().String(), clientOptions...)
	require.NoError(t, err)

	dir, er := os.MkdirTemp("", "server-test")
	require.NoError(t, err)

	cfg = &Config{
		CommitLog: clog,
	}

	if fn != nil {
		fn(cfg)
	}
	server, err := NewGRPCServer(cfg)
	require.NoError(t, err)

	go func() {
		server.Serve(l)
	}()

	client = api.NewLogClient(cc)

	return client, cfg, func() {
		cc.Close()
		server.Stop()
		l.Close()
		clog.Remove()
	}
}