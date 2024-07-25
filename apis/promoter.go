package apis

import (
	"encoding/json"
	"net/http"
)

type Promoter struct {
	Messages []string
}

func (api *Api) Promoter(w http.ResponseWriter, r *http.Request) {
	messages := make([]string, 1)
	messages[0] = "Limited time offer 15% off: You'll get the discount on your order!"

	data := &Promoter{Messages: messages}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}
