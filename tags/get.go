package tags

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

func Ec2Get(region *string, creds *Creds, instanceId string) (map[string]string, error) {
	client, err := getEc2Client(creds, region)
	if err != nil {
		return nil, err
	}
	tags := make(map[string]string)
	var nextToken *string
	for {
		out, err := client.DescribeInstances(context.Background(), &ec2.DescribeInstancesInput{
			InstanceIds: []string{instanceId},
			NextToken:   nextToken,
		})
		if err != nil {
			return nil, err
		}
		for _, r := range out.Reservations {
			for _, i := range r.Instances {
				for _, t := range i.Tags {
					tags[aws.ToString(t.Key)] = aws.ToString(t.Value)
				}
			}
		}
		if out.NextToken == nil {
			break
		}
		nextToken = out.NextToken
	}
	return tags, nil
}
