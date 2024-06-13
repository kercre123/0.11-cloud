package web

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/kercre123/0.11-cloud/pkg/vars"
)

func apiHandler(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/api/is_esn_valid":
		var esnReq ESNValidRequest
		req, _ := io.ReadAll(r.Body)
		json.Unmarshal(req, &esnReq)
		isValid, isNew := validateESN(esnReq.ESN)
		var esnResp ESNValidResponse
		esnResp.IsValid = isValid
		esnResp.IsNew = isNew
		if isValid {
			w.WriteHeader(500)
		}
		resp, _ := json.Marshal(esnResp)
		w.Write(resp)
	}
}

func AddSecretsWebroot() {
	http.Handle("/", http.FileServer(http.Dir(vars.SecretsWebroot)))

}

func AddSecretsAPI() {
	http.HandleFunc("/api", apiHandler)
}
