package listennotes

import (
	"net"
	"net/http"
	"time"
)

var defaultHTTPClient = &http.Client{
	Timeout: time.Second * 30,
	Transport: &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 10 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 10 * time.Second,
	},
}
