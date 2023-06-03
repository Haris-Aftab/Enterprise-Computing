package resources

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
	"bytes"
)


func cooltown(w http.ResponseWriter, r *http.Request) {
	sample := map[string]interface{}{}
	if err := json.NewDecoder(r.Body).Decode(&sample); err == nil {
		if sample, ok := sample["Audio"].(string); ok {
			if id, err := search(sample); err == nil {
				if audio, err := getAudio(id); err == nil {
					r := map[string]interface{} {"Audio" : audio}
					if r["Audio"].(string) != "" {
						if err := json.NewEncoder(w).Encode(r); err == nil {
							w.WriteHeader(http.StatusOK)
							return
						} else {
							w.WriteHeader(500) /* Internal Server Error */
						}
					} else {
						w.WriteHeader(404) /* Internal Server Error */
					}
				} else {
					w.WriteHeader(500) /* Internal Server Error */
				}
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

func search(sample string,) (string, error) {
	if marshBody, err := json.Marshal(map[string]interface{} {"Audio" : sample}); err == nil {
		if res, err := http.Post("http://localhost:3001/search", "application/json", bytes.NewBuffer(marshBody)); err == nil {
			defer res.Body.Close()
			if res.StatusCode == http.StatusOK {
				resBody := map[string]interface{}{}
				if err := json.NewDecoder(res.Body).Decode(&resBody); err == nil {
					if id, ok := resBody["Id"].(string); ok {
						return id, nil
					} else {
						return "", err
					}
				} else {
					return "", err
				}
			} else {
				return "", err
			}
		} else {
			return "", err
		}
	} else {
		return "", err
	}
}

func getAudio(id string) (string, error) {
	urlTrack := "http://localhost:3000/tracks/" + strings.Replace(id, " ", "+", -1)
	if res, err := http.Get(urlTrack); err == nil {
		defer res.Body.Close()
		if res.StatusCode == http.StatusOK {
			resBody := map[string]interface{}{}
			if err := json.NewDecoder(res.Body).Decode(&resBody); err == nil {
				if audio, ok := resBody["Audio"].(string); ok {
					return audio, nil
				} else {
					return "", err
				}
			}else {
				return "", err
			}
		} else {
			return "", err
		}
	} else {
		return "", err
	}
}


func Router() http.Handler {
	r := mux.NewRouter()
	/* controller */
	r.HandleFunc("/cooltown", cooltown).Methods("POST")
	return r
}

