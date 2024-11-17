package main

import (
	"awstags/tags"
	"encoding/json"
	"errors"
	"log"
	"os"

	"github.com/jessevdk/go-flags"
)

type cli struct {
	Set    *cliSet    `command:"set" description:"set tags from file/stdin json"`
	Get    *cliGet    `command:"get" description:"get tags of an instance, and print to stdout to json"`
	Delete *cliDelete `command:"delete" description:"delete all tags of an instance"`
}

type cliSet struct {
	items
	FileName *string `short:"f" long:"filename" description:"filename of the file to read tags from; will read os.Stdin if not specified"`
}

type cliGet struct {
	items
}

type cliDelete struct {
	items
}

type items struct {
	ProfileName *string `short:"p" long:"profile-name" description:"login using a specific shared credentials profile name"`
	KeyID       *string `short:"k" long:"key-id" description:"login using a specific keyId"`
	SecretKey   *string `short:"s" long:"secret-key" description:"login using a specific secretKey"`
	Region      *string `short:"r" long:"region" description:"use a specific AWS region"`
	InstanceId  *string `short:"i" long:"instance-id" description:"required: instance-id to query/set"`
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

func (s *cliSet) Execute(tail []string) error {
	if s.FileName == nil && !isInputFromPipe() {
		return errors.New("either specify a filename to read from, or cat/echo the json into this program")
	}
	creds, err := validate(&s.items)
	if err != nil {
		return err
	}
	t := make(map[string]string)
	if s.FileName != nil {
		data, err := os.ReadFile(*s.FileName)
		if err != nil {
			return err
		}
		err = json.Unmarshal(data, &t)
		if err != nil {
			return err
		}
	} else {
		err = json.NewDecoder(os.Stdin).Decode(&t)
		if err != nil {
			return err
		}
	}
	log.Printf("Setting tags for %s", *s.InstanceId)
	err = tags.Ec2Set(s.Region, creds, *s.InstanceId, t)
	if err != nil {
		return err
	}
	log.Print("Done")
	return nil
}

func (s *cliGet) Execute(tail []string) error {
	creds, err := validate(&s.items)
	if err != nil {
		return err
	}
	log.Printf("Getting tags for %s", *s.InstanceId)
	out, err := tags.Ec2Get(s.Region, creds, *s.InstanceId)
	if err != nil {
		return err
	}
	err = json.NewEncoder(os.Stdout).Encode(out)
	if err != nil {
		return err
	}
	log.Print("Done")
	return nil
}

func (s *cliDelete) Execute(tail []string) error {
	creds, err := validate(&s.items)
	if err != nil {
		return err
	}
	log.Printf("Deleting all tags for %s", *s.InstanceId)
	err = tags.Ec2Delete(s.Region, creds, *s.InstanceId)
	if err != nil {
		return err
	}
	log.Print("Done")
	return nil
}

func validate(i *items) (creds *tags.Creds, err error) {
	if i.InstanceId == nil {
		return nil, errors.New("instance-id is a required field")
	}
	if i.KeyID != nil && i.SecretKey == nil {
		return nil, errors.New("when specifying key-id, secret-key must also be set")
	}
	if i.SecretKey != nil && i.KeyID == nil {
		return nil, errors.New("when specifying secret-key, key-id must also be set")
	}
	if i.KeyID != nil && i.ProfileName != nil {
		return nil, errors.New("specify either a profile-name or static key-id credential, not both")
	}
	if i.KeyID != nil {
		creds = &tags.Creds{
			Static: &tags.CredsStatic{
				Key:    *i.KeyID,
				Secret: *i.SecretKey,
			},
		}
	}
	if i.ProfileName != nil {
		creds = &tags.Creds{
			ProfileName: i.ProfileName,
		}
	}
	return creds, nil
}