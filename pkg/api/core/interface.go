package core

import (
	"github.com/jinzhu/gorm"
	"time"
)

type User struct {
	gorm.Model
	Tokens        []*Token `json:"tokens"`
	Group         *Group   `json:"group"`
	GroupID       uint     `json:"group_id"`
	Name          string   `json:"name"`
	NameEn        string   `json:"name_en"`
	Email         string   `json:"email"`
	Pass          string   `json:"pass"`
	ExpiredStatus *uint    `json:"expired_status"`
	Level         uint     `json:"level"`
	MailVerify    *bool    `json:"mail_verify"`
	MailToken     string   `json:"mail_token"`
}

type Group struct {
	gorm.Model
	Users     []User `json:"users"`
	Org       string `json:"org"`
	Status    uint   `json:"status"`
	Comment   string `json:"comment"`
	Vlan      uint   `json:"vlan"`
	Lock      bool   `json:"lock"`
	MaxVM     uint   `json:"max_VM"`
	MaxCPU    uint   `json:"max_cpu"`
	MaxMemory uint   `json:"max_memory"`
}

type Token struct {
	gorm.Model
	ExpiredAt   time.Time `json:"expired_at"`
	UserID      uint      `json:"user_id"`
	User        User      `json:"user"`
	Status      uint      `json:"status"` //0: initToken(30m) 1: 30m 2:6h 3: 12h 10: 30d 11:180d
	Admin       *bool     `json:"admin"`
	UserToken   string    `json:"user_token"`
	TmpToken    string    `json:"tmp_token"`
	AccessToken string    `json:"access_token"`
	Debug       string    `json:"debug"`
}

type Notice struct {
	gorm.Model
	UserID     uint   `json:"user_id"`
	GroupID    uint   `json:"group_id"`
	Everyone   *bool  `json:"everyone"`
	StartTime  uint   `json:"start_time"`
	EndingTime uint   `json:"ending_time"`
	Important  *bool  `json:"important"`
	Fault      *bool  `json:"fault"`
	Info       *bool  `json:"info"`
	Title      string `json:"title"`
	Data       string `json:"data"`
}

type Ticket struct {
	gorm.Model
	GroupID uint   `json:"group_id"`
	UserID  uint   `json:"user_id"`
	Chat    []Chat `json:"chat"`
	Solved  *bool  `json:"solved"`
	Title   string `json:"title"`
	Group   Group  `json:"group"`
	User    User   `json:"user"`
}

type Chat struct {
	gorm.Model
	TicketID uint   `json:"ticket_id"`
	UserID   uint   `json:"user_id"`
	Admin    bool   `json:"admin"`
	Data     string `json:"data" gorm:"size:65535"`
	User     User   `json:"user"`
}

type Region struct {
	gorm.Model
	Name    string `json:"name"`
	Comment string `json:"comment"`
	Lock    *bool  `json:"lock"`
}

type Zone struct {
	gorm.Model
	RegionID uint   `json:"region_id"`
	Name     string `json:"name"`
	Postcode string `json:"postcode"`
	Address  string `json:"address"`
	Tel      string `json:"tel"`
	Mail     string `json:"mail"`
	Comment  string `json:"comment"`
	Lock     *bool  `json:"lock"`
}

type Node struct {
	gorm.Model
	ZoneID    uint      `json:"zone_id"`
	GroupID   uint      `json:"group_id"`
	AdminOnly *bool     `json:"admin_only"`
	Name      string    `json:"name"`
	IP        string    `json:"ip"`
	Port      uint      `json:"port"`
	WsPort    uint      `json:"ws_port"`
	ManageNet string    `json:"manage_net"`
	Mac       string    `json:"mac"`
	Machine   string    `json:"machine"`
	Emulator  string    `json:"emulator"`
	Comment   string    `json:"comment"`
	Storage   []Storage `json:"storage"`
	NIC       []NIC     `json:"nic"`
}

type ImaCon struct {
	gorm.Model
	ZoneID uint    `json:"zone_id"`
	Name   string  `json:"name"`
	IP     string  `json:"ip"`
	Port   uint    `json:"port"`
	Image  []Image `json:"image"`
}

type Image struct {
	gorm.Model
	ImaConID  uint       `json:"ima_con_id"`
	GroupID   uint       `json:"group_id"` //0: All 1~: Only Group
	Type      uint       `json:"type"`     //0: ISO 1:Image
	Path      string     `json:"path"`     //node側のパス
	UUID      string     `json:"uuid"`
	Name      string     `json:"name"`
	CloudInit *bool      `json:"cloud_init"` //cloud-init対応イメージであるか否か
	MinCPU    uint       `json:"min_cpu"`
	MinMem    uint       `json:"min_mem"`
	OS        string     `json:"os"`
	Admin     *bool      `json:"admin"` //管理者専用イメージであるか否か
	Lock      *bool      `json:"lock"`  //削除保護
	Template  []Template `json:"template"`
}

type Template struct {
	gorm.Model
	Name         string         `json:"name"`
	Tag          string         `json:"tag"`
	ImageID      string         `json:"image_id"`
	Image        Image          `json:"image"`
	TemplatePlan []TemplatePlan `json:"template_plan"`
}

type TemplatePlan struct {
	gorm.Model
	TemplateID uint  `json:"template_id"`
	CPU        uint  `json:"cpu"`
	Mem        uint  `json:"mem"`
	Storage    uint  `json:"storage"`
	Hide       *bool `json:"hide"` //管理者専用イメージであるか否か
}

type Storage struct {
	gorm.Model
	NodeID      uint   `json:"node_id"`
	AdminOnly   *bool  `json:"admin"`
	Name        string `json:"name"`
	Type        uint   `json:"type"`
	Path        string `json:"path"`
	MaxCapacity uint   `json:"max_capacity"`
	Comment     string `json:"comment"`
}

type NIC struct {
	gorm.Model
	NodeID    uint   `json:"node_id"`
	GroupID   uint   `json:"group_id"`
	AdminOnly *bool  `json:"admin"`
	Name      string `json:"name"`
	Enable    *bool  `json:"enable"`
	Virtual   *bool  `json:"virtual"`
	Type      uint   `json:"type"`
	Vlan      uint   `json:"vlan"`
	Speed     uint   `json:"speed"`
	MAC       string `json:"mac"`
	Comment   string `json:"comment"`
}

type VM struct {
	gorm.Model
	NodeID  uint   `json:"node_id"`
	GroupID uint   `json:"group_id"`
	Name    string `json:"name"`
	UUID    string `json:"uuid"`
	VNCPort uint   `json:"vnc_port"`
	Lock    *bool  `json:"lock"`
}
