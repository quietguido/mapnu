package rest

import (
	"encoding/json"
	"net/http"
)

type Err string

func (e Err) Error() string {
	return string(e)
}

const (
	badJson = Err("bad json")
)

func JsonBodyDecoding(r *http.Request, dest any) error {
	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(dest); err != nil {
		return badJson
	}
	return nil
}

func RespondWithError(w http.ResponseWriter, code int, message string) {
	RespondWithJson(w, code, map[string]string{
		"error": message,
	})
}

func RespondWithJson(w http.ResponseWriter, code int, payload any) {
	data, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}
