package cloudformation

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

const (
	typeProduct = "AWS::ServiceCatalog::CloudFormationProvisionedProduct"
	typeStack   = "AWS::CloudFormation::Stack"
)

type AwsResource struct {
	Type string
	Id   string
}

func ListResources(stackName string) ([]AwsResource, error) {
	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to load AWS config: %s", err.Error())
	}
	stsClient := sts.NewFromConfig(cfg)
	cfnClient := cloudformation.NewFromConfig(cfg)

	output, err := stsClient.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
	if err != nil {
		return nil, fmt.Errorf("unable to get current AWS account id: %s", err.Error())
	}
	accountId := *output.Account

	resourceList := []AwsResource{
		{
			Type: typeStack,
			Id:   stackName,
		},
	}

	for i := 0; i < len(resourceList); i++ {
		res := resourceList[i]
		if res.Type != typeStack && res.Type != typeProduct {
			fmt.Println(res.Type, res.Id)
			continue
		}

		stackName := res.Id
		if res.Type == typeProduct {
			stackName = "SC-" + accountId + "-" + res.Id
		}
		output, err := cfnClient.DescribeStackResources(
			ctx,
			&cloudformation.DescribeStackResourcesInput{
				StackName: aws.String(stackName),
			},
		)
		if err != nil {
			return nil, fmt.Errorf("failed to list stack [%s]: %s", stackName, err.Error())
		}
		for _, r := range output.StackResources {
			resourceList = append(resourceList, AwsResource{
				Type: *r.ResourceType,
				Id:   *r.PhysicalResourceId,
			})
		}
	}

	return resourceList, nil
}
