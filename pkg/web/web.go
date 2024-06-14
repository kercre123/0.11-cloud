package web

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/kercre123/0.11-cloud/pkg/vars"
)

func getRemoteIP(r *http.Request) string {
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		ips := strings.Split(xff, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}
	remoteAddr := r.RemoteAddr
	if colonIndex := strings.LastIndex(remoteAddr, ":"); colonIndex != -1 {
		return remoteAddr[:colonIndex]
	}
	return remoteAddr
}

func ipWhitelistMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		remoteIP := getRemoteIP(r)
		is, _ := vars.IsInWhitelist(remoteIP, "")
		if !is {
			http.Error(w, "Error 403. Your IP is not whitelisted. Try a voice command to whitelist your IP, then reload this page.", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	remoteIP := getRemoteIP(r)
	is, _ := vars.IsInWhitelist(remoteIP, "")
	if !is {
		http.Error(w, "Error 403. Your IP is not whitelisted. Try a voice command to whitelist your IP, then reload this page.", http.StatusForbidden)
		return
	}
	switch r.URL.Path {
	case "/api/is_esn_valid":
		var esnReq ESNValidRequest
		req, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading body. Try again.", 500)
			return
		}
		err = json.Unmarshal(req, &esnReq)
		if err != nil {
			http.Error(w, "Error unmarshaling JSON. Try again.", 500)
			return
		}
		esnCheck := strings.ToLower(strings.TrimSpace(esnReq.ESN))
		isValid, isNew := validateESN(esnCheck)
		var esnResp ESNValidResponse
		esnResp.IsValid = isValid
		esnResp.IsNew = isNew
		_, matches := vars.IsInWhitelist(remoteIP, esnCheck)
		esnResp.MatchesIP = matches
		resp, _ := json.Marshal(esnResp)
		w.WriteHeader(200)
		w.Write(resp)
	case "/api/get_settings":
		var esnReq ESNValidRequest
		req, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading body. Try again.", 500)
			return
		}
		err = json.Unmarshal(req, &esnReq)
		if err != nil {
			http.Error(w, "Error unmarshaling JSON. Try again.", 500)
			return
		}
		esnCheck := strings.ToLower(strings.TrimSpace(esnReq.ESN))
		_, matches := vars.IsInWhitelist(remoteIP, esnCheck)
		if !matches {
			http.Error(w, "Error 403. Your bot's ESN doesn't match with this IP. Try a voice command then try again. It is possible your IP changed.", http.StatusForbidden)
			return
		}
		_, info := vars.GetUserInfo(esnCheck)
		marshalled, err := json.Marshal(info)
		if err != nil {
			http.Error(w, "There was an error marshalling the data.", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(200)
		w.Write(marshalled)
	case "/api/set_settings":
		var infoReq vars.Vector_UserInfo
		req, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading body. Try again.", 500)
			return
		}
		err = json.Unmarshal(req, &infoReq)
		if err != nil {
			http.Error(w, "Error unmarshaling JSON. Try again.", 500)
			return
		}
		infoReq.KG.APIKey = strings.TrimSpace(infoReq.KG.APIKey)
		infoReq.KG.ClientKey = strings.TrimSpace(infoReq.KG.ClientKey)
		esnCheck := strings.ToLower(strings.TrimSpace(infoReq.ESN))
		_, matches := vars.IsInWhitelist(remoteIP, esnCheck)
		if !matches {
			http.Error(w, "Error 403. Your bot's ESN doesn't match with this IP. Try a voice command then try again. It is possible your IP changed.", http.StatusForbidden)
			return
		}
		vars.ChangeUserInfo(infoReq)
		w.WriteHeader(200)
		fmt.Fprintf(w, "success")
	}
}

func AddSecretsWebroot() {
	http.Handle("/", ipWhitelistMiddleware(http.FileServer(http.Dir(vars.SecretsWebroot))))
}

func AddSecretsAPI() {
	http.HandleFunc("/api/", apiHandler)
}
