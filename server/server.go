package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"miikka.xyz/sgoreboard/game"
	"miikka.xyz/sgoreboard/manager"
)

// Server ...
type Server struct {
	// ID
	counter int
	// This gets passed to Game for creating ID
	http  *http.Server
	games map[string]*game.Course
}

// StartingRequest holds data thats needed for starting new game
type StartingRequest struct {
	BasketCount int      `json:"basketCount"`
	Players     []string `json:"players"`
}

// Start ...
func Start(path string) {
	server := Server{}
	router := mux.NewRouter()
	server.http = &http.Server{
		Handler:      router,
		Addr:         "0.0.0.0:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	// Our games/courses
	server.games = make(map[string]*game.Course)

	// Init routes
	router.HandleFunc("/games/{id}/{active:[0-9]+}", server.GetGameHandle).Methods("GET")
	router.HandleFunc("/test_create", server.TestCreate).Methods("POST")
	router.HandleFunc("/test_edit", TestEdit).Methods("POST")
	router.PathPrefix("/").Handler(http.FileServer(http.Dir(path)))
	server.http.ListenAndServe()
}

// GetGameHandle ...
func (s *Server) GetGameHandle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["id"]
	active := vars["active"]

	log.Println("REQUEST", id, active)

	if _, exist := s.games[id]; exist {
		fmt.Fprintf(w, "{}")
		return
	}
	fmt.Fprintf(w, "No game found")
}

// TestCreate ...
func (s *Server) TestCreate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var query StartingRequest
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(bytes, &query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// TODO: Inc only if all is legal
	s.counter++
	course := manager.CreateCourse(query.Players, query.BasketCount, s.counter)
	bytes, err = json.Marshal(course)
	if err != nil {
		fmt.Fprintf(w, "{}")
		return
	}
	log.Println(string(bytes))
	s.games[course.ID] = course
	fmt.Fprintf(w, string(bytes))
}

// TestEdit ...
func TestEdit(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	bytes, err := ioutil.ReadAll(r.Body)
	log.Println(string(bytes))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	course := manager.JSONToCourse(string(bytes))
	log.Printf("%+v\n", course)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "{}")
}

//
//
//
// Old code below, left for example
//
//
//

// CreateGameHandle ...
func (s *Server) CreateGameHandle(w http.ResponseWriter, r *http.Request) {
	// TODO: Mutex here
	g := game.NewCourse()
	err := json.NewDecoder(r.Body).Decode(g)
	if err != nil {
		log.Println(err)
		text(w, http.StatusBadRequest, err.Error())
		return
	}
	s.games[g.ID] = g
	fmt.Fprintf(w, "New Game: %d, %+v, %+v", len(g.Baskets), g, g.Baskets[1])
}

// SetBasketScore ...
func (s *Server) SetBasketScore(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	fmt.Println(vars)
	b := game.Basket{}
	err := json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		log.Println(err)
		fmt.Fprintf(w, err.Error())
		return
	}
	fmt.Fprintf(w, "%+v", b)
}

// HomeHandler ...
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	text(w, 200, "OK")
}

// QueryGame ...
func QueryGame(w http.ResponseWriter, r *http.Request) {
	fmt.Println(mux.Vars(r))
	text(w, 200, "OK")
}

func text(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(code)
	fmt.Fprintln(w, msg)
}
