package proxy

import (
	"io"
	"net"
	"net/http"
)

func HijackTCP(w http.ResponseWriter, r *http.Request, backend string) error {
	hj, ok := w.(http.Hijacker)
	if !ok {
		return http.ErrNotSupported
	}
	clientConn, _, err := hj.Hijack()
	if err != nil {
		return err
	}

	backendConn, err := net.Dial("tcp", backend)
	if err != nil {
		clientConn.Close()
		return err
	}

	if err := r.Write(backendConn); err != nil {
		backendConn.Close()
		clientConn.Close()
		return err
	}
	go io.Copy(backendConn, clientConn)
	go io.Copy(clientConn, backendConn)
	return nil
}
