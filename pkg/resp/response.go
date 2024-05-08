package response

import (
	"encoding/json"
	"net/http"
)

type (
	SuccessReponse struct {
		Message string `json:"message"`
		Data    any    `json:"data"`
	}

	SuccessPageReponse struct {
		Message string `json:"message"`
		Data    any    `json:"data"`
		Meta    Meta   `json:"meta"`
	}

	Meta struct {
		Limit  int `json:"limit"`
		Offset int `json:"offset"`
		Total  int `json:"total"`
	}
)

type ErrorResponse struct {
	Error struct {
		Code    int         `json:"code"`
		Message interface{} `json:"message"`
	} `json:"error"`
}

// CustomResponse adalah struktur untuk respons JSON yang dapat disesuaikan
type CustomResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// RespondWithError mengirimkan respons JSON dengan pesan kesalahan yang disediakan
func RespondWithError(w http.ResponseWriter, code int, message interface{}) {
	errorMessage := ErrorResponse{
		Error: struct {
			Code    int         `json:"code"`
			Message interface{} `json:"message"`
		}{
			Code:    code,
			Message: message,
		},
	}

	respondJSON(w, code, errorMessage)
}

// RespondWithJSON mengirimkan respons JSON dengan data yang disediakan
func RespondWithJSON(w http.ResponseWriter, code int, message string, data interface{}) {
	response := CustomResponse{
		Message: message,
		Data:    data,
	}

	respondJSON(w, code, response)
}

func respondJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(payload)
}
