package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
)

func WriteJson(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}

type APIfunc func(http.ResponseWriter, *http.Request) error

type ApiError struct {
	Error string `json:"error"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

// handling error since handler function does not return error but our api function do
func httpHandleFunc(f APIfunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			// handle err
			WriteJson(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

type APIServer struct {
	listenAddr string
	store      Storage
}

func newAPIServer(listenAddr string, store Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/login", httpHandleFunc(s.handleLogin)).Methods("POST")
	router.HandleFunc("/account", httpHandleFunc(s.handleGetAccounts)).Methods("GET")
	router.HandleFunc("/account/{id}", httpHandleFunc(s.handleGetAccount)).Methods("GET")
	router.HandleFunc("/account", httpHandleFunc(s.handleCreateAccount)).Methods("POST")
	router.HandleFunc("/account/{id}", withJWTAuth(httpHandleFunc(s.handleDeleteAccount), s.store)).Methods("DELETE")
	router.HandleFunc("/transfer", withJWTAuth(httpHandleFunc(s.handleTransfer), s.store)).Methods("POST")

	log.Println("Server is running on Port ", s.listenAddr)
	log.Fatal(http.ListenAndServe(s.listenAddr, router))
}

func withJWTAuth(handlerFunc http.HandlerFunc, s Storage) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		tokenString := r.Header.Get("token")

		token, err := validateJWT(tokenString)

		if err != nil {
			log.Println("invalid token")
			WriteJson(w, http.StatusUnauthorized, ApiError{Error: "permission denied"})
			return
		}

		if !token.Valid {
			log.Println("invalid token")
			WriteJson(w, http.StatusUnauthorized, ApiError{Error: "permission denied"})
			return
		}

		userID, err := getID(r)

		if err != nil {
			log.Println("invalid id")
			WriteJson(w, http.StatusUnauthorized, ApiError{Error: "permission denied"})
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		
		account, err := s.GetAccountByID(userID)

		if account.Number != int64(claims["accountNumber"].(float64)) {
			log.Println("invalid token")
			WriteJson(w, http.StatusForbidden, ApiError{Error: "permission denied"})
			return
		}
		
		if err != nil {
			log.Println("invalid token")
			WriteJson(w, http.StatusForbidden, ApiError{Error: "invalid token"})
			return
		}

		handlerFunc(w, r)
	}
} 

func validateJWT(tokenString string) (*jwt.Token, error) {

	secret := os.Getenv("SECRET")

	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(secret), nil
	})
}

func createJWT(account *Account) (string, error) {
	 claims := &jwt.MapClaims{
		"expiresAt": 15000,
		"accountNumber": account.Number,
	}

	secret := os.Getenv("SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	
	return token.SignedString([]byte(secret))
}

func (s *APIServer) handleLogin(w http.ResponseWriter, r *http.Request) error {
	
	req := new(LoginRequest)
	err := json.NewDecoder(r.Body).Decode(req)

	if err != nil {
		return err
	}

	tokenString, err := s.store.LoginAccount(req)

	if err != nil {
		return err
	}

	return WriteJson(w, http.StatusOK, TokenResponse{Token: tokenString})
}

func (s *APIServer) handleGetAccounts(w http.ResponseWriter, r *http.Request) error {
	accounts, err := s.store.GetAccounts()

	if err != nil {
		return err
	}

	return WriteJson(w, http.StatusOK, accounts)
}

func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)

	if err != nil {
		return err
	}

	account, err := s.store.GetAccountByID(id)

	if err != nil {
		return err
	}

	return WriteJson(w, http.StatusOK, account)
}

func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {

	createAccount := new(CreateAccountRequest)
	if err := json.NewDecoder(r.Body).Decode(createAccount); err != nil {
		return err
	}

	if len(createAccount.Password) < 6 {
		return fmt.Errorf("password should have more than 6 characters")
	}

	account := NewAccount(createAccount.FirstName, createAccount.LastName, createAccount.Password)

	if account == nil {
		return fmt.Errorf("account could not be created")
	}

	accountCreated, err := s.store.CreateAccount(account)

	if err != nil {
		return err
	}

	tokenString, err := createJWT(accountCreated) 

	if err != nil {
		return err
	}

	log.Println("JWT: ", tokenString)
	r.Header.Add("token", tokenString)

	return WriteJson(w, http.StatusOK, TokenResponse{Token: tokenString})
}

func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {

	id, err := getID(r)

	if err != nil {
		return err
	}

	err = s.store.DeleteAccount(id)
	if err != nil {
		return nil
	}

	return nil
}

func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {

	transferData := new(TransferRequest)

	err := json.NewDecoder(r.Body).Decode(transferData)
	if err != nil {
		return err
	}

	return WriteJson(w, http.StatusOK, transferData)
}

func getID(r *http.Request) (int, error) {
	val := mux.Vars(r)["id"]
	id, err := strconv.Atoi(val)

	if err != nil {
		return id, fmt.Errorf("invalid id given")
	}

	return id, nil
}
