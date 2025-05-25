package sec

import (
	"net"
	"net/http"

	"github.com/utkukrl/safegate/internal/core"
)

type IPFilter struct {
	allowedNets []*net.IPNet
}

func NewIPFilter(cidrs []string) (*IPFilter, error) {
	filter := &IPFilter{}
	for _, cidr := range cidrs {
		_, network, err := net.ParseCIDR(cidr)
		if err != nil {
			return nil, err
		}
		filter.allowedNets = append(filter.allowedNets, network)
	}
	return filter, nil
}

func (f *IPFilter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := core.ClientIP(r)
		if ip == nil || !f.isAllowed(ip) {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (f *IPFilter) isAllowed(ip net.IP) bool {
	for _, netw := range f.allowedNets {
		if netw.Contains(ip) {
			return true
		}
	}
	return false
}
