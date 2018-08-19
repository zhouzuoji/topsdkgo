package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/zhouzuoji/topsdkgo"
)

var ClientSecret string
var TopApps map[string]*topsdk.App = map[string]*topsdk.App{}

func copyHeaders(src, dst http.Header) {
	for k, vs := range src {
		for _, value := range vs {
			dst.Add(k, value)
		}
	}
}

func httpPostProxy(w http.ResponseWriter, url, contentType, data string) {
	resp, err := http.Post(url, contentType, strings.NewReader(data))
	if err == nil {
		copyHeaders(resp.Header, w.Header())
		io.Copy(w, resp.Body)
	}
}

func invalidSignature(format string, w http.ResponseWriter) {
	var text string
	if format == "json" {
		text = `{"error_response":{"code":25,"msg":"Invalid signature"}}`
	} else {
		text = `<?xml version="1.0" encoding="utf-8" ?><error_response><code>25</code><msg>Invalid signature</msg></error_response>`
	}
	w.Write([]byte(text))
}

func topProxy(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	form := r.Form
	appKey := form.Get("app_key")
	format := form.Get("format")
	method := form.Get("method")
	fmtid := topsdk.ResponseFormatFromName(format)

	fmt.Println(form)

	app, exists := TopApps[appKey]

	if !exists {
		w.Write([]byte(topsdk.MissingAppKey(fmtid)))
		return
	}

	if len(method) == 0 {
		w.Write([]byte(topsdk.MissingMethod(topsdk.ResponseFormatFromName(format))))
		return
	}

	sign := form.Get("sign")
	delete(form, "sign")
	if len(sign) == 0 || topsdk.Sign(topsdk.FormValues(form), ClientSecret, nil) != sign {
		invalidSignature(format, w)
		return
	}

	uri := form.Get("_uri")
	if len(uri) == 0 {
		uri = topsdk.UriDefault
	}

	delete(form, "_uri")
	form.Set("timestamp", topsdk.GetTimestamp())
	form.Set("sign", topsdk.Sign(topsdk.FormValues(form), app.Secret, nil))
	httpPostProxy(w, uri, "application/x-www-form-urlencoded", form.Encode())
}

func topBatchProxy(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}

	r.ParseForm()
	form := r.Form
	appKey := form.Get("app_key")
	format := form.Get("format")
	fmtid := topsdk.ResponseFormatFromName(format)

	app, exists := TopApps[appKey]

	if !exists {
		w.Write([]byte(topsdk.MissingAppKey(fmtid)))
		return
	}

	sign := form.Get("sign")
	delete(form, "sign")
	if len(sign) == 0 || topsdk.Sign(topsdk.FormValues(form), ClientSecret, body) != sign {
		invalidSignature(format, w)
		return
	}

	uri := form.Get("_uri")
	if len(uri) == 0 {
		uri = topsdk.UriBatch
	}

	delete(form, "_uri")
	form.Set("timestamp", topsdk.GetTimestamp())
	form.Set("sign", topsdk.Sign(topsdk.FormValues(form), app.Secret, body))
	httpPostProxy(w, uri+"?"+form.Encode(), "text/plain;charset=UTF-8", string(body))
}

func main() {
	var port int
	var appkey, secret string
	flag.IntVar(&port, "port", 30001, "listen port")
	flag.StringVar(&appkey, "appkey", "", "TOP appkey")
	flag.StringVar(&secret, "secret", "", "TOP app secret")
	flag.StringVar(&ClientSecret, "ClientSecret", "github.com/zhouzuoji/topsdkgo", "client sign secret")
	flag.Parse()
	if port <= 0 {
		log.Println("invalid port:", port)
		return
	}

	if appkey == "" {
		log.Println("app key is empty")
		return
	}

	if secret == "" {
		log.Println("app secret is empty")
		return
	}

	fmt.Println(appkey, secret, port)

	TopApps[appkey] = &topsdk.App{appkey, secret, ""}

	http.HandleFunc("/router/rest", topProxy)
	http.HandleFunc("/router/batch", topBatchProxy)

	if err := http.ListenAndServe(fmt.Sprint(":", port), nil); err != nil {
		log.Println(err.Error())
	}
}
