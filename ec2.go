package main

import (
	"awstags/tags"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
)

func (s *ec2Set) Execute(tail []string) error {
	if s.FileName == nil && !isInputFromPipe() {
		return errors.New("either specify a filename to read from, or cat/echo the json into this program")
	}
	creds, err := ec2Validate(&s.ec2Items)
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

func (s *ec2Get) Execute(tail []string) error {
	creds, err := ec2Validate(&s.ec2Items)
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

func (s *ec2Delete) Execute(tail []string) error {
	creds, err := ec2Validate(&s.ec2Items)
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

func (s *ec2List) Execute(tail []string) error {
	creds, err := ec2ValidateLogin(&s.login)
	if err != nil {
		return err
	}
	log.Print("Listing all instances")
	if !s.WithTags {
		items, err := tags.Ec2List(s.Region, creds)
		if err != nil {
			return err
		}
		for _, item := range items {
			fmt.Println(item)
		}
		log.Print("Done")
		return nil
	}
	items, err := tags.Ec2ListWithTags(s.Region, creds)
	if err != nil {
		return err
	}
	err = json.NewEncoder(os.Stdout).Encode(items)
	if err != nil {
		return err
	}
	log.Print("Done")
	return nil
}

func ec2Validate(i *ec2Items) (creds *tags.Creds, err error) {
	if i.InstanceId == nil {
		return nil, errors.New("instance-id is a required field")
	}
	return ec2ValidateLogin(&i.login)
}

func ec2ValidateLogin(i *login) (creds *tags.Creds, err error) {
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
