package enrollment

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/JuD4Mo/go_rest_api/pkg/meta"
)

type (
	Controller func(w http.ResponseWriter, r *http.Request)

	Endpoints struct {
		Create Controller
	}

	CreateReq struct {
		UserId   string `json:"user_id"`
		CourseId string `json:"course_id"`
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
				Err:    fmt.Sprintf("invalid request format: %s", err.Error()),
			})
			return
		}

		if req.UserId == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(&Response{
				Status: http.StatusBadRequest,
				Err:    "user id is required",
			})
			return
		}

		if req.CourseId == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(&Response{
				Status: http.StatusBadRequest,
				Err:    "course id is required",
			})
			return
		}

		enroll, err := s.Create(req.UserId, req.CourseId)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(&Response{
				Status: http.StatusBadRequest,
				Err:    err.Error(),
			})
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(&Response{
			Status: http.StatusCreated,
			Data:   enroll,
		})
	}
}
