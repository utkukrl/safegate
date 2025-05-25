package grpcproxy

import (
	"net/http"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCProxy struct {
	target  string
	timeout time.Duration
}

func NewGRPCProxy(target string, timeout time.Duration) *GRPCProxy {
	return &GRPCProxy{
		target:  target,
		timeout: timeout,
	}
}

func (p *GRPCProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := grpc.NewClient(p.target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		http.Error(w, "Failed to connect to gRPC backend", http.StatusBadGateway)
		return
	}
	defer conn.Close()
	http.Error(w, "gRPC proxy is not fully implemented", http.StatusNotImplemented)
}
