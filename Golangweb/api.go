package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func WriteJSON(w http.ResponseWriter,status int,v any )error{
	w.WriteHeader(status)
	//tells client reposne is json data
	w.Header().Set("Content-type","application/json")
	//converts http to json
	return json.NewEncoder(w).Encode(v)
   }

//Custom function type 
type apiFunc func(http.ResponseWriter,*http.Request ) error

type ApiError struct{
	Error string
}

func checkerr(err error){
	if err!=nil{
		log.Fatal(err) 
	}
}


func makeHTTPHandlefunc (f apiFunc) http.HandlerFunc{
	return func(w http.ResponseWriter,r *http.Request){
 if err:=f(w,r);err!=nil{
WriteJSON(w,http.StatusBadRequest,ApiError{Error :err.Error()})
 }
	}
}

type APIServer struct {
	listenAddr string
}

func NewApiServer(listenAddr string) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
	}
}

func(s *APIServer) Run(){
	router:= mux.NewRouter()
	router.Handlefunc("/account",makeHTTPHandlefunc(s.handleAccount))
	log.Println("Json api running on :",s.listenAddr)
	http.ListenAndServe(s.listenAddr,router)
}

func (s *APIServer) handleAccount(w http.ResponseWriter, r *http.Request)error{
	return nil
}

func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request)error{
	return nil
}

func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request)error{
	return nil
}

func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request)error{
	return nil
}
func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request)error{
	return nil
}