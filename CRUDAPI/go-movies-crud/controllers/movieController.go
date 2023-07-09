package controllers

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/karalakrepp/golang/models"
)

var movies []models.Movie

func GetMovies(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(movies)
}

func GetMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)

	for _, movie := range movies {
		if movie.ID == params["id"] {
			json.NewEncoder(w).Encode(movie)
			return
		}
	}

}

func CreateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var movie models.Movie
	_ = json.NewDecoder(r.Body).Decode(&movie) //body kısmından movie ye yazılacak

	movie.ID = strconv.Itoa(rand.Intn(10000))

	// movie postmandan json olarak gelicek decode oyuzden yapıldı
	movies = append(movies, movie)

	json.NewEncoder(w).Encode(movie)
}

func UpdateMovie(w http.ResponseWriter, r *http.Request) {

	//set json content type
	w.Header().Set("Content-Type", "application/json")

	//params tanımla
	params := mux.Vars(r)
	//döngüde ara
	for index, item := range movies {
		if item.ID == params["id"] {
			movies = append(movies[:index], movies[index+1:]...)
			var movie models.Movie
			_ = json.NewDecoder(r.Body).Decode(&movie)
			movie.ID = params["id"]

			// movie postmandan json olarak gelicek decode oyuzden yapıldı
			movies = append(movies, movie)

			json.NewEncoder(w).Encode(movie)

			return

		}
	}

}

func DeleteMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	for index, movie := range movies {
		if movie.ID == params["id"] {
			movies = append(movies[:index], movies[index+1:]...) // yeni movies indexten sonrasını ve öncesini gösterir
			break
		}
	}
	json.NewEncoder(w).Encode(movies)
}
