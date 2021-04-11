package v0

import (
	"fmt"
	"github.com/vmmgr/controller/pkg/api/core/node/storage/vmImage"
)

func inputCheck(input vmImage.Post) error {
	if input.Name == "" {
		return fmt.Errorf("invalid name")
	}
	if input.URL == "" {
		return fmt.Errorf("invalid url")
	}

	return nil
}
