package resources

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"io"
	"bytes"
)

const (
	TOKEN = "69385399feadf6836cf6e7c3e5e5bcfe"
)

type response struct {
    Status string `json:"status"`
    Result struct {
		Id string `json:"title"`
	} `json:"result"`
}

func search(w http.ResponseWriter, r *http.Request) {
	// Get sample
	sample := map[string]interface{}{}
	if err := json.NewDecoder(r.Body).Decode(&sample); err == nil {
		if sample, ok := sample["Audio"].(string); ok {
			reqBody := map[string]interface{} {"api_token" : TOKEN, "audio" : sample}
			if marshBody, err := json.Marshal(reqBody); err == nil {
				sendData := bytes.NewBuffer(marshBody)
				if res, err := http.Post("https://api.audd.io/recognize", "application/json", sendData); err == nil {
					defer res.Body.Close()
					if resBodyMarsh, err := io.ReadAll(res.Body); err == nil {
						resBody := response{}
						if err := json.Unmarshal(resBodyMarsh, &resBody); err == nil {
							resId := map[string]interface{} {"Id" : resBody.Result.Id}
							if resId["Id"].(string) != "" {
								if err := json.NewEncoder(w).Encode(resId); err == nil {
									w.WriteHeader(http.StatusOK)
									return
								} else {
									w.WriteHeader(500) /* Internal Server Error */
								}
							} else {
								w.WriteHeader(404) /* Not Found */
							}			
						} else {
							w.WriteHeader(500) /* Internal Server Error */
						}
					} else {
						w.WriteHeader(500) /* Internal Server Error */
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
	}
	w.WriteHeader(400) /* Bad Request */
}

func Router() http.Handler {
	r := mux.NewRouter()
	/* controller */
	r.HandleFunc("/search", search).Methods("POST")
	return r
}
