package zapi

import (
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

type cors struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	AllowCredentials bool
}

type Public struct {
	url   string
	dir   string
	inUse bool
}

type App struct {
	router *router
	server *http.Server
	public Public
	CORS   cors
}

func New() *App {
	return &App{
		router: &router{},
		public: Public{inUse: false},
		CORS: cors{
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
			AllowedHeaders:   []string{"Content-Type", "Authorization"},
			AllowCredentials: true,
		},
	}
}

func (app *App) Public(url, dir string) {
	app.public.url = url
	app.public.dir = dir
	app.public.inUse = true
}

func (app *App) Get(path string, h HandlerFunc)    { app.router.get(path, h) }
func (app *App) Post(path string, h HandlerFunc)   { app.router.post(path, h) }
func (app *App) Put(path string, h HandlerFunc)    { app.router.put(path, h) }
func (app *App) Delete(path string, h HandlerFunc) { app.router.delete(path, h) }

func (app *App) Run(addr string) {
	app.server = &http.Server{
		Addr:    addr,
		Handler: app,
	}

	app.server.ListenAndServe()
}

func (app *App) ServeHTTP(res http.ResponseWriter, req *http.Request) {

	allowedOrigins := app.CORS.AllowedOrigins
	allowedMethods := app.CORS.AllowedMethods
	allowedHeaders := app.CORS.AllowedHeaders
	allowCredentials := app.CORS.AllowCredentials

	origin := req.Header.Get("Origin")

	allowedOrigin := ""
	if slices.Contains(allowedOrigins, origin) {
		allowedOrigin = origin
	}

	res.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
	res.Header().Set("Access-Control-Allow-Methods", strings.Join(allowedMethods, ", "))
	res.Header().Set("Access-Control-Allow-Headers", strings.Join(allowedHeaders, ", "))

	if allowCredentials {
		res.Header().Set("Access-Control-Allow-Credentials", "true")
	}

	if req.Method == http.MethodOptions {
		res.WriteHeader(http.StatusNoContent)
		return
	}

	// CORS ends here

	path := req.URL.Path
	method := req.Method

	if app.public.inUse && strings.HasPrefix(path, app.public.url+"/") {

		filePath := app.public.dir + strings.Replace(path, app.public.url, "", 1)
		_, err := os.Stat(filePath)
		if err == nil {
			data, _ := os.ReadFile(filePath)
			ext := filepath.Ext(filePath)
			mime, exists := mimeTypes[ext]
			if !exists {
				mime = "application/octet-stream"
			}

			if ext != "" && mime != "" {
				res.Header().Set("Content-Type", mime)
			}

			res.Write(data)
			return
		}
	}

	for _, route := range app.router.routes {
		if !route.regex.MatchString(path) {
			continue
		}

		matches := route.regex.FindStringSubmatch(path)
		if len(matches) == 0 {
			continue
		}

		if route.method == method {
			params := make(map[string]string)
			for i, name := range route.paramNames {
				if i+1 < len(matches) {
					params[name] = matches[i+1]
				}
			}

			wrappedReq := &Req{
				Request: req,
				Params:  params,
			}

			route.handler(wrappedReq, Res{res})
			return
		}

		res.WriteHeader(404)
		res.Write([]byte("404 Not Found !"))
	}
}
