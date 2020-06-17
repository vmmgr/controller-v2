package server

import (
	"context"
	"github.com/vmmgr/controller/data"
	"github.com/vmmgr/controller/db"
	pb "github.com/vmmgr/controller/proto/proto-go"
	"google.golang.org/grpc/metadata"
	"log"
	"strconv"
)

func (s *baseServer) AddGroup(ctx context.Context, in *pb.GroupData) (*pb.Result, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	log.Println("----AddGroup----")
	log.Println("Receive GroupName: " + in.GetName())
	log.Println("Receive Mode: " + strconv.Itoa(int(in.GetMode())))
	log.Println("Token: " + md.Get("authorization")[0])

	//if data.AdminUserCertification(in.Base.GetUser(), in.Base.GetPass(), in.Base.GetToken()) == false {
	//	return &pb.Result{Status: false, Info: "Authentication failed!!"}, nil
	//}
	if data.ExistUser(in.GetName()) {
		return &pb.Result{Status: false, Info: "exists username !!"}, nil
	}
	if db.AddDBGroup(db.Group{
		Name: in.GetName(),
	}) {
		return &pb.Result{Status: true, Info: "OK!"}, nil
	} else {
		return &pb.Result{Status: false, Info: "NG!!"}, nil
	}
}

func (s *baseServer) DeleteGroup(ctx context.Context, in *pb.GroupData) (*pb.Result, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	log.Println("----DeleteGroup----")
	log.Println("Receive GroupID: " + strconv.Itoa(int(in.GetId())))
	log.Println("Token: " + md.Get("authorization")[0])

	if data.ExistUser(in.GetName()) == false {
		return &pb.Result{Status: false, Info: "not exists username !!"}, nil
	}
	if db.DeleteDBGroup(db.Group{ID: int(in.GetId())}) {
		return &pb.Result{Status: true, Info: "OK!"}, nil
	} else {
		return &pb.Result{Status: false, Info: "NG!!"}, nil
	}
}

func (s *baseServer) GetAllGroup(d *pb.Null, stream pb.Controller_GetAllGroupServer) error {
	md, _ := metadata.FromIncomingContext(stream.Context())
	log.Println("----GetAllUser----")
	log.Println("Token: " + md.Get("authorization")[0])

	result := db.GetAllDBGroup()
	for _, r := range result {
		if err := stream.Send(&pb.GroupData{
			Id:   int64(r.ID),
			Name: r.Name,
			//Admin: r.AdminUser,
			//User:  r.StandardUser,
		}); err != nil {
			return err
		}
	}

	return nil
}

func (s *baseServer) UpdateGroup(ctx context.Context, in *pb.GroupData) (*pb.Result, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	log.Println("----UpdateGroup----")
	log.Println("Receive GroupID: " + strconv.Itoa(int(in.GetId())))
	log.Println("Receive GroupName: " + in.GetName())
	log.Println("Token: " + md.Get("authorization")[0])

	if db.UpdateDBGroup(db.Group{
		ID:   int(in.GetId()),
		Name: in.GetName(),
	}) {
		return &pb.Result{Status: true, Info: "OK!"}, nil
	} else {
		return &pb.Result{Status: false, Info: "NG!!"}, nil
	}

}

func (s *baseServer) JoinAddGroup(ctx context.Context, in *pb.GroupData) (*pb.Result, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	log.Println("----JoinAddGroup----")
	log.Println("Receive GroupID: " + strconv.Itoa(int(in.GetId())))
	log.Println("Receive GroupName: " + in.GetName())
	log.Println("Token: " + md.Get("authorization")[0])

	admin := false
	var user string

	if in.GetAdmin() != "" {
		admin = true
		user = in.GetAdmin()
	} else if in.GetUser() != "" {
		user = in.GetUser()
	} else {
		return &pb.Result{Status: false, Info: "NG!!"}, nil
	}

	if db.AddDBGroupUser(db.Group{
		ID: int(in.GetId()),
	}, user, admin) {
		return &pb.Result{Status: true, Info: "OK!"}, nil
	} else {
		return &pb.Result{Status: false, Info: "NG!!"}, nil
	}

}

func (s *baseServer) JoinDeleteGroup(ctx context.Context, in *pb.GroupData) (*pb.Result, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	log.Println("----JoinDeleteGroup----")
	log.Println("Receive GroupID: " + strconv.Itoa(int(in.GetId())))
	log.Println("Token: " + md.Get("authorization")[0])

	admin := false
	var user string

	if in.GetAdmin() != "" {
		admin = true
		user = in.GetAdmin()
	} else if in.GetUser() != "" {
		user = in.GetUser()
	} else {
		return &pb.Result{Status: false, Info: "NG!!"}, nil
	}

	if db.DeleteDBGroupUser(db.Group{
		ID: int(in.GetId()),
	}, user, admin) {
		return &pb.Result{Status: true, Info: "OK!"}, nil
	} else {
		return &pb.Result{Status: false, Info: "NG!!"}, nil
	}

}
