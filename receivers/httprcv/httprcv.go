package httprcv

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"

	app "github.com/ciphermarco/BOAST"
	"github.com/ciphermarco/BOAST/log"
)

// Receiver represents the HTTP protocol receiver.
type Receiver struct {
	Name        string
	Host        string
	Ports       []int
	TLSPorts    []int
	TLSCertPath string
	TLSKeyPath  string
	IPHeader    string
	Storage     app.Storage
}

// ListenAndServe sets the necessary conditions for the underlying http.Server
// to serve the HTTP and/or the HTTPS server for each configured port.
//
// Any errors are returned via the received channel.
func (r *Receiver) ListenAndServe(err chan error) {
	http.HandleFunc("/", catchAll(r.Storage, r.IPHeader))

	for _, port := range r.Ports {
		go func(p int) {
			r.serveHTTP(p, err)
		}(port)
	}

	for _, port := range r.TLSPorts {
		go func(p int) {
			r.serveHTTPS(p, err)
		}(port)
	}
}

func (r *Receiver) serveHTTP(port int, err chan error) {
	addr := r.Addr(port)
	srv := &http.Server{
		Addr:         addr,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Info("%s: Listening on http://%s\n", r.Name, addr)
	err <- srv.ListenAndServe()
}

func (r *Receiver) serveHTTPS(port int, err chan error) {
	tlsConfig := &tls.Config{
		PreferServerCipherSuites: true,
		CurvePreferences: []tls.CurveID{
			tls.CurveP256, tls.X25519,
		},
	}

	addr := r.Addr(port)
	srv := &http.Server{
		Addr:         addr,
		TLSConfig:    tlsConfig,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Info("%s: Listening on https://%s\n", r.Name, addr)
	err <- srv.ListenAndServeTLS(r.TLSCertPath, r.TLSKeyPath)
}

// Addr returns an address in the format expected by http.Server.
func (r *Receiver) Addr(port int) string {
	return r.Host + fmt.Sprintf(":%d", port)
}

func catchAll(strg app.Storage, ipHdr string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("HTTP event received")
		dump, err := httputil.DumpRequest(r, true)
		if err != nil {
			log.Info("Could not dump HTTP request event")
			log.Debug("Dump HTTP request error: %v", err)
			errCode := http.StatusInternalServerError
			errTxt := http.StatusText(errCode)
			http.Error(w, errTxt, errCode)
			return
		}

		// Does the request contain any known test ID (id)?
		searchDumpID := func(k, v string) bool {
			return strings.Contains(string(dump), k)
		}
		id, canary := strg.SearchTest(searchDumpID)
		if id == "" || canary == "" {
			log.Debug("HTTP event test not found: id=\"%s\" canary=\"%s\"",
				id, canary)
			u := "https://github.com/ciphermarco/BOAST"
			h := fmt.Sprintf("<html><body>BOAST (<a href=\"%s\">learn more</a>)</body></html>", u)
			fmt.Fprint(w, h)
			return
		}

		// HTTP or HTTPS event?
		rcv := "HTTP"
		if r.TLS != nil {
			rcv = "HTTPS"
		}

		// Real IP header?
		remoteAddr := r.RemoteAddr
		if realIP := r.Header.Get(ipHdr); realIP != "" {
			remoteAddr = realIP
		}

		// Try to create and store the event
		evt, err := app.NewEvent(id, rcv, remoteAddr, string(dump))
		if err != nil {
			log.Info("Error creating a new HTTP event")
			log.Debug("New HTTP event error: %v", err)
		} else {
			if err := strg.StoreEvent(evt); err != nil {
				log.Info("Error storing a new HTTP event")
				log.Debug("Store HTTP event error: %v", err)
			} else {
				log.Info("New HTTP event stored")
			}
			log.Debug("HTTP event object:\n%s", evt.String())
		}

		// Respond the canary to the client
		fmt.Fprintf(w, "<html><body>%s</body></html>", canary)
	}
}
