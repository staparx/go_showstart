package client

type OrderListReq struct {
	PageNo      int    `json:"pageNo"`
	PageSize    int    `json:"pageSize"`
	TotalAmount string `json:"totalAmount"`
	GoodsID     string `json:"goodsId"`
	GoodsType   int    `json:"goodsType"`
	TicketID    string `json:"ticketId"`
	StFlpv      string `json:"st_flpv"`
	Sign        string `json:"sign"`
	TrackPath   string `json:"trackPath"`
}

type OrderReq struct {
	OrderDetails      []*OrderDetail `json:"orderDetails"`
	CommonPerfomerIds []int          `json:"commonPerfomerIds"`
	AreaCode          string         `json:"areaCode"`
	Telephone         string         `json:"telephone"`
	AddressID         string         `json:"addressId"`
	TeamID            string         `json:"teamId"`
	CouponID          string         `json:"couponId"`
	CheckCode         string         `json:"checkCode"`
	Source            int            `json:"source"`
	Discount          int            `json:"discount"`
	SessionID         int            `json:"sessionId"`
	Freight           int            `json:"freight"`
	AmountPayable     string         `json:"amountPayable"`
	TotalAmount       string         `json:"totalAmount"`
	Partner           string         `json:"partner"`
	OrderSource       int            `json:"orderSource"`
	VideoID           string         `json:"videoId"`
	PayVideotype      string         `json:"payVideotype"`
	StFlpv            string         `json:"st_flpv"`
	Sign              string         `json:"sign"`
	TrackPath         string         `json:"trackPath"`
}

type OrderDetail struct {
	GoodsType  int     `json:"goodsType"`
	SkuType    int     `json:"skuType"`
	Num        string  `json:"num"`
	GoodsID    int     `json:"goodsId"`
	SkuID      string  `json:"skuId"`
	Price      float64 `json:"price"`
	GoodsPhoto string  `json:"goodsPhoto"`
	DyPOIType  int     `json:"dyPOIType"`
	GoodsName  string  `json:"goodsName"`
}
