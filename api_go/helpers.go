package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// tipo para converter nossas funções handler em HTTP Handle
type apiFunc func(http.ResponseWriter, *http.Request) error

type ApiError struct {
	Error string
}

// decorator da apiFUnc. Essa função recebe um apiFunc e retorna um http.HandlerFunc
func makeHTTPHandlerFuncHelper(f apiFunc) http.HandlerFunc {
	// a única diferença do HandlerFunc para nossas funções handler é que nossas funções retornam um erro. Precisamos converter isso
	// Vamos converter retornando uma função anônima:
	return func(w http.ResponseWriter, r *http.Request) {
		// administrando o erro
		if err := f(w, r); err != nil {
			WriteJSONHelper(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

func WriteJSONHelper(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	// especificando o status que recebemos no header
	w.WriteHeader(status)

	if w == nil {
		fmt.Println("Error validating data or data is null")
		return nil
	}

	// pŕecisamos fazer o encode do responseWriter
	return json.NewEncoder(w).Encode(v)
}

func getBrazilCurrentTimeHelper() (*time.Location, error) {
	loc, err := time.LoadLocation("America/Sao_Paulo")
	if err != nil {
		return nil, fmt.Errorf("An error occured getting Brazilian time: %v", err.Error())
	}

	return loc, nil
}

func withJWTAuthHelper(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		fmt.Println("Calling JWT Helper")
		handlerFunc(w, r)
	}
}
