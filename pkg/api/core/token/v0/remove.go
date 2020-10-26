package v0

import (
	"github.com/jinzhu/gorm"
	"github.com/vmmgr/controller/pkg/api/core/token"
	dbToken "github.com/vmmgr/controller/pkg/api/store/token/v0"
	"log"
	"time"
)

func TokenRemove() {
	go func() {
		// 15分おき
		t := time.NewTicker(15 * 60 * time.Second)
		for {
			select {
			case <-t.C:
				result := dbToken.Get(token.ExpiredTime, &token.Token{})
				if result.Err != nil {
					log.Println(result.Err)
				}
				for _, tmp := range result.Token {
					dbToken.Delete(&token.Token{Model: gorm.Model{ID: tmp.ID}})
				}
			}
		}
		t.Stop() //停止
	}()
}
