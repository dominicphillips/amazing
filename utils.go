package amazing

import (
	"net"
	"net/http"
	"time"
)

func TimeoutDialer(cTimeout time.Duration, rwTimeout time.Duration) func(net, addr string) (c net.Conn, err error) {
	return func(netw, addr string) (net.Conn, error) {
		conn, err := net.DialTimeout(netw, addr, cTimeout)
		if err != nil {
			return nil, err
		}
		conn.SetReadDeadline(time.Now().Add(rwTimeout))
		return conn, nil
	}
}

func NewTimeoutClient(connectTimeout time.Duration, readWriteTimeout time.Duration) *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			Dial:  TimeoutDialer(connectTimeout, readWriteTimeout),
		},
	}
}
