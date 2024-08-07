package helper

import (
	"encoding/json"
	"net/http"
)

func Response(code int, message string, data any) *WebResponse[interface{}] {
	res := new(WebResponse[interface{}])

	// if code == 0 {
	// 	message = "Sukses"
	// } else if code == 500 {
	// 	message = "Terjadi gangguan."
	// 	data = nil
	// }

	res.Code = code
	res.Message = message
	res.Data = data

	return res
}

func ReturnResponse(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}
