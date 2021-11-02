package rds

import (
	"net"
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

func GetPort() *int {
	validate := func(input string) error {
		num, err := strconv.ParseInt(input, 64)
		if err != nil {
			return errors.New("Invalid number")
		}
		if num < 1 || num > 65535 {
			return errors.New("Invalid Port")
		}
		return nil
	}
	while (1){
		prompt := promptui.Prompt{
			Label:    "Enter Local Port",
			Validate: validate,
		}
	
		result, err := prompt.Run()
		if err != nil {
			logger.Sugar().Error("Prompt failed")
			return nil
		}else{
			port, err := CheckPort(result)
			if err != nil {
				logger.Sugar().Error(err.Error())
			}
			if port == nil {
				logger.Sugar().Error("Recevied null pointer in port")
			}else{
				return port
			}
		}
	}
}

func CheckPort(port int) (*int, error){
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		// should not ever trigger because of validate
		return nil, err
	}
	listen, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return nil, err
	}
	defer listen.Close()
	return &l.Addr().(*net.TCPAddr).Port, nil
}