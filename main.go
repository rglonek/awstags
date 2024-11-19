package main

import (
	"fmt"
	"os"

	_ "embed"

	"github.com/jessevdk/go-flags"
)

//go:embed version.txt
var toolVersion string

type cli struct {
	Ec2 struct {
		Set    *ec2Set    `command:"set" description:"set tags from file/stdin json"`
		Get    *ec2Get    `command:"get" description:"get tags of an instance, and print to stdout to json"`
		Delete *ec2Delete `command:"delete" description:"delete all tags of an instance"`
		List   *ec2List   `command:"list" description:"list all ec2 instances"`
	} `command:"ec2" description:"operate on ec2 tags"`
	Efs struct {
		Set    *efsSet    `command:"set" description:"set tags from file/stdin json"`
		Get    *efsGet    `command:"get" description:"get tags of an efs, and print to stdout to json"`
		Delete *efsDelete `command:"delete" description:"delete all tags of an efs"`
		List   *efsList   `command:"list" description:"list all efs instances"`
	} `command:"efs" description:"operate on efs tags"`
	Version *version `command:"version" description:"print tool version"`
}

type version struct{}

type ec2Set struct {
	ec2Items
	FileName *string `short:"f" long:"filename" description:"filename of the file to read tags from; will read os.Stdin if not specified"`
}

type ec2Get struct {
	ec2Items
}

type ec2Delete struct {
	ec2Items
}

type ec2List struct {
	login
	WithTags bool `short:"t" long:"with-tags" description:"list instances with tags - produces long json"`
}

type login struct {
	ProfileName *string `short:"p" long:"profile-name" description:"login using a specific shared credentials profile name"`
	KeyID       *string `short:"k" long:"key-id" description:"login using a specific keyId"`
	SecretKey   *string `short:"s" long:"secret-key" description:"login using a specific secretKey"`
	Region      *string `short:"r" long:"region" description:"use a specific AWS region"`
}

type ec2Items struct {
	login
	InstanceId *string `short:"i" long:"instance-id" description:"required: instance-id to query/set"`
}

type efsSet struct {
	efsItems
	FileName *string `short:"f" long:"filename" description:"filename of the file to read tags from; will read os.Stdin if not specified"`
}

type efsGet struct {
	efsItems
}

type efsDelete struct {
	efsItems
}

type efsList struct {
	login
	WithTags bool `short:"t" long:"with-tags" description:"list efs with tags - produces long json"`
}

type efsItems struct {
	login
	EfsId *string `short:"e" long:"efs-id" description:"required: efs-id to query/set"`
}

func main() {
	c := &cli{}
	_, err := flags.Parse(c)
	if err != nil {
		os.Exit(1)
	}
}

func isInputFromPipe() bool {
	fileInfo, _ := os.Stdin.Stat()
	return fileInfo.Mode()&os.ModeCharDevice == 0
}

func (s *version) Execute(tail []string) error {
	fmt.Println(toolVersion)
	return nil
}
