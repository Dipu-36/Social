package main

import (
	"log"
	"net/http"
	"time"

	"github.com/Dipu-36/social/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)
//dbConfig struct is being used to configure the Database settings in runtime 
type dbConfig struct {
	addr string
	maxOpenConns int
	maxidleConns int 
	maxIdleTime string
}

//config struct holds the configuration for our application
type config struct {
	addr string
	db dbConfig
	env string
}
//application struct holds the dependencies of our application
type application struct {
	config config
	store store.Storage
}

func (app *application) mount() http.Handler{
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	//creates a route group for /v1
	//registers a GET route at /v1/health
	//Converts handler function
	r.Route("/v1", func(r chi.Router){
		r.Get("/health", app.HealthCheckhandler)//this is a sub router that inherits its parent router properties like configuration and prefixed with v1 routes to the specific handler

		r.Route("/posts", func(r chi.Router){

			r.Post("/", app.createPosthandler)

			r.Route("/{postID}", func(r chi.Router){
				
				r.Use(app.postsContextMiddleware)

				r.Get("/", app.getPosthandler)
				r.Delete("/", app.deletePosthandler)
				r.Patch("/", app.updatePosthandler)
			})
		})
	})
	
	return r
}
//Takes the return value of the mount(), then it is passed to the Handler, which then directs the incoming requests to the specific handlers
func (app *application) run(mux http.Handler) error {

	srv := &http.Server{
		Addr: app.config.addr,
		Handler: mux,
		WriteTimeout: time.Second *30,
		ReadTimeout: time.Second *10,
		IdleTimeout: time.Minute,
	}
	log.Printf("server has started at %s", app.config.addr)
	return srv.ListenAndServe()
}