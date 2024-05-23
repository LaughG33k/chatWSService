package app

import (
	"fmt"
	"net/http"
	"time"

	"golang.org/x/net/http2"
)

func initHttpServer(host, port string, rt, rht, wt time.Duration, maxConcStreams, maxHandlers int) (*http.Server, error) {

	httpServer := &http.Server{
		Addr:              fmt.Sprintf("%s:%s", host, port),
		ReadTimeout:       rt,
		ReadHeaderTimeout: rht,
		WriteTimeout:      wt,
	}

	http2Server := &http2.Server{
		MaxHandlers:          maxHandlers,
		MaxConcurrentStreams: uint32(maxConcStreams),
	}

	if err := http2.ConfigureServer(httpServer, http2Server); err != nil {
		return nil, err
	}

	return httpServer, nil

}
