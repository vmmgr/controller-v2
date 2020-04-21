package server

import (
	"context"
	"fmt"
	"github.com/vmmgr/controller/data"
	"github.com/vmmgr/controller/db"
	pb "github.com/vmmgr/controller/proto/proto-go"
	"google.golang.org/grpc/metadata"
	"log"
)

func (s *server) GenerateToken(ctx context.Context, in *pb.Base) (*pb.AuthResult, error) {
	log.Println("----GenerateToken----")
	log.Println("Receive AuthUser : " + in.GetUser() + ", AuthPass: " + in.GetPass())
	log.Println("Receive Token    : " + in.GetToken())

	test, _ := metadata.FromIncomingContext(ctx)
	fmt.Println(test)

	if data.UserCertification(in.GetUser(), in.GetPass()) == false {
		return &pb.AuthResult{Result: false, Token: "Auth Failed!!"}, nil
	}
	userid, result := db.GetDBUserID(in.GetUser())
	if result == false {
		userid = 0
	}

	uuid, result := data.NewToken(in.GetUser())
	if result {
		return &pb.AuthResult{Result: true, Token: uuid, Name: in.GetUser(), Id: int32(userid)}, nil
	} else {
		return &pb.AuthResult{Result: false, Token: uuid}, nil
	}
}

func (s *server) DeleteToken(ctx context.Context, in *pb.Base) (*pb.Result, error) {
	go data.DeleteExpiredToken()
	log.Println("----DeleteToken----")
	log.Println("Receive Token    : " + in.Token)

	data, result := db.GetDBToken(in.Token)
	if result == false {
		return &pb.Result{Status: false, Info: "DB Search failed!!"}, nil
	}
	info, result := db.RemoveDBToken(data.ID)
	if result {
		return &pb.Result{Status: true, Info: "OK!"}, nil
	} else {
		return &pb.Result{Status: false, Info: info}, nil
	}
}

func (s *server) CheckToken(ctx context.Context, in *pb.Base) (*pb.Result, error) {
	go data.DeleteExpiredToken()
	log.Println("----TokenCheck----")
	log.Println("Receive Token    : " + in.Token)

	a, result := db.GetDBToken(in.Token)
	if result {
		return &pb.Result{Status: true, Info: "OK!", Id: int32(a.Userid)}, nil
	} else {
		return &pb.Result{Status: false, Info: "NG"}, nil
	}
}

func (s *server) GetAllToken(d *pb.Base, stream pb.Grpc_GetAllTokenServer) error {
	go data.DeleteExpiredToken()
	log.Println("----GetAllToken----")
	log.Println("Receive AuthUser : " + d.GetUser() + ", AuthPass: " + d.GetPass())
	log.Println("Receive Token    : " + d.GetToken())
	if data.AdminUserCertification(d.GetUser(), d.GetPass(), d.GetToken()) == false {
		fmt.Println("Administrator certification failed!!!")
		return nil
	}
	result := db.GetDBAllToken()
	fmt.Printf("DBstruct: ")
	fmt.Println(result)
	for _, a := range result {
		if err := stream.Send(&pb.TokenData{
			Id:        int64(a.ID),
			Token:     a.Token,
			Userid:    int32(a.Userid),
			Groupid:   a.Groupid,
			Begintime: int64(a.Begintime),
			Endtime:   int64(a.Endtime),
		}); err != nil {
			return err
		}
	}
	return nil
}
