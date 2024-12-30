package utils

import (
	"encoding/json"
	"net/http"
)

func DecodeJSON(w http.ResponseWriter, r *http.Request, v interface{}) error {
	err := json.NewDecoder(r.Body).Decode(v)
	if err != nil {
		http.Error(w, "Entrada inválida. Certifique-se de que todos os campos estão corretos.", http.StatusBadRequest)
		return err
	}
	return nil
}
