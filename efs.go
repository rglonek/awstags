package main

import (
	"awstags/tags"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
)

func (s *efsSet) Execute(tail []string) error {
	if s.FileName == nil && !isInputFromPipe() {
		return errors.New("either specify a filename to read from, or cat/echo the json into this program")
	}
	creds, err := efsValidate(&s.efsItems)
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
	log.Printf("Setting tags for %s", *s.EfsId)
	err = tags.EfsSet(s.Region, creds, *s.EfsId, t)
	if err != nil {
		return err
	}
	log.Print("Done")
	return nil
}

func (s *efsGet) Execute(tail []string) error {
	creds, err := efsValidate(&s.efsItems)
	if err != nil {
		return err
	}
	log.Printf("Getting tags for %s", *s.EfsId)
	out, err := tags.EfsGet(s.Region, creds, *s.EfsId)
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

func (s *efsDelete) Execute(tail []string) error {
	creds, err := efsValidate(&s.efsItems)
	if err != nil {
		return err
	}
	log.Printf("Deleting all tags for %s", *s.EfsId)
	err = tags.EfsDelete(s.Region, creds, *s.EfsId)
	if err != nil {
		return err
	}
	log.Print("Done")
	return nil
}

func (s *efsList) Execute(tail []string) error {
	creds, err := efsValidateLogin(&s.login)
	if err != nil {
		return err
	}
	log.Print("Listing all efs")
	if !s.WithTags {
		items, err := tags.EfsList(s.Region, creds)
		if err != nil {
			return err
		}
		for _, item := range items {
			fmt.Println(item)
		}
		log.Print("Done")
		return nil
	}
	items, err := tags.EfsListWithTags(s.Region, creds)
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

func (s *efsRegions) Execute(tail []string) error {
	creds, err := efsValidateLogin(&s.login)
	if err != nil {
		return err
	}
	log.Print("Listing regions")
	items, err := tags.EfsRegions(s.Region, creds)
	if err != nil {
		return err
	}
	for _, item := range items {
		fmt.Println(item)
	}
	log.Print("Done")
	return nil
}

func efsValidate(i *efsItems) (creds *tags.Creds, err error) {
	if i.EfsId == nil {
		return nil, errors.New("efs-id is a required field")
	}
	return efsValidateLogin(&i.login)
}

func efsValidateLogin(i *login) (creds *tags.Creds, err error) {
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
