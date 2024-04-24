package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

type api struct {
	r mux.Router
}

func New(router *mux.Router) *api {
	return &api{r: *router}
}

func (api *api) FillEndpoints() {
	api.r.HandleFunc("/api/generateToken", GenerateHandler).Methods(http.MethodGet).Queries("GUID", "{GUID}")
	api.r.HandleFunc("/api/refreshToken", RefreshHandler).Methods(http.MethodGet).Queries("Access", "{Access}", "Refresh", "{Refresh}")
}

func (api *api) ListenAndServe(addr string) error {
	return http.ListenAndServe(addr, &api.r)
}
