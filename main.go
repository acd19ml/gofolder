package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-openapi/runtime/middleware"

	"github.com/nicholasjackson/env"

	"github.com/acd19ml/gofolder/handlers"

	"github.com/gorilla/mux"
)

var bindAddress = env.String("BIND_ADDRESS", false, ":9090", "Bind address for the server")

func main() {

	env.Parse()
	// This is how we create a new instance of our handler
	// os.Stdout: is a stream that we can write to
	// "product-api": is a prefix that will appear before any log message
	// log.LstdFlags: is a flag that specifies the logging properties
	l := log.New(os.Stdout, "product-api", log.LstdFlags) //logger

	// Create the handlers
	ph := handlers.NewProducts(l) //products handler

	// Create a new serve mux and register the handlers
	sm := mux.NewRouter() //serve mux

	getRouter := sm.Methods(http.MethodGet).Subrouter() //create a subrouter for GET requests
	getRouter.HandleFunc("/", ph.GetProducts)

	putRouter := sm.Methods(http.MethodPut).Subrouter() //create a subrouter for PUT requests
	putRouter.HandleFunc("/{id:[0-9]+}", ph.UpdateProducts)
	putRouter.Use(ph.MiddlewareProductValidation)

	postRouter := sm.Methods(http.MethodPost).Subrouter() //create a subrouter for POST requests
	postRouter.HandleFunc("/", ph.AddProduct)
	postRouter.Use(ph.MiddlewareProductValidation)

	deleteRouter := sm.Methods(http.MethodDelete).Subrouter() //create a subrouter for DELETE requests
	deleteRouter.HandleFunc("/{id:[0-9]+}", ph.DeleteProduct)

	opts := middleware.RedocOpts{SpecURL: "/swagger.yaml"}
	sh := middleware.Redoc(opts, nil)

	getRouter.Handle("/docs", sh)
	getRouter.Handle("/swagger.yaml", http.FileServer(http.Dir("./"))) //serve the swagger.yaml file
	

	// Create a new server
	s := http.Server{
		Addr:         *bindAddress, //":9090
		Handler:      sm,
		ErrorLog:     l,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// start the server
	go func() {
		err := s.ListenAndServe()
		if err != nil {
			l.Printf("Error starting server: %s\n", err)
			os.Exit(1) //1: means an error
		}
	}()

	//This will block until we get a signal
	//ctrl + c: will send an interrupt signal
	//'chan': is a channel
	//'Notify': will notify the channel when we get the specified signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan
	l.Println("Received terminate, graceful shutdown", sig)

	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(tc)

	http.ListenAndServe(":9090", sm)
}
