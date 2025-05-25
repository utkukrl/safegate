package handler

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/utkukrl/safegate/internal/core"
	grpc "github.com/utkukrl/safegate/internal/grpcproxy"
)

type Handler struct {
	target     *url.URL
	proxy      *httputil.ReverseProxy
	middleware []core.Middleware
	grpcProxy  *grpc.GRPCProxy
}

func NewHandler(target string, grpcTarget string, timeout time.Duration) (*Handler, error) {
	parsedURL, err := url.Parse(target)
	if err != nil {
		return nil, err
	}

	var grpcProxy *grpc.GRPCProxy
	if grpcTarget != "" {
		grpcProxy = grpc.NewGRPCProxy(grpcTarget, timeout)
	}

	proxy := httputil.NewSingleHostReverseProxy(parsedURL)
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
	}

	return &Handler{
		target:    parsedURL,
		proxy:     proxy,
		grpcProxy: grpcProxy,
	}, nil
}

func (h *Handler) Use(m core.Middleware) {
	h.middleware = append(h.middleware, m)
}
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if h.grpcProxy != nil && r.ProtoMajor == 2 && r.Header.Get("Content-Type") == "application/grpc" {
			h.grpcProxy.ServeHTTP(w, r)
			return
		}
		h.proxy.ServeHTTP(w, r)
	})

	for i := len(h.middleware) - 1; i >= 0; i-- {
		handler = h.middleware[i](handler)
	}

	handler.ServeHTTP(w, r)
}
