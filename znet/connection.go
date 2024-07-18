package znet

import (
	"fmt"
	"net"
	"zinx/ziface"
)

type Connection struct {
	Conn         *net.TCPConn
	ConnID       uint32
	isClosed     bool
	handleAPI    ziface.HandFunc
	ExitBuffChan chan bool
}

func NewConntion(conn *net.TCPConn, connID uint32, callback_api ziface.HandFunc) *Connection {
	c := &Connection{
		Conn:         conn,
		ConnID:       connID,
		isClosed:     false,
		handleAPI:    callback_api,
		ExitBuffChan: make(chan bool, 1),
	}
	return c
}

func (c *Connection) StartReader() {
	fmt.Println("[from Server] [StartReader] new connection catch, Reader Goroutine is running")
	defer fmt.Println(c.Conn.RemoteAddr().String(), "conn reader exit!")
	defer c.Stop()

	for {
		//读取我们最大的数据到buf中
		buf := make([]byte, 512)
		cnt, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("recv buf err ", err)
			c.ExitBuffChan <- true
			// continue
			break
		}
		//调用当前链接业务(这里执行的是当前conn的绑定的handle方法)
		if err := c.handleAPI(c.Conn, buf, cnt); err != nil {
			fmt.Println("connID ", c.ConnID, " handle is error")
			c.ExitBuffChan <- true
			return
		}
	}

}

// 启动连接，让当前连接开始工作
func (c *Connection) Start() {

	//开启处理该链接读取到客户端数据之后的请求业务
	go c.StartReader()

	for {
		select {
		case <-c.ExitBuffChan:
			//得到退出消息，不再阻塞
			return
		}
	}
}

// 停止连接，结束当前连接状态M
func (c *Connection) Stop() {
	//1. 如果当前链接已经关闭
	if c.isClosed {
		return
	}
	c.isClosed = true

	//TODO Connection Stop() 如果用户注册了该链接的关闭回调业务，那么在此刻应该显示调用

	// 关闭socket链接
	c.Conn.Close()

	//通知从缓冲队列读数据的业务，该链接已经关闭
	c.ExitBuffChan <- true

	//关闭该链接全部管道
	close(c.ExitBuffChan)
	fmt.Println("[from Server] Conn is closed.")
}

// 从当前连接获取原始的socket TCPConn
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

// 获取当前连接ID
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

// 获取远程客户端地址信息
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}
