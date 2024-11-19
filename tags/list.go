package tags

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/efs"
)

func Ec2List(region *string, creds *Creds) ([]string, error) {
	client, err := getEc2Client(creds, region)
	if err != nil {
		return nil, err
	}
	instances := []string{}
	var nextToken *string
	for {
		out, err := client.DescribeInstances(context.Background(), &ec2.DescribeInstancesInput{
			NextToken: nextToken,
		})
		if err != nil {
			return nil, err
		}
		for _, r := range out.Reservations {
			for _, i := range r.Instances {
				instances = append(instances, aws.ToString(i.InstanceId))
			}
		}
		if out.NextToken == nil {
			break
		}
		nextToken = out.NextToken
	}
	return instances, nil
}

func EfsList(region *string, creds *Creds) ([]string, error) {
	client, err := getEfsClient(creds, region)
	if err != nil {
		return nil, err
	}
	vols := []string{}
	var nextToken *string
	for {
		out, err := client.DescribeFileSystems(context.Background(), &efs.DescribeFileSystemsInput{
			Marker: nextToken,
		})
		if err != nil {
			return nil, err
		}
		for _, e := range out.FileSystems {
			vols = append(vols, aws.ToString(e.FileSystemId))
		}
		if out.Marker == nil {
			break
		}
		nextToken = out.Marker
	}
	return vols, nil
}

func Ec2ListWithTags(region *string, creds *Creds) (map[string]map[string]string, error) {
	client, err := getEc2Client(creds, region)
	if err != nil {
		return nil, err
	}
	instances := make(map[string]map[string]string)
	var nextToken *string
	for {
		out, err := client.DescribeInstances(context.Background(), &ec2.DescribeInstancesInput{
			NextToken: nextToken,
		})
		if err != nil {
			return nil, err
		}
		for _, r := range out.Reservations {
			for _, i := range r.Instances {
				instances[aws.ToString(i.InstanceId)] = make(map[string]string)
				for _, t := range i.Tags {
					instances[aws.ToString(i.InstanceId)][aws.ToString(t.Key)] = aws.ToString(t.Value)
				}
			}
		}
		if out.NextToken == nil {
			break
		}
		nextToken = out.NextToken
	}
	return instances, nil
}

func EfsListWithTags(region *string, creds *Creds) (map[string]map[string]string, error) {
	client, err := getEfsClient(creds, region)
	if err != nil {
		return nil, err
	}
	vols := make(map[string]map[string]string)
	var nextToken *string
	for {
		out, err := client.DescribeFileSystems(context.Background(), &efs.DescribeFileSystemsInput{
			Marker: nextToken,
		})
		if err != nil {
			return nil, err
		}
		for _, e := range out.FileSystems {
			vols[aws.ToString(e.FileSystemId)] = make(map[string]string)
			for _, t := range e.Tags {
				vols[aws.ToString(e.FileSystemId)][aws.ToString(t.Key)] = aws.ToString(t.Value)
			}
		}
		if out.Marker == nil {
			break
		}
		nextToken = out.Marker
	}
	return vols, nil
}
