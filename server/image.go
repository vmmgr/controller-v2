package server

import (
	"context"
	"fmt"
	"github.com/vmmgr/controller/data"
	"github.com/vmmgr/controller/db"
	pb "github.com/vmmgr/controller/proto/proto-go"
	"google.golang.org/grpc"
	"io"
	"log"
	"strconv"
	"time"
)

// ID 0:AddISO 1:AddDISK 10:JoinISO 11:JoinDISK
func (s *server) AddImage(ctx context.Context, in *pb.ImageData) (*pb.ImageResult, error) {
	log.Println("----AddImage----")
	log.Println("Receive FileName : " + in.GetFilename())
	log.Println("Receive ImaConID : " + strconv.Itoa(int(in.GetImaconid())))
	log.Println("Receive Type     : " + strconv.Itoa(int(in.GetType())))
	log.Println("Receive Name     : " + in.GetName())
	log.Println("Receive AuthUser : " + in.GetBase().User + ", AuthPass: " + in.GetBase().Pass)
	log.Println("Receive Token    : " + in.GetBase().GetToken())

	if data.AdminUserCertification(in.GetBase().GetUser(), in.GetBase().GetPass(), in.GetBase().GetToken()) == false {
		return &pb.ImageResult{Result: false, Info: "Authentication failed!!"}, nil
	}

	if in.GetImaconid() < 0 {
		return &pb.ImageResult{Result: false, Info: "imacon id is wrong!!"}, nil
	}
	ip := data.GetImaConHostIP(int(in.GetImaconid()))

	conn, err := grpc.Dial(ip, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		fmt.Printf("Not connect; ")
		fmt.Println(err)
	}
	defer conn.Close()
	c := pb.NewGrpcClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if in.GetType() == 10 || in.GetType() == 11 {
		r, err := c.AddImage(ctx, &pb.ImageData{
			Name:      in.GetName(),
			Tag:       in.GetTag(),
			Type:      in.GetType(),
			Minmem:    in.GetMinmem(),
			Authority: in.GetAuthority(),
			Filename:  in.GetFilename(),
		})
		if err != nil {
			fmt.Printf("Not connect; ")
			fmt.Println(err)
		}
		return &pb.ImageResult{Result: r.Result, Info: r.Info}, nil
	} else if in.GetType() == 1 {
		d := data.GetNodeDISK(in.GetFilename())
		if d.Path == "0" {
			return &pb.ImageResult{Result: false, Info: "Not Found!!"}, nil
		}
		r, err := c.AddImage(ctx, &pb.ImageData{
			Name:      in.GetName(),
			Path:      d.Path,
			Tag:       in.GetTag(),
			Type:      in.GetType(),
			Minmem:    in.GetMinmem(),
			Authority: in.GetAuthority(),
			Filename:  in.GetFilename(),
			Raddr:     d.Address,
		})
		if err != nil {
			fmt.Printf("Not connect; ")
			fmt.Println(err)
		}
		return &pb.ImageResult{Result: r.Result, Info: r.Info}, nil
	}

	return &pb.ImageResult{Result: false, Info: "Failed"}, nil
}

func (s *server) DeleteImage(ctx context.Context, in *pb.ImageData) (*pb.ImageResult, error) {
	log.Println("----DeleteImage----")
	log.Println("Receive FileName : " + in.GetFilename())
	log.Println("Receive ID     : " + strconv.Itoa(int(in.GetType())))
	log.Println("Receive AuthUser : " + in.GetBase().User + ", AuthPass: " + in.GetBase().Pass)
	log.Println("Receive Token    : " + in.GetBase().GetToken())
	log.Println("Receive ImaConID : " + strconv.Itoa(int(in.GetImaconid())))

	if data.AdminUserCertification(in.GetBase().GetUser(), in.GetBase().GetPass(), in.GetBase().GetToken()) == false {
		return &pb.ImageResult{Result: false, Info: "Authentication failed!!"}, nil
	}

	d := data.GetImage(in.GetName(), in.GetTag())

	conn, err := grpc.Dial(d.Address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		fmt.Printf("Not connect; ")
		fmt.Println(err)
	}
	defer conn.Close()
	c := pb.NewGrpcClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.DeleteImage(ctx, &pb.ImageData{Name: in.GetName(), Tag: in.GetTag(), Type: in.GetType()})
	if err != nil {
		fmt.Printf("Not connect; ")
		fmt.Println(err)
	}

	return &pb.ImageResult{Result: r.Result, Info: r.Info}, nil
}

func (s *server) ChangeNameImage(ctx context.Context, in *pb.ImageData) (*pb.ImageResult, error) {
	log.Println("----ChangeNameImage----")
	log.Println("Receive FileName : " + in.GetFilename())
	log.Println("Receive Name     : " + in.GetName())
	log.Println("Receive AuthUser : " + in.GetBase().User + ", AuthPass: " + in.GetBase().Pass)
	log.Println("Receive Token    : " + in.GetBase().GetToken())

	if data.AdminUserCertification(in.GetBase().GetUser(), in.GetBase().GetPass(), in.GetBase().GetToken()) == false {
		return &pb.ImageResult{Result: false, Info: "Authentication failed!!"}, nil
	}

	if in.GetImaconid() < 0 {
		return &pb.ImageResult{Result: false, Info: "imacon id is wrong!!"}, nil
	}
	ip := data.GetImaConHostIP(int(in.GetImaconid()))

	conn, err := grpc.Dial(ip, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		fmt.Printf("Not connect; ")
		fmt.Println(err)
	}
	defer conn.Close()
	c := pb.NewGrpcClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.ChangeNameImage(ctx, &pb.ImageData{
		Filename: in.GetFilename(),
		Name:     in.GetName(),
	})
	if err != nil {
		fmt.Printf("Not connect; ")
		fmt.Println(err)
	}

	return &pb.ImageResult{Result: r.Result, Info: r.Info}, nil
}

func (s *server) ChangeTagImage(ctx context.Context, in *pb.ImageData) (*pb.ImageResult, error) {
	log.Println("----ChangeNameImage----")
	log.Println("Receive FileName : " + in.GetFilename())
	log.Println("Receive Tag      : " + in.GetTag())
	log.Println("Receive AuthUser : " + in.GetBase().User + ", AuthPass: " + in.GetBase().Pass)
	log.Println("Receive Token    : " + in.GetBase().GetToken())

	if data.AdminUserCertification(in.GetBase().GetUser(), in.GetBase().GetPass(), in.GetBase().GetToken()) == false {
		return &pb.ImageResult{Result: false, Info: "Authentication failed!!"}, nil
	}

	if in.GetImaconid() < 0 {
		return &pb.ImageResult{Result: false, Info: "imacon id is wrong!!"}, nil
	}
	ip := data.GetImaConHostIP(int(in.GetImaconid()))

	conn, err := grpc.Dial(ip, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		fmt.Printf("Not connect; ")
		fmt.Println(err)
	}
	defer conn.Close()
	c := pb.NewGrpcClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.ChangeTagImage(ctx, &pb.ImageData{
		Filename: in.GetFilename(),
		Tag:      in.GetTag(),
	})
	if err != nil {
		fmt.Printf("Not connect; ")
		fmt.Println(err)
	}

	return &pb.ImageResult{Result: r.Result, Info: r.Info}, nil
}

func (s *server) ExistImage(ctx context.Context, in *pb.ImageData) (*pb.ImageResult, error) {
	log.Println("----ExistImage----")
	log.Println("Receive ID : " + strconv.Itoa(int(in.GetType())))
	log.Println("Receive Name : " + in.GetName())
	log.Println("Receive Tag  : " + in.GetTag())
	log.Println("Receive AuthUser : " + in.GetBase().User + ", AuthPass: " + in.GetBase().Pass)
	log.Println("Receive Token    : " + in.GetBase().GetToken())

	if data.AdminUserCertification(in.GetBase().GetUser(), in.GetBase().GetPass(), in.GetBase().GetToken()) == false {
		return &pb.ImageResult{Result: false, Info: "Authentication failed!!"}, nil
	}

	if in.GetImaconid() < 0 {
		return &pb.ImageResult{Result: false, Info: "imacon id is wrong!!"}, nil
	}
	ip := data.GetImaConHostIP(int(in.GetImaconid()))

	conn, err := grpc.Dial(ip, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		fmt.Printf("Not connect; ")
		fmt.Println(err)
	}
	defer conn.Close()
	c := pb.NewGrpcClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.ExistImage(ctx, &pb.ImageData{
		Name: in.GetName(),
		Tag:  in.GetTag(),
	})
	if err != nil {
		fmt.Printf("Not connect; ")
		fmt.Println(err)
	}

	return &pb.ImageResult{Result: r.Result, Info: r.Info}, nil
}

func (s *server) GetAllImage(d *pb.Base, stream pb.Grpc_GetAllImageServer) error {
	go data.DeleteExpiredToken()
	log.Println("----GetAllImage----")
	log.Println("Receive AuthUser : " + d.GetUser() + ", AuthPass: " + d.GetPass())
	log.Println("Receive Token    : " + d.GetToken())

	isAdmin := true

	if data.AdminUserCertification(d.GetUser(), d.GetPass(), d.GetToken()) == false {
		_, _, r := data.TokenCertification(d.GetToken())
		if r {
			isAdmin = false
		} else {
			return nil
		}
	}
	fmt.Printf("isAdmin: ")
	fmt.Println(isAdmin)

	var data []pb.ImageData

	for _, a := range db.GetDBAllImaCon() {
		conn, err := grpc.Dial(a.IP, grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			fmt.Printf("Not connect; ")
			fmt.Println(err)
		}
		defer conn.Close()

		c := pb.NewGrpcClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		stream, err := c.GetAllImage(ctx, &pb.Base{})
		if err != nil {
			fmt.Println(err)
		}
		for {
			article, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				fmt.Println(err)
			}

			d2 := pb.ImageData{Id: article.Id, Path: article.Path, Name: article.Name,
				Tag: article.Tag, Type: article.Type, Capacity: article.Capacity, Addtime: article.Addtime,
				Minmem: article.Minmem, Authority: article.Authority, Filename: article.Filename, Imaconid: int32(a.ID)}

			data = append(data, d2)
		}
	}

	if isAdmin {
		for _, a := range data {
			if err := stream.Send(&a); err != nil {
				return err
			}
		}
	} else {
		for _, a := range data {
			if a.Authority == 0 {
				if err := stream.Send(&a); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
