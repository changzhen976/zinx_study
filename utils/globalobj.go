package utils

import (
	"encoding/json"
	"os"
	"zinx/ziface"
)

type GlobalObj struct {
	TcpServer     ziface.IServer
	Host          string
	TcpPort       int
	Name          string
	Version       string
	MaxPacketSize uint32
	MaxConn       int
}

var GlobalObject *GlobalObj

// 读取用户的配置文件
func (g *GlobalObj) Reload() {
	data, err := os.ReadFile("../conf/zinx.json")
	if err != nil {
		panic(err)
	}
	//将json数据解析到struct中
	//fmt.Printf("json :%s\n", data)
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
}

/*
提供init方法，默认加载
*/
func init() {
	//初始化GlobalObject变量，设置一些默认值
	GlobalObject = &GlobalObj{
		Name:          "ZinxServerApp",
		Version:       "V0.4",
		TcpPort:       7777,
		Host:          "0.0.0.0",
		MaxConn:       12000,
		MaxPacketSize: 4096,
	}

	//从配置文件中加载一些用户配置的参数
	GlobalObject.Reload()
}
