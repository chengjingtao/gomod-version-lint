package main

import (
	"context"
	"github.com/chengjingtao/gomod-version-lint/cmd"
	pkgctx "github.com/chengjingtao/gomod-version-lint/pkg/context"
	"go.uber.org/zap"
	"os"
)

func main() {
	ctx := context.Background()
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	log := logger.Sugar()

	ctx = pkgctx.WithLogger(ctx, log)

	rootCmd := cmd.NewRootCmd(ctx)
	err := rootCmd.Execute()
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
}
