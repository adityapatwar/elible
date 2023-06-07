package pkg

import (
	"encoding/json"
	"net/http"
)

type ResponseData struct {
	Data struct {
		Status  Status      `json:"status"`
		Results interface{} `json:"results"`
	} `json:"data"`
}

type Status struct {
	Code  int    `json:"code"`
	Pesan string `json:"pesan"`
}

func NewResponseData(code int, pesan string, results interface{}) ResponseData {
	response := ResponseData{}
	response.Data.Status.Code = code
	response.Data.Status.Pesan = pesan
	response.Data.Results = results

	return response
}

// ResponseError adalah tipe data untuk menyimpan informasi tentang kesalahan API
type ResponseError struct {
	Data struct {
		Status struct {
			Code  int    `json:"code"`
			Pesan string `json:"pesan"`
		} `json:"status"`
	} `json:"data"`
}

// NewResponseError membuat instance baru dari ResponseError dengan kode status dan pesan yang diberikan
func NewResponseError(code int, pesan string) ResponseError {
	var respErr ResponseError
	respErr.Data.Status.Code = code
	respErr.Data.Status.Pesan = pesan
	return respErr
}

func WriteErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	respErr := NewResponseError(statusCode, message)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(respErr)
}
