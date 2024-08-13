package apis

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type RequestPayload struct {
	PhoneNumber string `json:"phone_number"`
	Message     string `json:"message"`
	RoomId      string `json:"room_id"`
	FromName    string `json:"from_name"`
	Intent      string `json:"intent"`
}

func (api *Api) TestRequest(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		payload := RequestPayload{}
		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		fmt.Println(payload)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{\"success\": true}"))
}
