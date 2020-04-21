package data

import (
	"fmt"
	"github.com/vmmgr/controller/db"
	"strconv"
	"strings"
)

func AddGroupUser(name, group string) (string, bool) {
	data, info, result := VerifyGroup(name, group)
	if result == false {
		return info, false
	}
	fmt.Println(data)
	d, result := ProcessStringToArray(data.User, name, 0)
	if result == false {
		fmt.Println("Error: User is exists this group")
		return "Error: User is exists this group", false
	}

	if db.ChangeDBGroupUser(data.ID, d) {
		return "ok", true
	}
	return "DBChangeError!!", false
}

func AddGroupAdmin(name, group string) (string, bool) {
	data, info, result := VerifyGroup(name, group)
	if result == false {
		return info, false
	}
	fmt.Println("Admin")

	d, result := ProcessStringToArray(data.Admin, name, 0)
	if result == false {
		fmt.Println("Error: User is exists this group")
		return "Error: User is exists this group", false
	}

	if db.ChangeDBGroupAdmin(data.ID, d) {
		return "ok", true
	}
	return "DBChangeError!!", false
}

func RemoveGroupAdmin(name, group string) (string, bool) {
	data, info, result := VerifyGroup(name, group)
	if result == false {
		return info, false
	}
	d, result := ProcessStringToArray(data.Admin, name, 1)
	if result == false {
		fmt.Println("Error: User is exists this group")
		return "Error: User is exists this group", false
	}

	if db.ChangeDBGroupAdmin(data.ID, d) {
		return "ok", true
	}
	return "DBChangeError!!", false
}

func RemoveGroupUser(name, group string) (string, bool) {
	data, info, result := VerifyGroup(name, group)
	if result == false {
		return info, false
	}
	d, result := ProcessStringToArray(data.User, name, 1)
	if result == false {
		fmt.Println("Error: User is exists this group")
		return "Error: User is exists this group", false
	}

	if db.ChangeDBGroupUser(data.ID, d) {
		return "ok", true
	}
	return "DBChangeError!!", false
}

func VerifyGroup(name, group string) (db.Group, string, bool) {
	_, result := db.GetDBUserID(name)
	if result == false {
		fmt.Println("Error: Not exists this User")
		return db.Group{Name: name}, "Not exists this User", false
	}
	id, result := db.GetDBGroupID(group)
	if result == false {
		fmt.Println("Error: Not exists this group")
		return db.Group{Name: name}, "Not exists this group", false
	}
	data, result := db.GetDBGroup(id)
	if result == false {
		fmt.Println("Error: Failed GetGroup")
		return db.Group{Name: name}, "Failed GetGroup", false
	}
	return data, "ok", true
}

//mode 0: create 1:delete
func ProcessStringToArray(basedata, data string, mode int) (string, bool) {
	fmt.Println(basedata)
	basedataarray := strings.Split(basedata, ",")
	fmt.Println("length: " + strconv.Itoa(len(basedataarray)))
	fmt.Println(basedataarray)
	if mode == 0 {
		if len(basedataarray) == 0 {
			var result []string
			result[0] = data
			return strings.Join(result, ","), true
		} else {
			//verify same user data
			for _, a := range basedataarray {
				if a == data {
					return "0", false
				}
			}
			basedataarray := append(basedataarray, data)
			return strings.Join(basedataarray, ","), true
		}
	} else if mode == 1 {
		if len(basedataarray) == 0 {
			fmt.Println("GroupData: No data!!")
			return "0", false
		}
		var dataarray []string
		for _, a := range basedataarray {
			if a != data {
				dataarray = append(dataarray, a)
			}
		}
		return strings.Join(dataarray, ","), true
	}
	return "0", false
}

func SearchUserForAllGroup(user string) ([]int, bool) {
	var r []int
	admindata, result := SearchUserForAdminGroup(user)
	if result == false {
		r[0] = 0
		return r, false
	}
	userdata, result := SearchUserForUserGroup(user)
	if result == false {
		r[0] = 0
		return r, false
	}

	tmp := admindata[:]
	tmp = append(tmp, userdata...)

	resultarray := make([]int, 0)

	m := make(map[int]struct{}, 0)
	for _, el := range tmp {
		if _, ok := m[el]; ok == false {
			m[el] = struct{}{}
			resultarray = append(resultarray, el)
		}
	}
	return resultarray, true
}

func SearchUserForAdminGroup(user string) ([]int, bool) {
	data := db.GetDBAllGroup()
	fmt.Println(data)
	var result []int
	for _, a := range data {
		dataarray := strings.Split(a.Admin, ",")
		for _, d := range dataarray {
			if d == user {
				result = append(result, a.ID)
				break
			}
		}
		fmt.Println(a)
	}
	fmt.Println("AdminGroup")
	fmt.Println(result)
	return result, true
}
func SearchUserForUserGroup(user string) ([]int, bool) {
	data := db.GetDBAllGroup()
	var result []int
	for _, a := range data {
		dataarray := strings.Split(a.User, ",")
		for _, d := range dataarray {
			if d == user {
				result = append(result, a.ID)
				break
			}
		}
	}
	fmt.Println("UserGroup")
	fmt.Println(result)
	return result, true
}

//basedata: group admin and user data
func SearchAllGroupUser(basedata, searchnamedata string) bool {
	dataarray := strings.Split(basedata, ",")
	fmt.Printf("groupuser&admin: ")
	fmt.Println(dataarray)
	for _, d := range dataarray {
		if d == searchnamedata {
			return true
		}
	}
	return false
}

//mode 0:admin 1:user
func SearchGroupUser(name, group string, mode int) bool {
	id, result := db.GetDBGroupID(group)
	if result == false {
		fmt.Println("Error: Not found group")
	}
	fmt.Println("GroupID: " + strconv.Itoa(id))
	data, result := db.GetDBGroup(id)
	if result == false {
		fmt.Println("Error: Not found group")
	}
	if mode == 0 {
		d := strings.Split(data.Admin, ",")
		for _, a := range d {
			if a == name {
				return true
			}
		}
	} else if mode == 1 {
		d := strings.Split(data.User, ",")
		for _, a := range d {
			if a == name {
				return true
			}
		}
	}
	return false
}
