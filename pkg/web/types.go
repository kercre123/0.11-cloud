package web

import "github.com/kercre123/0.11-cloud/pkg/vars"

type ESNValidRequest struct {
	ESN string `json:"esn"`
}

type ESNValidResponse struct {
	IsValid bool `json:"esn_isvalid"`
	IsNew   bool `json:"esn_isnew"`
}

/*
type Vector_UserInfo struct {
	ESN      string        `json:"esn"`
	Location string        `json:"location"`
	KG       Vector_KGInfo `json:"kg"`
}

type Vector_KGInfo struct {
	Enabled   bool   `json:"enabled"`
	Prompt    string `json:"prompt"`
	Service   string `json:"service"`
	APIKey    string `json:"apikey"`
	ClientKey string `json:"clientkey"`
}
*/

func validateESN(esn string) (isValid, isNew bool) {
	if len([]rune(esn)) != 8 {
		return false, false
	}
	isNew, _ = vars.GetUserInfo(esn)
	return true, isNew
}
