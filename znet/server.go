package znet

import (
	"errors"
	"fmt"
	"net"
	"time"
	"zinx/utils"
	"zinx/ziface"
)

type Server struct {
	Name      string
	IPVersion string
	IP        string
	Port      int
	Router    ziface.IRouter
}

// ============== 定义当前客户端链接的handle api ===========
func CallBackToClient(conn *net.TCPConn, data []byte, cnt int) error {
	//回显业务
	fmt.Println("[from Server] [Conn Handle] CallBackToClient ... ")
	if _, err := conn.Write(data[:cnt]); err != nil {
		fmt.Println("write back buf err ", err)
		return errors.New("CallBackToClient error")
	}
	return nil
}

func (s *Server) Start() {
	fmt.Printf("[from Server] [START] Server listenner at IP: %s, Port %d, is starting\n", s.IP, s.Port)
	go func() {
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr err: ", err)
			return
		}

		listenner, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen", s.Port, "  err: ", err)
			return
		}

		fmt.Println("[from Server] [Running] start Zinx server  ", s.Name, " succ!")

		var cid uint32 = 0
		// cid = 0

		for {
			conn, err := listenner.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err ", err)
				return
			}

			//3.3 处理该新连接请求的 业务 方法， 此时应该有 handler 和 conn是绑定的
			dealConn := NewConntion(conn, cid, s.Router)
			cid++
			fmt.Println("current cid :", cid)

			//3.4 启动当前链接的处理业务
			go dealConn.Start()

		}

	}()
}

func (s *Server) Stop() {
	// s.Stop()
}

func (s *Server) Serve() {
	s.Start()
	for {
		time.Sleep(10 * time.Second)
	}
}

func (s *Server) AddRouter(router ziface.IRouter) {
	s.Router = router
}

func NewServer(name string) ziface.IServer {
	// func NewServer(name string) *Server {

	utils.GlobalObject.Reload()
	s := &Server{
		Name:      utils.GlobalObject.Name,
		IPVersion: "tcp4",
		IP:        utils.GlobalObject.Host,
		Port:      utils.GlobalObject.TcpPort,
		Router:    nil,
	}
	return s
}
