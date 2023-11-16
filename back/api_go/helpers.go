package main

import (
	"encoding/json"
	"fmt"
	"net/http"
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
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
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
