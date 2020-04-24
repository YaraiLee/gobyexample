/**
 * @Author: liyalei
 * @Description: jsonrpc实现跨语言的RPC
 * @Version:
 * @Date: 2020/4/24 1:30 下午
 */
package main

import (
	"gobyexample/midware/rpc/hello/service"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

type HelloService struct {
}
//服务具体实现
func (p *HelloService) Hello(req string, rep *string) error {
	*rep = "hello:" + req
	return nil
}

func main() {
	service.RegisterHelloService(new(HelloService))
	//反射获取类型名称
	//log.Println(reflect.Indirect(reflect.ValueOf(new(HelloService))).Type().Name())

	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal("ListenTcp error:", err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("Accept error:", err)
		}

		//go rpc.ServeConn(conn)
		//json code
		go rpc.ServeCodec(jsonrpc.NewServerCodec(conn))
	}
}
