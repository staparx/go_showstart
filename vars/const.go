package vars

import "time"

var (
	TimeLocal *time.Location
)

var (
	SaleStatusMap = map[int]string{
		1:  "立即购买",
		4:  "即将开售",
		9:  "APP候补购票",
		10: "APP缺票登记",
	}

	// NeedCpMap 是否需要填写观演人信息
	NeedCpMap = map[int]bool{
		2: true,
		3: false,
		4: false,
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
