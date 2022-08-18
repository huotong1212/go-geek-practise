package ch1

import (
	"fmt"
	"testing"
)

type Family struct {
	name, value any
}

func TestFamily(t *testing.T) {
	f := &Family{}
	f.name = "小明"
	f.value = "HongKong"
	fmt.Println(f)
}

// 装饰器模式：装饰器模式一般用在基类功能封装不错，但使用的时候需要对功能进行一些加强，而这些加强版的功能也会被其它加强版需要，这种就比较适合。
/**
飞行器接口
*/
type Aircraft interface {
	fly()
	landing()
}

/**
普通直升机
*/
type Helicopter struct{}

func (h *Helicopter) fly() {
	fmt.Println("飞行功能...")
}

func (h Helicopter) landing() {
	fmt.Println("着陆功能...")
}

/**
武装直升机
*/
type WeaponAircraft struct {
	Aircraft
}

func (h *WeaponAircraft) fly() {
	h.Aircraft.fly()
	fmt.Println("添加武器功能...")
}

/**
救援直升机
*/
type RescueAircraft struct {
	Aircraft
}

func (a *RescueAircraft) fly() {
	a.Aircraft.fly()
	fmt.Println("添加救援功能...")
}

func TestDecorator(t *testing.T) {
	//w := WeaponAircraft{}
	//w.landing() // runtime error: invalid memory address or nil pointer dereference

	fmt.Println("武装直升机==============")
	w := WeaponAircraft{&Helicopter{}}
	w.fly()
	w.landing()

	fmt.Println("武装救援直升机==============")
	r := RescueAircraft{&w}
	r.fly()
	r.landing()
}

/**
总结：
1.需要一个抽象接口
2.基类实现这个接口
3.装饰类组合这个接口
4.装饰类重新接口中的需要装饰的方法，然后调用接口中的方法
*/
