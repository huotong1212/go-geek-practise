package v1

import (
	"net"
	"net/http"
)

func ServerStart() {
	var s Server = &HttpServer{}
	http.ListenAndServe(":8081", s)
}

type Context struct {
	Req  *http.Request
	Resp http.ResponseWriter
}

type HandleFunc interface {
	http.Handler
}

type Server interface {
	http.Handler
	Start() error
	AddRoute(method, path string, handleFunc HandleFunc)
	Server(ctx Context)
}

type HttpServer struct {
}

func (m *HttpServer) AddRoute(method, path string, handleFunc HandleFunc) {
	//TODO implement me
	panic("implement me")
}

func (m *HttpServer) Server(ctx Context) {
	//TODO implement me
	panic("implement me")
}

func (m *HttpServer) Start() error {
	// 启动前做点事
	listener, err := net.Listen("tcp", ":8081")
	if err != nil {
		return err
	}

	// 端口启动后
	// 注册本服务到管理平台，比如说注册到etcd然后打开管理界面，就能看到这个实例
	return http.Serve(listener, m)
	//http.ListenAndServe(":8081", m)
	// 启动后做点事
}

func (m *HttpServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	ctx := Context{
		Req:  request,
		Resp: writer,
	}
	m.Server(ctx)
	writer.Write([]byte("hello world"))
}

// 装饰器模式
type HttpsServer struct {
	HttpServer
	CertFile string
	KeyFile  string
}

func (m *HttpsServer) Start() error {
	// 启动前做点事
	listener, err := net.Listen("tcp", ":8081")
	if err != nil {
		return err
	}

	// 端口启动后
	// 注册本服务到管理平台，比如说注册到etcd然后打开管理界面，就能看到这个实例
	return http.ServeTLS(listener, m, m.CertFile, m.KeyFile)
	//http.ListenAndServe(":8081", m)
	// 启动后做点事
}
