package ch3

type MyInterface interface {
	Abc()
}

// ok
var _ MyInterface = Player{}

// ok
var _ MyInterface = &Player{}

type Player struct{}

func (receiver Player) Abc() {}

// ok
var _ MyInterface = &Buyer{}

// failed
//var _ MyInterface = Buyer{}

type Buyer struct{}

func (receiver *Buyer) Abc() {}
