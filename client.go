package topsdk

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

type Client struct {
	App
}

type NewSession struct {
	R1ExpiresIn  int    `json:"r1_expires_in,string"`
	W1ExpiresIn  int    `json:"w1_expires_in,string"`
	R2ExpiresIn  int    `json:"r2_expires_in,string"`
	W2ExpiresIn  int    `json:"w2_expires_in,string"`
	ExpiresIn    int    `json:"expires_in,string"`
	ReExpiresIn  int    `json:"re_expires_in,string"`
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"top_session"`
	Sign         string `json:"sign"`
	ErrCode      int    `json:"error,string"`
	ErrDesc      string `json:"error_description"`
}

func RefreshSession(appkey, secret, AccessToken, RefreshToken string) (ans NewSession, err error) {
	h := md5.New()
	h.Write([]byte("appkey"))
	h.Write([]byte(appkey))
	h.Write([]byte("refresh_token"))
	h.Write([]byte(RefreshToken))
	h.Write([]byte("sessionkey"))
	h.Write([]byte(AccessToken))
	h.Write([]byte(secret))
	sign := strings.ToUpper(hex.EncodeToString(h.Sum(nil)))
	url := fmt.Sprint("http://container.open.taobao.com/container/refresh?appkey=",
		appkey, "&refresh_token=", RefreshToken, "&sessionkey=", AccessToken, "&sign=", sign)
	var text string
	if text, err = HttpGet(url); err == nil {
		err = json.Unmarshal([]byte(text), &ans)
	}
	return
}

func MissingAppKey(formt ResponseFormat) string {
	var text string
	if formt == JsonResponse {
		text = `{"error_response":{"code":29,"msg":"Invalid app Key",
		"sub_code":"isv.appkey-not-exists","request_id":"3bsxgbx6seom"}}`
	} else {
		text = `<?xml version="1.0" encoding="utf-8" ?><error_response><code>29</code>
          <msg>Invalid app Key</msg><sub_code>isv.appkey-not-exists</sub_code>
          <request_id>14sqam8deckzh</request_id></error_response><!--top010179080184.s.et2-->`
	}
	return text
}

func MissingMethod(format ResponseFormat) string {
	var text string
	if format == JsonResponse {
		text = `{"error_response":{"code":21,"msg":"Missing method","sub_msg":"http传入的参数加入method字段","request_id":"14ibtxvb4vqjh"}}`
	} else {
		text = `<?xml version="1.0" encoding="utf-8" ?><error_response><code>21</code><msg>Missing method</msg><sub_msg>http传入的参数加入method字段</sub_msg><request_id>14ibtxvb4vqjh</request_id></error_response><!--top011009063212.na61-->`
	}
	return text
}

func (self *Client) CallMethod(uri, AccessToken, method string, paramsMap IParameterMap,
	signMeth SignMethod, respFmt ResponseFormat) (body string, err error) {
	if len(self.Key) == 0 {
		return MissingAppKey(respFmt), nil
	}

	if len(method) == 0 {
		return MissingMethod(respFmt), nil
	}

	formData := BuildParams(self.Key, self.Secret, AccessToken, method, nil, paramsMap, signMeth, respFmt)
	return HttpPost(uri, "application/x-www-form-urlencoded", formData)
}

func (self *Client) CallMethodEx(uri, method string, ssn *Session, params IParameterMap,
	signMeth SignMethod, respFmt ResponseFormat) (body string, err error) {
	body, err = self.CallMethod(uri, ssn.AccessToken, method, params, signMeth, respFmt)
	return
}

func (self *Client) DoRequest(uri, AccessToken string, request interface{},
	signMeth SignMethod, respFmt ResponseFormat) (body string, err error) {
	params := Struct2Map(request)
	method := methodName(request)
	body, err = self.CallMethod(uri, AccessToken, (method), params, signMeth, respFmt)
	return
}

func encodeParameterMap(writer *bytes.Buffer, params IParameterMap) {
	if params == nil {
		return
	}

	for _, key := range params.Names() {
		value, _ := params.Get(key)
		if len(value) > 0 {
			if writer.Len() > 0 {
				writer.WriteByte('&')
			}
			writer.WriteString(url.QueryEscape(key))
			writer.WriteByte('=')
			writer.WriteString(url.QueryEscape(value))
		}
	}
}

func (self *Client) DoBatch(uri, AccessToken string, requests []IParameterMap,
	signMeth SignMethod, respFmt ResponseFormat) (body string, err error) {
	var writer = &bytes.Buffer{}
	encodeParameterMap(writer, requests[0])
	for i := 1; i < len(requests); i++ {
		writer.WriteString("\r\n-S-\r\n")
		encodeParameterMap(writer, requests[i])
	}
	payload := writer.String()
	queryString := BuildParams(self.Key, self.Secret, AccessToken, "",
		[]byte(payload), nil, signMeth, respFmt)
	return HttpPost(uri+"?"+queryString, "text/plain;charset=UTF-8", payload)
}

func (self *Client) DoBatchEx(uri, AccessToken string, requests []interface{},
	signMeth SignMethod, respFmt ResponseFormat) (body string, err error) {
	paramMaps := make([]IParameterMap, len(requests))
	for i, v := range requests {
		paramMaps[i] = Struct2Map(v)
		paramMaps[i].Set("method", methodName(v))
	}
	return self.DoBatch(uri, AccessToken, paramMaps, signMeth, respFmt)
}
