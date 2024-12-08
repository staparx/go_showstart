package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"sync"
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

var ErrorChannel = make(chan error)

var orderJobKeyAcquired bool = false

// 控制orderJobKeyAcquired的并发访问锁
var orderJobKeyAcquiredLock = sync.Mutex{}

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

	//是否需要填写地址
	if vars.NeedAdress[confirm.Result.OrderInfoVo.TicketPriceVo.TicketType] {
		log.Logger.Info(fmt.Sprintf("🏃地址票务类型为:%d ，匹配地址信息中...", confirm.Result.OrderInfoVo.TicketPriceVo.TicketType))
		//查询地址信息
		adressList, err := c.AdressList(ctx)
		if err != nil {
			log.Logger.Error("❌ 查询地址信息失败：", zap.Error(err))
			return err
		}

		if len(adressList.Result) > 0 {
			for _, v := range adressList.Result {
				if v.IsDefault == 1 {
					orderReq.AddressID = strconv.Itoa(v.ID)
					log.Logger.Info(fmt.Sprintf("🏠地址信息匹配成功！地址：%s", v.Address))
					break
				}
			}
			if orderReq.AddressID == "" {
				log.Logger.Error("❌ 地址信息匹配失败，请设置默认地址")
				return errors.New("地址信息匹配失败，请设置默认地址")
			}
		} else {
			log.Logger.Error("❌ 地址信息匹配失败，请设置默认地址")
			return errors.New("地址信息匹配失败，请设置默认地址")
		}
	} else {
		log.Logger.Info(fmt.Sprintf("🏃地址票务类型为:%d ，无需选择地址 ", confirm.Result.OrderInfoVo.TicketPriceVo.TicketType))
	}

	t, err := time.ParseInLocation("2006-01-02 15:04:05.000", cfg.Ticket.StartTime, vars.TimeLocal)
	if err != nil {
		log.Logger.Error("⏰时间格式" + cfg.Ticket.StartTime + "错误，正确格式为：2006-01-02 15:04:05.000 ")
		return err
	}

	log.Logger.Info(fmt.Sprintf("🕒 抢票启动时间为：%s", t.Format("2006-01-02 15:04:05.000")))

	// time.Millisecond，精确到毫秒
	startTime := t.UnixNano() / int64(time.Millisecond)

	// 开始抢票进程
	StartOrder := func() {
		// since 精确到毫秒
		since := (startTime - time.Now().UnixNano()/int64(time.Millisecond))
		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Duration(since) * time.Millisecond):
			log.Logger.Info("🚀活动即将开始，开始监听抢票！！！")
			for i := 0; i < cfg.System.MaxGoroutine; i++ {
				go GoOrder(ctx, i, c, orderReq, cfg, order)
			}
		}
	}

	// 倒计时进程
	Countdown := func() {
		// since 精确到毫秒
		since := (startTime - time.Now().UnixNano()/int64(time.Millisecond))
		// since 减去 10s
		since -= 10000
		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Duration(since) * time.Millisecond):
			since = (startTime - time.Now().UnixNano()/int64(time.Millisecond))
			// 加入 ctx.Done() 退出
			for since > 0 && ctx.Err() == nil {
				log.Logger.Info(fmt.Sprintf("🕒 距离抢票开始还有：%d秒", since/1000))
				time.Sleep(1 * time.Second)
				since -= 1000
			}
		}
	}

	// token 重新获取进程
	GetTokenAgain := func() {
		// since 精确到毫秒
		since := (startTime - time.Now().UnixNano()/int64(time.Millisecond))
		// since 减去 3min
		since -= 1000 * 60 * 3
		// 如果距离开始时间小于3min，不再重新获取token
		if since < 0 {
			return
		}
		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Duration(since) * time.Millisecond):
			// token 重新获取
			err := c.GetToken(ctx)
			if err != nil {
				log.Logger.Error("token重新获取失败：", zap.Error(err))
				// 再次获取
				err = c.GetToken(ctx)
				if err != nil {
					log.Logger.Error("token重新获取失败：", zap.Error(err))
					ErrorChannel <- err
					return
				}
			}
		}
	}

	// 启动
	go StartOrder()
	go Countdown()
	go GetTokenAgain()

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

	// 除线程0，初始循环仍然加入随机等待
	firstLoop := index == 0

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

			//获取orderJobKey锁
			orderJobKeyAcquiredLock.Lock()
			if orderJobKeyAcquired { //已经有线程获取到orderJobKey
				orderJobKeyAcquiredLock.Unlock()
				continue
			}
			orderJobKeyAcquiredLock.Unlock() //释放orderJobKey锁

			//下单
			orderResp, err := c.Order(ctx, orderReq)
			if err != nil {
				log.Logger.Error(logPrefix+"下单失败：", zap.Error(err))
				continue
			}

			orderJobKey := orderResp.Result.OrderJobKey
			if orderJobKey == "" {
				log.Logger.Error(logPrefix + "orderJobKey为空")
				continue
			}

			log.Logger.Info(fmt.Sprintf(logPrefix+"获取orderJobKey成功！orderJobKey：%s", orderJobKey))

			//获取orderJobKey锁
			orderJobKeyAcquiredLock.Lock()
			orderJobKeyAcquired = true // 有线程获取到orderJobKey
			orderJobKeyAcquiredLock.Unlock()

			OrderResult, orderResultCancel := context.WithCancel(ctx)
			defer orderResultCancel()

			// 每隔200ms发送查询订单结果
			for {
				select {
				case <-OrderResult.Done():
					//停止循环查询订单结果
					return
				default:
					//查询订单结果
					go func() {
						GetOrderResp, err := c.GetOrderResult(ctx, orderJobKey)

						// 如果OrderResult.Done()则不再继续查询订单结果
						if OrderResult.Err() != nil {
							return
						}

						if err != nil {
							log.Logger.Error(logPrefix+"查询订单结果失败：", zap.Error(err))
							// 如果err中包含“小手指点得太快啦，休息一下”，则不停止循环查询订单结果
							if strings.Contains(err.Error(), "小手指点得太快啦，休息一下") {
								return
							}
							//释放orderJobKeyAcquired
							orderJobKeyAcquiredLock.Lock()
							orderJobKeyAcquired = false
							orderJobKeyAcquiredLock.Unlock()
							//停止循环查询订单结果
							orderResultCancel()
							return
						}
						log.Logger.Info(fmt.Sprintf(logPrefix+"查询订单结果成功！订单号：%s", GetOrderResp.Result.OrderSn))
						//停止循环查询订单结果
						orderResultCancel()
						channel <- order
					}()
					// 间隔200ms查询订单结果
					time.Sleep(200 * time.Millisecond)
				}
			}
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
