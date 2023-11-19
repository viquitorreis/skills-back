package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

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
	router.HandleFunc("/account/{id}", makeHTTPHandlerFuncHelper(s.handleGetAccountById))

	log.Println("Escutando API JSON na porta:", s.listenAddr)

	// "escutando" e servirndo a API do server
	http.ListenAndServe(s.listenAddr, router)
}

// criando handler => vai ser um método de APIServer. Vamos poder acessar essa função ao chamar nossa struct APIServer
func (s *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) error {
	// ResponseWriter => vai escrever cabeçalhos / header, corpo da resposta
	// Request => request recebida pelo servidor. Vai ter informações do método, cabeçalho, header etc
	if r.Method == "GET" {
		return s.handleGetAccounts(w, r)
	}

	if r.Method == "POST" {
		return s.handleCreateAccount(w, r)
	}

	return fmt.Errorf("Método não suportado %s", r.Method)
}

func (s *APIServer) handleGetAccountById(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		id, err := GetUserId(r)
		if err != nil {
			return err
		}

		account, err := s.store.GetAccountById(id)
		if err != nil {
			return err
		}

		return WriteJSONHelper(w, http.StatusOK, account)
	}

	if r.Method == "DELETE" {
		return s.handleDeleteAccount(w, r)
	}
	
	return fmt.Errorf("Method not supported %s", r.Method)
}

func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	createAccountReq := CreateAccountRequest{}
	
	if err := json.NewDecoder(r.Body).Decode(&createAccountReq); err != nil {
		return err
	}

	account, err := NewAccount(
		createAccountReq.Email,
		createAccountReq.FullName,
		createAccountReq.Password,
		false,
		string(*createAccountReq.Sex),
		createAccountReq.Country,
		string(*createAccountReq.Language),
	)
	if err != nil {
		return err
	}

	if err := s.store.CreateAccount(account); err != nil {
		return err
	}

	return WriteJSONHelper(w, http.StatusOK, account)
}

func (s *APIServer) handleGetAccounts(w http.ResponseWriter, r *http.Request) error {
	accounts, err := s.store.GetAccounts()
	if err != nil {
		return err
	}

	return WriteJSONHelper(w, http.StatusOK, accounts)
}

func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	id, err := GetUserId(r)
	if err != nil {
		return err
	}

	if err := s.store.DeleteAccount(id); err != nil {
		return err
	}

	return WriteJSONHelper(w, http.StatusOK, map[string]int{"user deleted:": id})
}


func GetUserId(r *http.Request) (int, error) {
	idStr := mux.Vars(r)["id"] // o vars retorna as variáveis de rota que estão na request, se existir algum (pega os parâmetros da request)
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return id, fmt.Errorf("Given ID is invalid %s", idStr)
	}

	return id, nil
}
