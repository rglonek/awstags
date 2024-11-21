package tags

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

func Ec2Regions(region *string, creds *Creds) ([]string, error) {
	client, err := getEc2Client(creds, region)
	if err != nil {
		return nil, err
	}
	regions := []string{}
	out, err := client.DescribeRegions(context.Background(), &ec2.DescribeRegionsInput{
		AllRegions: aws.Bool(false),
	})
	if err != nil {
		return nil, err
	}
	for _, r := range out.Regions {
		regions = append(regions, aws.ToString(r.RegionName))
	}
	return regions, nil
}

func EfsRegions(region *string, creds *Creds) ([]string, error) {
	client, err := getEc2Client(creds, region)
	if err != nil {
		return nil, err
	}
	regions := []string{}
	out, err := client.DescribeRegions(context.Background(), &ec2.DescribeRegionsInput{
		AllRegions: aws.Bool(false),
	})
	if err != nil {
		return nil, err
	}
	for _, r := range out.Regions {
		regions = append(regions, aws.ToString(r.RegionName))
	}
	return regions, nil
}
