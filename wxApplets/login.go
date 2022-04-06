package wxApplets

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
	_ "github.com/tidwall/gjson"
	"net/http"
	_ "net/url"
)

//需要传入wx.login得到的code，以及getPhoneNumber得到的e.detail.code，返回结果：手机号，openid，error
func Login(wxLoginCode string, getPhoneCode string) (string, string, error) {
	type RequestBody struct {
		Code string `json:"code"`
	}
	v := viper.New()
	v.SetConfigName("application")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	err := v.ReadInConfig()
	if err != nil {
		return "", "", err
	}
	v.WatchConfig()
	appid := v.GetString("wechat.appid")
	secret := v.GetString("wechat.secret")
	response, err := http.Get("https://api.weixin.qq.com/sns/jscode2session?appid=" + appid + "&secret=" + secret + "&js_code=" + wxLoginCode + "&grant_type=authorization_code")
	if err != nil {
		return "", "", err
	}
	body := response.Body
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(body)
	if err != nil {
		return "", "", err
	}
	openid := gjson.Get(buf.String(), "openid")
	errCode := gjson.Get(buf.String(), "errcode")
	errMsg := gjson.Get(buf.String(), "errmsg")
	if errCode.Int() != 0 {
		return "", "", errors.New(errMsg.String())
	}
	get, err := http.Get("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=" + appid + "&secret=" + secret)
	if err != nil {
		return "", "", err
	}
	closer := get.Body
	buffer := new(bytes.Buffer)
	_, err = buffer.ReadFrom(closer)
	if err != nil {
		return "", "", err
	}
	accessToken := gjson.Get(buffer.String(), "access_token")
	errCode = gjson.Get(buffer.String(), "errcode")
	errMsg = gjson.Get(buffer.String(), "errmsg")
	if errCode.Int() != 0 {
		return "", "", errors.New(errMsg.String())
	}
	requestBody := RequestBody{
		Code: getPhoneCode,
	}
	js, err := json.MarshalIndent(&requestBody, "", "\t")
	if err != nil {
		return "", "", err
	}
	request, err := http.NewRequest("post", "https://api.weixin.qq.com/wxa/business/getuserphonenumber?access_token="+accessToken.String()+"&", bytes.NewBuffer(js))
	if err != nil {
		return "", "", err
	}
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err = client.Do(request)
	if err != nil {
		return "", "", err
	}
	body = response.Body
	buf = new(bytes.Buffer)
	_, err = buf.ReadFrom(body)
	if err != nil {
		return "", "", err
	}
	phoneNumFormJson := gjson.Get(buf.String(), "phone_info.phoneNumber")
	phoneNum := phoneNumFormJson.String()
	errorCodeFormJson := gjson.Get(buf.String(), "errcode")
	errorMsgFormJson := gjson.Get(buf.String(), "errmsg")
	wxErrCode := errorCodeFormJson.Int()
	if wxErrCode != 0 {
		return "", "", errors.New(errorMsgFormJson.String())
	}
	return phoneNum, openid.String(), nil
}
