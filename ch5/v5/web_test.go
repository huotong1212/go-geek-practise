package v5

import (
	"fmt"
	"net/http"
	"testing"
)

func TestWeb(t *testing.T) {
	s := NewHttpServer()

	s.Get("/", func(ctx Context) {
		ctx.Resp.Write([]byte("hello world"))
	})

	s.Get("/user", func(ctx Context) {
		ctx.Resp.Write([]byte("hello user"))
	})

	s.Get("/user/*", func(ctx Context) {
		ctx.Resp.Write([]byte("通配符匹配"))
	})

	s.Get("/user/home/:id", func(ctx Context) {
		s2 := ctx.PathParams["id"]
		ctx.Resp.Write([]byte(fmt.Sprintf("路径匹配 %s", s2)))
	})

	g := s.Group("cart")
	g.AddRoute(http.MethodGet, "info", func(ctx Context) {
		ctx.Resp.Write([]byte("路由分组"))
	})

	s.Start(":8081")
}
