/**
 * @Author: liyalei
 * @Description:jsonrpc实现跨语言的RPC
 * @Version:
 * @Date: 2020/4/24 1:54 下午
 */
package main

import (
	"fmt"
	"gobyexample/midware/rpc/hello/service"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

type HelloServiceClient struct {
	*rpc.Client
}

//var _ service.HelloServiceInterface = (*HelloServiceClient)(nil)

func DialHelloService(network, address string) (*HelloServiceClient, error) {
	conn, err := net.Dial(network, address)
	if err != nil {
		return nil, err
	}
	c := rpc.NewClientWithCodec(jsonrpc.NewClientCodec(conn))
	return &HelloServiceClient{c}, nil
}
//服务具体调用
func (p *HelloServiceClient) Hello(req string, rep *string) error {
	return p.Client.Call(service.HelloServiceName+".Hello", req, rep)
}

func main() {
	client, err := DialHelloService("tcp", ":1234")
	if err != nil {
		log.Fatal("dialing:", err)
	}
	var args = "golang"
	var rep string
	err = client.Hello(args, &rep)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(rep)
}
