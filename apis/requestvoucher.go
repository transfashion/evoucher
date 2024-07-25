package apis

import (
	"net/http"
)

func (api *Api) RequestVoucher(w http.ResponseWriter, r *http.Request) {
	messages := make([]string, 1)
	messages[0] = "Limited time offer 15% off: You'll get the discount on your order!"

	// data := ""
	// w.Header().Set("Content-Type", "application/json")
	// w.WriteHeader(http.StatusOK)
	// json.NewEncoder(w).Encode(data)

}
