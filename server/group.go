package server

import (
	"context"
	"fmt"
	"github.com/vmmgr/controller/data"
	"github.com/vmmgr/controller/db"
	pb "github.com/vmmgr/controller/proto/proto-go"
	"log"
	"strconv"
)

func (s *server) AddGroup(ctx context.Context, in *pb.GroupData) (*pb.Result, error) {
	log.Println("----AddGroup----")
	log.Println("Receive GroupName       : " + in.GetName())
	log.Println("Receive GroupMaxVM      : " + strconv.Itoa(int(in.GetSepc().GetMaxvm())))
	log.Println("Receive GroupMaxCPU     : " + strconv.Itoa(int(in.GetSepc().GetMaxcpu())))
	log.Println("Receive GroupMaxMem     : " + strconv.Itoa(int(in.GetSepc().GetMaxmem())))
	log.Println("Receive GroupMaxStorage : " + strconv.Itoa(int(in.GetSepc().GetMaxstorage())))
	log.Println("Receive GroupNet        : " + in.GetSepc().GetNet())
	log.Println("Receive AuthUser        : " + in.Base.User + ", AuthPass: " + in.Base.Pass)
	log.Println("Receive Token           : " + in.Base.GetToken())

	if data.AdminUserCertification(in.Base.GetUser(), in.Base.GetPass(), in.Base.GetToken()) == false {
		return &pb.Result{Status: false, Info: "Authentication failed!!"}, nil
	}
	if data.ExistGroupCheck(in.GetName()) {
		return &pb.Result{Status: false, Info: "Exists Group!!"}, nil
	}
	if db.AddDBGroup(db.Group{Name: in.GetName(), Admin: "", User: "", MaxVM: int(in.GetSepc().GetMaxvm()), MaxCPU: int(in.GetSepc().GetMaxcpu()), MaxMem: int(in.GetSepc().GetMaxmem()), MaxStorage: int(in.GetSepc().GetMaxstorage()), Net: in.GetSepc().GetNet()}) {
		return &pb.Result{Status: true, Info: "OK!"}, nil
	} else {
		return &pb.Result{Status: false, Info: "DB Error!!"}, nil
	}
}

func (s *server) RemoveGroup(ctx context.Context, in *pb.GroupData) (*pb.Result, error) {
	log.Println("----RemoveGroup----")
	log.Println("Receive GroupName       : " + in.GetName())
	log.Println("Receive AuthUser        : " + in.Base.User + ", AuthPass: " + in.Base.Pass)
	log.Println("Receive Token           : " + in.Base.GetToken())

	if data.AdminUserCertification(in.Base.GetUser(), in.Base.GetPass(), in.Base.GetToken()) == false {
		return &pb.Result{Status: false, Info: "Authentication failed!!"}, nil
	}
	if data.ExistGroupCheck(in.GetName()) == false {
		return &pb.Result{Status: false, Info: "Not exists Group!!"}, nil
	}
	id, err := db.GetDBGroupID(in.GetName())
	if err == false {
		fmt.Println("Get ID Error!!")
	}
	if db.RemoveDBGroup(id) {
		return &pb.Result{Status: true, Info: "OK!"}, nil
	} else {
		return &pb.Result{Status: false, Info: "DB Error!!"}, nil
	}
}

func (s *server) UserAddGroup(ctx context.Context, in *pb.GroupData) (*pb.Result, error) {
	log.Println("----UserAddGroup----")
	log.Println("Receive GroupName   : " + in.GetName())
	log.Println("Receive AddUserName : " + in.GetUser())
	log.Println("Receive AuthUser    : " + in.Base.User + ", AuthPass: " + in.Base.Pass)
	log.Println("Receive Token     : " + in.Base.GetToken())

	if in.GetMode() == 0 {
		log.Println("Receive Mode    : Admin")
	} else if in.GetMode() == 1 {
		log.Println("Receive Mode    : User")
	} else {
		return &pb.Result{Status: false, Info: "Mode error!!!"}, nil
	}

	if data.AdminUserCertification(in.Base.GetUser(), in.Base.GetPass(), in.Base.GetToken()) == false {
		return &pb.Result{Status: false, Info: "Authentication failed!!"}, nil
	}
	if data.ExistGroupCheck(in.GetName()) == false {
		return &pb.Result{Status: false, Info: "Not exists Group!!"}, nil
	}
	if in.GetMode() == 0 {
		info, result := data.AddGroupAdmin(in.GetUser(), in.GetName())
		if result {
			return &pb.Result{Status: true, Info: info}, nil
		}
		return &pb.Result{Status: false, Info: info}, nil
	} else if in.GetMode() == 1 {
		info, result := data.AddGroupUser(in.GetUser(), in.GetName())
		if result {
			return &pb.Result{Status: true, Info: info}, nil
		}
		return &pb.Result{Status: false, Info: info}, nil
	} else {
		return &pb.Result{Status: false, Info: "Mode error!!!"}, nil
	}
}

func (s *server) UserRemoveGroup(ctx context.Context, in *pb.GroupData) (*pb.Result, error) {
	log.Println("----UserRemoveGroup----")
	log.Println("Receive GroupName   : " + in.GetName())
	log.Println("Receive AddUserName : " + in.GetUser())
	log.Println("Receive AuthUser    : " + in.Base.User + ", AuthPass: " + in.Base.Pass)
	log.Println("Receive Token     : " + in.Base.GetToken())

	if in.GetMode() == 0 {
		log.Println("Receive Mode    : Admin")
	} else if in.GetMode() == 1 {
		log.Println("Receive Mode    : User")
	} else {
		return &pb.Result{Status: false, Info: "Mode error!!!"}, nil
	}

	if data.AdminUserCertification(in.Base.GetUser(), in.Base.GetPass(), in.Base.GetToken()) == false {
		return &pb.Result{Status: false, Info: "Authentication failed!!"}, nil
	}
	if data.ExistGroupCheck(in.GetName()) == false {
		return &pb.Result{Status: false, Info: "Not exists Group!!"}, nil
	}
	if in.GetMode() == 0 {
		info, result := data.RemoveGroupAdmin(in.GetUser(), in.GetName())
		if result {
			return &pb.Result{Status: true, Info: info}, nil
		}
		return &pb.Result{Status: false, Info: info}, nil
	} else if in.GetMode() == 1 {
		info, result := data.RemoveGroupUser(in.GetUser(), in.GetName())
		if result {
			return &pb.Result{Status: true, Info: info}, nil
		}
		return &pb.Result{Status: false, Info: info}, nil
	} else {
		return &pb.Result{Status: false, Info: "Mode error!!!"}, nil
	}
}

func (s *server) GetGroup(d *pb.GroupData, stream pb.Grpc_GetGroupServer) error {
	log.Println("----GetGroup----")
	log.Println("Receive AuthUser  : " + d.Base.GetUser() + ", AuthPass: " + d.Base.GetPass())
	log.Println("Receive Token     : " + d.Base.GetToken())

	isAdmin := false

	log.Printf("Receive GetAllGroup")
	if data.AdminUserCertification(d.Base.GetUser(), d.Base.GetPass(), d.Base.GetToken()) {
		isAdmin = true
	}

	_, _, r := data.TokenCertification(d.Base.GetToken())
	if r == false {
		fmt.Println("Auth Failed...")
		return nil
	}

	data, r := data.SearchUserForAdminGroup(data.GetUserName(d.Base.GetUser(), d.Base.GetToken()))
	if r == false {
		fmt.Println("Search Error!!")
		return nil
	}

	var groupdata []string

	for _, a := range data {
		tmp, r := db.GetDBGroup(a)
		if r == false {
			fmt.Println("Search Error!!")
			return nil
		}
		groupdata = append(groupdata, tmp.Name)
	}

	fmt.Println(db.GetDBAllGroup())
	result := db.GetDBAllGroup()
	for _, a := range result {
		if isAdmin {
			if err := stream.Send(&pb.GroupData{
				Id:    int32(a.ID),
				Name:  a.Name,
				Admin: a.Admin,
				User:  a.User,
				Sepc: &pb.SpecData{
					Maxvm:      int32(a.MaxVM),
					Maxcpu:     int32(a.MaxCPU),
					Maxmem:     int32(a.MaxMem),
					Maxstorage: int32(a.MaxStorage),
					Net:        a.Net,
				},
			}); err != nil {
				return err
			}
		} else if isAdmin == false {
			for _, b := range groupdata {
				if a.Name == b {
					if err := stream.Send(&pb.GroupData{
						Id:    int32(a.ID),
						Name:  a.Name,
						Admin: a.Admin,
						User:  a.User,
						Sepc: &pb.SpecData{
							Maxvm:      int32(a.MaxVM),
							Maxcpu:     int32(a.MaxCPU),
							Maxmem:     int32(a.MaxMem),
							Maxstorage: int32(a.MaxStorage),
							Net:        a.Net,
						},
					}); err != nil {
						return err
					}
					break
				}
			}

		}
	}

	return nil
}
