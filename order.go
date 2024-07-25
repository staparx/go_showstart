package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/staparx/go_showstart/client"
	"github.com/staparx/go_showstart/config"
	"github.com/staparx/go_showstart/log"
	"github.com/staparx/go_showstart/util"
	"github.com/staparx/go_showstart/vars"
	"go.uber.org/zap"
	"math/rand"
	"strconv"
	"time"
)

type OrderDetail struct {
	ActivityID int
	GoodType   int
	TicketID   string
}

var channel = make(chan struct{})

func ConfirmOrder(ctx context.Context, order *OrderDetail, cfg *config.Config) error {
	c := client.NewShowStartClient(ctx, cfg.Showstart)

	num := len(cfg.Ticket.People)
	//订单信息确认
	confirm, err := c.Confirm(ctx, order.ActivityID, order.TicketID, fmt.Sprintf("%d", num))
	if err != nil {
		log.Logger.Error("❌ 订单信息确认失败：", zap.Error(err))
		return err
	}

	log.Logger.Info("👀订单信息确认成功！", zap.Any("ticket_id", order.TicketID))

	pay := strconv.FormatFloat(confirm.Result.OrderInfoVo.TicketPriceVo.Price*float64(num), 'f', 2, 64)
	//下单
	orderReq := &client.OrderReq{
		OrderDetails: []*client.OrderDetail{
			{
				GoodsType:  order.GoodType,
				SkuType:    confirm.Result.OrderInfoVo.TicketPriceVo.TicketType,
				Num:        fmt.Sprintf("%d", num),
				GoodsID:    confirm.Result.OrderInfoVo.ActivityID,
				SkuID:      confirm.Result.OrderInfoVo.TicketPriceVo.TicketID,
				Price:      confirm.Result.OrderInfoVo.TicketPriceVo.Price,
				GoodsPhoto: confirm.Result.OrderInfoVo.Poster,
				DyPOIType:  confirm.Result.OrderInfoVo.TicketPriceVo.DyPOIType,
				GoodsName:  confirm.Result.OrderInfoVo.Title,
			},
		},
		CommonPerfomerIds: []int{},
		AreaCode:          confirm.Result.OrderInfoVo.AreaCode,
		Telephone:         confirm.Result.OrderInfoVo.Telephone,
		AddressID:         "",
		TeamID:            "",
		CouponID:          "",
		CheckCode:         "",
		Source:            0,
		Discount:          0,
		SessionID:         confirm.Result.OrderInfoVo.SessionID,
		Freight:           0,
		AmountPayable:     pay,
		TotalAmount:       pay,
		Partner:           "",
		OrderSource:       1,
		VideoID:           "",
		PayVideotype:      "",
		StFlpv:            "",
		Sign:              "",
		TrackPath:         "",
	}
	//是否需要查询观演人
	if vars.NeedCpMap[confirm.Result.OrderInfoVo.BuyType] {
		log.Logger.Info(fmt.Sprintf("🏃票务类型为:%d ，匹配观演人信息中...", confirm.Result.OrderInfoVo.BuyType))
		//查询观演人信息
		cpResp, err := c.CpList(ctx, order.TicketID)
		if err != nil {
			log.Logger.Error("❌ 查询观演人信息失败：", zap.Error(err))
			return err
		}

		var perfomerIds []int
		for _, v := range cpResp.Result {
			for _, user := range cfg.Ticket.People {
				if v.Name == user {
					perfomerIds = append(perfomerIds, v.ID)
				}
			}
		}

		if len(perfomerIds) > 0 && len(perfomerIds) == len(cfg.Ticket.People) {
			log.Logger.Info("🙎观演人信息匹配成功！!")
			orderReq.CommonPerfomerIds = perfomerIds
		} else {
			log.Logger.Error("❌ 观演人信息匹配失败")
			return errors.New("观演人信息匹配失败")
		}
	} else {
		log.Logger.Info(fmt.Sprintf("🏃票务类型为:%d ，无需选择观演人 ", confirm.Result.OrderInfoVo.BuyType))

	}

	t, err := time.ParseInLocation("2006-01-02 15:04:05", cfg.Ticket.StartTime, vars.TimeLocal)
	if err != nil {
		log.Logger.Error("⏰时间格式" + cfg.Ticket.StartTime + "错误，正确格式为：2006-01-02 15:04:05 ")
		return err
	}

	startTime := t.Unix()
	now := time.Now().Unix()

	// 计算等待时间
	waitTime := startTime - now - 2

	// 等待开票
	if waitTime > 0 {
		day, hour, minute, second := util.ConvertSeconds(waitTime)
		log.Logger.Info(fmt.Sprintf("⏰活动还未开始，预计等待时间为：%d天%d时%d分%d秒 \n", day, hour, minute, second))
		// 转换为 Duration 类型
		waitDuration := time.Duration(waitTime) * time.Second

		// 设置定时器
		timer := time.NewTimer(waitDuration)

		// 等待定时器到期
		<-timer.C
	}

	log.Logger.Info("👂活动即将开始，开始监听抢票！！！")
	for i := 0; i < cfg.System.MaxGoroutine; i++ {
		go GoOrder(ctx, i, c, orderReq, cfg)
	}

	return nil
}

func GoOrder(ctx context.Context, index int, c client.ShowStartIface, orderReq *client.OrderReq, cfg *config.Config) {
	logPrefix := fmt.Sprintf("[%d]", index)

	for {
		select {
		case <-ctx.Done():
			return
		default:
			TimeSleep(cfg.System)
			//下单
			orderResp, err := c.Order(ctx, orderReq)
			if err != nil {
				log.Logger.Error(logPrefix+"下单失败：", zap.Error(err))
				continue
			}

			log.Logger.Info(fmt.Sprintf(logPrefix+"下单成功！核心订单Key：%s", orderResp.Result.CoreOrderKey))

			coreOrder, err := c.CoreOrder(ctx, orderResp.Result.CoreOrderKey)
			if err != nil {
				log.Logger.Error(logPrefix+"查询核心订单失败：", zap.Error(err))
				continue
			}

			var orderJobKey string
			if coreOrderMap, ok := coreOrder.Result.(map[string]interface{}); ok {
				if _, okk := coreOrderMap["orderJobKey"].(string); okk {
					orderJobKey = coreOrderMap["orderJobKey"].(string)
				}
			}

			if orderJobKey == "" {
				log.Logger.Error(logPrefix + "核心订单Key为空")
				continue
			}

			log.Logger.Info(fmt.Sprintf(logPrefix+"查询核心订单成功！订单任务Key：%s", orderJobKey))

			//查询订单结果
			_, err = c.GetOrderResult(ctx, orderJobKey)
			if err != nil {
				log.Logger.Error(logPrefix+"查询订单结果失败：", zap.Error(err))
				continue
			}

			channel <- struct{}{}
		}

	}
}

func TimeSleep(cfg *config.System) {
	// 生成随机休眠时间
	minInterval := cfg.MinInterval
	maxInterval := cfg.MaxInterval
	interval := rand.Intn(maxInterval-minInterval+1) + minInterval
	time.Sleep(time.Duration(interval) * time.Millisecond)
}
