package main

import (
  "fmt"
  "net/http"
  "reflect"

  "github.com/gorilla/mux"
  "github.com/urfave/negroni"
)

func main() {
  r := mux.NewRouter()
  r.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
    fmt.Fprintf(w, "Welcome to the home page!",reflect.TypeOf(r))
  })

  n := negroni.Classic() // Includes some default middlewares

  // add jwt authentication
  JWTMiddleware := JWTMiddlewareNew()
  n.Use(negroni.HandlerFunc(JWTMiddleware.ServeHTTP))

  // add recovery middleware 
  n.Use(negroni.NewRecovery())

  RegisterRequests(r)
  /* add middlewares here, since router is the last one */
  n.UseHandler(r)

  http.ListenAndServe(":3000", n)
  fmt.Printf("Web server exits")
}

func getBackendServerUrl (w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  w.WriteHeader(http.StatusOK)
  fmt.Fprintf(w, "backend " + vars["backendid"] + " server " + vars["serverid"] + " url is:")
  return
}

func RegisterRequests (r *mux.Router) {
  r.HandleFunc("/traefik/backends/{backendid}/servers/{serverid}/url", getBackendServerUrl).Methods("GET")
  r.HandleFunc("/traefik/backends/{backendid}/servers/{serverid}/weight", getBackendServerUrl).Methods("GET")
}