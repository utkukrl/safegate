package core

import (
	"net"
	"net/http"
	"strings"
)

type Middleware func(http.Handler) http.Handler

type MiddlewareChain struct {
	middlewares []Middleware
}

func NewMiddlewareChain() *MiddlewareChain {
	return &MiddlewareChain{
		middlewares: []Middleware{},
	}
}

func (m *MiddlewareChain) Use(mw Middleware) {
	m.middlewares = append(m.middlewares, mw)
}

func (m *MiddlewareChain) Wrap(h http.Handler) http.Handler {
	for i := len(m.middlewares) - 1; i >= 0; i-- {
		h = m.middlewares[i](h)
	}
	return h
}

func FirewallMiddleware(rules *RuleSet) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			method := r.Method
			path := r.URL.Path

			if rule, ok := rules.GetRule("*", path); ok && rule.Action == ActionBlock {
				http.Error(w, "Blocked", http.StatusForbidden)
				return
			}

			if rule, ok := rules.GetRule(method, path); ok && rule.Action == ActionBlock {
				http.Error(w, "Blocked", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func IPWhitelistMiddleware(allowedCIDRs []string) Middleware {
	nets := []*net.IPNet{}
	for _, cidr := range allowedCIDRs {
		_, network, err := net.ParseCIDR(cidr)
		if err == nil {
			nets = append(nets, network)
		}
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := ClientIP(r)
			allowed := false
			for _, net := range nets {
				if net.Contains(ip) {
					allowed = true
					break
				}
			}
			if !allowed {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func ClientIP(r *http.Request) net.IP {
	ipStr := r.Header.Get("X-Forwarded-For")
	if ipStr == "" {
		ipStr = r.RemoteAddr
	}
	ip := strings.Split(ipStr, ":")[0]
	return net.ParseIP(ip)
}
