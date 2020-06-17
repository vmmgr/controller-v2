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

func (s *baseServer) GenerateToken(ctx context.Context, in *pb.UserData) (*pb.Result, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	log.Println("----GenerateToken----")
	log.Println("Receive UserID: " + strconv.Itoa(int(in.GetId())))
	log.Println("Token: " + md.Get("authorization")[0])

	//if data.AdminUserCertification(in.Base.GetUser(), in.Base.GetPass(), in.Base.GetToken()) == false {
	//	return &pb.Result{Status: false, Info: "Authentication failed!!"}, nil
	//}
	if data.ExistUser(in.GetName()) {
		return &pb.Result{Status: false, Info: "exists username !!"}, nil
	}
	if db.AddDBToken(db.Token{
		Token:  etc.GenerateUUID(),
		UserID: int(in.GetId()),
	}) {
		return &pb.Result{Status: true, Info: "OK!"}, nil
	} else {
		return &pb.Result{Status: false, Info: "NG!!"}, nil
	}
}

func (s *baseServer) DeleteToken(ctx context.Context, in *pb.Null) (*pb.Result, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	log.Println("----DeleteUser----")
	log.Println("Token: " + md.Get("authorization")[0])

	if db.DeleteDBToken(db.Token{Token: md.Get("authorization")[0]}) {
		return &pb.Result{Status: true, Info: "OK!"}, nil
	} else {
		return &pb.Result{Status: false, Info: "NG!!"}, nil
	}
}

func (s *baseServer) GetAllToken(d *pb.Null, stream pb.Controller_GetAllTokenServer) error {
	md, _ := metadata.FromIncomingContext(stream.Context())
	log.Println("----GetAllUser----")
	log.Println("Token: " + md.Get("authorization")[0])

	result := db.GetAllDBToken()
	for _, r := range result {
		if err := stream.Send(&pb.TokenData{
			Token:  r.Token,
			Userid: int64(r.UserID),
		}); err != nil {
			return err
		}
	}

	return nil
}
