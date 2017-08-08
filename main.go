package main

import (
	"io"
	"log"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"bytes"
	"os"

	proxyproto "github.com/exosite/proxyprotov2"
)

type ProxyListener struct {
	listener net.Listener
}

func NewProxyListener(listenAddr string) (*ProxyListener, error) {
	ln, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return nil, err
	}

	return &ProxyListener{
		listener: ln,
	}, nil
}

func (self *ProxyListener) Accept() (net.Conn, error) {
	conn, err := self.listener.Accept()
	if err != nil {
		return conn, err
	}
	connInfo := fmt.Sprintf("%s->%s", conn.RemoteAddr().String(), conn.LocalAddr().String())

	proxyInfo, bytesToWrite, err := proxyproto.HandleProxy(conn)
	if err != nil {
		if err == io.EOF {
			// EOF?  Just return.  Screw it.
			return conn, nil
		}
		log.Printf("[%s] Failed to handle proxy protocol: %s", connInfo, err.Error())
		return conn, nil
	}
	if bytesToWrite != nil && len(bytesToWrite) > 0 {
		return conn, fmt.Errorf("Read too much!")
	}
	if proxyInfo != nil {
		for _, tlv := range proxyInfo.TLVs {
			// log.Printf("[%s] TLV 0x%x: %#v", connInfo, tlv.Type, tlv.Value)
			tlsInfo, isTls := tlv.(*proxyproto.TlsTLV)
			if isTls {
				log.Printf("[%s] TLS Version: %s", connInfo, tlsInfo.Version())
				log.Printf("[%s] CN: %s", connInfo, tlsInfo.CN())
				log.Printf("[%s] SNI: %s", connInfo, tlsInfo.SNI())
				if tlsInfo.Certs != nil {
					certs, err := tlsInfo.Certs()
					if err != nil {
						log.Printf("[%s] Failed to parse certificates: %s", err.Error())
					} else {
						for i, cert := range certs {
							log.Printf("[%s] Certificate %d: %s", connInfo, i, cert.Subject.CommonName)
						}
					}
				}
				fp := tlsInfo.Fingerprint()
				if fp != nil {
					var shabuf bytes.Buffer
					for _, part := range fp {
						shabuf.WriteString(fmt.Sprintf("%x", part))
					}
					log.Printf("[%s] Cert SHA1: %s", connInfo, shabuf.String())
				}
			}
		}
	}

	return conn, nil
}

func (self *ProxyListener) Close() error {
	return self.listener.Close()
}

func (self *ProxyListener) Addr() net.Addr {
	return self.listener.Addr()
}

type HandleAll struct {
	Verbose bool
}

func (self *HandleAll) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if self.Verbose {
		log.Printf("%s %s", req.Method, req.URL.String())
	}

	reqDump, err := httputil.DumpRequest(req, true)
	if err != nil {
		log.Printf("Failed to dump request: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, err.Error())
	}
	log.Print(string(reqDump))
	
	w.WriteHeader(http.StatusOK)
	w.Write(reqDump)
}

func main() {
	hAll := HandleAll{
		Verbose: os.Getenv("VERBOSE") == "y",
	}
	var err error
	var ln net.Listener
	if os.Getenv("USE_PROXY_PROTO") == "y" {
		ln, err = NewProxyListener(":8080")
	} else {
		ln, err = net.Listen("tcp", ":8080")
	}
	if err != nil {
		log.Fatalf("Error setting up listener: %s", err.Error())
	}
	log.Fatal(http.Serve(ln, &hAll))
}
