package v0

import (
	"fmt"
	"github.com/vmmgr/controller/pkg/api/core/support"
)

func check(input support.FirstInput) error {
	if input.TicketID == 0 {
		return fmt.Errorf("no data: TicketID")
	}
	if input.Data == "" {
		return fmt.Errorf("no data: data")
	}

	return nil
}
