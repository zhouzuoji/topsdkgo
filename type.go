package topsdk

import (
	"time"
)

// gateway, some method must be called through https
const UriDefault string = "http://gw.api.taobao.com/router/rest"
const UriBatch string = "http://gw.api.taobao.com/router/batch"
const UriHttpsDefault string = "https://eco.taobao.com/router/rest"
const UriSandbox string = "http://mockgw.hz.taeapp.com/gw"

// sign method, optional md5 or hmac
type SignMethod int

const (
	UnknownSign = iota
	Md5Sign
	HmacSign
)

func SignMethodFromName(name string) SignMethod {
	if name == "md5" || len(name) == 0 {
		return Md5Sign
	} else if name == "hmac" {
		return HmacSign
	} else {
		return UnknownSign
	}
}

// response content-type, optional xml or json
type ResponseFormat int

const (
	UnknownResponse = iota
	JsonResponse
	XmlResponse
)

func ResponseFormatFromName(name string) ResponseFormat {
	if name == "json" || len(name) == 0 {
		return JsonResponse
	} else if name == "hmac" {
		return XmlResponse
	} else {
		return UnknownResponse
	}
}

// top app parameters
type App struct {
	Key, Secret, RandomNum string
}

type Session struct {
	AccessToken, RefreshToken string
	ExpireAt, ReExpireAt      time.Time
}

type Parameter struct {
	Key, Value string
}

type Parameters []Parameter

func (self Parameters) Len() int {
	return len(self)
}

func (self Parameters) Less(i, j int) bool {
	return self[i].Key < self[j].Key
}

func (self Parameters) Swap(i, j int) {
	self[i], self[j] = self[j], self[i]
}

type Request struct {
	//Timestamp string `json:"timestamp"`
	//Sign string `json:"sign"`
	//AppKey       string `json:"app_key"`
	Method       string
	TargetAppKey string `json:"target_app_key"`
	SignMethod   string `json:"sign_method"`
	AccessToken  string `json:"session"`
	Format       string `json:"format"`
	Version      string `json:"v"`
	PartnerId    string `json:"partner_id"`
	Simplify     bool   `json:"simplify"`
}

type ErrorResponseObject struct {
	Code      int    `json:"code"`
	Msg       string `json:"msg"`
	SubCode   string `json:"sub_code"`
	SubMsg    string `json:"sub_msg"`
	RequestId string `json:"request_id"`
}

type ErrorResponse struct {
	ErrorResponseObject `json:"error_response"`
}

type UserSellerGetRequest struct {
	Fields string `json:"fields"`
}
