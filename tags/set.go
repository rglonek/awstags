package tags

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

func Ec2Set(region *string, creds *Creds, instanceId string, tags map[string]string) error {
	client, err := getEc2Client(creds, region)
	if err != nil {
		return err
	}
	t := []types.Tag{}
	for k, v := range tags {
		t = append(t, types.Tag{
			Key:   aws.String(k),
			Value: aws.String(v),
		})
	}
	_, err = client.CreateTags(context.Background(), &ec2.CreateTagsInput{
		Resources: []string{instanceId},
		Tags:      t,
	})
	return err
}
