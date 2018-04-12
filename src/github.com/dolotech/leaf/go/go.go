package g

import (
	"github.com/dolotech/leaf/conf"
	"github.com/golang/glog"
	"runtime"
	"github.com/dolotech/lib/grpool"
)

// one Go per goroutine (goroutine not safe)
//type Go struct {
//	//ChanCb    chan func()
//	pendingGo int
//}

/*
type LinearGo struct {
	f  func()
	cb func()
}
*/

/*type LinearContext struct {
	g              *Go
	linearGo       *list.List
	mutexLinearGo  sync.Mutex
	mutexExecution sync.Mutex
}*/

//func New(l int) *Go {
//	g := new(Go)
	//g.ChanCb = make(chan func(), l)
	//return g
//}

var pool = grpool.NewPool(runtime.NumCPU()*2, 1024*10)
//func (g *Go) Go(f func(), cb func()) {
func Go(f func(), cb func()) {
	//g.pendingGo++
	pool.JobQueue <- func() {
		defer func() {
			pool.JobQueue  <- cb
			if r := recover(); r != nil {
				if conf.LenStackBuf > 0 {
					buf := make([]byte, conf.LenStackBuf)
					l := runtime.Stack(buf, false)
					glog.Error("%v: %s", r, buf[:l])
				} else {
					glog.Error("%v", r)
				}
			}
		}()
		f()
	}
}

/*func  Cb(cb func()) {
	defer func() {
		//g.pendingGo--
		if r := recover(); r != nil {
			if conf.LenStackBuf > 0 {
				buf := make([]byte, conf.LenStackBuf)
				l := runtime.Stack(buf, false)
				glog.Error("%v: %s", r, buf[:l])
			} else {
				glog.Error("%v", r)
			}
		}
	}()

	if cb != nil {
		cb()
	}
}*/
/*

func (g *Go) Close() {
	for g.pendingGo > 0 {
		g.Cb(<-g.ChanCb)
	}
}

func (g *Go) Idle() bool {
	return g.pendingGo == 0
}
*/

/*

func (g *Go) NewLinearContext() *LinearContext {
	c := new(LinearContext)
	c.g = g
	c.linearGo = list.New()
	return c
}

func (c *LinearContext) Go(f func(), cb func()) {
	c.g.pendingGo++

	c.mutexLinearGo.Lock()
	c.linearGo.PushBack(&LinearGo{f: f, cb: cb})
	c.mutexLinearGo.Unlock()

	go func() {
		c.mutexExecution.Lock()
		defer c.mutexExecution.Unlock()

		c.mutexLinearGo.Lock()
		e := c.linearGo.Remove(c.linearGo.Front()).(*LinearGo)
		c.mutexLinearGo.Unlock()

		defer func() {
			c.g.ChanCb <- e.cb
			if r := recover(); r != nil {
				if conf.LenStackBuf > 0 {
					buf := make([]byte, conf.LenStackBuf)
					l := runtime.Stack(buf, false)
					glog.Error("%v: %s", r, buf[:l])
				} else {
					glog.Error("%v", r)
				}
			}
		}()

		e.f()
	}()
}
*/