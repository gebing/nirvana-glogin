package plat

import (
	"fmt"
	"glogin/config"
	"glogin/pbs/glogin"
	"io/ioutil"
	"net/http"

	"github.com/tidwall/gjson"
)

const (
	googleAuthKey = "google_oauth_url"
)

var Google google

type google struct{}

// Auth 登录返回第三方账号id 和 错误信息
func (g google) Auth(request *glogin.ThirdLoginReq) (string, string, error) {
	baseUrl := authURL(request.Game.BundleId, googleAuthKey)
	url := baseUrl + request.ThirdToken
	resp, err := http.Get(url)
	if err != nil {
		resErr := fmt.Errorf("failed communicating with server: %v", err)
		elkAlarm("error", url, resErr)
		return "", "", resErr
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		elkAlarm(resp.Status, url, "")
		return "", "", fmt.Errorf(resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		resErr := fmt.Errorf("failed reading from metadata server: %v", err)
		elkAlarm("error", url, resErr)
		return "", "", resErr
	}

	uid := gjson.GetBytes(body, "sub").String()
	return uid, uid, nil
}

func (g google) String() string {
	return "google"
}

func (g google) DbFieldName() string {
	return "google"
}

func authURL(bundleId string, platKey string) string {
	//data := gmoss.DynamicCfg("glogin", bundleId, nil)
	//if len(data) == 0 {
	//	return ""
	//}
	//return gjson.GetBytes(data, platKey).String()
	return config.Field(platKey).String()
}
