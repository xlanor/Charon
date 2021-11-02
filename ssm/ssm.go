package ssm

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2instanceconnect"
	"github.com/aws/aws-sdk-go/aws"
	"os"
	charonConf "xlanor/charon/config"
	"xlanor/charon/logger"
)

func ConnectPublicKey(ctx context.Context) {
	pk := charonConf.GetPublicKey()
	cfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithRegion(charonConf.GetAwsRegion()),
		config.WithSharedConfigProfile(charonConf.GetAwsProfile()),
	)
	client := ec2.NewFromConfig(cfg)

	jumphost, err := GetJumphost(ctx, client)
	if err != nil {
		logger.Sugar().Error(err.Error())
		os.Exit(1)
	}
	connect := ec2instanceconnect.NewFromConfig(cfg)

	status, err := client.DescribeInstanceStatus(ctx, &ec2.DescribeInstanceStatusInput{
		InstanceIds: []string{*jumphost.InstanceId},
	})

	if err != nil {
		logger.Sugar().Error(err.Error())
		os.Exit(1)
	}

	svc, err := connect.SendSSHPublicKey(ctx, &ec2instanceconnect.SendSSHPublicKeyInput{
		AvailabilityZone: status.InstanceStatuses[0].AvailabilityZone,
		InstanceId:       jumphost.InstanceId,
		InstanceOSUser:   aws.String(charonConf.GetJumpHostUser()),
		SSHPublicKey:     aws.String(pk),
	})
	if err != nil {
		logger.Sugar().Error(err.Error())
	}

	if !svc.Success {
		logger.Sugar().Error("SSM unable to upload key")
		os.Exit(1)
	}

}
