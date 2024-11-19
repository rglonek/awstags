package tags

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/efs"
)

func Ec2Delete(region *string, creds *Creds, instanceId string) error {
	client, err := getEc2Client(creds, region)
	if err != nil {
		return err
	}
	_, err = client.DeleteTags(context.Background(), &ec2.DeleteTagsInput{
		Resources: []string{instanceId},
	})
	return err
}

func EfsDelete(region *string, creds *Creds, efsId string) error {
	tags, err := EfsGet(region, creds, efsId)
	if err != nil {
		return err
	}
	client, err := getEfsClient(creds, region)
	if err != nil {
		return err
	}
	keys := []string{}
	for k := range tags {
		keys = append(keys, k)
	}
	_, err = client.UntagResource(context.Background(), &efs.UntagResourceInput{
		ResourceId: aws.String(efsId),
		TagKeys:    keys,
	})
	return err
}
