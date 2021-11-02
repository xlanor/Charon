package rds

import (
	//"net"
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
	"os"
	charonConf "xlanor/charon/config"
	logger "xlanor/charon/logger"
)

func GetRds(ctx context.Context) []types.DBInstance {
	cfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithRegion(charonConf.GetAwsRegion()),
		config.WithSharedConfigProfile(charonConf.GetAwsProfile()),
	)
	client := rds.NewFromConfig(cfg)
	result, err := client.DescribeDBInstances(ctx, &rds.DescribeDBInstancesInput{})
	if err != nil {
		logger.Sugar().Error(err.Error())
		os.Exit(1)
	}

	if len(result.DBInstances) == 0 {
		logger.Sugar().Error("No Databases avaliable")
		os.Exit(1)
	}
	return result.DBInstances
}
