package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	app "github.com/ciphermarco/BOAST"
	"github.com/ciphermarco/BOAST/api/httplogger"
	"github.com/ciphermarco/BOAST/log"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/prometheus/procfs"
)

type env struct {
	strg   app.Storage
	proc   procfs.Proc
	domain string
}

func api(domain string, statusPath string, strg app.Storage) (http.Handler, error) {
	e := &env{strg: strg, domain: domain}
	r := chi.NewRouter()

	if e.domain != "" {
		r.Use(e.hostCheck)
	}
	r.Use(httplogger.Logger)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Get("/", e.home)
	r.With(e.authorize).Get("/events", e.events)

	if statusPath != "" && statusPath != "/" {
		p, err := procfs.Self()
		if err != nil {
			return nil, err
		}
		e.proc = p
		r.Get(statusPath, e.status)
	}

	return r, nil
}

func (env *env) hostCheck(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hostHdr := []string{
			// RFC 7239
			"Forwarded",
			// Popular but non-standard
			"X-Forwarded-Host",
		}
		host := r.Host
		for _, hdr := range hostHdr {
			if h := r.Header.Get(hdr); h != "" {
				host = h
			}
		}
		d := strings.Split(host, ":")[0]
		if strings.ToLower(d) != env.domain {
			http.Error(w, http.StatusText(404), 404)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (env *env) home(w http.ResponseWriter, r *http.Request) {
	u := "https://github.com/ciphermarco/BOAST"
	h := fmt.Sprintf("<html><body>BOAST API (<a href=\"%s\">learn more</a>)</body></html>", u)
	fmt.Fprint(w, h)
}

func (env *env) status(w http.ResponseWriter, r *http.Request) {
	statusErr := errors.New("could not access process status")
	check := func(i string, d string, err error) {
		if err != nil {
			log.Info(i)
			log.Debug("%s: %s", d, err)
			render.Render(w, r, errInternalServerError(statusErr))
			return
		}
	}

	stat, err := env.proc.Stat()
	check("could not access process stat", "process stat error", err)

	fdLen, err := env.proc.FileDescriptorsLen()
	check("could not access open file descriptors", "file descriptors len error", err)

	limits, err := env.proc.Limits()
	check("could not access process limits", "process limits error", err)

	res := &statusResponse{
		StoredTests:  env.strg.TotalTests(),
		StoredEvents: env.strg.TotalEvents(),
		RSS:          stat.ResidentMemory(),
		FDLen:        fdLen,
		FDLimit:      limits.OpenFiles,
	}
	render.Render(w, r, res)
}

type statusResponse struct {
	StoredTests  int    `json:"storedTests"`
	StoredEvents int    `json:"storedEvents"`
	RSS          int    `json:"residentSetSizeBytes"`
	FDLen        int    `json:"openFileDescriptors"`
	FDLimit      uint64 `json:"openFileDescriptorsLimit"`
}

func (res *statusResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type ctxKey string

var idCtxKey = ctxKey("id")
var canaryCtxKey = ctxKey("canary")

func (env *env) events(w http.ResponseWriter, r *http.Request) {
	id, idOk := r.Context().Value(idCtxKey).(string)
	canary, canaryOk := r.Context().Value(canaryCtxKey).(string)

	if !idOk || !canaryOk || id == "" || canary == "" {
		log.Info("API /events could not get authorization context keys from context")
		log.Debug("API /events got id from context of type %T", id)
		log.Debug("API /events got canary from context of type %T", canary)

		err := errors.New("internal authentication error")
		render.Render(w, r, errUnauthorized(err))
		return
	}

	if events, exists := env.strg.LoadEvents(id); exists {
		res := &eventsResponse{ID: id, Canary: canary, Events: events}
		render.Render(w, r, res)
	} else {
		res := &eventsResponse{ID: id, Canary: canary, Events: []app.Event{}}
		render.Render(w, r, res)
	}
}

type eventsResponse struct {
	ID     string      `json:"id"`
	Canary string      `json:"canary"`
	Events []app.Event `json:"events"`
}

func (res *eventsResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type errResponse struct {
	Err            error  `json:"-"`
	HTTPStatusCode int    `json:"-"`
	StatusText     string `json:"status"`
	ErrorText      string `json:"error,omitempty"`
}

func (e *errResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func errUnauthorized(err error) render.Renderer {
	return &errResponse{
		Err:            err,
		HTTPStatusCode: http.StatusUnauthorized,
		StatusText:     "Unauthorized",
		ErrorText:      err.Error(),
	}
}

func errInternalServerError(err error) render.Renderer {
	return &errResponse{
		Err:            err,
		HTTPStatusCode: http.StatusInternalServerError,
		StatusText:     "Internal Server Error",
		ErrorText:      err.Error(),
	}
}
