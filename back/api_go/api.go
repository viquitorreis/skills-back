package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type APIServer struct {
	listenAddr string
	store Storage
}

// pegando o listenAddr como parâmetro e retornar o valor como um novo APIServer
func NewApiServer(listenAddr string, store Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store: store,
	}
}

 // Inicializando o servidor
func (s *APIServer) Run() {
	// criando router
	router := mux.NewRouter()

	// endpoints
	router.HandleFunc("/account", makeHTTPHandlerFuncHelper(s.handleAccount))
	router.HandleFunc("/account/{id}", makeHTTPHandlerFuncHelper(s.handleGetAccount))

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
//	account, err := NewAccount("victorreis@reis.com", "Victor Reis", "123456", true, "Male", "en")
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusBadRequest)
//		return nil
//	}
	vars := mux.Vars(r) // o vars retorna as variáveis de rota que estão na request, se existir algum (pega os parâmetros da request)

	return WriteJSONHelper(w, http.StatusOK, vars)
}

func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	createAccountReq := CreateAccountRequest{}
	
	if err := json.NewDecoder(r.Body).Decode(&createAccountReq); err != nil {
		return err
	}

	fmt.Println(createAccountReq)
	account, err := NewAccount(
		createAccountReq.Email,
		createAccountReq.FullName,
		createAccountReq.Password,
		false,
		string(*createAccountReq.Sex),
		string(*createAccountReq.Language),
	)
	if err != nil {
		return err
	}

	if err := s.store.CreateAccount(account); err != nil {
		return err
	}
fmt.Println(account)
	return WriteJSONHelper(w, http.StatusOK, account)
}

func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	return nil	
}


