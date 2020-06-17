package server

import (
	"context"
	"github.com/vmmgr/controller/db"
	spb "github.com/vmmgr/controller/proto/proto-go"
	pb "github.com/vmmgr/node/proto/proto-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"io"
	"log"
	"strconv"
	"time"
)

func (s *vmServer) AddStorage(in *pb.StorageData, stream pb.Node_AddStorageServer) error {
	md, _ := metadata.FromIncomingContext(stream.Context())
	log.Println("----AddStorage----")
	log.Println("Receive GroupName: " + in.GetName())
	log.Println("Receive Mode: " + strconv.Itoa(int(in.GetMode())))
	log.Println("Token: " + md.Get("authorization")[0])

	token := md.Get("authorization")
	if in.GetID() < 100000 {
		if err := stream.Send(&pb.Result{
			Info:   "ID Error",
			Status: false,
		}); err != nil {
			return err
		}
		return nil
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

	cStream, err := client.AddStorage(cCtx, in)
	if err != nil {
		log.Fatal(err)
	}
	for {
		d, err := cStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		if err := stream.Send(&pb.Result{
			Info:   d.Info,
			Status: d.Status,
		}); err != nil {
			return err
		}
		log.Println("Info: " + d.GetInfo())
		log.Println("Status: " + strconv.FormatBool(d.GetStatus()))
	}
	return nil
}

func (s *vmServer) DeleteStorage(ctx context.Context, in *pb.StorageData) (*spb.Result, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	log.Println("----DeleteStorage----")
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

	r, err := client.DeleteStorage(cCtx, in)
	if err != nil {
		log.Fatal(err)
	}
	return &spb.Result{Status: r.Status, Info: r.Info}, nil
}

func (s *vmServer) UpdateStorage(ctx context.Context, in *pb.StorageData) (*spb.Result, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	log.Println("----UpdateStorage----")
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

	r, err := client.UpdateStorage(cCtx, in)
	if err != nil {
		log.Fatal(err)
		return &spb.Result{Status: false, Info: r.Info}, nil
	}
	return &spb.Result{Status: r.Status, Info: r.Info}, nil
}
