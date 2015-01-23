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
	matcher       map[string]*Matcher
	replacer      map[string]*Replacer
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
		glog.Warning(p.prefix+s, err)
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
	glog.Info("Opened %s >>> %s", p.lconn.RemoteAddr().String(), p.rconn.RemoteAddr().String())
	//bidirectional copy
	go p.pipe(p.lconn, p.rconn)
	go p.pipe(p.rconn, p.lconn)
	//wait for close...
	<-p.errsig
	glog.Info("Closed (%d bytes sent, %d bytes recieved)", p.sentBytes, p.receivedBytes)
}

func (p *Proxy) pipe(src, dst *net.TCPConn) {
	//data direction
	var f, h string
	islocal := src == p.lconn
	if islocal {
		f = ">>> %d bytes sent%s"
	} else {
		f = "<<< %d bytes recieved%s"
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
		//execute match
		for _, m := range p.matcher {
			m.match(b)
		}
		//execute replace
		for _, r := range p.replacer {
			b = r.replace(b)
		}
		p.mutex.Unlock()

		//show output
		glog.Info(f, n, "\n"+fmt.Sprintf(h, b))

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

func (p *Proxy) AddMatcher(key *string, m *Matcher) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.matcher[*key] = m
}

func (p *Proxy) DelMatcher(key *string) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.matcher[*key] = nil
}

func (p *Proxy) AddReplacer(key *string, r *Replacer) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.replacer[*key] = r
}

func (p *Proxy) DelReplacer(key *string) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.replacer[*key] = nil
}
