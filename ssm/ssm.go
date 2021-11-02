package ssm

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2instanceconnect"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	//ssmtype "github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"os"
	charonConf "xlanor/charon/config"
	"xlanor/charon/logger"
)

func ConnectPublicKey(ctx context.Context, database_url string) {
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
	user := charonConf.GetJumpHostUser()
	svc, err := connect.SendSSHPublicKey(ctx, &ec2instanceconnect.SendSSHPublicKeyInput{
		AvailabilityZone: status.InstanceStatuses[0].AvailabilityZone,
		InstanceId:       jumphost.InstanceId,
		InstanceOSUser:   &user,
		SSHPublicKey:     &pk,
	})
	if err != nil {
		logger.Sugar().Error(err.Error())
	}

	if !svc.Success {
		logger.Sugar().Error("SSM unable to upload key")
		os.Exit(1)
	}
	_, err = openSSM(ctx, cfg, *jumphost.InstanceId)
	if err != nil {
		logger.Sugar().Error(err.Error())
	}
	logger.Sugar().Info("Opened port")

}

func openSSM(ctx context.Context, cfg aws.Config, instance string) (*ssm.StartSessionOutput, error){
	docName := "AWS-StartSSHSession"

	client := ssm.NewFromConfig(cfg)

	input := &ssm.StartSessionInput {
		DocumentName: &docName,
		Target: &instance,
		Parameters: map[string][]string{
			"portNumber": []string{"3306"},
		},
	}
	return client.StartSession(ctx, input)


}