package main

import (
	"log"
	"net"
)

//A proxy represents a pair of connections and their state
type Proxy struct {
	sentBytes     uint64
	receivedBytes uint64
	laddr, raddr  *net.TCPAddr
}

type ProxyConn struct {
	sentBytes     uint64
	receivedBytes uint64
	laddr, raddr  *net.TCPAddr
	lconn, rconn  *net.TCPConn
	proxy         *Proxy
}

func NewProxy(laddr, raddr string) (*Proxy, error) {
	la, err := net.ResolveTCPAddr("tcp", laddr)
	if err != nil {
		return nil, err
	}
	ra, err := net.ResolveTCPAddr("tcp", raddr)
	if err != nil {
		return nil, err
	}
	return &Proxy{
		laddr: la,
		raddr: ra,
	}, nil
}

func (p *Proxy) ListenAndServe() error {
	listener, err := net.ListenTCP("tcp", p.laddr)
	if err != nil {
		return err
	}
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			log.Println(err)
			continue
		}
		pc := ProxyConn{
			laddr: p.laddr,
			raddr: p.raddr,
			lconn: conn,
			proxy: p,
		}
		go pc.start()
	}
}

func (p *ProxyConn) start() {
	defer p.lconn.Close()
	//connect to remote
	rconn, err := net.DialTCP("tcp", nil, p.raddr)
	if err != nil {
		log.Printf("Remote connection failed: %s", err)
		return
	}
	p.rconn = rconn
	defer p.rconn.Close()

	// FIXME: may need to set a flag
	p.lconn.SetNoDelay(true)
	p.rconn.SetNoDelay(true)

	//display both ends
	// log.Printf("Opened %s >>> %s", p.lconn.RemoteAddr().String(), p.rconn.RemoteAddr().String())
	//bidirectional copy
	ch1 := p.pipe(p.lconn, p.rconn)
	ch2 := p.pipe(p.rconn, p.lconn)
	//wait for close...
	<-ch1
	<-ch2
	// log.Printf("Closed (%d bytes sent, %d bytes recieved)", p.sentBytes, p.receivedBytes)
}

func (p *ProxyConn) pipe(src, dst *net.TCPConn) chan error {
	//data direction
	errch := make(chan error, 1)
	islocal := src == p.lconn

	//directional copy (64k buffer)
	buff := make([]byte, 0xffff)
	go func() {
		for {
			n, err := src.Read(buff)
			if err != nil {
				errch <- err
				return
			}
			b := buff[:n]

			//write out result
			n, err = dst.Write(b)
			if err != nil {
				errch <- err
				log.Printf("Write failed '%s'\n", err)
				return
			}
			if islocal {
				p.sentBytes += uint64(n)
				p.proxy.sentBytes += uint64(n)
			} else {
				p.receivedBytes += uint64(n)
				p.proxy.receivedBytes += uint64(n)
			}
		}
	}()
	return errch
}
