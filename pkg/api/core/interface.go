package core

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	Notice        []*Notice `json:"notice" gorm:"many2many:user_notice;"`
	Tokens        []*Token  `json:"tokens"`
	Group         *Group    `json:"group"`
	GroupID       *uint     `json:"group_id"`
	Name          string    `json:"name"`
	NameEn        string    `json:"name_en"`
	Email         string    `json:"email"`
	Pass          string    `json:"pass"`
	ExpiredStatus *uint     `json:"expired_status"`
	Level         uint      `json:"level"`
	MailVerify    *bool     `json:"mail_verify"`
	MailToken     string    `json:"mail_token"`
}

type Group struct {
	gorm.Model
	VMs       []*VM  `json:"vms"`
	Users     []User `json:"users"`
	Org       string `json:"org"`
	Status    uint   `json:"status"`
	Comment   string `json:"comment"`
	Vlan      uint   `json:"vlan"`
	Enable    *bool  `json:"enable"`
	MaxVM     uint   `json:"max_VM"`
	MaxCPU    uint   `json:"max_cpu"`
	MaxMemory uint   `json:"max_memory"`
}

type Token struct {
	gorm.Model
	ExpiredAt   time.Time `json:"expired_at"`
	UserID      *uint     `json:"user_id"`
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
	User      []User    `json:"user" gorm:"many2many:notice_user;"`
	Everyone  *bool     `json:"everyone"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Important *bool     `json:"important"`
	Fault     *bool     `json:"fault"`
	Info      *bool     `json:"info"`
	Title     string    `json:"title"`
	Data      string    `json:"data" gorm:"size:15000"`
}

type Ticket struct {
	gorm.Model
	GroupID       *uint  `json:"group_id"`
	UserID        *uint  `json:"user_id"`
	Chat          []Chat `json:"chat"`
	Request       *bool  `json:"request"`
	RequestReject *bool  `json:"request_reject"`
	Solved        *bool  `json:"solved"`
	Admin         *bool  `json:"admin"`
	Title         string `json:"title"`
	Group         Group  `json:"group"`
	User          User   `json:"user"`
}

type Chat struct {
	gorm.Model
	Ticket   Ticket `json:"ticket"`
	TicketID uint   `json:"ticket_id"`
	UserID   *uint  `json:"user_id"`
	Admin    bool   `json:"admin"`
	Data     string `json:"data" gorm:"size:10000"`
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
	Region   Region `json:"region"`
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
	Zone       Zone      `json:"zone"`
	Group      Group     `json:"group"`
	ZoneID     uint      `json:"zone_id"`
	GroupID    *uint     `json:"group_id"`
	AdminOnly  *bool     `json:"admin_only"`
	Name       string    `json:"name"`
	HostName   string    `json:"host_name"`
	IP         string    `json:"ip"`
	Port       uint      `json:"port"`
	User       string    `json:"user"`
	Pass       string    `json:"pass"`
	WsPort     uint      `json:"ws_port"`
	ManageNet  string    `json:"manage_net"`
	Mac        string    `json:"mac"`
	Machine    string    `json:"machine"`
	Emulator   string    `json:"emulator"`
	Comment    string    `json:"comment"`
	Enable     *bool     `json:"enable"`
	PrimaryNIC string    `json:"primary_nic"`
	Storage    []Storage `json:"storage"`
	NIC        []NIC     `json:"nic"`
}

type ImaCon struct {
	gorm.Model
	Zone       Zone    `json:"zone"`
	ZoneID     uint    `json:"zone_id"`
	Name       string  `json:"name"`
	HostName   string  `json:"host_name"`
	IP         string  `json:"ip"`
	User       string  `json:"user"`
	Pass       string  `json:"pass"`
	Port       uint    `json:"port"`
	Enable     *bool   `json:"enable"`
	AppPath    string  `json:"app_path"`
	ConfigPath string  `json:"config_path"`
	Image      []Image `json:"image"`
}

type Image struct {
	gorm.Model
	ImaConID  uint        `json:"ima_con_id"`
	GroupID   *uint       `json:"group_id"` //nil: All 1~: Only Group
	Type      uint        `json:"type"`     //0: ISO 1:Image
	Path      string      `json:"path"`     //node側のパス
	UUID      string      `json:"uuid"`
	Name      string      `json:"name"`
	CloudInit *bool       `json:"cloud_init"` //cloud-init対応イメージであるか否か
	MinCPU    uint        `json:"min_cpu"`
	MinMem    uint        `json:"min_mem"`
	OS        string      `json:"os"`
	Admin     *bool       `json:"admin"` //管理者専用イメージであるか否か
	Lock      *bool       `json:"lock"`  //削除保護
	ImaCon    *ImaCon     `json:"imacon"`
	Group     *Group      `json:"group"`
	Template  *[]Template `json:"template"`
}

type Template struct {
	gorm.Model
	Name         string          `json:"name"`
	Tag          string          `json:"tag"`
	ImageID      string          `json:"image_id"`
	Image        *Image          `json:"image"`
	TemplatePlan []*TemplatePlan `json:"template_plan"`
}

type TemplatePlan struct {
	gorm.Model
	TemplateID uint      `json:"template_id"`
	CPU        uint      `json:"cpu"`
	Mem        uint      `json:"mem"`
	Storage    uint      `json:"storage"`
	Hide       *bool     `json:"hide"` //管理者専用イメージであるか否か
	Template   *Template `json:"template"`
}

// Type:  1:SSD 2:HDD 3:NVMe 11:SSD(iSCSI) 12:HDD(iSCSI) 13:NVme(iSCSI)
type Storage struct {
	gorm.Model
	NodeID      uint   `json:"node_id"`
	AdminOnly   *bool  `json:"admin"`
	Name        string `json:"name"`
	Type        uint   `json:"type"`
	Path        string `json:"path"`
	MaxCapacity uint   `json:"max_capacity"`
	Comment     string `json:"comment"`
	Node        Node   `json:"node"`
}

type NIC struct {
	gorm.Model
	NodeID    uint   `json:"node_id"`
	GroupID   uint   `json:"group_id"`
	AdminOnly *bool  `json:"admin"`
	Name      string `json:"name"`
	Enable    *bool  `json:"enable"`
	Virtual   *bool  `json:"virtual"`
	Type      uint   `json:"type"` //0-9: Global 10-19: Private
	Vlan      uint   `json:"vlan"`
	Speed     uint   `json:"speed"`
	MAC       string `json:"mac"`
	Comment   string `json:"comment"`
}

type VM struct {
	gorm.Model
	Node          Node   `json:"node"`
	IP            IP     `json:"ip"`
	NodeID        uint   `json:"node_id"`
	GroupID       *uint  `json:"group_id"`
	Name          string `json:"name"`
	UUID          string `json:"uuid"`
	VNCPort       *uint  `json:"vnc_port"`
	WebSocketPort *uint  `json:"web_socket_port"`
	Lock          *bool  `json:"lock"`
}

type IP struct {
	gorm.Model
	VMID     *uint  `json:"vm_id"`
	IP       string `json:"ip"`
	Subnet   string `json:"subnet"`
	Gateway  string `json:"gateway"`
	DNS      string `json:"dns"`
	Reserved *bool  `json:"reserved"`
}
