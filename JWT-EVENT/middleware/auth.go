package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	database "github.com/karalakrepp/Golang/JWT/Database"
	"github.com/karalakrepp/Golang/JWT/auth"
)

// Permision Denied
func permisionDenied(w http.ResponseWriter, status int) error {
	w.WriteHeader(status)

	res := "Permision Denied"
	return json.NewEncoder(w).Encode(res)
}

// Middleware For token authorization
func RequireAuth(next http.HandlerFunc, s database.Storer) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		//get cookies
		cookie, err := r.Cookie("Authorization")
		if err != nil {
			//Check Error
			permisionDenied(w, http.StatusBadRequest)
			return
		}
		//get token from cookie
		tokenstring := cookie.Value

		//validate token for auth

		token, err := auth.ValidateJWT(tokenstring)

		if err != nil {

			permisionDenied(w, http.StatusBadRequest)

			return
		}

		if !token.Valid {
			permisionDenied(w, http.StatusBadRequest)
			return
		}

		//Getting ID
		id, err := GetID(r)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode("Can not get ID")
			return
		}

		//Check email from claims
		claims := token.Claims.(jwt.MapClaims)
		ctx := context.Background()
		fmt.Println(claims)
		account, err := s.GetUserById(ctx, id)
		if err != nil {
			permisionDenied(w, http.StatusBadRequest)
			return
		}

		if account.Email != claims["email"] {
			permisionDenied(w, http.StatusBadRequest)
			log.Println("İts not your account")
			return
		}
		next(w, r)
	}
}

// get ıd from request
func GetID(r *http.Request) (int, error) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return id, fmt.Errorf("invalid id given %s", idStr)
	}
	return id, nil
}
