package tags

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/efs"
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

func EfsGet(region *string, creds *Creds, efsId string) (map[string]string, error) {
	client, err := getEfsClient(creds, region)
	if err != nil {
		return nil, err
	}
	tags := make(map[string]string)
	var nextToken *string
	for {
		out, err := client.DescribeFileSystems(context.Background(), &efs.DescribeFileSystemsInput{
			FileSystemId: aws.String(efsId),
			Marker:       nextToken,
		})
		if err != nil {
			return nil, err
		}
		for _, e := range out.FileSystems {
			for _, t := range e.Tags {
				tags[aws.ToString(t.Key)] = aws.ToString(t.Value)
			}
		}
		if out.Marker == nil {
			break
		}
		nextToken = out.Marker
	}
	return tags, nil
}
