package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/staparx/go_showstart/config"
	"github.com/staparx/go_showstart/log"
	"github.com/staparx/go_showstart/vars"
	"go.uber.org/zap"
)

func main() {
	// 用于结束程序
	defer func() {
		fmt.Println("Press Enter to exit...")
		fmt.Scanln()
	}()
	ctx := context.Background()

	//初始化日志
	log.InitLogger()

	var err error

	vars.ShowPortal()

	//初始化时间地区
	vars.TimeLocal, err = time.LoadLocation(vars.TimeLoadLocation)
	if err != nil {
		log.Logger.Error("⚠️ 初始化时间地区失败，正在使用手动定义的时区信息", zap.Error(err))
		vars.TimeLocal = time.FixedZone("CST", 8*3600)
		log.Logger.Info("✅ 手动定义的时区信息成功！!")
	} else {
		log.Logger.Info(fmt.Sprintf("✅ 时间地区 %s 初始化成功！!", vars.TimeLoadLocation))
	}

	// 打印当前系统时间
	log.Logger.Info(fmt.Sprintf("⏰ 当前系统时间：%s", time.Now().Format("2006-01-02 15:04:05")))

	cfg, err := config.InitCfg()
	if err != nil {
		log.Logger.Error("❌ 配置信息读取失败：", zap.Error(err))
		return
	}
	log.Logger.Info("✅ 系统初始化配置完成！")

	log.Logger.Info("👍开始进入到票务系统抢票流程！！！")
	validate := NewValidateService(ctx, cfg)
	buyTicketList, err := validate.ValidateSystem(ctx)
	if err != nil {
		log.Logger.Error("❌ 配置信息校验失败！！！程序结束", zap.Error(err))
		return
	}

	cancelCtx, cancel := context.WithCancel(ctx)
	for _, ticket := range buyTicketList {
		err = ConfirmOrder(cancelCtx, &OrderDetail{
			ActivityName: ticket.ActivityName,
			SessionName:  ticket.SessionName,
			Price:        ticket.Ticket.SellingPrice,
			ActivityID:   cfg.Ticket.ActivityId,
			GoodType:     ticket.Ticket.GoodType,
			TicketID:     ticket.Ticket.TicketID,
		}, cfg)
		if err != nil {
			log.Logger.Error("❌ 抢票失败！！！程序结束")
			return
		}
	}

	// 捕获终止信号
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case order := <-channel:
		cancel()
		log.Logger.Info("🎉抢票成功！赶紧去订单页面支付吧！！🎉")
		// 下单成功，发送邮件提醒
		if cfg.SmtpEmail.Enable {
			subject := vars.GetEmailTitle()

			body := vars.GetEmailFormat(order.ActivityName, order.SessionName, order.Price)

			if err := sendEmail(subject, body, cfg); err != nil {
				log.Logger.Error("发送邮件失败：", zap.Error(err))
			} else {
				log.Logger.Info("下单成功，邮件已发送")
			}
		}
	case Error := <-ErrorChannel:
		cancel()
		log.Logger.Error("❌ 抢票失败！！！程序结束")
		// 下单失败，发送邮件提醒
		if cfg.SmtpEmail.Enable {
			subject := "抢票初始化失败，请查看错误，并及时处理重启程序！！！"

			body := fmt.Sprintf("错误信息：%s", Error.Error())

			if err := sendEmail(subject, body, cfg); err != nil {
				log.Logger.Error("发送邮件失败：", zap.Error(err))
			} else {
				log.Logger.Info("下单失败，邮件已发送")
			}
		}
	case <-stopChan:
		log.Logger.Info("⚠️ 接收到关闭信号，程序关闭")
		cancel()
		return
	}
}
