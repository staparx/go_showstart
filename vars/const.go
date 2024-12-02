package vars

import "time"

var (
	TimeLocal        *time.Location
	TimeLoadLocation = "Asia/Shanghai"
)

var (
	SaleStatusMap = map[int]string{
		1:  "立即购买",
		2:  "立即购买",
		3:  "票已售罄",
		4:  "即将开售",
		5:  "活动结束",
		6:  "售票结束",
		7:  "仅供展示",
		8:  "APP开票提醒",
		9:  "APP候补购票",
		10: "APP缺票登记",
		11: "APP购票",
		12: "已超过候补限制",
		13: "APP候补购票",
		14: "已参与候补",
	}

	// NeedCpMap 是否需要填写观演人信息
	NeedCpMap = map[int]bool{
		2: true,
		3: false,
		4: false,
	}

	// NeedAdress 是否需要填写地址信息
	NeedAdress = map[int]bool{
		1: false,
		2: true,
	}
)

var (
	EncryptPathMap = map[string]bool{
		"/nj/coupon/order_list":    true,
		"/nj/order/order":          true,
		"/nj/order/coreOrder":      true,
		"/nj/order/getOrderResult": true,
	}
)
