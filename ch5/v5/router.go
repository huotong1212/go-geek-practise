package v5

import (
	"strings"
)

type HandleFunc func(ctx Context)

type matchInfo struct {
	node   *node
	params map[string]string
}

func newRouter() *router {
	return &router{trees: map[string]*node{}}
}

type router struct {
	// key HTTP METHOD node 根节点
	trees map[string]*node
}

func (r *router) addRoute(method, path string, handle HandleFunc) {
	root, ok := r.trees[method]
	if !ok {
		// 如果没有根节点，新建一个根节点，并继续向下执行
		root = &node{path: "/"}
		r.trees[method] = root
	}
	if path == "/" {
		// 如果是根节点，提前返回
		root.handler = handle
		return
	}

	// 前后去掉 /
	path = strings.Trim(path, "/")
	// 支持 /user
	segs := strings.Split(path[:], "/")

	curNode := root
	for _, seg := range segs {
		curNode = curNode.childOrCreate(seg)
	}
	// 只有在最后的时候才将 handle 方法交给当前节点
	curNode.handler = handle
}

// 获取到seg所代表的子节点，如果没有，创建一个
func (n *node) childOrCreate(seg string) *node {
	// 通配符匹配
	if seg == "*" {
		if n.startNode == nil {
			n.startNode = &node{
				path: "*",
			}
		}

		return n.startNode
	}

	// 路径匹配
	if string(seg[0]) == ":" {
		if n.paramNode == nil {
			n.paramNode = &node{
				path: seg,
			}
		}

		return n.paramNode
	}

	// 获取到children，如果没有，创建一个
	if n.children == nil {
		n.children = make(map[string]*node)
	}

	// 判断当前节点是否包含这个节点
	childNode, ok := n.children[seg]
	if !ok {
		childNode = &node{
			path: seg,
		}
		// 将这个节点加入到父节点中
		n.children[seg] = childNode
		// 如果路由有优先级，按照路由优先级查找下级节点
	}

	return childNode
}

func (r *router) findRoute(method, path string) (*matchInfo, bool) {
	root, ok := r.trees[method]
	if !ok {
		return nil, false
	}
	if path == "/" {
		return &matchInfo{node: root}, true
	}

	// 前后去掉 /
	path = strings.Trim(path, "/")
	// 支持 /user
	segs := strings.Split(path[:], "/")

	curNode := root
	for _, seg := range segs {
		// 如果当前节点没有子节点
		if curNode.children == nil {
			// 看看通配符路由中是否有
			if curNode.startNode != nil {
				return &matchInfo{node: curNode.startNode}, true
			}

			// 看看路径匹配中是否有
			if curNode.paramNode != nil {
				params := make(map[string]string)
				params[curNode.paramNode.path[1:]] = seg
				return &matchInfo{node: curNode.paramNode, params: params}, true
			}

			return nil, false
		}

		childNode, ok := curNode.children[seg]
		if !ok {
			// 如果找不到当前子节点

			// 这里再判断一次的原因是 如果注册了 /user/bus/xxx 路由，但是找 /user/* 路由，就会走到这里， user有children，却不是静态匹配
			// 看看通配符路由中是否有
			if curNode.startNode != nil {
				return &matchInfo{node: curNode.startNode}, true
			}

			// 看看路径匹配中是否有
			if curNode.paramNode != nil {
				params := make(map[string]string)
				params[curNode.paramNode.path[1:]] = seg
				return &matchInfo{node: curNode.paramNode, params: params}, true
			}

			return nil, false
		}

		curNode = childNode
	}

	return &matchInfo{node: curNode}, true
}

type node struct {
	path      string
	handler   HandleFunc
	children  map[string]*node // 子节点
	startNode *node            // 通配符匹配 *
	paramNode *node
	//param     string
}
