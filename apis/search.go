package apis

import (
	"encoding/json"
	"net/http"
)

type SearchPayload struct {
	Text string `json:"text"`
}

type Item struct {
	Id    string
	Sku   string
	Name  string
	Descr string
	Price float32
	Image string
}

func (api *Api) Search(w http.ResponseWriter, r *http.Request) {
	var payload SearchPayload

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		ApiError(w, 1, err, http.StatusBadRequest)
		return
	}

	// testing aja
	if payload.Text == "bag" || payload.Text == "bags" {
		data := make([]*Item, 4)
		data[0] = &Item{Name: "EMEA HOBO BAG LOGO M", Price: 14000000}
		data[1] = &Item{Name: "CELESTE SHOULDER BAG LOGO S", Price: 11000000}
		data[2] = &Item{Name: "EVVIVA SHOULDER BAG S", Price: 12900000}
		data[3] = &Item{Name: "IVY SHOULDER BAG M", Price: 11000000}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(data)
	} else {
		data := make([]*Item, 0)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(data)
	}

}
