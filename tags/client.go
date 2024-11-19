package tags

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/efs"
)

type Creds struct {
	Static      *CredsStatic
	ProfileName *string
}

type CredsStatic struct {
	Key    string
	Secret string
}

func getCfgForClient(creds *Creds, region *string) (aws.Config, error) {
	opts := []func(*config.LoadOptions) error{}
	if creds != nil {
		if creds.Static != nil {
			opts = append(opts, config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(creds.Static.Key, creds.Static.Secret, "")))
		}
		if creds.ProfileName != nil {
			opts = append(opts, config.WithSharedConfigProfile(*creds.ProfileName))
		}
	}
	if region != nil {
		opts = append(opts, config.WithRegion(*region))
	}
	cfg, err := config.LoadDefaultConfig(context.Background(), opts...)
	if err != nil {
		return aws.Config{}, err
	}
	return cfg, nil
}

func getEc2Client(creds *Creds, region *string) (client *ec2.Client, err error) {
	cfg, err := getCfgForClient(creds, region)
	if err != nil {
		return nil, err
	}
	client = ec2.NewFromConfig(cfg)
	return client, nil
}

func getEfsClient(creds *Creds, region *string) (client *efs.Client, err error) {
	cfg, err := getCfgForClient(creds, region)
	if err != nil {
		return nil, err
	}
	client = efs.NewFromConfig(cfg)
	return client, nil
}
