package main

import (
	"context"
	"fmt"
	"math"
	"strings"

	"github.com/staparx/go_showstart/client"
	"github.com/staparx/go_showstart/config"
	"github.com/staparx/go_showstart/log"
	"go.uber.org/zap"
)

type ValidateService struct {
	cfg *config.Config
}

func NewValidateService(ctx context.Context, cfg *config.Config) *ValidateService {
	return &ValidateService{
		cfg: cfg,
	}
}

type buyTicket struct {
	ActivityName                string             `json:"activityName"`
	SessionName                 string             `json:"sessionName"`
	SessionID                   int                `json:"sessionId"`
	IsConfirmedStartTime        int                `json:"isConfirmedStartTime"`
	CommonPerformerDocumentType string             `json:"commonPerformerDocumentType"`
	IsSupportTransform          int                `json:"isSupportTransform"`
	Ticket                      *client.TicketInfo `json:"ticket"`
}

// ValidateSystem 前置检查操作
func (s *ValidateService) ValidateSystem(ctx context.Context) ([]*buyTicket, error) {
	c := client.NewShowStartClient(ctx, s.cfg.Showstart)

	activityId := s.cfg.Ticket.ActivityId

	err := c.GetToken(ctx)
	if err != nil {
		log.Logger.Error("获取登陆token失败", zap.Error(err))
		return nil, err
	}
	log.Logger.Info("👌获取登陆token成功")

	log.Logger.Info("🏃正在查询活动详情信息...")
	//获取活动详情
	detail, err := c.ActivityDetail(ctx, activityId)
	if err != nil {
		log.Logger.Error("❌ 查询活动详情信息失败", zap.Error(err))
		return nil, err
	}
	log.Logger.Info("🎯查询到activity_id对应的活动名称为:")
	log.Logger.Info("==============================================")
	log.Logger.Info(detail.Result.ActivityName)
	log.Logger.Info("==============================================")

	//查询票务信息
	log.Logger.Info("🏃正在查询活动的票务信息...")
	ticketList, err := c.ActivityTicketList(ctx, activityId)
	if err != nil {
		log.Logger.Error("❌ 查询活动票务信息失败", zap.Error(err))
		return nil, err
	}

	//按顺序查找票务信息
	var buyTicketList []*buyTicket
	for _, ticket := range s.cfg.Ticket.List {
		for _, result := range ticketList.Result {
			//找到对应的场次
			if DelectStringBlank(result.SessionName) == DelectStringBlank(ticket.Session) {
				//找到对应的票价
				for _, ticketPrice := range result.TicketPriceList {
					if ticket.Price == ticketPrice.Price {
						//将场次票价信息保存下来
						buyTicketList = append(buyTicketList, &buyTicket{
							ActivityName:                detail.Result.ActivityName,
							SessionName:                 result.SessionName,
							SessionID:                   result.SessionID,
							IsConfirmedStartTime:        result.IsConfirmedStartTime,
							CommonPerformerDocumentType: result.CommonPerformerDocumentType,
							IsSupportTransform:          result.IsSupportTransform,
							Ticket:                      ticketPrice.TicketList[0],
						})
					}
				}
			}
		}
	}
	if len(buyTicketList) == 0 {
		log.Logger.Error("❌ 配置匹配票档失败！在场次中未找寻到对应票价的信息")
		log.Logger.Info("🎯进入手动匹配模式，请根据以下信息进行匹配:")
		// return nil, errors.New("匹配票档失败！在场次中未找寻到对应票价的信息")

		if len(ticketList.Result) == 1 { // 单场次
			log.Logger.Info("🎯仅有一个场次，默认匹配，场次名为:" + ticketList.Result[0].SessionName)
			if len(ticketList.Result[0].TicketPriceList) == 1 { // 单场次单票价
				log.Logger.Info("🎯仅有一个票价，默认匹配，票价为:" + ticketList.Result[0].TicketPriceList[0].Price)
				err := config.SaveCfg(ticketList.Result[0].SessionName, ticketList.Result[0].TicketPriceList[0].Price) // 保存配置到config.yaml
				if err != nil {
					log.Logger.Error("❌ 保存手动匹配配置信息失败", zap.Error(err))
				} else {
					log.Logger.Info("🎯保存手动匹配配置信息成功")
				}
				buyTicketList = append(buyTicketList, &buyTicket{
					ActivityName:                detail.Result.ActivityName,
					SessionName:                 ticketList.Result[0].SessionName,
					SessionID:                   ticketList.Result[0].SessionID,
					IsConfirmedStartTime:        ticketList.Result[0].IsConfirmedStartTime,
					CommonPerformerDocumentType: ticketList.Result[0].CommonPerformerDocumentType,
					IsSupportTransform:          ticketList.Result[0].IsSupportTransform,
					Ticket:                      ticketList.Result[0].TicketPriceList[0].TicketList[0],
				})
			} else { // 单场次多票价
				log.Logger.Info("🎯有多个票价，请手动匹配")
				for index, ticketPrice := range ticketList.Result[0].TicketPriceList {
					log.Logger.Info(fmt.Sprintf("🎯票价%d：%s", index+1, ticketPrice.Price))
				}
				log.Logger.Info("🎯请输入票价序号:")
				var ticketIndex int
				fmt.Scanln(&ticketIndex)
				err := config.SaveCfg(ticketList.Result[0].SessionName, ticketList.Result[0].TicketPriceList[ticketIndex-1].Price) // 保存配置到config.yaml
				if err != nil {
					log.Logger.Error("❌ 保存手动匹配配置信息失败", zap.Error(err))
				} else {
					log.Logger.Info("🎯保存手动匹配配置信息成功")
				}
				buyTicketList = append(buyTicketList, &buyTicket{
					ActivityName:                detail.Result.ActivityName,
					SessionName:                 ticketList.Result[0].SessionName,
					SessionID:                   ticketList.Result[0].SessionID,
					IsConfirmedStartTime:        ticketList.Result[0].IsConfirmedStartTime,
					CommonPerformerDocumentType: ticketList.Result[0].CommonPerformerDocumentType,
					IsSupportTransform:          ticketList.Result[0].IsSupportTransform,
					Ticket:                      ticketList.Result[0].TicketPriceList[ticketIndex-1].TicketList[0],
				})
			}
		} else { // 多场次
			log.Logger.Info("🎯有多个场次，请手动匹配")
			for index, session := range ticketList.Result {
				log.Logger.Info(fmt.Sprintf("🎯场次%d：%s", index+1, session.SessionName))
			}
			log.Logger.Info("🎯请输入场次序号:")
			var sessionIndex int
			fmt.Scanln(&sessionIndex)
			if len(ticketList.Result[sessionIndex-1].TicketPriceList) == 1 { // 多场次单票价
				log.Logger.Info("🎯仅有一个票价，默认匹配，票价为:" + ticketList.Result[sessionIndex-1].TicketPriceList[0].Price)
				err := config.SaveCfg(ticketList.Result[sessionIndex-1].SessionName, ticketList.Result[sessionIndex-1].TicketPriceList[0].Price) // 保存配置到config.yaml
				if err != nil {
					log.Logger.Error("❌ 保存手动匹配配置信息失败", zap.Error(err))
				} else {
					log.Logger.Info("🎯保存手动匹配配置信息成功")
				}
				buyTicketList = append(buyTicketList, &buyTicket{
					ActivityName:                detail.Result.ActivityName,
					SessionName:                 ticketList.Result[sessionIndex-1].SessionName,
					SessionID:                   ticketList.Result[sessionIndex-1].SessionID,
					IsConfirmedStartTime:        ticketList.Result[sessionIndex-1].IsConfirmedStartTime,
					CommonPerformerDocumentType: ticketList.Result[sessionIndex-1].CommonPerformerDocumentType,
					IsSupportTransform:          ticketList.Result[sessionIndex-1].IsSupportTransform,
					Ticket:                      ticketList.Result[sessionIndex-1].TicketPriceList[0].TicketList[0],
				})
			} else { // 多场次多票价
				log.Logger.Info("🎯有多个票价，请手动匹配")
				for index, ticketPrice := range ticketList.Result[sessionIndex-1].TicketPriceList {
					log.Logger.Info(fmt.Sprintf("🎯票价%d：%s", index+1, ticketPrice.Price))
				}
				log.Logger.Info("🎯请输入票价序号:")
				var ticketIndex int
				fmt.Scanln(&ticketIndex)
				err := config.SaveCfg(ticketList.Result[sessionIndex-1].SessionName, ticketList.Result[sessionIndex-1].TicketPriceList[ticketIndex-1].Price) // 保存配置到config.yaml
				if err != nil {
					log.Logger.Error("❌ 保存手动匹配配置信息失败", zap.Error(err))
				} else {
					log.Logger.Info("🎯保存手动匹配配置信息成功")
				}
				buyTicketList = append(buyTicketList, &buyTicket{
					ActivityName:                detail.Result.ActivityName,
					SessionName:                 ticketList.Result[sessionIndex-1].SessionName,
					SessionID:                   ticketList.Result[sessionIndex-1].SessionID,
					IsConfirmedStartTime:        ticketList.Result[sessionIndex-1].IsConfirmedStartTime,
					CommonPerformerDocumentType: ticketList.Result[sessionIndex-1].CommonPerformerDocumentType,
					IsSupportTransform:          ticketList.Result[sessionIndex-1].IsSupportTransform,
					Ticket:                      ticketList.Result[sessionIndex-1].TicketPriceList[ticketIndex-1].TicketList[0],
				})
			}
		}
	}

	log.Logger.Info("🎫获取票务信息成功，系统将按照以下优先级进行抢购:")
	log.Logger.Info("==============================================")
	var startTime int64 = math.MaxInt64
	for _, v := range buyTicketList {
		log.Logger.Info(fmt.Sprintf("%s - %s - %s", v.SessionName, v.Ticket.TicketType, v.Ticket.CostPrice))
		startTime = int64(math.Min(float64(startTime), float64(v.Ticket.StartTime)))
	}
	log.Logger.Info("==============================================")

	return buyTicketList, nil
}

// DelectStringBlank 函数移除字符串中的所有空格
func DelectStringBlank(s string) string {
	return strings.ReplaceAll(s, " ", "")
}
