package main

import (
	"emu/core"
	"emu/rpc"
	"fmt"
	"math/rand"
	"reflect"
	"time"
	"unsafe"
)

func testVarHash() {
	m := make(map[string]int)
	b := []uint8{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	fmt.Println(b)
	fmt.Println(cap(b), len(b))

	m["a"] = 12
	m["b"] = 13
	s := string(b)
	fmt.Printf(" %d,%v", len(s), s)

	m[string(b)] = 17
	_, prs := m[string(b)]
	fmt.Println(prs)
	fmt.Println(m)
}

func testVarHash2() {
	/* convert to string but not efficient*/
	m := make(map[[2]byte]int)
	b := [2]uint8{1, 2}

	m[b] = 17
	_, prs := m[b]
	fmt.Println(prs)
	fmt.Println(m)
}

func testCnt1() {
	var cnt uint64
	var cnt1 float64
	cnt = 17
	cnt1 = 18.1

	c1 := &core.CCounterRec{
		Counter:  &cnt,
		Name:     "A",
		Help:     "an example",
		Unit:     "pkts",
		DumpZero: false,
		Info:     core.ScINFO}
	c2 := &core.CCounterRec{
		Counter:  &cnt1,
		Name:     "B",
		Help:     "an example",
		Unit:     "pkts",
		DumpZero: false,
		Info:     core.ScINFO}
	fmt.Println("val:" + string(c1.MarshalValue()))
	fmt.Println("meta:" + string(c1.MarshalMetaAndVal()))
	db := core.NewCCounterDb("my db")
	db.Add(c1)
	db.Add(c2)
	db.Dump()
	cnt = 18
	cnt1 = 19.31

	fmt.Println()
	fmt.Printf(string(db.MarshalMeta()))
	fmt.Println()
	fmt.Printf(string(db.MarshalValues()))
	fmt.Println()
}

func desc(i interface{}) {

	switch v := i.(type) {
	case *int:
		val := reflect.ValueOf(i)
		elm := val.Elem()
		a := (*int)(unsafe.Pointer(elm.Addr().Pointer()))
		fmt.Printf("int  %d\n", *a)
	case *float32:
		fmt.Printf("float32\n")
	case *float64:
		fmt.Printf("float64\n")
	default:
		fmt.Printf("I don't know about type %T!\n", v)
	}
}

type MapPortT map[uint16]bool

func testTunnelKey() {
	var portMap MapPortT
	portMap = make(MapPortT)
	portMap[1] = true
	portMap[2] = true
	fmt.Println(portMap)
	v, ok := portMap[3]
	if ok {
		fmt.Println(v)
	} else {
		fmt.Println("3 is not there")
	}

}

const (
	eARP_TIMER1 = 1
	eARP_TIMER2 = 2
)

type ArpExample struct {
	ctx    *core.TimerCtx
	timer1 core.CHTimerObj
	timer2 core.CHTimerObj
}

func (o *ArpExample) StartTimers() {
	o.timer1.SetCB(o, eARP_TIMER1, nil)
	o.timer2.SetCB(o, eARP_TIMER2, nil)
	o.ctx.Start(&o.timer1, time.Second*1)
	o.ctx.Start(&o.timer2, time.Second*2)
}

func (o *ArpExample) OnEvent(a, b interface{}) {

	vala, ok := a.(int)
	if !ok {
		fmt.Printf("error ")
		return
	}
	if vala == eARP_TIMER1 {
		fmt.Printf("timer a")
		o.ctx.Start(&o.timer1, time.Second*1)
	}
	if vala == eARP_TIMER2 {
		fmt.Printf("timer b")
		o.ctx.Start(&o.timer2, time.Second*2)
	}
}

func TestNs1() {
	var arp ArpExample
	tctx := core.NewTimerCtx()
	arp.ctx = tctx
	arp.StartTimers()
	fmt.Printf("start  \n")
	for {
		select {
		case <-tctx.Timer.C:
			if tctx.Ticks%100 == 0 {
				fmt.Printf("\n %d \n", tctx.Ticks)
			}
			tctx.RestartTimer()
		}
	}

	fmt.Printf("end  \n")
}

//func (o *DList) RemoveSelffunc (o *DList) RemoveSelf(n *DList) {
//(n *DList) {
type MyObjectTest struct {
	val   uint32
	dlist core.DList
}

//There is no better solotion in go right now! maybe go2.0
func covert(dlist *core.DList) *MyObjectTest {
	var s MyObjectTest
	return (*MyObjectTest)(unsafe.Pointer(uintptr(unsafe.Pointer(dlist)) - unsafe.Offsetof(s.dlist)))
}

func testdList() {
	var head core.DList
	head.SetSelf()

	var it core.DListIterHead

	for i := 0; i < 10; i++ {
		o := new(MyObjectTest)
		o.val = uint32(i)
		head.AddLast(&o.dlist)
	}

	//head.RemoveNode(&ptr.dlist)
	for it.Init(&head); it.IsCont(); it.Next() {
		fmt.Println(covert(it.Val()).val)
	}
}

/*
type A struct {
	a uint32
}

type B struct {
	b uint32
}

type PluginEvents interface {
	OnEvent(msg string, a, b interface{})
}

type PluginBase struct {
	a     *A           // pointer to ns
	b     *B           // pointer to the thread
	i     PluginEvents //
	Ext   interface{}  // extention memory that can be converted to specific memory
	NsExt interface{}
}

type PluginArp struct {
	PluginBase
	arpEnable bool
}

func (o *PluginArp) OnEvent(a, b interface{}) {
	fmt.Printf(" OnEvent %v %v \n", a, b)
	fmt.Printf(" values %v %v arp:%v \n", o.a.a, o.b.b, o.arpEnable)
}

type InterfaceMem struct {
	a *interface{}
	b *interface{}
	c *interface{}
}

func RunArp(o *PluginBase) {
	o.i.OnEvent(27, 28)
	p := o.Ext.(PluginArp)
	fmt.Printf(" %v  \n", p)
}

func testPlugin() {
	var a A
	var b B
	a.a = 11
	b.b = 12
	var plugArp PluginArp
	plugArp.a = &a
	plugArp.b = &b
	plugArp.i = &plugArp
	plugArp.Ext = plugArp // point to itself
	plugArp.arpEnable = true
	RunArp(&plugArp.PluginBase)
}

*/

func main() {
	//testPlugin()
	//testdList()
	//TestNs1()
	return
	var i interface{}
	var cnt int
	cnt = 17
	i = cnt
	val, _ := i.(int)

	fmt.Printf("(%v, %T)\n", val, val)
	//desc(i)

	//desc(i)
	//fmt.Printf("(%v, %T)\n", i, i)

	//i = 42.21
	//desc(i)
	//fmt.Printf("(%v, %T)\n", i, i)

	//testVarHash2()
	//testPkt2()
	//testpcapWrite3()
	//testpcapWrite2()
	//testPkt1()
	//testSlice()
	//testMemory1()
	//testMemory2()
	//testpcapWrite()
	//testMpool2()
	//testStats()
	return
	var data *[]byte
	d := make([]byte, 10)
	data = &d
	fmt.Println(*data)
	(*data)[0] = 1
	fmt.Println(*data)

	//b := []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	//fmt.Println(b)
	//fmt.Println(b[0:1])

	return
	rand.Seed(time.Now().UnixNano())
	rpc.RcpCtx.Create(4510)
	rpc.RcpCtx.StartRxThread()

	for {
		select {
		case req := <-rpc.RcpCtx.GetC():
			rpc.RcpCtx.HandleReqToChan(req)
		}
	}

	rpc.RcpCtx.Delete()
	return
}
