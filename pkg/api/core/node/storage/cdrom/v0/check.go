package v0

import (
	"fmt"
	cdrom "github.com/vmmgr/controller/pkg/api/core/node/storage/cdrom"
)

func inputCheck(input cdrom.Post) error {
	if input.Name == "" {
		return fmt.Errorf("invalid name")
	}
	if input.URL == "" {
		return fmt.Errorf("invalid url")
	}

	return nil
}
