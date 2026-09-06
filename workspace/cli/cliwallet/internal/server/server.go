package server

import "net/http"

func newServer(app *Application) *http.Server {
	mux := http.NewServeMux()
	addRoutes(mux, app)

	var handler http.Handler = mux
	handler = loggerMiddleware(app.logger)(handler)
	handler = requestIDMiddleware(handler)

	server := &http.Server{
		Addr:           app.config.ServerAddress,
		ReadTimeout:    app.config.ServerReadTimeout,
		WriteTimeout:   app.config.ServerWriteTimeout,
		IdleTimeout:    app.config.ServerIdleTimeout,
		MaxHeaderBytes: app.config.ServerMaxHeaderBytes,
		Handler:        handler,
	}
	return server

}
