package topsdk

type LogisticsDummySendRequest struct {
	Request
	Tid      uint64 `json:"tid,omitempty"`
	Feature  string `json:"feature,omitempty"`
	SellerIp string `json:"seller_ip,omitempty"`
}

type Shipping struct {
	IsSuccess bool `json:"is_success,omitempty"`
}

type LogisticsDummySendResponseRoot struct {
	Shipping `json:"shipping"`
}

type LogisticsDummySendResponse struct {
	LogisticsDummySendResponseRoot `json:"logistics_dummy_send_response"`
}
