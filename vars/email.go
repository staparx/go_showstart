package vars

import "fmt"

var emailTitle = "go_showstart 抢票成功通知"

var emailFormat = `
尊敬的用户，您好：

感谢您使用 go_showstart 抢票工具！我们注意到您成功抢到了演出的心仪票务，为了确保您的订单有效，请尽快完成支付。以下是订单的详细信息：

	•	演出名称：<%s>
	•	演出场次：<%s>
	•	票价：<%s>

为了避免票务超时释放，请您尽快前往app端完成支付。


若工具对你有所帮助，别忘记给 go_showstart[https://github.com/staparx/go_showstart] 点个star哦！


祝您观演愉快！

`

func GetEmailTitle() string {
	return emailTitle
}

func GetEmailFormat(activityName, session, price string) string {
	return fmt.Sprintf(emailFormat, activityName, session, price)
}
