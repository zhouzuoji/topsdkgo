package topsdk

type QimenEventProduceRequest struct {
	Status   string `json:"status,omitempty"`
	Tid      string `json:"tid,omitempty"`
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
