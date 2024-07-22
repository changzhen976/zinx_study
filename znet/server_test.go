package znet

import (
	"fmt"
	"net"
	"testing"
	"time"
	"zinx/ziface"
)

func ClientTest() {
	fmt.Println("[from Client] [START] Client test started")
	time.Sleep(3 * time.Second)

	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("write error err ", err)
		return
	}

	for i := 0; i < 3; i++ {
		_, err := conn.Write([]byte("hello ZINX"))
		if err != nil {
			fmt.Println("write error err ", err)
			return
		}

		buf := make([]byte, 512)
		cnt, err := conn.Read(buf)
		if err != nil {
			fmt.Println("read buf error ")
			return
		}

		fmt.Printf("[from Client] server call back : %s, cnt = %d\n", buf, cnt)

		time.Sleep(1 * time.Second)
	}
	conn.Close()
	time.Sleep(3 * time.Second)

	conn, err = net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("write error err ", err)
		return
	}

	for {
		_, err := conn.Write([]byte("hello ZINX"))
		if err != nil {
			fmt.Println("write error err ", err)
			return
		}

		buf := make([]byte, 512)
		cnt, err := conn.Read(buf)
		if err != nil {
			fmt.Println("read buf error ")
			return
		}

		fmt.Printf("[from Client] server call back : %s, cnt = %d\n", buf, cnt)

		time.Sleep(1 * time.Second)
	}

}

// ping test 自定义路由
type PingRouter struct {
	BaseRouter //一定要先基础BaseRouter
}

// Test PreHandle
func (pr *PingRouter) PreHandle(request ziface.IRequest) {
	fmt.Println("[From Router Handle] Call Router PreHandle")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("before ping ....\n"))
	if err != nil {
		fmt.Println("call back ping ping ping error")
	}
}

// Test Handle
func (pr *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("[From Router Handle] Call PingRouter Handle")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("ping...ping...ping\n"))
	if err != nil {
		fmt.Println("call back ping ping ping error")
	}
}

// Test PostHandle
func (pr *PingRouter) PostHandle(request ziface.IRequest) {
	fmt.Println("[From Router Handle] Call Router PostHandle")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("After ping .....\n"))
	if err != nil {
		fmt.Println("call back ping ping ping error")
	}
}

func TestServer(t *testing.T) {
	s := NewServer("[zinx V0.3]")
	s.AddRouter(&PingRouter{})
	go ClientTest()
	s.Serve()

}
