package v5

import (
	"fmt"
	"net"
	"net/http"
)

func ServerStart() {
	var s Server = &HttpServer{}
	http.ListenAndServe(":8081", s)
}

type Context struct {
	Req        *http.Request
	Resp       http.ResponseWriter
	PathParams map[string]string
}

type Server interface {
	http.Handler
	Start(port string) error
	Server(ctx Context)
	AddRoute(method, path string, handle HandleFunc)
	Get(path string, handle HandleFunc)
	Post(path string, handle HandleFunc)
	findRoute(method, path string) (*matchInfo, bool)
	Group(prefix string) *Group
}

var server Server = &HttpServer{}

type HttpServer struct {
	router *router
}

func (m *HttpServer) Group(prefix string) *Group {
	return &Group{prefix, m}
}

func NewHttpServer() *HttpServer {
	return &HttpServer{router: newRouter()}
}

func (m *HttpServer) Get(path string, handle HandleFunc) {
	m.AddRoute(http.MethodGet, path, handle)
}

func (m *HttpServer) Post(path string, handle HandleFunc) {
	m.AddRoute(http.MethodPost, path, handle)
}

func (m *HttpServer) findRoute(method, path string) (*matchInfo, bool) {
	return m.router.findRoute(method, path)
}

func (m *HttpServer) AddRoute(method, path string, handle HandleFunc) {
	m.router.addRoute(method, path, handle)
}

func (m *HttpServer) Server(ctx Context) {
	//
	matchInfo, ok := m.findRoute(ctx.Req.Method, ctx.Req.URL.Path)
	if !ok || matchInfo.node.handler == nil {
		// 未找到该路由 或者路由没有 handler 方法
		ctx.Resp.Write([]byte("404 Not Found"))
		return
	}

	for k, v := range matchInfo.params {
		ctx.PathParams[k] = v
	}

	matchInfo.node.handler(ctx)
}

func (m *HttpServer) Start(port string) error {
	// 启动前做点事
	listener, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}

	// 端口启动后
	// 注册本服务到管理平台，比如说注册到etcd然后打开管理界面，就能看到这个实例
	err = http.Serve(listener, m)
	fmt.Println("server started")
	//http.ListenAndServe(":8081", m)
	// 启动后做点事
	return err
}

func (m *HttpServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	ctx := Context{
		Req:        request,
		Resp:       writer,
		PathParams: map[string]string{},
	}
	m.Server(ctx)
	//writer.Write([]byte("hello world"))
}

type Group struct {
	prefix string
	s      Server
}

func (g *Group) AddRoute(method, path string, handle HandleFunc) {
	g.s.AddRoute(method, fmt.Sprintf("%s/%s", g.prefix, path), handle)
}
