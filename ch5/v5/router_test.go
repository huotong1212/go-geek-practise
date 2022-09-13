package v5

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"reflect"
	"testing"
)

func (r router) equal(y *router) (string, bool) {
	for k, v := range r.trees {
		yv, ok := y.trees[k]
		if !ok {
			return fmt.Sprintf("目标 router 里面没有方法 %s 的路由树", k), false
		}
		str, ok := v.equal(yv)
		if !ok {
			return k + "-" + str, ok
		}
	}
	return "", true
}

func (n *node) equal(y *node) (string, bool) {
	if y == nil {
		return "目标节点为 nil", false
	}
	if n.path != y.path {
		return fmt.Sprintf("%s 节点 path 不相等 x %s, y %s", n.path, n.path, y.path), false
	}

	nhv := reflect.ValueOf(n.handler)
	yhv := reflect.ValueOf(y.handler)
	if nhv != yhv {
		return fmt.Sprintf("%s 节点 handler 不相等 x %s, y %s", n.path, nhv.Type().String(), yhv.Type().String()), false
	}

	if len(n.children) != len(y.children) {
		return fmt.Sprintf("%s 子节点长度不等", n.path), false
	}
	if len(n.children) == 0 {
		return "", true
	}

	for k, v := range n.children {
		yv, ok := y.children[k]
		if !ok {
			return fmt.Sprintf("%s 目标节点缺少子节点 %s", n.path, k), false
		}
		str, ok := v.equal(yv)
		if !ok {
			return n.path + "-" + str, ok
		}
	}
	return "", true
}

func Test_router_addRoute(t *testing.T) {

	tests := []struct {
		name string

		method string
		path   string
	}{
		{
			method: http.MethodGet,
			path:   "/",
		},
		{
			method: http.MethodGet,
			path:   "/user",
		},
		{
			method: http.MethodGet,
			path:   "/user/detail/profile",
		},
		{
			method: http.MethodGet,
			path:   "/user/detail/image",
		},
		{
			method: http.MethodPost,
			path:   "/order/cancel",
		},
		{
			method: http.MethodGet,
			path:   "/user/info/*",
		},
		{
			method: http.MethodGet,
			path:   "/user/:phone",
		},
	}
	var handle HandleFunc = func(ctx Context) {

	}

	wantRouter := &router{
		trees: map[string]*node{
			http.MethodGet: &node{
				path:    "/",
				handler: handle,
				children: map[string]*node{
					"user": &node{
						path:    "user",
						handler: handle,
						children: map[string]*node{
							"detail": {
								path: "detail",
								children: map[string]*node{
									"profile": {
										path:    "profile",
										handler: handle,
									},
									"image": {
										path:    "image",
										handler: handle,
									},
								},
							},
							"info": {
								path: "info",
								startNode: &node{
									path:    "*",
									handler: handle,
								},
							},
						},
						paramNode: &node{
							path:    ":phone",
							handler: handle,
						},
					},
				},
			},
			http.MethodPost: &node{
				path: "/",
				children: map[string]*node{
					"order": &node{
						path: "order",
						children: map[string]*node{
							"cancel": &node{
								path:    "cancel",
								handler: handle,
							},
						},
					},
				},
			},
		},
	}

	res := &router{
		trees: map[string]*node{},
	}

	for _, tc := range tests {
		res.addRoute(tc.method, tc.path, handle)
	}

	errStr, ok := wantRouter.equal(res)
	//fmt.Println(errStr, ok)
	// 判断ok为true errStr 为false时的错误打印
	assert.True(t, ok, errStr)

	findCases := []struct {
		name   string
		method string

		path       string
		found      bool
		wantPath   string
		handleFunc bool
	}{
		{
			name:       "/",
			method:     http.MethodGet,
			path:       "/",
			found:      true,
			wantPath:   "/",
			handleFunc: true,
		},
		{
			name:       "/user/detail",
			method:     http.MethodGet,
			path:       "/user/detail",
			found:      false,
			wantPath:   "detail",
			handleFunc: false,
		},
	}

	for _, tc := range findCases {
		t.Run(tc.name, func(t *testing.T) {
			matchInfo, ok := res.findRoute(http.MethodGet, tc.path)
			assert.True(t, ok)
			if !ok {
				return
			}
			assert.Equal(t, tc.wantPath, matchInfo.node.path)
			assert.Equal(t, tc.handleFunc, matchInfo.node.handler != nil)
		})
	}

}
