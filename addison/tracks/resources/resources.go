package resources

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"tracks/repository"
)

func updateTrack(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	var c repository.Track
	if err := json.NewDecoder(r.Body).Decode(&c); err == nil {
		if id == c.Id {
			if n := repository.Update(c); n > 0 {
				w.WriteHeader(204) /* No Content */
			} else if n := repository.Insert(c); n > 0 {
				w.WriteHeader(201) /* Created */
			} else {
				w.WriteHeader(500) /* Internal Server Error */
			}
		} else {
			w.WriteHeader(400) /* Bad Request */
		}
	} else {
		w.WriteHeader(400) /* Bad Request */
	}
}

func readTrack(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if c, n := repository.Read(id); n > 0 {
		d := repository.Track{Id: c.Id, Audio: c.Audio}
		w.WriteHeader(200) /* OK */
		json.NewEncoder(w).Encode(d)
	} else if n == 0 {
		w.WriteHeader(404) /* Not Found */
	} else {
		w.WriteHeader(500) /* Internal Server Error */
	}
}


func listTrack(w http.ResponseWriter, r *http.Request) {
	if c, n := repository.List(); n >= 0 {
		w.WriteHeader(200) /* OK */
		json.NewEncoder(w).Encode(c)
	} else {
		w.WriteHeader(500) /* Internal Server Error */
	}
}

func deleteTrack(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if n := repository.Delete(id); n > 0 {
		w.WriteHeader(200) /* OK */
	} else if n == 0 {
		w.WriteHeader(404) /* Not Found */
	} else {
		w.WriteHeader(500) /* Internal Server Error */
	}
}

func Router() http.Handler {
	r := mux.NewRouter()
	/* Create */
	r.HandleFunc("/tracks/{id}", updateTrack).Methods("PUT")
	/* Read */
	r.HandleFunc("/tracks/{id}", readTrack).Methods("GET")
	/* List */
	r.HandleFunc("/tracks", listTrack).Methods("GET")
	/* Delete */
	r.HandleFunc("/tracks/{id}", deleteTrack).Methods("DELETE")

	return r
}
