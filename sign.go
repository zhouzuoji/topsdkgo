package topsdk

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"encoding/hex"
	"net/url"
	"sort"
	"strings"
)

func encodeParameters(params Parameters) string {
	if params == nil {
		return ""
	}

	var buf bytes.Buffer
	for _, param := range params {
		if buf.Len() > 0 {
			buf.WriteByte('&')
		}
		buf.WriteString(url.QueryEscape(param.Key))
		buf.WriteByte('=')
		buf.WriteString(url.QueryEscape(param.Value))
	}
	return buf.String()
}

func SignMd5(params Parameters, secret string, payload []byte) string {
	h := md5.New()
	h.Write([]byte(secret))
	for _, param := range params {
		h.Write([]byte(param.Key))
		h.Write([]byte(param.Value))
	}
	if len(payload) > 0 {
		h.Write(payload)
	}
	h.Write([]byte(secret))
	return strings.ToUpper(hex.EncodeToString(h.Sum(nil)))
}

func SignHmac(params Parameters, secret string, payload []byte) string {
	h := hmac.New(md5.New, []byte(secret))
	for _, param := range params {
		h.Write([]byte(param.Key))
		h.Write([]byte(param.Value))
	}
	if len(payload) > 0 {
		h.Write(payload)
	}
	return strings.ToUpper(hex.EncodeToString(h.Sum(nil)))
}

func Sign(form IParameterMap, secret string, payload []byte) string {
	params := make(Parameters, form.Len())
	i := 0
	for _, key := range form.Names() {
		value, _ := form.Get(key)
		if len(value) > 0 {
			params[i].Key = key
			params[i].Value = value
			i++
		}
	}
	params = params[0:i]
	sort.Sort(params)

	if signMethName, _ := form.Get("sign_method"); SignMethodFromName(signMethName) == HmacSign {
		return SignHmac(params, secret, payload)
	} else {
		return SignMd5(params, secret, payload)
	}
}

func BuildParams(appkey, secret, AccessToken, method string, payload []byte, params IParameterMap,
	respFmt ResponseFormat) string {
	n := 5
	if params != nil {
		n += params.Len()
	}
	if len(method) > 0 {
		n++
	}
	if len(AccessToken) > 0 {
		n++
	}
	allParams := make(Parameters, n)
	allParams[0].Key = "app_key"
	allParams[0].Value = appkey

	allParams[1].Key = "timestamp"
	allParams[1].Value = GetTimestamp()

	allParams[2].Key = "v"
	allParams[2].Value = "2.0"

	allParams[3].Key = "format"
	if respFmt == JsonResponse {
		allParams[3].Value = "json"
	} else {
		allParams[3].Value = "xml"
	}

	i := 4
	signMethName, _ := params.Get("sign_method")

	var signMeth SignMethod = HmacSign
	if signMethName != "" {
		signMeth = SignMethodFromName(signMethName)
	} else {
		allParams[i].Key = "sign_method"
		allParams[i].Value = "hmac"
		i++
	}

	if len(method) > 0 {
		allParams[i].Key = "method"
		allParams[i].Value = method
		i++
	}
	if len(AccessToken) > 0 {
		allParams[i].Key = "session"
		allParams[i].Value = AccessToken
		i++
	}

	if params != nil {
		for _, key := range params.Names() {
			value, _ := params.Get(key)
			if len(value) > 0 {
				allParams[i].Key = key
				allParams[i].Value = value
				i++
			}
		}
	}

	allParams = allParams[0:i]
	sort.Sort(allParams)
	var sign string
	if signMeth == Md5Sign {
		sign = SignMd5(allParams, secret, payload)
	} else {
		sign = SignHmac(allParams, secret, payload)
	}

	return encodeParameters(allParams) + "&sign=" + sign
}
