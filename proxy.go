package sipt

import (
	"fmt"
	"github.com/golang/glog"
	"io"
	"net"
	"sync"
)

//A proxy represents a pair of connections and their state
type Proxy struct {
	sentBytes     uint64
	receivedBytes uint64
	laddr, raddr  *net.TCPAddr
	lconn, rconn  *net.TCPConn
	Connid        uint64
	nagles        bool
	erred         bool
	errsig        chan bool
	prefix        string
	rules         map[string]*Rule
	mutex         *sync.Mutex
}

//func (p *Proxy) log(s string, args ...interface{}) {
//	if *verbose {
//		glog.Info(p.prefix+s, args...)
//	}
//}

func (p *Proxy) err(s string, err error) {
	if p.erred {
		return
	}
	if err != io.EOF {
		glog.Warningf(p.prefix+s, err)
	}
	p.errsig <- true
	p.erred = true
}

func (p *Proxy) Start() {
	defer p.lconn.Close()
	//connect to remote
	rconn, err := net.DialTCP("tcp", nil, p.raddr)
	if err != nil {
		p.err("Remote connection failed: %s", err)
		return
	}
	p.rconn = rconn
	defer p.rconn.Close()

	if p.nagles {
		p.lconn.SetNoDelay(true)
		p.rconn.SetNoDelay(true)
	}
	//display both ends
	glog.Infof("Opened %s >>> %s", p.lconn.RemoteAddr().String(), p.rconn.RemoteAddr().String())
	//bidirectional copy
	go p.pipe(p.lconn, p.rconn)
	go p.pipe(p.rconn, p.lconn)
	//wait for close...
	<-p.errsig
	glog.Infof("Closed (%d bytes sent, %d bytes recieved)", p.sentBytes, p.receivedBytes)
}

func (p *Proxy) Stop() {
	p.errsig <- true
	p.erred = true
}

func (p *Proxy) pipe(src, dst *net.TCPConn) {
	//data direction
	var f, h string
	islocal := src == p.lconn
	if islocal {
		f = ">>> %d bytes sent%s\n"
	} else {
		f = "<<< %d bytes recieved%s\n"
	}
	h = "%s"

	//directional copy (64k buffer)
	buff := make([]byte, 0xffff)
	for {
		n, err := src.Read(buff)
		if err != nil {
			p.err("Read failed '%s'\n", err)
			return
		}
		b := buff[:n]

		p.mutex.Lock()
		for _, r := range p.rules {
			b = r.Do(b)
		}
		p.mutex.Unlock()

		//show output
		glog.Infof(f, n, "\n"+fmt.Sprintf(h, b))

		//write out result
		n, err = dst.Write(b)
		if err != nil {
			p.err("Write failed '%s'\n", err)
			return
		}
		if islocal {
			p.sentBytes += uint64(n)
		} else {
			p.receivedBytes += uint64(n)
		}
	}
}

func (p *Proxy) AddRule(key string, r *Rule) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.rules[key] = r
}

func (p *Proxy) RmRule(key string) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	delete(p.rules, key)
}
