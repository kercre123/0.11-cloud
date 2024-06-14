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
		if !vars.IsInWhitelist(remoteIP) {
			http.Error(w, "Error 403. Your IP is not whitelisted. Try a voice command to whitelist your IP, then reload this page.", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	remoteIP := getRemoteIP(r)
	if !vars.IsInWhitelist(remoteIP) {
		http.Error(w, "Error 403. Your IP is not whitelisted. Try a voice command to whitelist your IP, then reload this page.", http.StatusForbidden)
		return
	}
	fmt.Println(r.URL.Path)
	switch r.URL.Path {
	case "/api/is_esn_valid":
		var esnReq ESNValidRequest
		req, _ := io.ReadAll(r.Body)
		json.Unmarshal(req, &esnReq)
		isValid, isNew := validateESN(esnReq.ESN)
		var esnResp ESNValidResponse
		esnResp.IsValid = isValid
		esnResp.IsNew = isNew
		// if isValid {
		// 	w.WriteHeader(500)
		// }
		resp, _ := json.Marshal(esnResp)
		w.Write(resp)
	}
}

func AddSecretsWebroot() {
	http.Handle("/", ipWhitelistMiddleware(http.FileServer(http.Dir(vars.SecretsWebroot))))
}

func AddSecretsAPI() {
	http.HandleFunc("/api/", apiHandler)
}
