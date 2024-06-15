package web

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/kercre123/0.11-cloud/pkg/vars"
)

func bringToError(w http.ResponseWriter, r *http.Request) {
	// the only user-facing error is unauthorized
	http.Redirect(w, r, "/errors/unauthorized.html", http.StatusTemporaryRedirect)
}

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
			bringToError(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	remoteIP := getRemoteIP(r)
	is, _ := vars.IsInWhitelist(remoteIP, "")
	if !is {
		if r.URL.Path == "/api/check_if_whitelisted" {
			var wListO IsWhitelisted
			wListO.Whitelisted = false
			marshalledwList, _ := json.Marshal(wListO)
			w.Write(marshalledwList)
			return
		}
		http.Error(w, "403 unauthorized (likely not whitelisted)", http.StatusForbidden)
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
	case "/api/check_if_whitelisted":
		var wListO IsWhitelisted
		wListO.Whitelisted = true
		marshalledwList, _ := json.Marshal(wListO)
		w.Write(marshalledwList)
	}
}

func AddSecretsWebroot() {
	http.Handle("/", ipWhitelistMiddleware(http.FileServer(http.Dir(vars.SecretsWebroot))))

	// errors need access to style files as well
	http.Handle("/res/", http.StripPrefix("/res/", http.FileServer(http.Dir(filepath.Join(vars.SecretsWebroot, "resources")))))

	// don't check IP for errors
	http.Handle("/errors/", http.StripPrefix("/errors/", http.FileServer(http.Dir(filepath.Join(vars.SecretsWebroot, "errors")))))
}

func AddSecretsAPI() {
	http.HandleFunc("/api/", apiHandler)
}
