package ch4

import (
	_ "database/sql"
	_ "github.com/silenceper/pool"
	"testing"
)

type name struct {
}

func (n name) Scan(src any) error {
	//TODO implement me
	panic("implement me")
}

func TestPool(t *testing.T) {
	//pool.Pool()
	//sql.DB{}
}

type channelPool struct {
	// 空闲的连接
	//conns chan *idleConn
	// 阻塞的请求
	//connReqs []chan connReq
}
