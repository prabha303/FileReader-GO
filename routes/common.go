package routes

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	lg "log"
	"net/http"
)

func renderJSON(w http.ResponseWriter, status int, res interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	if status == http.StatusNoContent {
		return
	}

	if err := json.NewEncoder(w).Encode(res); err != nil {
		lg.Printf("ERROR: renderJson - %q\n", err)
	}
}

func parseJSON(_ http.ResponseWriter, body io.ReadCloser, model interface{}) bool {
	defer body.Close()
	b, _ := ioutil.ReadAll(body)
	err := json.Unmarshal(b, model)
	if err != nil {
		fmt.Printf("ERROR: param - %q\n", err)
		return false
	}

	return true
}
