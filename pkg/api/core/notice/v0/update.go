package v0

import (
	"github.com/vmmgr/controller/pkg/api/core"
)

func updateAdminUser(input, replace core.Notice) (core.Notice, error) {

	//Title
	if input.Title != "" {
		replace.Title = input.Title
	}
	//Data
	if input.Data != "" {
		replace.Data = input.Data
	}

	// uint boolean
	//StartTime
	if input.StartTime != replace.StartTime {
		replace.StartTime = input.StartTime
	}
	//Everyone
	if input.Everyone != replace.Everyone {
		replace.Everyone = input.Everyone
	}
	//Important
	if input.Important != replace.Important {
		replace.Important = input.Important
	}
	//Fault
	if input.Fault != replace.Fault {
		replace.Fault = input.Fault
	}
	//Info
	if input.Info != replace.Info {
		replace.Info = input.Info
	}

	return replace, nil
}
