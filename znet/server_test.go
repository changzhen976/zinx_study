package znet

import (
	"fmt"
	"net"
	"testing"
	"time"
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

func TestServer(t *testing.T) {
	s := NewServer("[zinx V0.1]")

	go ClientTest()
	s.Serve()

}
