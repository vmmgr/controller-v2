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
	"strconv"
	"time"
)

func (s *vmServer) AddVM(ctx context.Context, in *pb.VMData) (*spb.Result, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	log.Println("----AddVM----")
	log.Println("Receive VMID    : " + strconv.Itoa(int(in.GetID())))
	log.Println("Receive GroupID : " + strconv.Itoa(int(in.GetGroupID())))
	log.Println("Receive Mode    : " + strconv.Itoa(int(in.GetMode())))
	log.Println("Token: " + md.Get("authorization")[0])

	if in.GetID() < 100000 {
		return &spb.Result{Status: false, Info: "ID Error!!"}, nil
	}
	dataNode := db.SearchDBNode(int(in.GetID() / 100000))
	in.ID = in.GetID() % 100000

	if 2 <= data.VerifySameGroup(md.Get("authorization")[0], int(in.GetGroupID())) {
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

	r, err := client.AddVM(cCtx, in)
	if err != nil {
		log.Fatal(err)
	}
	return &spb.Result{Status: r.Status, Info: r.Info}, nil
}

func (s *vmServer) DeleteVM(ctx context.Context, in *pb.VMData) (*spb.Result, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	log.Println("----DeleteVM----")
	log.Println("Receive VMID    : " + strconv.Itoa(int(in.GetID())))
	log.Println("Receive GroupID : " + strconv.Itoa(int(in.GetGroupID())))
	log.Println("Receive Mode    : " + strconv.Itoa(int(in.GetMode())))
	log.Println("Token: " + md.Get("authorization")[0])

	if in.GetID() < 100000 {
		return &spb.Result{Status: false, Info: "ID Error!!"}, nil
	}
	dataNode := db.SearchDBNode(int(in.GetID() / 100000))
	in.ID = in.GetID() % 100000

	if 2 <= data.VerifySameGroup(md.Get("authorization")[0], int(in.GetGroupID())) {
		log.Println("Error: Authentication Error")
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

	r, err := client.DeleteVM(cCtx, in)
	if err != nil {
		log.Fatal(err)
	}
	return &spb.Result{Status: r.Status, Info: r.Info}, nil
}

func (s *vmServer) UpdateVM(ctx context.Context, in *pb.VMData) (*spb.Result, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	log.Println("----DeleteVM----")
	log.Println("Receive VMID    : " + strconv.Itoa(int(in.GetID())))
	log.Println("Receive GroupID : " + strconv.Itoa(int(in.GetGroupID())))
	log.Println("Receive Mode    : " + strconv.Itoa(int(in.GetMode())))
	log.Println("Token: " + md.Get("authorization")[0])

	if in.GetID() < 100000 {
		return &spb.Result{Status: false, Info: "ID Error!!"}, nil
	}
	dataNode := db.SearchDBNode(int(in.GetID() / 100000))
	in.ID = in.GetID() % 100000

	if 2 <= data.VerifySameGroup(md.Get("authorization")[0], int(in.GetGroupID())) {
		log.Println("Error: Authentication Error")
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

	r, err := client.UpdateVM(cCtx, in)
	if err != nil {
		log.Fatal(err)
		return &spb.Result{Status: false, Info: r.Info}, nil
	}
	return &spb.Result{Status: r.Status, Info: r.Info}, nil
}

func getVMGroup(ip string, id uint64) int {
	conn, err := grpc.Dial(ip, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(2*time.Second))
	if err != nil {
		log.Fatalf("Not connect; %v", err)
	}
	defer conn.Close()

	client := pb.NewNodeClient(conn)
	header := metadata.New(map[string]string{"node": "true"})
	ctx := metadata.NewOutgoingContext(context.Background(), header)

	r, err := client.GetVM(ctx, &pb.VMData{ID: id})
	if err != nil {
		log.Fatal(err)
	}
	return int(r.GroupID)
}
