package server

import (
	"context"
	"github.com/vmmgr/controller/db"
	spb "github.com/vmmgr/controller/proto/proto-go"
	pb "github.com/vmmgr/node/proto/proto-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log"
	"strconv"
	"time"
)

func (s *vmServer) AddVM(ctx context.Context, in *pb.VMData) (*spb.Result, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	log.Println("----AddVM----")
	log.Println("Receive GroupName: " + in.GetName())
	log.Println("Receive Mode: " + strconv.Itoa(int(in.GetMode())))
	log.Println("Token: " + md.Get("authorization")[0])

	token := md.Get("authorization")
	if in.GetID() < 100000 {
		return &spb.Result{Status: false, Info: "ID Error!!"}, nil
	}
	node := in.GetID() / 100000
	data := db.SearchDBNode(int(node))
	in.ID = in.GetID() % 100000

	conn, err := grpc.Dial(data.IP, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(2*time.Second))
	if err != nil {
		log.Fatalf("Not connect; %v", err)
	}
	defer conn.Close()

	client := pb.NewNodeClient(conn)
	header := metadata.New(map[string]string{"authorization": token[0]})
	cCtx := metadata.NewOutgoingContext(context.Background(), header)

	r, err := client.AddVM(cCtx, in)
	if err != nil {
		log.Fatal(err)
	}
	return &spb.Result{Status: r.Status, Info: r.Info}, nil
}

func (s *vmServer) DeleteVM(ctx context.Context, in *pb.VMData) (*spb.Result, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	log.Println("----DeleteVM----")
	log.Println("Receive GroupName: " + in.GetName())
	log.Println("Receive Mode: " + strconv.Itoa(int(in.GetMode())))
	log.Println("Token: " + md.Get("authorization")[0])

	token := md.Get("authorization")
	if in.GetID() < 100000 {
		return &spb.Result{Status: false, Info: "ID Error!!"}, nil
	}
	node := in.GetID() / 100000
	data := db.SearchDBNode(int(node))
	in.ID = in.GetID() % 100000

	conn, err := grpc.Dial(data.IP, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(2*time.Second))
	if err != nil {
		log.Fatalf("Not connect; %v", err)
	}
	defer conn.Close()

	client := pb.NewNodeClient(conn)
	header := metadata.New(map[string]string{"authorization": token[0]})
	cCtx := metadata.NewOutgoingContext(context.Background(), header)

	r, err := client.DeleteVM(cCtx, in)
	if err != nil {
		log.Fatal(err)
	}
	return &spb.Result{Status: r.Status, Info: r.Info}, nil
}

func (s *vmServer) UpdateVM(ctx context.Context, in *pb.VMData) (*spb.Result, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	log.Println("----UpdateVM----")
	log.Println("Receive GroupName: " + in.GetName())
	log.Println("Receive Mode: " + strconv.Itoa(int(in.GetMode())))
	log.Println("Token: " + md.Get("authorization")[0])

	token := md.Get("authorization")
	if in.GetID() < 100000 {
		return &spb.Result{Status: false, Info: "ID Error!!"}, nil
	}
	node := in.GetID() / 100000
	data := db.SearchDBNode(int(node))
	in.ID = in.GetID() % 100000

	conn, err := grpc.Dial(data.IP, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(2*time.Second))
	if err != nil {
		log.Fatalf("Not connect; %v", err)
	}
	defer conn.Close()

	client := pb.NewNodeClient(conn)
	header := metadata.New(map[string]string{"authorization": token[0]})
	cCtx := metadata.NewOutgoingContext(context.Background(), header)

	r, err := client.UpdateVM(cCtx, in)
	if err != nil {
		log.Fatal(err)
		return &spb.Result{Status: false, Info: r.Info}, nil
	}
	return &spb.Result{Status: r.Status, Info: r.Info}, nil
}
