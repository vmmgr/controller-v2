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

// ID 0:AddISO 1:AddDISK 10:JoinISO 11:JoinDISK
func (s *server) AddImacon(ctx context.Context, in *pb.ImaconData) (*pb.Result, error) {
	log.Println("----AddImacon----")
	log.Println("Receive ImaConID : " + strconv.Itoa(int(in.GetId())))
	log.Println("Receive HostName : " + in.GetHostname())
	log.Println("Receive IP       : " + in.GetIP())
	log.Println("Receive AuthUser : " + in.GetBase().User + ", AuthPass: " + in.GetBase().Pass)
	log.Println("Receive Token    : " + in.GetBase().GetToken())

	if data.AdminUserCertification(in.GetBase().GetUser(), in.GetBase().GetPass(), in.GetBase().GetToken()) == false {
		return &pb.Result{Status: false, Info: "Authentication failed!!"}, nil
	}

	result := db.AddDBImaCon(db.ImaCon{
		ID:       int(in.GetId()),
		HostName: in.GetHostname(),
		IP:       in.GetIP(),
		Status:   int(in.GetStatus()),
	})

	return &pb.Result{Status: result}, nil
}

func (s *server) RemoveImacon(ctx context.Context, in *pb.ImaconData) (*pb.Result, error) {
	log.Println("----DeleteImacon----")
	log.Println("Receive ImaConID : " + strconv.Itoa(int(in.GetId())))
	log.Println("Receive AuthUser : " + in.GetBase().User + ", AuthPass: " + in.GetBase().Pass)
	log.Println("Receive Token    : " + in.GetBase().GetToken())

	if data.AdminUserCertification(in.GetBase().GetUser(), in.GetBase().GetPass(), in.GetBase().GetToken()) == false {
		return &pb.Result{Status: false, Info: "Authentication failed!!"}, nil
	}

	result := db.RemoveDBImaCon(int(in.GetId()))

	return &pb.Result{Status: result}, nil
}

func (s *server) GetAllImacon(d *pb.Base, stream pb.Grpc_GetNodeServer) error {
	log.Println("----GetAllImage----")

	result := db.GetDBAllImaCon()
	fmt.Printf("DBstruct: ")
	fmt.Println(result)
	for _, a := range result {

		if err := stream.Send(&pb.NodeData{
			NodeID:   int32(a.ID),
			Hostname: a.HostName,
			IP:       a.IP,
			Status:   int32(a.Status),
		}); err != nil {
			return err
		}
	}
	return nil
}
