package topsdk

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/transform"
)

type IParameterMap interface {
	Len() int
	Names() []string
	Get(name string) (string, bool)
	Set(name, value string)
}

type ParameterMap map[string]string

func NewParameterMap(capacity int) ParameterMap {
	return make(ParameterMap, capacity)
}

func (self ParameterMap) Len() int {
	return len(self)
}

func (self ParameterMap) Names() []string {
	ans := make([]string, len(self))
	i := 0
	for key, _ := range self {
		ans[i] = key
		i++
	}
	return ans
}

func (self ParameterMap) Get(name string) (value string, exists bool) {
	value, exists = self[name]
	return
}

func (self ParameterMap) Set(name, value string) {
	self[name] = value
}

type FormValues map[string][]string

func (self FormValues) Len() int {
	return len(self)
}

func (self FormValues) Names() []string {
	ans := make([]string, len(self))
	i := 0
	for key, _ := range self {
		ans[i] = key
		i++
	}
	return ans
}

func (self FormValues) Get(name string) (value string, exists bool) {
	values, ok := self[name]
	if ok && len(values) > 0 {
		value = values[0]
	}
	return
}

func (self FormValues) Set(name, value string) {
	values := self[name]
	if values == nil {
		values = make([]string, 1)
		self[name] = values
		values[0] = value
	} else {
		self[name] = append(values, value)
	}
}

func Value2Str(v *reflect.Value) (ans string) {
	switch v.Type().Kind() {
	case reflect.String:
		ans = v.String()
	case reflect.Invalid:
		ans = ""
	case reflect.Bool:
		ans = strconv.FormatBool(v.Bool())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		ans = strconv.FormatInt(v.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		ans = strconv.FormatUint(v.Uint(), 10)
	case reflect.Float32, reflect.Float64:
		ans = strconv.FormatFloat(v.Float(), 'g', -1, 32)
	default:
		ans = fmt.Sprint(v.Interface())
	}
	return ans
}

func ValueIsEmpty(v *reflect.Value) (ans bool) {
	switch v.Type().Kind() {
	case reflect.String:
		ans = v.String() == ""
	case reflect.Invalid:
		ans = true
	case reflect.Bool:
		ans = !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		ans = v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		ans = v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		ans = v.Float() == 0
	default:
		ans = true
	}
	return ans
}

func Struct2Map(obj interface{}) (ans ParameterMap) {
	ans = make(ParameterMap)
	if obj == nil {
		return
	}
	t := reflect.TypeOf(obj)
	if t.Kind() != reflect.Struct {
		return
	}
	n := t.NumField()
	v := reflect.ValueOf(obj)
	ans = make(ParameterMap, n+1)
	for i := 0; i < n; i++ {
		ft := t.Field(i)
		if ft.Anonymous {
			continue
		}
		fv := v.Field(i)

		fldName := ft.Name
		if tag := ft.Tag.Get("json"); len(tag) > 0 {
			tags := strings.Split(tag, ",")
			if len(tags) > 1 && tags[1] == "omitempty" && ValueIsEmpty(&fv) {
				continue
			}
			fldName = tags[0]
		}
		ans[fldName] = Value2Str(&fv)
	}
	return ans
}

func MethodName(request interface{}) string {
	t := reflect.TypeOf(request)
	if t.Kind() != reflect.Struct {
		return ""
	}
	var sb strings.Builder
	sb.WriteString("taobao")
	for _, c := range t.Name() {
		if c >= 'A' && c <= 'Z' {
			sb.WriteByte('.')
			c += 32
		}
		sb.WriteRune(c)
	}
	sb.String()
	ans := sb.String()
	return ans[0:strings.LastIndexByte(ans, '.')]
}

func GetTimestamp() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func determineEncoding(resp *http.Response) encoding.Encoding {
	//bytes, _ := bufio.NewReader(resp.Body).read Peek(1024)
	//fmt.Println(resp.Header)
	contentType := resp.Header.Get("Content-Type")
	//fmt.Println("Content-Type:", contentType)
	//fmt.Println("Content-Length:", resp.ContentLength)
	e, _, _ := charset.DetermineEncoding(nil, contentType)
	//fmt.Println("codepage:", name)
	return e
}

func getResponseText(resp *http.Response) (string, error) {
	defer resp.Body.Close()
	e := determineEncoding(resp)
	reader := transform.NewReader(resp.Body, e.NewDecoder())
	//fmt.Println("status code:", resp.StatusCode)
	bodyBytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", err
	}
	return string(bodyBytes), err
}

func HttpGet(url string) (body string, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	return getResponseText(resp)
}

func HttpPost(url, contentType, data string) (body string, err error) {
	resp, err := http.Post(url, contentType, strings.NewReader(data))
	if err != nil {
		return "", err
	}
	return getResponseText(resp)
}
