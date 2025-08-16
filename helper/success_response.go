package helper

import (
	"encoding/json"
	"net/http"

	"github.com/syrlramadhan/api-bendahara-inovdes/dto"
)

// WriteJSONSuccess digunakan untuk mengirim response sukses dalam format JSON
func WriteJSONSuccess(w http.ResponseWriter, data interface{}, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	response := dto.ListResponseOK{
		Code:    code,
		Status:  http.StatusText(http.StatusOK),
		Data:    data,
		Message: message,
	}

	json.NewEncoder(w).Encode(response)
}
