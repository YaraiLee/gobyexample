/**
 * @Author: liyalei
 * @Description: 启动该provider通过curl请求，如下
				 curl localhost:1234/jsonrpc -X POST -d '{"method":"gobyexample/midware/rpc/hello/service/helloservice.Hello","params":["golang"],"id":1}'
 * @Version:
 * @Date: 2020/4/24 1:30 下午
 */
package main

import (
	"gobyexample/midware/rpc/hello/service"
	"io"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"
)

func main() {
	service.RegisterHelloService(new(HelloService))

	http.HandleFunc("/jsonrpc", func(w http.ResponseWriter, r *http.Request) {
		var conn io.ReadWriteCloser = struct {
			io.Writer
			io.ReadCloser
		}{
			w, r.Body,
		}
		rpc.ServeRequest(jsonrpc.NewServerCodec(conn))
	})

	http.ListenAndServe(":1234", nil)
}
