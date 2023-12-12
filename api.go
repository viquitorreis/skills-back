package main

import (
	// "context"

	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	validator "github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type APIServer struct {
	listenAddr string
	store      Storage
}

// pegando o listenAddr como parâmetro e retornar o valor como um novo APIServer
func NewApiServer(listenAddr string, store Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
	}
}

// Inicializando o servidor
func (s *APIServer) Run() {
	// criando router
	router := mux.NewRouter()

	router.HandleFunc("/login", MakeHTTPHandlerFuncHelper(s.handleLogin))
	router.HandleFunc("/account", MakeHTTPHandlerFuncHelper(s.handleAccount))
	router.HandleFunc("/account/{id}", WithJWTAuthHelper(MakeHTTPHandlerFuncHelper(s.handleGetAccountById), s.store))
	router.HandleFunc("/auth/google", MakeHTTPHandlerFuncHelper(s.handleGoogleAuthLogin))

	log.Println("Escutando API JSON na porta:", s.listenAddr)

	http.ListenAndServe(s.listenAddr, router)
}

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
		id, err := GetAccountIdInRequestHelper(r)
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

	validate := validator.New()
	err := validate.Struct(createAccountReq)
	if err != nil {
		errMsg := fmt.Errorf("Not all fields were given %s: ", err)
		fmt.Println(errMsg)
		return WriteJSONHelper(w, http.StatusBadRequest, ApiError{Error: errMsg.Error()})
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
	id, err := GetAccountIdInRequestHelper(r)
	if err != nil {
		return err
	}

	if err := s.store.DeleteAccount(id); err != nil {
		return err
	}

	return WriteJSONHelper(w, http.StatusOK, map[string]int{"user deleted:": id})
}

func (s *APIServer) handleLogin(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "POST" {
		return fmt.Errorf("Method not supported %s", r.Method)
	}

	req := LoginRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}

	acc, err := s.GetAccountByEmailHelper(req.Email)
	if err != nil {
		return err
	}

	if !acc.ValidateHashedPasswordHelper(req.Password) {
		return fmt.Errorf("User not authenticated")
	}

	token, err := CreateJWTHelper(acc)
	if err != nil {
		return err
	}
	resp := LoginResponse{
		Email: acc.Email,
		Token: token,
	}

	return WriteJSONHelper(w, http.StatusOK, resp)
}

func (s *APIServer) handleGoogleAuthLogin(w http.ResponseWriter, r *http.Request) error {

	fmt.Print("Called here")
	if r.Method != "POST" {
		// return
		return WriteJSONHelper(w, http.StatusBadRequest, ApiError{Error: "Method not supported"})
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return WriteJSONHelper(w, http.StatusBadRequest, ApiError{Error: err.Error()})
	}

	bodyString := string(body)
	fmt.Println(bodyString)

	err = r.ParseForm()
	if err != nil {
		return WriteJSONHelper(w, http.StatusBadRequest, ApiError{Error: err.Error()})
	}

	code := r.Form.Get("code")
	fmt.Println("code =>", code)

	http.Redirect(w, r, "http://localhost:4200/dashboard", http.StatusFound)
	return nil
}
