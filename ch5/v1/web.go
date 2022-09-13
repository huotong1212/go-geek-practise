package v1

import "net/http"

func Start() {
	var s Server = &MyServer{}
	http.ListenAndServe(":8081", s)
}

type Server interface {
	http.Handler
}

type MyServer struct {
}

func (m *MyServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte("hello world"))
}
