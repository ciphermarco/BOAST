package api

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"time"

	app "github.com/ciphermarco/BOAST"
	"github.com/ciphermarco/BOAST/log"
)

// Server represents the API server.
type Server struct {
	Host        string
	Domain      string
	Port        int
	TLSPort     int
	TLSCertPath string
	TLSKeyPath  string
	StatusPath  string
	Storage     app.Storage
}

// ListenAndServe sets the necessary conditions for the underlying http.Server
// to serve the API via HTTPS.
//
// Any errors are returned via the received channel.
func (s *Server) ListenAndServe(err chan error) {
	tlsConfig := &tls.Config{
		MinVersion:               tls.VersionTLS12,
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		},
		CurvePreferences: []tls.CurveID{
			tls.CurveP256, tls.X25519,
		},
	}

	addr := s.Addr(s.TLSPort)
	statusPath := ensureLeadingSlash(url.PathEscape(s.StatusPath))
	r, e := api(s.Domain, statusPath, s.Storage)
	if e != nil {
		err <- e
	}

	srv := &http.Server{
		Addr:         addr,
		Handler:      r,
		TLSConfig:    tlsConfig,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	if statusPath != "" && statusPath != "/" {
		log.Info("Web API Server: status URL is https://%s%s", addr, statusPath)
	}
	log.Info("Web API Server: Listening on https://%s\n", addr)
	err <- srv.ListenAndServeTLS(s.TLSCertPath, s.TLSKeyPath)
}

// Addr returns an address in the format expected by http.Server.
func (s *Server) Addr(port int) string {
	return s.Host + fmt.Sprintf(":%d", port)
}

func ensureLeadingSlash(s string) string {
	l := len(s)
	if l > 0 && s[0] != '/' {
		return "/" + s
	}
	return "/"
}
