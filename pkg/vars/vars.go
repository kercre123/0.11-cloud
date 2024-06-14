package vars

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

const (
	UserInfoFile    = "secrets/user-info.json"
	IPWhitelistFile = "secrets/ip-whitelist.json"
	SecretsWebroot  = "webroot/"
	// poll cert/key combo every hour
	PollHrs          = 1
	CertFileEnv      = "TLS_CERT_PATH"
	KeyFileEnv       = "TLS_KEY_PATH"
	WeatherAPIKeyEnv = "WEATHER_API_KEY"
	VoskModelPathEnv = "VOSK_MODEL_PATH"
	ChipperPortEnv   = "CHIPPER_PORT"
	WebPortEnv       = "WEB_PORT"
)

var TLSCert []byte
var TLSKey []byte

var UserInfo []Vector_UserInfo
var IPWhitelists []IPWhitelist

var UserDataMutex sync.Mutex
var IPWhitelistMutex sync.Mutex

type IPWhitelist struct {
	ESN string `json:"esn"`
	IP  string `json:"ip"`
}

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

func Init() {
	// deal with secrets
	fmt.Println("loading secrets")
	os.Mkdir("./secrets", 0777)
	uInfo, err := os.ReadFile(UserInfoFile)
	if err != nil {
		fmt.Println("creating secrets")
		os.Create(UserInfoFile)
	}
	err = json.Unmarshal(uInfo, &UserInfo)
	if err != nil {
		fmt.Println("error unmarshaling user info: " + err.Error())
	}

	ipW, err := os.ReadFile(IPWhitelistFile)
	if err != nil {
		fmt.Println("creating IP whitelist file")
		os.Create(UserInfoFile)
	}
	err = json.Unmarshal(ipW, &IPWhitelists)
	if err != nil {
		fmt.Println("error unmarshaling IP whitelist: " + err.Error())
	}

	// open the cert and key
	cFile := os.Getenv(CertFileEnv)
	if cFile == "" {
		fmt.Println("no TLS_CERT_PATH env var given, exiting")
		os.Exit(1)
	}
	TLSCert, err = os.ReadFile(cFile)
	if err != nil {
		fmt.Println("can't read TLS cert file: " + err.Error())
		os.Exit(1)
	}
	kFile := os.Getenv(KeyFileEnv)
	if kFile == "" {
		fmt.Println("no TLS_KEY_PATH env var given, exiting")
		os.Exit(1)
	}
	TLSKey, err = os.ReadFile(kFile)
	if err != nil {
		fmt.Println("can't read TLS key file: " + err.Error())
		os.Exit(1)
	}
}

func GetUserInfo(esn string) (bool, Vector_UserInfo) {
	UserDataMutex.Lock()
	defer UserDataMutex.Unlock()
	for _, info := range UserInfo {
		if info.ESN == esn {
			return false, info
		}
	}
	return true, Vector_UserInfo{}
}

func ChangeUserInfo(uInfo Vector_UserInfo) {
	UserDataMutex.Lock()
	defer UserDataMutex.Unlock()
	for i, info := range UserInfo {
		if info.ESN == uInfo.ESN {
			fmt.Println("changing " + uInfo.ESN + "'s userinfo")
			UserInfo[i] = uInfo
			SaveUserInfo()
			return
		}
	}
	fmt.Println("adding " + uInfo.ESN + " to userinfo")
	UserInfo = append(UserInfo, uInfo)
	SaveUserInfo()
}

func SaveUserInfo() error {
	jsonBytes, err := json.Marshal(UserInfo)
	if err != nil {
		fmt.Println("error saving user info " + err.Error())
		return err
	}
	os.WriteFile(UserInfoFile, jsonBytes, 0777)
	return nil
}

func SaveIPWhitelist() error {
	jsonBytes, err := json.Marshal(IPWhitelists)
	if err != nil {
		fmt.Println("error saving user info " + err.Error())
		return err
	}
	os.WriteFile(IPWhitelistFile, jsonBytes, 0777)
	return nil
}

func AddToIPWhitelist(ipAddr, esn string) {
	UserDataMutex.Lock()
	defer UserDataMutex.Unlock()
	for i, wl := range IPWhitelists {
		if wl.ESN == esn {
			if wl.IP != ipAddr {
				IPWhitelists[i].IP = ipAddr
				SaveIPWhitelist()
			}
			return
		}
	}
	var newWl IPWhitelist
	newWl.ESN = esn
	newWl.IP = ipAddr
	IPWhitelists = append(IPWhitelists, newWl)
	fmt.Println("adding " + ipAddr + " to IP whitelist")
	SaveIPWhitelist()
}

func IsInWhitelist(ipAddr, esn string) (is, ESNmatches bool) {
	UserDataMutex.Lock()
	defer UserDataMutex.Unlock()
	for _, wl := range IPWhitelists {
		if wl.IP == ipAddr {
			if wl.ESN == esn {
				return true, true
			}
			return true, false
		}
	}
	return false, false
}
