package topsdk

type TradesSoldGetRequest struct {
	Request
	Fields       string `json:"fields,omitempty"`
	StartCreated string `json:"start_created,omitempty"`
	EndCreated   string `json:"end_created,omitempty"`
	Status       string `json:"status,omitempty"`
	BuyerNick    string `json:"buyer_nick,omitempty"`
	Type         string `json:"type,omitempty"`
	ExtType      string `json:"ext_type,omitempty"`
	RateStatus   string `json:"rate_status,omitempty"`
	Tag          string `json:"tag,omitempty"`
	PageNo       int    `json:"page_no,omitempty"`
	PageSize     int    `json:"page_size,omitempty"`
	UseHasNext   bool   `json:"use_has_next,omitempty"`
}

type Order struct {
	BuyerRate    bool   `json:"buyer_rate"`
	SellerRate   bool   `json:"seller_rate"`
	Cid          uint64 `json:"cid"`
	Num          uint   `json:"num"`
	NumIid       uint64 `json:"num_iid"`
	Id           string `json:"oid"`
	OuterIid     string `json:"outer_iid"`
	Payment      string `json:"payment"`
	Price        string `json:"price"`
	RefundStatus string `json:"refund_status"`
	Status       string `json:"status"`
	Title        string `json:"title"`
	TotalFee     string `json:"total_fee"`
}

type Orders struct {
	Order []Order `json:"order"`
}

type Trade struct {
	Tid              string `json:"tid"`
	SellerNick       string `json:"seller_nick"`
	BuyerNick        string `json:"buyer_nick"`
	Created          string `json:"created"`
	HasBuyerMessage  bool   `json:"has_buyer_message"`
	Orders           Orders `json:"orders"`
	PayTime          string `json:"pay_time"`
	ReceiverAddress  string `json:"receiver_address"`
	ReceiverCity     string `json:"receiver_city"`
	ReceiverDistrict string `json:"receiver_district"`
	ReceiverMobile   string `json:"receiver_mobile"`
	ReceiverName     string `json:"receiver_name"`
	ReceiverState    string `json:"receiver_state"`
	ReceiverZip      string `json:"receiver_zip"`
	SellerFlag       int    `json:"seller_flag"`
	TotalFee         string `json:"total_fee"`
	Type             string `json:"type"`
}

type Trades struct {
	Trade []Trade `json:"trade"`
}

type TradesSoldGetResponseRoot struct {
	TotalResults int    `json:"total_results"`
	Trades       Trades `json:"trades"`
}

type TradesSoldGetResponse struct {
	TradesSoldGetResponseRoot `json:"trades_sold_get_response"`
}
