package core

import (
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type Proxy struct {
	target    *url.URL
	rules     *RuleSet
	transport http.RoundTripper
}

func NewProxy(target string, rules *RuleSet) (*Proxy, error) {
	parsedURL, err := url.Parse(target)
	if err != nil {
		return nil, err
	}
	return &Proxy{
		target:    parsedURL,
		rules:     rules,
		transport: http.DefaultTransport,
	}, nil
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	path := r.URL.Path

	if rule, ok := p.rules.GetRule("*", path); ok && rule.Action == ActionBlock {
		w.WriteHeader(http.StatusForbidden)
		io.WriteString(w, "Blocked by firewall rule")
		return
	}

	if rule, ok := p.rules.GetRule(method, path); ok {
		if rule.Action == ActionBlock {
			w.WriteHeader(http.StatusForbidden)
			io.WriteString(w, "Blocked by firewall rule")
			return
		}
	}

	proxy := httputil.NewSingleHostReverseProxy(p.target)
	originalDirector := proxy.Director

	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		req.Host = p.target.Host
	}

	proxy.Transport = p.transport
	proxy.ServeHTTP(w, r)
}
