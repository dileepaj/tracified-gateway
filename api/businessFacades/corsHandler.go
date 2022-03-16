package businessFacades

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

/*
	EnableCorsAndResponse
*/
func EnableCorsAndResponse(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseForm(); err != nil {
		// Handle error
	}
	data, err := json.Marshal(r.Form)
	if err != nil {
		// Handle error
	}

	var raw map[string][]string
	json.Unmarshal(data, &raw)

	resData, err := http.Get(raw["web"][0])

	if err != nil {
		w.WriteHeader(400)
		return

	} else {
		w.WriteHeader(200)
		body, _ := ioutil.ReadAll(resData.Body)
		w.Write(body)
		return
	}

}
