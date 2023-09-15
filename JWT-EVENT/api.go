package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	database "github.com/karalakrepp/Golang/JWT/Database"
	"github.com/karalakrepp/Golang/JWT/Types"
	"github.com/karalakrepp/Golang/JWT/auth"
	"github.com/karalakrepp/Golang/JWT/middleware"
	"golang.org/x/crypto/bcrypt"
)

type APIServer struct {
	listenAddr string
	storage    database.Storer
}

type ApiError struct {
	Error string `json:"error"`
}

// use for thirthparty packages
type apiFunc func(context.Context, http.ResponseWriter, *http.Request) error

func makeHTTPHandlerFunc(fn apiFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		if err := fn(ctx, w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

// Constructor for APISERVER
func NewServer(listenAddr string, storer database.Storer) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		storage:    storer,
	}
}

// ROUTER AND RUN
func (s *APIServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/account", makeHTTPHandlerFunc(s.handleAccount))
	router.HandleFunc("/login", makeHTTPHandlerFunc(s.handleAccLogin))
	router.HandleFunc("/event/{id}", middleware.RequireAuth(makeHTTPHandlerFunc(s.handlerValidate), s.storage))

	router.HandleFunc("/account/{id}", middleware.RequireAuth(makeHTTPHandlerFunc(s.handleDeleteAccount), s.storage)).Methods(http.MethodDelete)

	fmt.Println("Server is running", s.listenAddr)
	log.Fatal(http.ListenAndServe(s.listenAddr, router))

}

// HANDLE VALIDATING USER EVENT
func (s *APIServer) handlerValidate(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	if r.Method == http.MethodGet {
		return s.handleGetEvent(ctx, w, r)
	}

	if r.Method == http.MethodPost {
		return s.handleCreateEvent(ctx, w, r)
	}
	return WriteJSON(w, http.StatusOK, "Error")

}

func (s *APIServer) handleAccount(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	if r.Method == "POST" {
		return s.handleCreateAccount(ctx, w, r)
	}

	return fmt.Errorf("method not allowed %s", r.Method)

}
func (s *APIServer) handleAccLogin(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	if r.Method == "POST" {
		return s.handleLogin(ctx, w, r)
	}

	return fmt.Errorf("method not allowed %s", r.Method)

}

func (s *APIServer) handleCreateAccount(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	//REQ FOR CREATE USER
	var req = new(Types.CreateAccountRequest)

	//DECODE IT
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		res := Types.Response{
			Message: "Unmarshall Error",
			Error:   err,
		}

		return WriteJSON(w, http.StatusBadRequest, res)
	}

	//password encrype
	encpw, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	if err != nil {
		log.Fatal(err.Error())

		return err

	}

	//DEFİNE USER
	acc, err := Types.NewAccount(req.FirstName, req.LastName, req.Email, string(encpw))

	//GETTING USER ERROR HANDLING
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, err)
	}
	//DB handling

	if err := s.storage.CreateUser(ctx, acc); err != nil {
		res := Types.Response{
			Message: "Creating User Error",
			Error:   err,
		}
		return WriteJSON(w, http.StatusBadRequest, res)
	}

	fmt.Println(acc)

	// Write and Send Json
	return WriteJSON(w, http.StatusOK, acc)

}

func (s *APIServer) handleLogin(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	//create request
	var req = new(Types.LoginRequest)

	//Decode it
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		res := Types.Response{Message: "Unmarshall Error", Error: err}
		return WriteJSON(w, http.StatusBadRequest, res)
	}

	//db handling

	acc, err := s.storage.GetAccountByMail(ctx, req.Email)
	if err != nil {
		res := Types.Response{
			Message: "Email not found",
			Error:   err,
		}

		WriteJSON(w, http.StatusBadRequest, res)
	}

	//pasword validating
	if !acc.ValidatePassword(req.Password) {
		res := Types.Response{
			Message: "Email or Password Error",
			Error:   err,
		}

		WriteJSON(w, http.StatusBadRequest, res)
	}

	//create token
	token, err := auth.CreateJWT(acc)

	if err != nil {
		res := Types.Response{
			Message: "Creating Token Err",
			Error:   err,
		}
		return WriteJSON(w, http.StatusBadGateway, res)

	}

	//response
	resp := Types.LoginResponse{
		Token:  token,
		Email:  acc.Email,
		UserId: acc.ID,
	}

	//send cookie
	cookie := http.Cookie{
		Name:     "Authorization",
		Value:    resp.Token,
		Path:     "",
		Domain:   "",
		MaxAge:   3600 * 24 * 30,
		Secure:   false,
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
	//token gönder

	return WriteJSON(w, http.StatusOK, resp)

}

func (s *APIServer) handleDeleteAccount(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	//GET ID FROM REQ
	id, err := middleware.GetID(r)
	if err != nil {
		return err
	}

	//DB HANDLING
	if err := s.storage.DeleteAccount(ctx, id); err != nil {
		return err
	}

	//ENCODE
	return WriteJSON(w, http.StatusOK, Types.Response{
		Message: fmt.Sprintf("Deleting User id %d", id),
	})
}

// Get event
func (s *APIServer) handleGetEvent(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	//Get Id from request
	id, _ := middleware.GetID(r)
	//db handling
	accounts, err := s.storage.GetEvents(ctx, id)
	if err != nil {
		return err
	}
	//Encode
	return WriteJSON(w, http.StatusOK, accounts)
}
func (s *APIServer) handleCreateEvent(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	//Create Request
	var req = new(Types.CreatedEventReq)

	//Decode It
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return WriteJSON(w, http.StatusBadRequest, "Unmarshall err")
	}
	//define date
	req.Date = time.Now()

	//define acc
	acc := Types.CreateEvent(
		req.Title, req.Description, req.Date, req.EndTime)

	//get id
	userId, err := middleware.GetID(r)

	// error handling
	if err != nil {
		return WriteJSON(w, http.StatusBadRequest, Types.Response{
			Message: "Get Id Err",
			Error:   err,
		})
	}

	//define user
	acc.User_id = userId
	now := time.Now()
	desiredTime := time.Date(now.Year(), now.Month(), now.Day(), 12, 30, 0, 0, now.Location())
	acc.StartTime = desiredTime
	acc.Created_At = time.Now()
	acc.Updated_At = time.Now()

	//DB Handling
	if err := s.storage.CreateEvent(ctx, acc); err != nil {
		return WriteJSON(w, http.StatusBadRequest, "Database register err")
	}

	return WriteJSON(w, http.StatusOK, req)
}

// ENCODING
func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(v)
}
