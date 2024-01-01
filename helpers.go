package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

// tipo para converter nossas funções handler em HTTP Handle
type apiFunc func(http.ResponseWriter, *http.Request) error

type ApiError struct {
	Error string
}

// decorator da apiFUnc. Essa função recebe um apiFunc e retorna um http.HandlerFunc
func MakeHTTPHandlerFuncHelper(f apiFunc) http.HandlerFunc {
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

	if w == nil {
		fmt.Println("Error validating data or data is null")
		return nil
	}

	// especificando o status que recebemos no header
	w.WriteHeader(status)

	// pŕecisamos fazer o encode do responseWriter
	return json.NewEncoder(w).Encode(v)
}

func validateSexHelper(sex string) (*Sex, error) {
	validSex := map[string]Sex{
		"male":   Male,
		"female": Female,
		"other":  Other,
	}

	if value, ok := validSex[sex]; ok {
		return &value, nil
	}

	return nil, fmt.Errorf("Invalid value for sex: %s", sex)
}

func validateLanguageHelper(language string) (*Language, error) {
	validLanguage := map[string]Language{
		"en":    En,
		"pt-BR": Br,
	}

	if value, ok := validLanguage[language]; ok {
		return &value, nil
	}
	return nil, fmt.Errorf("Invalid value for language: %s", language)
}

func GetBrazilCurrentTimeHelper() (*time.Location, error) {
	loc, err := time.LoadLocation("America/Sao_Paulo")
	if err != nil {
		return nil, fmt.Errorf("An error occured getting Brazilian time: %v", err.Error())
	}

	return loc, nil
}

func GetAccountIdInRequestHelper(r *http.Request) (int, error) {
	idStr := mux.Vars(r)["id"] // o vars retorna as variáveis de rota que estão na request, se existir algum (pega os parâmetros da request)
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return id, fmt.Errorf("Given ID is invalid %s", idStr)
	}

	return id, nil
}

func (s *APIServer) GetAccountByEmailHelper(email string) (*Account, error) {
	acc, err := s.store.GetAccountByEmail(email)
	if err != nil {
		return nil, err
	}

	return acc, nil
}

func WithJWTAuthHelper(handlerFunc http.HandlerFunc, s Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Calling JWT Helper")

		tokenString := r.Header.Get("x-jwt-token")
		token, err := ValidateJWTHelper(tokenString)
		if err != nil {
			PermissionDeniedHelper(w)
			return
		}

		if !token.Valid {
			PermissionDeniedHelper(w)
			return
		}

		userId, err := GetAccountIdInRequestHelper(r)
		if err != nil {
			PermissionDeniedHelper(w)
			return
		}

		account, err := s.GetAccountById(userId)
		if err != nil {
			PermissionDeniedHelper(w)
			return
		}

		// reivindicando / claim do token JWT
		claims := token.Claims.(jwt.MapClaims)
		if account.ID != int(claims["id"].(float64)) {
			PermissionDeniedHelper(w)
			return
		}

		handlerFunc(w, r)
	}
}

func ValidateJWTHelper(tokenString string) (*jwt.Token, error) {
	secret := os.Getenv("JWT_SECRET")
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})
}

func CreateJWTHelper(account *Account) (string, error) {
	claims := &jwt.MapClaims{
		"expiresAt":     15000,
		"accountNumber": account.ID,
	}

	secret := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}

func PermissionDeniedHelper(w http.ResponseWriter) {
	WriteJSONHelper(w, http.StatusForbidden, ApiError{Error: "Access denied"})
	return
}

func (a *Account) ValidateHashedPasswordHelper(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(a.Password), []byte(password)) == nil // precisamos retornar nil pois a função CompareHash... retorna um erro quando a senha não for igual a do hash
}
