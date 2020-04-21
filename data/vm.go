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

type VMDataStruct struct {
	NodeID      int
	Group       string
	GroupID     int
	ID          int
	Name        string
	CPU         int
	Mem         int
	Storage     string
	StoragePath string
	CDROM       string
	StorageType int
	Image       int
	Net         string
	VNC         int
	VNCURL      string
	AutoStart   bool
	Status      int
}

type SpecData struct {
	cpu     int
	mem     int
	storage int
}

type StorageData struct {
	Path   string
	Size   string
	Result bool
}

//mode 0: Administrator 1: SuperUser
func CheckNodeID(isadmin bool, nodeid int) (string, bool) {
	d, result := db.GetDBNodeID(nodeid)
	fmt.Println(d)
	if result == false {
		return "", false
	}
	if isadmin == false && d.OnlyAdmin == 1 {
		return "", false
	}

	return d.IP, true
}

func CheckMaxSpec(d *pb.VMData, s []VMDataStruct) bool {
	data := TotalSpec(s)
	id, r := db.GetDBGroupID(d.Base.GetGroup())
	if r == false {
		fmt.Println("NodeDB Error!!")
		return false
	}
	group, r := db.GetDBGroup(id)
	if r == false {
		fmt.Println("NodeDB Error!!")
		return false
	}
	fmt.Printf("maxmemory: ")
	fmt.Println(group.MaxMem)
	fmt.Printf("maxvm: ")
	fmt.Println(group.MaxVM)
	fmt.Printf("memory: ")
	fmt.Println(strconv.Itoa(data.mem))
	fmt.Printf("vm: ")
	fmt.Println(strconv.Itoa(len(s)))

	if group.MaxMem >= data.mem || group.MaxVM >= len(s) {
		return true
	}
	return false
}

//
// Issue #28
//
func GetVNCVMData(groupid int, vmname string) (string, int) {
	for _, node := range db.GetDBAllNode() {
		conn, err := grpc.Dial(node.IP, grpc.WithInsecure(), grpc.WithBlock())
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
			if err != nil {
				fmt.Println(err)
			}
			ip := strings.Split(node.IP, ":")
			s := strings.Split(article.Vmname, "-")
			if s[0] == strconv.Itoa(groupid) && article.Vmname == vmname {
				return ip[0], int(article.Option.GetVnc())
			}
		}
	}
	return "0", 0
}

func GetAllVMData(group string) []VMDataStruct {
	var d []VMDataStruct

	groupid, result := db.GetDBGroupID(group)
	if result == false {
		fmt.Println("DB Not Found Group!!")
	}

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
			if err != nil {
				fmt.Println(err)
			}
			s := strings.Split(article.Vmname, "-")
			if s[0] == strconv.Itoa(groupid) {

				d = append(d, VMDataStruct{
					NodeID:    a.ID,
					ID:        int(article.Option.Id) + (1000 * a.ID),
					Name:      article.Vmname,
					CPU:       int(article.Vcpu),
					Mem:       int(article.Vmem),
					Net:       article.Vnet,
					VNC:       int(article.Option.Vnc),
					Storage:   article.Storage,
					AutoStart: article.Option.Autostart,
					Status:    int(article.Option.Status),
				})
			}
		}
	}
	return d
}

func TotalSpec(d []VMDataStruct) SpecData {
	var cpu, mem, storage int

	for _, a := range d {
		s := strings.Split(a.Storage, ",")
		for _, b := range s {
			tmp, err := strconv.Atoi(b)
			if err != nil {
				fmt.Println("Error!! string to int")
			}
			storage = storage + tmp
		}
		cpu = cpu + a.CPU
		mem = mem + a.Mem
	}
	return SpecData{cpu: cpu, mem: mem, storage: storage}
}

func GenerateDiskPath(d *pb.VMData) StorageData {
	var size, path string
	var tmpsize, tmppath []string

	data := strings.Split(d.Storage, ",")
	if len(data) == 1 {
		path = GetNodePath(int(d.Node), 1)
		size = d.Storage
	} else {
		for i, a := range data {
			path, err := strconv.Atoi(a)
			if err != nil {
				fmt.Println("string to int false...")
				return StorageData{Result: false}
			}
			if i%2 == 0 {
				tmppath = append(tmppath, "1")
				tmppath = append(tmppath, GetNodePath(int(d.Node), path))
			} else {
				tmpsize = append(tmpsize, a)
			}
		}
		path = strings.Join(tmppath, ",")
		size = strings.Join(tmpsize, ",")
	}
	return StorageData{
		Path: path,
		Size: size,
	}
}

func GetNetworkName(d *pb.VMData) string {
	id, result := db.GetDBGroupID(d.Base.GetGroup())
	if result == false {
		fmt.Println("Error!! Group Not Found.....")
		return ""
	}
	group, result := db.GetDBGroup(id)
	if result == false {
		fmt.Println("Error!! Group Not Found.....")
		return ""
	}
	return "1," + group.Net
}

func GetNodePath(node, path int) string {
	var j int

	d, result := db.GetDBNodeID(node)
	if result == false {
		fmt.Println("GetNodeDB Error!!!")
	}
	data := strings.Split(d.Path, ",")
	for i, a := range data {
		if a == strconv.Itoa(path) {
			j = i + 1
		}
		if i == j {
			return a
		}
	}
	return ""
}
