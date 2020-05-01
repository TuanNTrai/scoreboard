package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"miikka.xyz/scoreboard/game"
	"miikka.xyz/scoreboard/geo"
)

const (
	maxBaskets   = 36
	maxPlayers   = 5
	maxPlayerLen = 10
	maxGames     = 10000
)

// New creates new server
func New(path string) *Server {
	server := &Server{}
	router := mux.NewRouter()
	server.HTTP = &http.Server{
		Handler:      router,
		Addr:         "0.0.0.0:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	server.games = make(map[string]*game.Course)

	file, err := ioutil.ReadFile(path + "courses.json")
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal([]byte(file), &server.courses)
	if err != nil {
		log.Fatal(err)
	}

	router.HandleFunc("/games/create", server.CreateGameHandle).Methods("POST")
	router.HandleFunc("/games/edit", server.EditGameHandle).Methods("POST")
	router.HandleFunc("/games/{id}", server.GetGameHandle).Methods("GET")
	router.PathPrefix("/").Handler(http.FileServer(http.Dir(path + "public")))

	return server
}

// GetGameHandle returns game by id
func (s *Server) GetGameHandle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["id"]

	if _, exist := s.games[id]; !exist {
		http.Error(w, "Error", http.StatusInternalServerError)
		return
	}

	bytes, err := json.Marshal(s.games[id])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, string(bytes))
}

func isValid(createQuery CreateRequest, maxPlayers int, maxBaskets int, maxPlayerLen int) bool {
	if len(createQuery.Players) > maxPlayers || createQuery.BasketCount > maxBaskets {
		return false
	}

	for _, player := range createQuery.Players {
		if len(player) > maxPlayerLen {
			return false
		}
	}
	return true
}

// CreateGameHandle creates new game
func (s *Server) CreateGameHandle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if len(s.games) > maxGames {
		http.Error(w, "Server if full", http.StatusTooManyRequests)
		return
	}

	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var query CreateRequest
	err = json.Unmarshal(bytes, &query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !isValid(query, maxPlayers, maxBaskets, maxPlayerLen) {
		http.Error(w, "Invalid data", http.StatusInternalServerError)
		return
	}

	if s.counter >= 10000 {
		s.counter = 0
	}
	s.mu.Lock()
	s.counter++
	s.mu.Unlock()

	var course *game.Course
	for _, info := range s.courses {
		m := geo.Distance(query.Lat, query.Lon, info.Lat, info.Lon)
		if m < 1000 && m > 0 {
			course = game.CreateExistingCourse(query.Players, query.BasketCount, s.counter, info.Pars, info.ShortName)
			fmt.Println("created", info.ShortName)
			break
		}
	}

	if course == nil {
		course = game.CreateCourse(query.Players, query.BasketCount, s.counter)
		fmt.Println("created default (all par 3)")
	}

	bytes, err = json.Marshal(course)
	var c *game.Course
	json.Unmarshal(bytes, &c)
	if err != nil {
		fmt.Fprintf(w, "{}")
		return
	}
	s.games[course.ID] = course
	fmt.Fprintf(w, string(bytes))
}

// EditGameHandle updates game on server
func (s *Server) EditGameHandle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var c *game.Course
	err = json.Unmarshal(bytes, &c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id := c.ID
	if _, found := s.games[id]; !found {
		http.Error(w, "Game not found", http.StatusInternalServerError)
		return
	}

	if s.games[id].Active > s.games[id].BasketCount || s.games[id].Active < 1 {
		fmt.Fprintf(w, string(bytes))
		return
	}

	// If editedAt is fraud, we can still delete game with orginal createdAt
	temp := s.games[id].CreatedAt
	s.games[id] = c
	s.games[id].CreatedAt = temp

	resp, err := json.Marshal(s.games[id])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, string(resp))
}

// CleanGames removes old games
func (s *Server) CleanGames() {
	for {
		time.Sleep(20 * time.Minute)
		s.remove()
	}
}

func (s *Server) remove() {
	for id, game := range s.games {
		if time.Since(game.EditedAt) > time.Hour*1 || time.Since(game.CreatedAt) > (time.Hour*5) {
			log.Println(id, "deleted")
			s.mu.Lock()
			delete(s.games, id)
			s.mu.Unlock()
		}
	}
}
