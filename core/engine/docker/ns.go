package docker

import (
	"agent/api"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"time"
)

// 镜像结构
type Image struct {
	Created     uint64
	Id          string
	ParentId    string
	RepoTags    []string
	Size        uint64
	VirtualSize uint64
}

// 容器结构
type Container struct {
	Id              string                 `json:"Id"`
	Names           []string               `json:"Names"`
	Image           string                 `json:"Image"`
	ImageID         string                 `json:"ImageID"`
	Command         string                 `json:"Command"`
	Created         uint64                 `json:"Created"`
	State           string                 `json:"State"`
	Status          string                 `json:"Status"`
	Ports           []Port                 `json:"Ports"`
	Labels          map[string]string      `json:"Labels"`
	HostConfig      map[string]string      `json:"HostConfig"`
	NetworkSettings map[string]interface{} `json:"NetworkSettings"`
	Mounts          []Mount                `json:"Mounts"`
}

// docker 端口映射
type Port struct {
	IP          string `json:"IP"`
	PrivatePort int    `json:"PrivatePort"`
	PublicPort  int    `json:"PublicPort"`
	Type        string `json:"Type"`
}

// docker 挂载
type Mount struct {
	Type        string `json:"Type"`
	Source      string `json:"Source"`
	Destination string `json:"Destination"`
	Mode        string `json:"Mode"`
	RW          bool   `json:"RW"`
	Propatation string `json:"Propagation"`
}

// 连接列表
var SockAddr = "/var/run//docker.sock"                                //这可不是随便写的，是docker官网文档的套接字默认值，当然守护进程通讯方式还有tcp,fd等方式，各自都有适用场景。。。
var imagesSock = "GET /images/json HTTP/1.0\r\n\r\n"                  //docker对外的镜像api操作
var containerSock = "GET /containers/json?all=true HTTP/1.0\r\n\r\n"  //docker对外的容器查看api
var startContainerSock = "POST /containers/%s/start HTTP/1.0\r\n\r\n" //docker对外的容器启动api

func (self *DockerKnow) RunDockerKnow() {
	lname, _ := os.Readlink("/proc/1/ns/pid")
	log.Printf("宿主机命名空间ID: %s\n", lname)
	namespaceID := lname[5 : len(lname)-1]
	dockerinfo := api.DockerInfo{
		ContainerID:   "localhost",
		ContainerName: "宿主机",
	}
	self.Set(namespaceID, dockerinfo)

	for {
		// 轮询docker
		err := self.listenDocker()
		if err != nil {
		//	log.Println(err.Error())
			time.Sleep(100 * time.Second)
		}
		time.Sleep(10 * time.Second)
	}
}

func (self *DockerKnow) listenDocker() error {
	// 获取容器列表,拿到所有的容器信息
	containers, err := readContainer()
	if err != nil {
		return errors.New("read container error: " + err.Error())
	}
	for _, container := range containers {

		if container.Status[:2] == "Up" {
			inspectinfo, err := readInspect(container.Id[:12])

			if err != nil {
				fmt.Println(err)
				break
			}
			lname, err := os.Readlink(fmt.Sprintf("/proc/%d/ns/pid", inspectinfo.State.Pid))
			namespaceID := lname[5 : len(lname)-1]
			dockerinfo := api.DockerInfo{
				ContainerID:   container.Id[:12],
				ContainerName: container.Names[0],
			}
			self.Set(namespaceID, dockerinfo)
			//log.Printf("id=%s, name=%s, state=%s pid = %d ===>namespace ID %s", container.Id[:12], container.Names, container.Status,inspectinfo.State.Pid,namespaceID)
		}

	}
	return nil
}

// 获取 unix sock 连接
func connectDocker() (*net.UnixConn, error) {
	addr := net.UnixAddr{SockAddr, "unix"} // SockAddr 这个变量的值被设定为docker的/var/run/docker 套接字路径值，也就是说此处就是拨通与docker的daemon通讯建立的关键处,其他处的代码就是些正常的逻辑处理了
	return net.DialUnix("unix", nil, &addr)
}

// 获取Inspect
func readInspect(id string) (*Inspect, error) {
	conn, err := connectDocker() //建立一个unix连接,这其实是一个关键点，需要你了解unix 套接字 建立连接
	if err != nil {
		return nil, errors.New("connect error: " + err.Error())
	}
	var inspectSock = fmt.Sprintf("GET /containers/%s/json HTTP/1.0\r\n\r\n", id)
	_, err = conn.Write([]byte(inspectSock))
	if err != nil {
		return nil, err
	}
	result, err := ioutil.ReadAll(conn)
	if err != nil {
		return nil, err
	}
	body := getBody(result)
	var inspect = &Inspect{}
	err = json.Unmarshal(body, &inspect)
	if err != nil {
		return nil, err
	}
	//fmt.Println(inspect)
	return inspect, nil

}

// 获取容器列表
func readContainer() ([]Container, error) {
	conn, err := connectDocker() //建立一个unix连接,这其实是一个关键点，需要你了解unix 套接字 建立连接
	if err != nil {
		return nil, errors.New("connect error: " + err.Error())
	}
	_, err = conn.Write([]byte(containerSock))
	if err != nil {
		return nil, err
	}
	result, err := ioutil.ReadAll(conn)
	if err != nil {
		return nil, err
	}
	body := getBody(result)
	var containers []Container
	err = json.Unmarshal(body, &containers)
	if err != nil {
		return nil, err
	}
	//log.Println("len of containers: ", containers)
	if len(containers) == 0 {
		return nil, errors.New("no containers")
	}
	return containers, nil
}

// 启动容器
func startContainer(id string) error {
	conn, err := connectDocker()
	if err != nil {
		return errors.New("connect error: " + err.Error())
	}
	start := fmt.Sprintf(startContainerSock, id)
	fmt.Println(start)
	cmd := []byte(start)
	code, err := conn.Write(cmd)
	if err != nil {
		return err
	}
	log.Println("start container response code: ", code)
	// 启动容器等待20秒，防止数据重发
	time.Sleep(20 * time.Second)
	return nil
}

// 获取镜像列表
func readImage(conn *net.UnixConn) ([]Image, error) {
	_, err := conn.Write([]byte(imagesSock))
	if err != nil {
		return nil, err
	}
	result, err := ioutil.ReadAll(conn)
	if err != nil {
		return nil, err
	}
	body := getBody(result[:])
	var images []Image
	err = json.Unmarshal(body, &images)
	if err != nil {
		return nil, err
	}
	return images, nil
}

// 从返回的 http 响应中提取 body
func getBody(result []byte) (body []byte) {
	for i := 0; i <= len(result)-4; i++ {
		if result[i] == 13 && result[i+1] == 10 && result[i+2] == 13 && result[i+3] == 10 {
			body = result[i+4:]
			break
		}
	}
	return
}
