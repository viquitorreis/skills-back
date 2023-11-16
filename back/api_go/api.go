package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// função que response em json
func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	// especificando o status que recebemos no header
	w.WriteHeader(status)
	// especificando no header que vamo retornar em JSON
//	w.Header().Set("Content-Type", "application/json")

	if w == nil {
		fmt.Println("Error validating data or data is null")
		return nil
	}

	// pŕecisamos fazer o encode do responseWriter
	return json.NewEncoder(w).Encode(v)
}

// tipo para converter nossas funções handler em HTTP Handler
type apiFunc func(http.ResponseWriter, *http.Request) error

// tipo pro erro da API
type ApiError struct {
	Error string
}

type APIServer struct {
	listenAddr string
}

// decorator da apiFUnc. Essa função recebe um apiFunc e retorna um http.HandlerFunc
func makeHTTPHandlerFunc(f apiFunc) http.HandlerFunc {
	// a única diferença do HandlerFunc para nossas funções handler é que nossas funções retornam um erro. Precisamos converter isso
	// Vamos converter retornando uma função anônima:
	return func(w http.ResponseWriter, r *http.Request) {
		// administrando o erro
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

// pegando o listenAddr como parâmetro e retornar o valor como um novo APIServer
func NewApiServer(listenAddr string) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
	}
}

 // Inicializando o servidor
func (s *APIServer) Run() {
	// criando router
	router := mux.NewRouter()

	// endpoints
	router.HandleFunc("/account", makeHTTPHandlerFunc(s.handleAccount))

	log.Println("Escutando API JSON na porta:", s.listenAddr)

	// "escutando" e servirndo a API do server
	http.ListenAndServe(s.listenAddr, router)
}

// criando handler => vai ser um método de APIServer. Vamos poder acessar essa função ao chamar nossa struct APIServer
func (s *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) error {
	// ResponseWriter => vai escrever cabeçalhos / header, corpo da resposta
	// Request => request recebida pelo servidor. Vai ter informações do método, cabeçalho, header etc
	if r.Method == "GET" {
		return s.handleGetAccount(w, r)
	}

	if r.Method == "POST" {
		return s.handleCreateAccount(w, r)
	}

	if r.Method == "DELETE" {
		return s.handleDeleteAccount(w, r)
	}

	return fmt.Errorf("Método não suportado %s", r.Method)
}

func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	account, err := NewAccount("victorreis@reis.com", "Victor Reis", "123456", true, "Maleeee", "en")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil
	}

	return WriteJSON(w, http.StatusOK, account)
}

func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	return nil	
}

func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	return nil	
}


