package data

import (
	"context"
	"fmt"
	"github.com/vmmgr/controller/db"
	pb "github.com/vmmgr/controller/proto/proto-go"
	"google.golang.org/grpc"
	"io"
	"strconv"
	"strings"
	"time"
)

type ImageStruct struct {
	Address string
	Result  bool
}

type NodeDISKData struct {
	Address string
	Path    string
}

func GetImaConHostIP(id int) string {
	d, result := db.GetDBImaCon(id)
	if result == false {
		fmt.Println("ImaCon IP not found....")
		return ""
	}
	return d.IP
}

func GetImage(name, tag string) ImageStruct {
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
		r, err := c.ExistImage(ctx, &pb.ImageData{Name: name, Tag: tag})
		if err != nil {
			fmt.Println(err)
		}
		if r.Result {
			return ImageStruct{Address: a.IP, Result: true}
		}
	}
	return ImageStruct{Result: false}
}

func CheckImage(name, tag string) bool {

	return false
}

func GetNodeDISK(data string) NodeDISKData {
	for _, a := range db.GetDBAllNode() {
		conn, err := grpc.Dial(a.IP, grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			fmt.Printf("Not connect; ")
			fmt.Println(err)
		}
		defer conn.Close()

		c := pb.NewGrpcClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		stream, err := c.GetAllVM(ctx, &pb.Base{})
		if err != nil {
			fmt.Println(err)
		}
		for {
			article, err := stream.Recv()
			if err == io.EOF {
				break
			}
			for i, b := range strings.Split(article.Storage, ",") {
				if i%2 != 0 {
					if data == b {
						return NodeDISKData{Address: a.IP, Path: b}
					}
				}
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	}
	return NodeDISKData{Address: "0", Path: "0"}
}

func RegistImage(name, group string) (string, bool) {
	data, info, result := VerifyGroup(name, group)
	if result == false {
		return info, false
	}
	fmt.Println(data)
	d, result := ProcessStringToArray(data.User, name, 0)
	if result == false {
		fmt.Println("Error: User is exists this group")
		return "Error: User is exists this group", false
	}

	if db.ChangeDBGroupUser(data.ID, d) {
		return "ok", true
	}
	return "DBChangeError!!", false
}

func UnRegistImgae(name, group string) (string, bool) {
	data, info, result := VerifyGroup(name, group)
	if result == false {
		return info, false
	}
	fmt.Println(data)
	d, result := ProcessStringToArray(data.User, name, 0)
	if result == false {
		fmt.Println("Error: User is exists this group")
		return "Error: User is exists this group", false
	}

	if db.ChangeDBGroupUser(data.ID, d) {
		return "ok", true
	}
	return "DBChangeError!!", false
}

func GetImagePath(d *pb.VMData) string {
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
		r, err := c.ExistImage(ctx, &pb.ImageData{Name: d.Image.GetName(), Tag: d.Image.GetTag()})
		if err != nil {
			fmt.Println(err)
		}
		if r.Result {
			fmt.Println("Generate: " + strconv.Itoa(a.ID) + "/" + r.Path)
			return strconv.Itoa(a.ID) + "/" + r.Path
		}
	}
	fmt.Println("GetImagePath_____|   Not Found!!")
	return ""
}
