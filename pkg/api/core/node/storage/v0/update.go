package v0

import (
	"github.com/vmmgr/controller/pkg/api/core/node/storage"
)

func updateAdminUser(input, replace storage.Storage) (storage.Storage, error) {

	//Name
	if input.Name != "" {
		replace.Name = input.Name
	}
	if input.Comment != "" {
		replace.Comment = input.Comment
	}
	if input.Path != "" {
		replace.Path = input.Path
	}
	if replace.AdminOnly != input.AdminOnly {
		replace.AdminOnly = input.AdminOnly
	}
	if replace.NodeID != input.NodeID {
		replace.NodeID = input.NodeID
	}
	if replace.MaxCapacity != input.MaxCapacity {
		replace.MaxCapacity = input.MaxCapacity
	}
	if replace.Type != input.Type {
		replace.Type = input.Type
	}

	return replace, nil
}
