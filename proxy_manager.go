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
	proxys     map[uint64]*Proxy
	mutex      *sync.Mutex
	connid     uint64
}

func NewProxyManager(localAddr, remoteAddr *string, nagles bool) *ProxyManager {
	return &ProxyManager{
		localAddr:  localAddr,
		remoteAddr: remoteAddr,
		proxys:     make(map[uint64]*Proxy),
		mutex:      &sync.Mutex{},
		connid:     0,
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
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			fmt.Printf("Failed to accept connection '%s'\n", err)
			continue
		}

		pm.mutex.Lock()
		p := &Proxy{
			lconn:    conn,
			laddr:    laddr,
			raddr:    raddr,
			erred:    false,
			errsig:   make(chan bool),
			Connid:   pm.connid,
			prefix:   fmt.Sprintf("Connection #%03d ", pm.connid),
			matcher:  nil,
			replacer: nil,
		}
		pm.proxys[pm.connid] = p
		pm.connid++
		pm.mutex.Unlock()

		go func() {
			p.Start()
			pm.mutex.Lock()
			pm.proxys[p.Connid] = nil
			pm.mutex.Unlock()
		}()
	}
}
