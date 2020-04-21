package data

import (
	"fmt"
	"github.com/vmmgr/controller/db"
	"github.com/vmmgr/controller/proto/proto-go"
	"strconv"
)

func CreateVMCheck(d *grpc.VMData) bool {
	if d.Vcpu < 0 || d.Vmem < 0 || d.Option.Vnc < 0 {
		return false
	}
	if d.Vmname == "" || d.Option.StoragePath == "" {
		return false
	}
	return true
}

func ExistNodeCheck(name, ip string) (string, bool) {
	data := db.GetDBAllNode()
	for _, a := range data {
		if a.HostName == name {
			return "Exists HostName", true
		}
		if a.IP == ip {
			return "Exists IP", true
		}
	}
	return "not found", false
}

func ExistGroupCheck(name string) bool {
	data := db.GetDBAllGroup()
	for _, a := range data {
		if a.Name == name {
			return true
		}
	}
	return false
}

func ExistUserCheck(name string) bool {
	data := db.GetDBAllUser()
	fmt.Println("UserCount: " + strconv.Itoa(len(data)))
	for _, a := range data {
		if a.Name == name {
			fmt.Println("Check OK(Exists User)")
			return true
		}
	}
	return false
}

//Exists All Group User
func GroupAllUserCheck(name string) bool {
	data := db.GetDBAllGroup()
	for i, _ := range data {
		if SearchAllGroupUser(data[i].Admin, name) {
			if SearchAllGroupUser(data[i].User, name) {
				fmt.Println("Exists group user.(Admin,User)")
				return true
			}
		}
	}
	return false
}

func CheckOnlyAdmin(nodeid int) bool {
	d, result := db.GetDBNodeID(nodeid)
	if result == false {
		fmt.Println("DB Connection Error!!")
		return false
	}
	if d.OnlyAdmin == 0 {
		return false
	} else {
		return true
	}
}
