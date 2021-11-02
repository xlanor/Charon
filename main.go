package main

import (
	"github.com/briandowns/spinner"
	"time"
	"os"
	"context"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	localLog "xlanor/charon/logger"
	localConfig "xlanor/charon/config"
	"xlanor/charon/rds"
	"xlanor/charon/ssm"
	"github.com/manifoldco/promptui"
)

func selectRDS(ctx context.Context){
	rds_list := rds.GetRds(ctx)
	// Re-order to string
	db_string := make([]string, 0)

	for _, d := range rds_list {
		db_string = append(db_string, *d.Endpoint.Address)
	}

	prompt := promptui.Select{
		Label: "Select Database",
		Items:  db_string,
	}

	_, result, err := prompt.Run()
	if err != nil {
		localLog.Sugar().Error(err.Error())
		os.Exit(1)
	}
	localLog.Sugar().Info("Selected ", result)

	ssm.ConnectPublicKey(ctx, result)
}

func main() {
	s := spinner.New(spinner.CharSets[7], 100*time.Millisecond)
	s.Suffix = "test1"
	s.Start()
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	logger, _ = config.Build()

	undo := zap.ReplaceGlobals(logger)
	defer undo()
	localConfig.Load("config/koanf.toml")

	ctx := context.Background()
	selectRDS(ctx)  
	s.Stop()
}
