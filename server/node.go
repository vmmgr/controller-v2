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

func (s *baseServer) AddNode(ctx context.Context, in *pb.NodeData) (*pb.Result, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	log.Println("----AddNode----")
	log.Println("Receive Host : " + in.GetHostname())
	log.Println("Receive IP   : " + in.GetIp())
	log.Println("Token: " + md.Get("authorization")[0])

	if data.VerifyGroup(md.Get("authorization")[0]) != 0 {
		return &pb.Result{Status: false, Info: "Authorization NG!!"}, nil
	}

	admin := 1
	if in.GetOnlyAdmin() {
		admin = 0
	}

	if db.AddDBNode(db.Node{
		HostName:  in.GetHostname(),
		IP:        in.GetIp(),
		Path:      in.GetPath(),
		OnlyAdmin: admin,
		MaxCPU:    int(in.GetCpu()),
		MaxMem:    int(in.GetMem()),
		Active:    1,
	}) == nil {
		return &pb.Result{Status: true, Info: "OK!"}, nil
	} else {
		return &pb.Result{Status: false, Info: "NG!!"}, nil
	}
}

func (s *baseServer) DeleteNode(ctx context.Context, in *pb.NodeData) (*pb.Result, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	log.Println("----DeleteNode----")
	log.Println("Receive NodeID: " + strconv.Itoa(int(in.GetId())))
	log.Println("Token: " + md.Get("authorization")[0])

	if data.VerifyGroup(md.Get("authorization")[0]) != 0 {
		return &pb.Result{Status: false, Info: "Authorization NG!!"}, nil
	}

	if db.DeleteDBNode(db.Node{Model: gorm.Model{ID: uint(in.GetId())}}) == nil {
		return &pb.Result{Status: true, Info: "OK!"}, nil
	} else {
		return &pb.Result{Status: false, Info: "NG!!"}, nil
	}
}

func (s *baseServer) UpdateNode(ctx context.Context, in *pb.NodeData) (*pb.Result, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	log.Println("----UpdateNode----")
	log.Println("Receive Host : " + in.GetHostname())
	log.Println("Receive IP   : " + in.GetIp())
	log.Println("Token: " + md.Get("authorization")[0])

	if data.VerifyGroup(md.Get("authorization")[0]) != 0 {
		return &pb.Result{Status: false, Info: "Authorization NG!!"}, nil
	}

	admin := 1
	if in.GetOnlyAdmin() {
		admin = 0
	}

	if db.UpdateDBNode(db.Node{
		Model:     gorm.Model{ID: uint(in.GetId())},
		HostName:  in.GetHostname(),
		IP:        in.GetIp(),
		Path:      in.GetPath(),
		OnlyAdmin: admin,
		MaxCPU:    int(in.GetCpu()),
		MaxMem:    int(in.GetMem()),
		Active:    1,
	}) == nil {
		return &pb.Result{Status: true, Info: "OK!"}, nil
	} else {
		return &pb.Result{Status: false, Info: "NG!!"}, nil
	}

}

func (s *baseServer) GetAllNode(d *pb.Null, stream pb.Controller_GetNodeServer) error {
	md, _ := metadata.FromIncomingContext(stream.Context())
	log.Println("----GetAllUser----")
	log.Println("Token: " + md.Get("authorization")[0])

	if data.VerifyGroup(md.Get("authorization")[0]) != 0 {
		if err := stream.Send(&pb.NodeData{
			Id: 0,
		}); err != nil {
			return err
		}
		return nil
	}

	for _, r := range db.GetAllDBNode() {
		if err := stream.Send(&pb.NodeData{
			Id:        uint64(r.ID),
			Ip:        r.IP,
			Hostname:  r.HostName,
			Cpu:       uint32(r.MaxCPU),
			Mem:       uint32(r.MaxMem),
			Path:      r.Path,
			OnlyAdmin: r.OnlyAdmin == 1,
			Active:    r.Active == 1,
		}); err != nil {
			return err
		}
	}

	return nil
}
