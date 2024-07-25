package main

import (
	"context"
	"github.com/staparx/go_showstart/config"
	"github.com/staparx/go_showstart/log"
	"github.com/staparx/go_showstart/vars"
	"go.uber.org/zap"
	"time"
)

func main() {
	ctx := context.Background()
	//初始化日志
	log.InitLogger()

	var err error
	//初始化时间地区
	vars.TimeLocal, err = time.LoadLocation("Asia/Shanghai")
	if err != nil {
		log.Logger.Error("❌ 初始化时间地区失败：", zap.Error(err))
		return
	}

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
			ActivityID: cfg.Ticket.ActivityId,
			GoodType:   ticket.Ticket.GoodType,
			TicketID:   ticket.Ticket.TicketID,
		}, cfg)
		if err != nil {
			log.Logger.Error("❌ 抢票失败！！！程序结束")
			return
		}
	}

	select {
	case <-channel:
		log.Logger.Info("🎉抢票成功！赶紧去订单页面支付吧！！🎉")
		cancel()
	}

}
