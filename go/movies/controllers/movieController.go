package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"goMovies/common"
	"goMovies/data"
	"gopkg.in/mgo.v2"

	"fmt"
	// "os"
	"github.com/garyburd/redigo/redis"
)

// Handler for HTTP Get - "/movies"
// Returns all Movie documents
func GetMovies(w http.ResponseWriter, r *http.Request) {
	context := NewContext()
	defer context.Close()
	c := context.DbCollection("movies")
	repo := &data.MovieRepository{c}
	movies := repo.GetAll()
	j, err := json.Marshal(MoviesResource{Data: movies})
	if err != nil {
		common.DisplayAppError(w, err, "An unexpected error has occurred", 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

// Handler for HTTP Post - "/movies"
// Insert a new Movie document
func CreateMovie(w http.ResponseWriter, r *http.Request) {
	var dataResourse MovieResource
	// Decode the incoming Movie json
	err := json.NewDecoder(r.Body).Decode(&dataResourse)
	if err != nil {
		common.DisplayAppError(w, err, "Invalid Movie data", 500)
		return
	}
	movie := &dataResourse.Data

	// create new context
	context := NewContext()
	defer context.Close()
	c := context.DbCollection("movies")
	// Insert a movie document
	repo := &data.MovieRepository{c}
	repo.Create(movie)
	j, err := json.Marshal(dataResourse)
	if err != nil {
		common.DisplayAppError(w, err, "An unexpected error has occurred", 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

// Handler for HTTP Get - "/movies/{id}"
// Get movie by id
func GetMovieById(w http.ResponseWriter, r *http.Request) {
	// Get id from incoming url
	vars := mux.Vars(r)
	id := vars["id"]

	// create new context
	context := NewContext()
	defer context.Close()
	c := context.DbCollection("movies")
	repo := &data.MovieRepository{c}

	// Get movie by id
	movie, err := repo.GetById(id)
	if err != nil {
		if err == mgo.ErrNotFound {
			w.WriteHeader(http.StatusNotFound)
			return
		} else {
			common.DisplayAppError(w, err, "An unexpected error has occurred", 500)
			return
		}
	}

	j, err := json.Marshal(MovieResource{Data: movie})
	if err != nil {
		common.DisplayAppError(w, err, "An unexpected error has occurred", 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

// Handler for HTTP Delete - "/movies/{id}"
// Delete movie by id
func DeleteMovie(rw http.ResponseWriter, req *http.Request) {
	// Get id from incoming url
	vars := mux.Vars(req)
	id := vars["id"]

	// Create new context
	context := NewContext()
	defer context.Close()
	c := context.DbCollection("movies")
	repo := &data.MovieRepository{c}

	err := repo.Delete(id)
	if err != nil {
		if err == mgo.ErrNotFound {
			rw.WriteHeader(http.StatusNotFound)
			return
		} else {
			common.DisplayAppError(rw, err, "An unexpected error has occurred", 500)
			return
		}
	}
	rw.WriteHeader(http.StatusNoContent)
}

func Demo(w http.ResponseWriter, r *http.Request) {
	host := "movies"
	fmt.Fprintf(w, "<p>Hi there, from <b>%s</b>!", host)
	c, err := redis.Dial("tcp", "redis:6379")
	if err != nil {
		panic(err)
	}
	defer c.Close()
	c.Do("INCR", host)
	keys, _ := redis.Strings(c.Do("KEYS", "*"))
	fmt.Fprintf(w, "<hr/>")
	fmt.Fprintf(w, "<table style='width: 10em; border-collapse: collapse;'><tr><th style='border: 2px dotted green;'>Container</th><th style='padding: 5px; border: 2px dotted green;'>#</th></tr>")
	for _, key := range keys {
		value, _ := redis.Int(c.Do("GET", key))
		fmt.Fprintf(w, "<tr><td style='border: 1px solid green;'>%s</td>", key)
		fmt.Fprintf(w, "<td style='border: 1px solid green; text-align: center;'>%d</td></tr>", value)
	}
	fmt.Fprintf(w, "</table>")
}
