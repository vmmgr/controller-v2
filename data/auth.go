package data

import (
	"context"
	"fmt"
	"github.com/vmmgr/controller/db"
	pb "github.com/vmmgr/controller/proto/proto-go"
	"google.golang.org/grpc"
	"strconv"
	"strings"
	"time"
)

type UserCertData struct {
	User    string
	Pass    string
	Group   string
	Token   string
	VMID    int
	GroupID int
	NodeID  int
}

//Token Authentication
func TokenCertification(token string) (int, []int, bool) {
	go DeleteExpiredToken()
	data, result := db.GetDBToken(token)
	if result == false {
		return 0, nil, false
	}
	if VerifyTime(data.Endtime) == false {
		return 1, nil, false
	}
	return data.Userid, nil, true
}

//User Certification Tool is testing now !!
func AdminUserCertification(name, pass, token string) bool {
	if token != "" {
		tokendata, result := db.GetDBToken(token)
		if result == false {
			fmt.Println("Error: Token Get Error 1")
			fmt.Println("Certification NG!! (Administrator)")
			return false
		}

		userdata, result := db.GetDBUser(tokendata.Userid)
		if result == false {
			fmt.Println("Error: Token Get Error 2")
			fmt.Println("Certification NG!! (Administrator)")
			return false
		}
		name = userdata.Name
	} else {
		if db.PassAuthDBUser(name, pass) == false {
			fmt.Println("Certification NG!! (Administrator)")
			return false
		}
	}

	if SearchGroupUser(name, "admin", 0) {
		fmt.Println("Certification OK!! (Administrator)")
		return true
	}
	fmt.Println("Certification NG!! (Administrator)")
	return false
}

func UserCertification(name, pass string) bool {
	if db.PassAuthDBUser(name, pass) {
		fmt.Println("Certification OK!! (User)")
		return true
	}
	fmt.Println("Certification NG!! (User)")
	return false
}

func GroupAdminCertification(name, pass, token, group string) bool {
	if token != "" {
		tokendata, result := db.GetDBToken(token)
		if result == false {
			fmt.Println("Error: Token Get Error 1")
			fmt.Println("Certification NG!! (GroupAdmin)")
			return false
		}

		userdata, result := db.GetDBUser(tokendata.Userid)
		if result == false {
			fmt.Println("Error: Token Get Error 2")
			fmt.Println("Certification NG!! (GroupAdmin)")
			return false
		}
		name = userdata.Name

	} else {
		if db.PassAuthDBUser(name, pass) == false {
			fmt.Println("Certification NG!! (GroupAdmin)")
			return false
		}
	}
	if SearchGroupUser(name, group, 0) {
		fmt.Println("Certification OK!! (GroupAdmin)")
		return true
	}

	fmt.Println("Certification NG!! (GroupAdmin)")
	return false
}

func GroupUserCertification(name, pass, token, group string) bool {
	if token != "" {
		tokendata, result := db.GetDBToken(token)
		if result == false {
			fmt.Println("Error: Token Get Error 1")
			fmt.Println("Certification NG!! (GroupUser)")
			return false
		}

		userdata, result := db.GetDBUser(tokendata.Userid)
		if result == false {
			fmt.Println("Error: Token Get Error 2")
			fmt.Println("Certification NG!! (GroupUser)")
			return false
		}
		name = userdata.Name

	} else {
		if db.PassAuthDBUser(name, pass) == false {
			fmt.Println("Certification NG!! (GroupUser)")
			return false
		}
	}
	if SearchGroupUser(name, group, 1) {
		fmt.Println("Certification OK!! (GroupUser)")
		return true
	}

	fmt.Println("Certification NG!! (GroupUser)")
	return false
}

//VM Certification
//
// 10-0-testVM
// VM verify group
func VMCertification(vmid, groupid int, address string) (string, bool) {
	if vmid < 1 {
		fmt.Println("Value False" + strconv.Itoa(vmid))
		fmt.Printf("Debug: value is ")
		return "Error VMID", false
	}

	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		fmt.Println("Not connect; %v", err)
	}
	defer conn.Close()
	c := pb.NewGrpcClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.GetVM(ctx, &pb.VMID{Id: int64(vmid)})
	if err != nil {
		fmt.Println("could not greet: %v", err)
	}
	if r.GetVmname() == "" {
		return "None", false
	}
	d := strings.Split(r.GetVmname(), "-")
	if d[0] == strconv.Itoa(groupid) {
		fmt.Println("Certification OK!!")
		return "ok", true
	}
	return "None", false
}

func SuperUserCertification(d *UserCertData) bool {
	if AdminUserCertification(d.User, d.Pass, d.Token) == false {
		if GroupAdminCertification(d.User, d.Pass, d.Token, d.Group) == false {
			return false
		}
	}
	return true
}

func StandardUserCertification(d *UserCertData) bool {
	if AdminUserCertification(d.User, d.Pass, d.Token) == false {
		if GroupAdminCertification(d.User, d.Pass, d.Token, d.Group) == false {
			if GroupUserCertification(d.User, d.Pass, d.Token, d.Group) == false {
				return false
			}
		}
	}
	return true
}

//Only VM Certification

func SuperUserVMCertification(d *UserCertData) (string, bool) {
	var username string
	if d.Token == "" {
		username = d.User
		if UserCertification(d.User, d.Pass) == false {
			return "", false
		}
	} else {
		userid, _, result := TokenCertification(d.Token)
		if result == false {
			return "", false
		}
		userdata, result := db.GetDBUser(userid)
		if result == false {
			return "", false
		}
		username = userdata.Name
	}

	//NodeIP
	n, result := db.GetDBNodeID(d.NodeID)
	if result == false {
		fmt.Println("Error: Get NodeData Error")
		return "", false
	}

	conn, err := grpc.Dial(n.IP, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		fmt.Printf("Not connect; ")
		fmt.Println(err)
	}
	defer conn.Close()
	c := pb.NewGrpcClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	g, err := c.GetVM(ctx, &pb.VMID{Id: int64(d.VMID)})
	if err != nil {
		fmt.Println("could not greet: %v", err)
	}
	d2 := strings.Split(g.GetVmname(), "-")
	groupid := d2[0]

	if AdminUserCertification(d.User, d.Pass, d.Token) == false {
		groupdata, result := SearchUserForAdminGroup(username)
		if result == false {
			return "", false
		}

		for _, a := range groupdata {
			if strconv.Itoa(a) == groupid {
				return n.IP, true
			}
		}
	} else {
		return n.IP, true
	}

	return "", true
}

func StandardUserVMCertification(d *UserCertData) (string, bool) {
	var username string
	if d.Token == "" {
		username = d.User
		if UserCertification(d.User, d.Pass) == false {
			return "", false
		}
	} else {
		userid, _, result := TokenCertification(d.Token)
		if result == false {
			return "", false
		}
		userdata, result := db.GetDBUser(userid)
		if result == false {
			return "", false
		}
		username = userdata.Name
	}

	//NodeIP
	n, result := db.GetDBNodeID(d.NodeID)
	if result == false {
		fmt.Println("Error: Get NodeData Error")
		return "NodeDataError", false
	}

	conn, err := grpc.Dial(n.IP, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		fmt.Printf("Not connect; ")
		fmt.Println(err)
	}
	defer conn.Close()
	c := pb.NewGrpcClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	g, err := c.GetVM(ctx, &pb.VMID{Id: int64(d.VMID)})
	if err != nil {
		fmt.Println("could not greet: %v", err)
	}
	d2 := strings.Split(g.GetVmname(), "-")
	groupid := d2[0]

	if AdminUserCertification(d.User, d.Pass, d.Token) == false {
		groupdata, result := SearchUserForAllGroup(username)
		if result == false {
			return "", false
		}

		for _, a := range groupdata {
			if strconv.Itoa(a) == groupid {
				return n.IP, true
			}
		}
	} else {
		return n.IP, true
	}
	return "", false
}

func GetUserName(name, token string) string {
	if token != "" {
		tokendata, result := db.GetDBToken(token)
		if result == false {
			fmt.Println("Error: Token Get Error 1")
			return ""
		}

		userdata, result := db.GetDBUser(tokendata.Userid)
		if result == false {
			fmt.Println("Error: Token Get Error 2")
			return ""
		}
		name = userdata.Name
	}
	return name
}
