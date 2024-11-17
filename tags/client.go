package tags

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

type Creds struct {
	Static      *CredsStatic
	ProfileName *string
}

type CredsStatic struct {
	Key    string
	Secret string
}

func getEc2Client(creds *Creds, region *string) (client *ec2.Client, err error) {
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
		return nil, err
	}
	client = ec2.NewFromConfig(cfg)
	return client, nil
}
