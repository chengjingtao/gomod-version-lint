package main

import (
	"context"
	"go.uber.org/zap"
	"gomod.alauda.cn/gomod-version-lint/cmd"
	pkgctx "gomod.alauda.cn/gomod-version-lint/pkg/context"
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
