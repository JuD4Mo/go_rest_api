package course

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/JuD4Mo/go_rest_api/pkg/meta"
	"github.com/gorilla/mux"
)

type (
	Controller func(w http.ResponseWriter, r *http.Request)

	Endpoints struct {
		Create Controller
		Get    Controller
		GetAll Controller
	}

	CreateReq struct {
		Name      string `json:"name"`
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
	}

	Response struct {
		Status int         `json:"status"`
		Data   interface{} `json:"data,omitempty"`
		Err    string      `json:"error,omitempty"`
		Meta   *meta.Meta  `json:"meta,omitempty"`
	}
)

func MakeEndpoints(s Service) Endpoints {
	return Endpoints{
		Create: makeCreateEndpoint(s),
		Get:    makeGetEndpoint(s),
		GetAll: makeGetAllEndpoint(s),
	}
}

func makeCreateEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateReq

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(&Response{
				Status: http.StatusBadRequest,
				Err:    fmt.Sprintf("invalid request format: %v", err.Error()),
			})
			return
		}

		if req.Name == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(&Response{
				Status: http.StatusBadRequest,
				Err:    "name is required",
			})
			return
		}

		if req.StartDate == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(&Response{
				Status: http.StatusBadRequest,
				Err:    "start date is required",
			})
			return
		}

		if req.EndDate == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(&Response{
				Status: http.StatusBadRequest,
				Err:    "end date is required",
			})
			return
		}

		course, err := s.Create(req.Name, req.StartDate, req.EndDate)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(&Response{
				Status: http.StatusInternalServerError,
				Err:    err.Error(),
			})
			return
		}

		json.NewEncoder(w).Encode(&Response{
			Status: http.StatusCreated,
			Data:   course,
		})
	}
}

func makeGetEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		path := mux.Vars(r)
		id, ok := path["id"]
		if !ok || id == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: "id is required"})
			return
		}

		course, err := s.Get(id)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(&Response{
				Status: http.StatusNotFound,
				Err:    "course does not exist",
			})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(&Response{
			Status: http.StatusOK,
			Data:   course,
		})
	}
}

func makeGetAllEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {

		v := r.URL.Query()
		filters := Filters{
			Name: v.Get("name"),
		}

		limit, _ := strconv.Atoi(v.Get("limit"))

		page, _ := strconv.Atoi(v.Get("page"))

		num, err := s.Count(filters)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(&Response{
				Status: http.StatusInternalServerError,
				Err:    err.Error(),
			})
			return
		}

		meta, err := meta.New(page, limit, num)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(&Response{
				Status: http.StatusInternalServerError,
				Err:    err.Error(),
			})
			return
		}

		courses, err := s.GetAll(filters, meta.Offset(), meta.Limit())
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(&Response{
				Status: http.StatusBadRequest,
				Err:    err.Error(),
			})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(&Response{
			Status: http.StatusOK,
			Data:   courses,
			Meta:   meta,
		})
	}
}
