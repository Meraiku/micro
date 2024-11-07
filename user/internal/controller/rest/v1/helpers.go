package v1

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (_ *API) JSON(w http.ResponseWriter, code int, payload any) error {
	w.Header().Set("Content-Type", "application/json")

	data, err := json.Marshal(&payload)
	if err != nil {
		return err
	}

	w.WriteHeader(code)
	_, err = w.Write(data)
	return err
}

func decodeIntoStruct(r *http.Request, v any) error {
	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&v)
	if err != nil {
		return fmt.Errorf("error decoding parameters: %s", err)
	}

	return nil
}
