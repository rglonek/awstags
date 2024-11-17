package tags

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
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
