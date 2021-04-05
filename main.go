package main

import (
	"context"
	"log"
	"os"

	"github.com/sourcegraph/jsonrpc2"
	"go.uber.org/zap"
)

func main() {
	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{
		"/Users/adrian/github.com/a-h/qt-lsp/log.txt",
	}
	logger, err := cfg.Build()
	if err != nil {
		log.Printf("failed to create logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()
	logger.Info("Starting up...")
	handler := NewHandler(logger)
	<-jsonrpc2.NewConn(context.Background(), jsonrpc2.NewBufferedStream(stdrwc{}, jsonrpc2.VSCodeObjectCodec{}), handler).DisconnectNotify()
}

type stdrwc struct{}

func (stdrwc) Read(p []byte) (int, error) {
	return os.Stdin.Read(p)
}

func (stdrwc) Write(p []byte) (int, error) {
	return os.Stdout.Write(p)
}

func (stdrwc) Close() error {
	if err := os.Stdin.Close(); err != nil {
		return err
	}
	return os.Stdout.Close()
}
