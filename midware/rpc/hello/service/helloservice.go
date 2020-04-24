/**
 * @Author: liyalei
 * @Description:服务声明
 * @Version:
 * @Date: 2020/4/24 1:26 下午
 */
package service

import "net/rpc"

const HelloServiceName = "gobyexample/midware/rpc/hello/service/helloservice"

//声明服务
type HelloServiceInterface interface {
	Hello(request string, reply *string) error
}

//服务注册
func RegisterHelloService(svc HelloServiceInterface) error {
	return rpc.RegisterName(HelloServiceName, svc)
}
