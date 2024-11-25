package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/staparx/go_showstart/client"
	"github.com/staparx/go_showstart/config"
	"github.com/staparx/go_showstart/log"
	"github.com/staparx/go_showstart/vars"
	"go.uber.org/zap"
	"gopkg.in/gomail.v2"
)

type OrderDetail struct {
	ActivityName string
	SessionName  string
	Price        string
	ActivityID   int
	GoodType     int
	TicketID     string
}

var channel = make(chan *OrderDetail)

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

	log.Logger.Info(fmt.Sprintf("👪观演人数：%d（请注意活动的购票数量限制！）", num))

	t, err := time.ParseInLocation("2006-01-02 15:04:05", cfg.Ticket.StartTime, vars.TimeLocal)
	if err != nil {
		log.Logger.Error("⏰时间格式" + cfg.Ticket.StartTime + "错误，正确格式为：2006-01-02 15:04:05 ")
		return err
	}

	startTime := t.Unix() - 2
	//时间戳转为时间日期字符串
	log.Logger.Info(fmt.Sprintf("🕒 抢票启动时间为：%s", time.Unix(startTime, 0).Format("2006-01-02 15:04:05")))

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				since := startTime - time.Now().Unix()

				if since <= 0 {
					log.Logger.Info("🚀活动即将开始，开始监听抢票！！！")
					for i := 0; i < cfg.System.MaxGoroutine; i++ {
						go GoOrder(ctx, i, c, orderReq, cfg, order)
					}
					return
				} else if since < 10 {
					log.Logger.Info(fmt.Sprintf("🕒 距离抢票开始还有：%d秒", since))
				}
				time.Sleep(time.Second)

			}
		}
	}()

	return nil
}

// 发送邮件
func sendEmail(subject, body string, cfg *config.Config) error {
	m := gomail.NewMessage()
	m.SetHeader("From", cfg.SmtpEmail.Username)
	m.SetHeader("To", cfg.SmtpEmail.To)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	d := gomail.NewDialer(cfg.SmtpEmail.Host, 587, cfg.SmtpEmail.Username, cfg.SmtpEmail.Password)

	// 发送邮件
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}

func GoOrder(ctx context.Context, index int, c client.ShowStartIface, orderReq *client.OrderReq, cfg *config.Config, order *OrderDetail) {
	logPrefix := fmt.Sprintf("[%d]", index)
	firstLoop := true

	for {
		select {
		case <-ctx.Done():
			return
		default:
			if !firstLoop {
				TimeSleep(cfg.System)
			} else {
				firstLoop = false
			}
			//下单
			orderResp, err := c.Order(ctx, orderReq)
			if err != nil {
				log.Logger.Error(logPrefix+"下单失败：", zap.Error(err))
				continue
			}

			// log.Logger.Info(fmt.Sprintf(logPrefix+"下单成功！核心订单Key：%s", orderResp.Result.CoreOrderKey))

			// coreOrder, err := c.CoreOrder(ctx, orderResp.Result.CoreOrderKey)
			// if err != nil {
			// 	log.Logger.Error(logPrefix+"查询核心订单失败：", zap.Error(err))
			// 	continue
			// }

			// var orderJobKey string
			// if coreOrderMap, ok := coreOrder.Result.(map[string]interface{}); ok {
			// 	if _, okk := coreOrderMap["orderJobKey"].(string); okk {
			// 		orderJobKey = coreOrderMap["orderJobKey"].(string)
			// 	}
			// }

			orderJobKey := orderResp.Result.OrderJobKey
			if orderJobKey == "" {
				log.Logger.Error(logPrefix + "orderJobKey为空")
				continue
			}

			log.Logger.Info(fmt.Sprintf(logPrefix+"查询订单成功！orderJobKey：%s", orderJobKey))

			//查询订单结果
			_, err = c.GetOrderResult(ctx, orderJobKey)
			if err != nil {
				log.Logger.Error(logPrefix+"查询订单结果失败：", zap.Error(err))
				continue
			}

			channel <- order
			return
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
