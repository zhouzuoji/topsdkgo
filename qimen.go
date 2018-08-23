package topsdk

import (
	"encoding/json"
	"fmt"
)

type QimenEventStatus int

const (
	QimenErpTransfer = iota
	QimenErpCheck
	QimenCpNotify
	QimenCpOut
)

func QimenEventName(status QimenEventStatus) (name string) {
	switch status {
	case QimenErpTransfer:
		name = "QIMEN_ERP_TRANSFER"
	case QimenErpCheck:
		name = "QIMEN_ERP_CHECK"
	case QimenCpNotify:
		name = "QIMEN_CP_NOTIFY"
	case QimenCpOut:
		name = "QIMEN_CP_OUT"
	}
	return name
}

type QimenEventProduceRequest struct {
	Status   string `json:"status,omitempty"`
	Tid      uint64 `json:"tid,omitempty"`
	Ext      string `json:"ext,omitempty"`
	Platform string `json:"platform,omitempty"`
	Create   string `json:"create,omitempty"`
	Nick     string `json:"nick,omitempty"`
}

type QimenEventProduceResponseRoot struct {
	IsSuccess bool `json:"is_success,omitempty,string"`
}

type QimenEventProduceResponse struct {
	QimenEventProduceResponseRoot `json:"qimen_event_produce_response"`
}

type QimenResult struct {
	IsSuccess    bool   `json:"is_success,omitempty,string"`
	ErrorCode    string `json:"error_code,omitempty"`
	ErrorMessage string `json:"error_message,omitempty"`
}

type QimenResults struct {
	QimenResult []QimenResult `json:"qimen_result,omitempty"`
}

type QimenEventsProduceResponseRoot struct {
	IsAllSuccess bool         `json:"is_all_success,omitempty,string"`
	Results      QimenResults `json:"results,omitempty"`
}

type QimenEventsProduceResponse struct {
	QimenEventsProduceResponseRoot `json:"qimen_events_produce_response"`
}

type qimenEvent struct {
}

type QimenEvent struct {
	Event QimenEventProduceRequest `json:"event"`
}

type QimenEvents []QimenEvent

func (self *Client) QimenEventsProduce(uri, AccessToken string, events QimenEvents) (err error) {
	if len(events) == 1 {
		_, err = self.DoRequest(uri, AccessToken, events[0].Event, JsonResponse)
		return err
	}
	var bin []byte
	bin, err = json.Marshal(events)
	if err != nil {
		return err
	}
	params := NewParameterMap(1)
	params["messages"] = string(bin)
	//fmt.Println(string(bin))
	var text string
	text, err = self.CallMethod(uri, AccessToken, "taobao.qimen.events.produce", params, JsonResponse)
	if err != nil {
		fmt.Println(text)
	}
	return err
}
