package ssm

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go/aws"
	charonConf "xlanor/charon/config"
	"xlanor/charon/logger"
)

func GetJumphost(ctx context.Context, client *ec2.Client) (*types.Instance, error) {

	filters := make([]types.Filter, 0)

	filters = append(filters, types.Filter{
		Name: aws.String(fmt.Sprintf("tag:%s", charonConf.GetJumpHostTagName())),
		Values: []string{
			charonConf.GetJumpHostTagValue(),
		},
	})
	if charonConf.GetJumpHostSecondaryValue() != "" {
		filters = append(filters, types.Filter{
			Name: aws.String(fmt.Sprintf("tag:%s", charonConf.GetJumpHostSecondaryName())),
			Values: []string{
				charonConf.GetJumpHostSecondaryValue(),
			},
		},
		)
	}
	svc, err := client.DescribeInstances(
		ctx, &ec2.DescribeInstancesInput{
			Filters: []types.Filter{
				{
					Name: aws.String(fmt.Sprintf("tag:%s", charonConf.GetJumpHostTagName())),
					Values: []string{
						charonConf.GetJumpHostTagValue(),
					},
				},
			},
		})

	if err != nil {
		return nil, err
	}

	if len(svc.Reservations) == 0 {
		return nil, errors.New(fmt.Sprintf("No ec2 tagged with %s found", charonConf.GetJumpHostTagName()))
	} else {
		logger.Sugar().Infof("Retrieved %d instances\n", len(svc.Reservations))
		return &svc.Reservations[0].Instances[0], nil
	}
}
