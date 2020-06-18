package server

import (
	"context"
	"github.com/vmmgr/controller/data"
	"github.com/vmmgr/controller/db"
	"github.com/vmmgr/controller/etc"
	pb "github.com/vmmgr/controller/proto/proto-go"
	"google.golang.org/grpc/metadata"
	"log"
	"strconv"
)

func (s *baseServer) AddUser(ctx context.Context, in *pb.UserData) (*pb.Result, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	log.Println("----AddUser----")
	log.Println("Receive UserName: " + in.GetName())
	log.Println("Receive Pass: " + in.GetPass())
	log.Println("Receive Auth: " + strconv.Itoa(int(in.GetAuth())))
	log.Println("Token: " + md.Get("authorization")[0])

	if data.VerifyGroup(md.Get("authorization")[0]) != 0 {
		return &pb.Result{Status: false, Info: "Authorization NG!!"}, nil
	}
	//if data.AdminUserCertification(in.Base.GetUser(), in.Base.GetPass(), in.Base.GetToken()) == false {
	//	return &pb.Result{Status: false, Info: "Authentication failed!!"}, nil
	//}
	if data.ExistUser(in.GetName()) {
		return &pb.Result{Status: false, Info: "exists username !!"}, nil
	}
	if db.AddDBUser(db.User{
		Name: in.GetName(),
		Pass: etc.HashGenerate(in.GetPass()),
		Auth: int(in.GetAuth()),
	}) {
		return &pb.Result{Status: true, Info: "OK!"}, nil
	} else {
		return &pb.Result{Status: false, Info: "NG!!"}, nil
	}
}

func (s *baseServer) DeleteUser(ctx context.Context, in *pb.UserData) (*pb.Result, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	log.Println("----DeleteUser----")
	log.Println("Receive ID: " + strconv.Itoa(int(in.GetId())))
	log.Println("Token: " + md.Get("authorization")[0])

	if data.VerifyGroup(md.Get("authorization")[0]) != 0 {
		return &pb.Result{Status: false, Info: "Authorization NG!!"}, nil
	}

	if data.ExistUser(in.GetName()) == false {
		return &pb.Result{Status: false, Info: "not exists username !!"}, nil
	}
	if db.DeleteDBUser(db.User{ID: int(in.GetId())}) {
		return &pb.Result{Status: true, Info: "OK!"}, nil
	} else {
		return &pb.Result{Status: false, Info: "NG!!"}, nil
	}
}

func (s *baseServer) GetAllUser(d *pb.Null, stream pb.Controller_GetAllUserServer) error {
	md, _ := metadata.FromIncomingContext(stream.Context())
	log.Println("----GetAllUser----")
	log.Println("Token: " + md.Get("authorization")[0])

	if data.VerifyGroup(md.Get("authorization")[0]) != 0 {
		if err := stream.Send(&pb.UserData{
			Id: 0,
		}); err != nil {
			return err
		}
		return nil
	}

	for _, r := range db.GetAllDBUser() {
		if err := stream.Send(&pb.UserData{
			Id:   uint64(r.ID),
			Name: r.Name,
			Pass: r.Pass,
			Auth: uint32(r.Auth),
		}); err != nil {
			return err
		}
	}

	return nil
}

func (s *baseServer) UpdateUser(ctx context.Context, in *pb.UserData) (*pb.Result, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	log.Println("----UpdateUser----")
	log.Println("Receive ID: " + strconv.Itoa(int(in.GetId())))
	log.Println("Receive UserName: " + in.GetName())
	log.Println("Token: " + md.Get("authorization")[0])

	if data.VerifyGroup(md.Get("authorization")[0]) != 0 {
		return &pb.Result{Status: false, Info: "Authorization NG!!"}, nil
	}

	if db.UpdateDBUser(db.User{
		ID:   int(in.GetId()),
		Name: in.GetName(),
		Pass: in.GetPass(),
		Auth: int(in.GetAuth()),
	}) {
		return &pb.Result{Status: true, Info: "OK!"}, nil
	} else {
		return &pb.Result{Status: false, Info: "NG!!"}, nil
	}
}
