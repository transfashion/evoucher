package apis

import (
	"encoding/json"
	"net/http"

	"github.com/fgtago/fgweb/appsmodel"
)

type Api struct {
	Webservice *appsmodel.Webservice
}

type ApiResult struct {
	Code    int
	Message string
	Data    interface{}
}

var api *Api

func New(ws *appsmodel.Webservice) *Api {
	api = &Api{
		Webservice: ws,
	}
	return api
}

func ApiError(w http.ResponseWriter, code int, err error, status int) {
	res := &ApiResult{
		Code:    400,
		Message: err.Error(),
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(res)
}
