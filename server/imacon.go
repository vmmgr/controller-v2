package server

import (
	"context"
	"github.com/jinzhu/gorm"
	"github.com/vmmgr/controller/data"
	"github.com/vmmgr/controller/db"
	pb "github.com/vmmgr/controller/proto/proto-go"
	"google.golang.org/grpc/metadata"
	"log"
	"strconv"
)

func (s *baseServer) AddImaCon(ctx context.Context, in *pb.ImaConData) (*pb.Result, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	log.Println("----AddImaCon----")
	log.Println("Receive Host : " + in.GetHostname())
	log.Println("Receive IP   : " + in.GetIp())
	log.Println("Token: " + md.Get("authorization")[0])

	if data.VerifyGroup(md.Get("authorization")[0]) != 0 {
		return &pb.Result{Status: false, Info: "Authorization NG!!"}, nil
	}

	if db.AddDBImaCon(db.ImaCon{
		HostName: in.Hostname,
		IP:       in.Ip,
		Status:   0,
	}) == nil {
		return &pb.Result{Status: true, Info: "OK!"}, nil
	} else {
		return &pb.Result{Status: false, Info: "NG!!"}, nil
	}
}

func (s *baseServer) DeleteImaCon(ctx context.Context, in *pb.ImaConData) (*pb.Result, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	log.Println("----DeleteImaCon----")
	log.Println("Receive ImaConID: " + strconv.Itoa(int(in.GetId())))
	log.Println("Token: " + md.Get("authorization")[0])

	if data.VerifyGroup(md.Get("authorization")[0]) != 0 {
		return &pb.Result{Status: false, Info: "Authorization NG!!"}, nil
	}

	if db.DeleteDBImaCon(db.ImaCon{Model: gorm.Model{ID: uint(in.GetId())}}) == nil {
		return &pb.Result{Status: true, Info: "OK!"}, nil
	} else {
		return &pb.Result{Status: false, Info: "NG!!"}, nil
	}
}

func (s *baseServer) UpdateImaCon(ctx context.Context, in *pb.ImaConData) (*pb.Result, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	log.Println("----UpdateImaCon----")
	log.Println("Receive Host : " + in.GetHostname())
	log.Println("Receive IP   : " + in.GetIp())
	log.Println("Token: " + md.Get("authorization")[0])

	if data.VerifyGroup(md.Get("authorization")[0]) != 0 {
		return &pb.Result{Status: false, Info: "Authorization NG!!"}, nil
	}

	if db.UpdateDBImaCon(db.ImaCon{
		HostName: in.Hostname,
		IP:       in.Ip,
		Status:   0,
	}) == nil {
		return &pb.Result{Status: true, Info: "OK!"}, nil
	} else {
		return &pb.Result{Status: false, Info: "NG!!"}, nil
	}
}

func (s *baseServer) GetAllImaCon(d *pb.Null, stream pb.Controller_GetImaConServer) error {
	md, _ := metadata.FromIncomingContext(stream.Context())
	log.Println("----GetAllUser----")
	log.Println("Token: " + md.Get("authorization")[0])

	if data.VerifyGroup(md.Get("authorization")[0]) != 0 {
		if err := stream.Send(&pb.ImaConData{
			Id: 0,
		}); err != nil {
			return err
		}
		return nil
	}

	for _, r := range db.GetAllDBImaCon() {
		if err := stream.Send(&pb.ImaConData{
			Id:       uint64(r.ID),
			Ip:       r.IP,
			Hostname: r.HostName,
		}); err != nil {
			return err
		}
	}

	return nil
}
