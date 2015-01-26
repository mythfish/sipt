package sipt

import (
	"fmt"
	"github.com/golang/glog"
	"net"
	"sync"
)

type ProxyManager struct {
	localAddr  *string
	remoteAddr *string
	listener   *net.TCPListener
	proxys     map[uint64]*Proxy
	mutex      *sync.Mutex
	connid     uint64
	quit       chan bool
	rules      map[string]*Rule
}

func NewProxyManager(localAddr, remoteAddr *string, nagles bool) *ProxyManager {
	return &ProxyManager{
		localAddr:  localAddr,
		remoteAddr: remoteAddr,
		listener:   nil,
		proxys:     make(map[uint64]*Proxy),
		mutex:      &sync.Mutex{},
		connid:     0,
		quit:       make(chan bool),
		rules:      make(map[string]*Rule),
	}
}

func (pm *ProxyManager) Run() error {
	laddr, err := net.ResolveTCPAddr("tcp", *pm.localAddr)
	if err != nil {
		glog.Error(err.Error())
		return err
	}
	raddr, err := net.ResolveTCPAddr("tcp", *pm.remoteAddr)
	if err != nil {
		glog.Error(err.Error())
		return err
	}

	listener, err := net.ListenTCP("tcp", laddr)
	if err != nil {
		glog.Error(err.Error())
		return err
	}
	pm.listener = listener
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			glog.Errorf("Failed to accept connection '%s'\n", err)

			// When Accept() returns with a non-nill error, we check the quit
			// channel to see if we should continue or quit. If quit, then we quit.
			// Otherwise we continue
			select {
			case <-pm.quit:
				return nil
			default:
				// thanks to martingx on reddit for noticing I am missing a default
				// case. without the default case the select will block.
			}

			continue
		}

		pm.mutex.Lock()
		p := &Proxy{
			lconn:  conn,
			laddr:  laddr,
			raddr:  raddr,
			erred:  false,
			errsig: make(chan bool),
			Connid: pm.connid,
			prefix: fmt.Sprintf("Connection #%03d ", pm.connid),
			rules:  pm.rules,
			mutex:  &sync.Mutex{},
		}
		pm.proxys[pm.connid] = p
		pm.connid++
		pm.mutex.Unlock()

		go func() {
			p.Start()
			pm.mutex.Lock()
			delete(pm.proxys, p.Connid)
			pm.mutex.Unlock()
		}()
	}
}

func (pm *ProxyManager) Stop() {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	for _, p := range pm.proxys {
		p.Stop()
	}

	close(pm.quit)

	pm.listener.Close()
	pm.listener = nil
}

func (pm *ProxyManager) AddRule(key string, r *Rule) error {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()
	pm.rules[key] = r
	return nil
}

func (pm *ProxyManager) RmRule(key string) error {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()
	delete(pm.rules, key)
	return nil
}

//func (pm *ProxyManager) AddRule(connid uint64, key string, r *Rule) error {
//	pm.mutex.Lock()
//	defer pm.mutex.Unlock()
//	p := pm.proxys[connid]
//	if p == nil {
//		err := fmt.Errorf("Fail to add matcher %v, connection #%d not found.", r, connid)
//		return err
//	}
//	p.AddRule(key, r)
//	return nil
//}
//
//func (pm *ProxyManager) RmRule(connid uint64, key string) error {
//	pm.mutex.Lock()
//	defer pm.mutex.Unlock()
//	p := pm.proxys[connid]
//	if p == nil {
//		err := fmt.Errorf("Fail to remove matcher %s, connection #%d not found.", key, connid)
//		return err
//	}
//	p.RmRule(key)
//	return nil
//}
//
//func (pm *ProxyManager) AddRuleToAllConn(key string, r *Rule) {
//	pm.mutex.Lock()
//	defer pm.mutex.Unlock()
//	for _, p := range pm.proxys {
//		p.AddRule(key, r)
//	}
//}
//
//func (pm *ProxyManager) RmRuleFromAllConn(key string) {
//	pm.mutex.Lock()
//	defer pm.mutex.Unlock()
//	for _, p := range pm.proxys {
//		p.RmRule(key)
//	}
//}
