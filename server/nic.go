package server

import (
	"context"
	"github.com/vmmgr/controller/data"
	"github.com/vmmgr/controller/db"
	spb "github.com/vmmgr/controller/proto/proto-go"
	pb "github.com/vmmgr/node/proto/proto-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log"
	"time"
)

func (s *vmServer) AddNIC(ctx context.Context, in *pb.NICData) (*spb.Result, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	log.Println("----AddNIC----")
	log.Println("Receive GroupName: " + in.GetName())
	log.Println("Token: " + md.Get("authorization")[0])

	if in.GetID() < 100000 {
		return &spb.Result{Status: false, Info: "ID Error!!"}, nil
	}
	dataNode := db.SearchDBNode(int(in.GetID() / 100000))
	in.ID = in.GetID() % 100000

	result := data.VerifyGroup(md.Get("authorization")[0], int(in.GetGroupID()))
	if result < 0 && 2 <= result {
		return &spb.Result{Status: false, Info: "Authentication Error "}, nil
	}

	conn, err := grpc.Dial(dataNode.IP, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(2*time.Second))
	if err != nil {
		log.Fatalf("Not connect; %v", err)
	}
	defer conn.Close()

	client := pb.NewNodeClient(conn)
	header := metadata.New(map[string]string{"node": "true"})
	cCtx := metadata.NewOutgoingContext(context.Background(), header)

	r, err := client.AddNIC(cCtx, in)
	if err != nil {
		log.Fatal(err)
	}
	return &spb.Result{Status: r.Status, Info: r.Info}, nil
}

func (s *vmServer) DeleteNIC(ctx context.Context, in *pb.NICData) (*spb.Result, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	log.Println("----DeleteNIC----")
	log.Println("Receive GroupName: " + in.GetName())

	if in.GetID() < 100000 {
		return &spb.Result{Status: false, Info: "ID Error!!"}, nil
	}
	dataNode := db.SearchDBNode(int(in.GetID() / 100000))
	in.ID = in.GetID() % 100000

	result := data.VerifyGroup(md.Get("authorization")[0], int(in.GetGroupID()))
	if result < 0 && 2 <= result {
		return &spb.Result{Status: false, Info: "Authentication Error "}, nil
	}

	conn, err := grpc.Dial(dataNode.IP, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(2*time.Second))
	if err != nil {
		log.Fatalf("Not connect; %v", err)
	}
	defer conn.Close()

	client := pb.NewNodeClient(conn)
	header := metadata.New(map[string]string{"node": "true"})
	cCtx := metadata.NewOutgoingContext(context.Background(), header)

	r, err := client.DeleteNIC(cCtx, in)
	if err != nil {
		log.Fatal(err)
	}
	return &spb.Result{Status: r.Status, Info: r.Info}, nil
}

func (s *vmServer) UpdateNIC(ctx context.Context, in *pb.NICData) (*spb.Result, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	log.Println("----UpdateNIC----")
	log.Println("Receive GroupName: " + in.GetName())
	log.Println("Token: " + md.Get("authorization")[0])

	if in.GetID() < 100000 {
		return &spb.Result{Status: false, Info: "ID Error!!"}, nil
	}
	dataNode := db.SearchDBNode(int(in.GetID() / 100000))
	in.ID = in.GetID() % 100000

	result := data.VerifyGroup(md.Get("authorization")[0], int(in.GetGroupID()))
	if result < 0 && 2 <= result {
		return &spb.Result{Status: false, Info: "Authentication Error "}, nil
	}

	conn, err := grpc.Dial(dataNode.IP, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(2*time.Second))
	if err != nil {
		log.Fatalf("Not connect; %v", err)
	}
	defer conn.Close()

	client := pb.NewNodeClient(conn)
	header := metadata.New(map[string]string{"node": "true"})
	cCtx := metadata.NewOutgoingContext(context.Background(), header)

	r, err := client.UpdateNIC(cCtx, in)
	if err != nil {
		log.Fatal(err)
		return &spb.Result{Status: false, Info: r.Info}, nil
	}
	return &spb.Result{Status: r.Status, Info: r.Info}, nil
}
