package client

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"

	jsoniter "github.com/json-iterator/go"
	"github.com/staparx/go_showstart/log"
)

type ShowStartIface interface {
	// GetToken 获取token
	GetToken(ctx context.Context) error
	// ActivityDetail 活动详情
	ActivityDetail(ctx context.Context, activityId int) (*ActivityDetailResp, error)
	// ActivityTicketList 获取票务场次信息
	ActivityTicketList(ctx context.Context, activityId int) (*ActivityTicketListResp, error)
	// Confirm 确认购买
	Confirm(ctx context.Context, activityId int, ticketId, ticketNum string) (*ConfirmResp, error)
	// CpList 观演人列表
	CpList(ctx context.Context, ticketId string) (*CpListResp, error)
	OrderList(ctx context.Context, req *OrderListReq) (*OrderListResp, error)
	// Order 下单
	Order(ctx context.Context, req *OrderReq) (*OrderResp, error)
	// CoreOrder 核心订单确认
	CoreOrder(ctx context.Context, coreOrderKey string) (*OrderCoreResp, error)
	// GetOrderResult 获取订单结果
	GetOrderResult(ctx context.Context, orderJobKey string) (*GetOrderResultResp, error)
}

// GetToken 获取token
func (c *ShowStartClient) GetToken(ctx context.Context) error {
	path := "/waf/gettoken"
	data := fmt.Sprintf(`{"st_flpv":"%s","sign":"%s","trackPath":""}`, c.StFlpv, c.Sign)

	result, err := c.Post(ctx, path, data)
	if err != nil {
		return err
	}

	var resp *GetTokenResp
	err = jsoniter.Unmarshal(result, &resp)
	if err != nil {
		return err
	}

	if !resp.Success {
		return errors.New(resp.Msg)
	}

	c.Cusat = resp.Result.AccessToken.AccessToken
	c.Cusit = resp.Result.IDToken.IDToken

	return nil
}

// ActivityDetail 获取活动详情
func (c *ShowStartClient) ActivityDetail(ctx context.Context, activityId int) (*ActivityDetailResp, error) {
	path := "/wap/activity/details"
	data := fmt.Sprintf(`{"activityId":"%d","coupon":"","shareId":"","st_flpv":"%s","sign":"%s","trackPath":""}`,
		activityId, c.StFlpv, c.Sign)

	result, err := c.Post(ctx, path, data)
	if err != nil {
		return nil, err
	}

	var resp *ActivityDetailResp
	err = jsoniter.Unmarshal(result, &resp)
	if err != nil {
		return nil, err
	}
	if resp.State == "token-expire-at" || resp.Msg == "登录过期了，请重新登录！" {
		err = c.GetToken(ctx)
		if err != nil {
			return nil, err
		}
		return c.ActivityDetail(ctx, activityId)
	}

	if resp.State != "1" {
		return nil, errors.New(resp.Msg)
	}

	return resp, nil
}

// ActivityTicketList 获取票务场次信息
func (c *ShowStartClient) ActivityTicketList(ctx context.Context, activityId int) (*ActivityTicketListResp, error) {
	path := "/wap/activity/V2/ticket/list"
	data := fmt.Sprintf(`{"activityId":"%d","coupon":"","st_flpv":"%s","sign":"%s","trackPath":""}`,
		activityId, c.StFlpv, c.Sign)

	result, err := c.Post(ctx, path, data)
	if err != nil {
		return nil, err
	}

	var resp *ActivityTicketListResp
	err = jsoniter.Unmarshal(result, &resp)
	if err != nil {
		return nil, err
	}
	if resp.State == "token-expire-at" || resp.Msg == "登录过期了，请重新登录！" {
		err = c.GetToken(ctx)
		if err != nil {
			return nil, err
		}
		return c.ActivityTicketList(ctx, activityId)
	}

	if resp.State != "1" {
		return nil, errors.New(resp.Msg)
	}

	return resp, nil
}

// Confirm 确认购买
func (c *ShowStartClient) Confirm(ctx context.Context, activityId int, ticketId, ticketNum string) (*ConfirmResp, error) {
	path := "/order/wap/order/confirm"
	body := fmt.Sprintf(`{"ticketId":"%s","sequence":"%d","ticketNum":"%s","st_flpv":"%s","sign":"%s","trackPath":""}`,
		ticketId, activityId, ticketNum, c.StFlpv, c.Sign)

	result, err := c.Post(ctx, path, body)
	if err != nil {
		return nil, err
	}

	var resp *ConfirmResp
	err = jsoniter.Unmarshal(result, &resp)
	if err != nil {
		return nil, err
	}
	if resp.State == "token-expire-at" || resp.Msg == "登录过期了，请重新登录！" {
		err = c.GetToken(ctx)
		if err != nil {
			return nil, err
		}
		return c.Confirm(ctx, activityId, ticketId, ticketNum)
	}

	return resp, nil
}

// CpList 观演人列表
func (c *ShowStartClient) CpList(ctx context.Context, ticketId string) (*CpListResp, error) {
	path := "/wap/cp/list"
	data := fmt.Sprintf(`{"ticketPriceId":"%s","audienceWhitelistPolicy":0,"st_flpv":"%s","sign":"%s","trackPath":""}`,
		ticketId, c.StFlpv, c.Sign)

	result, err := c.Post(ctx, path, data)
	if err != nil {
		return nil, err
	}

	var resp *CpListResp
	err = jsoniter.Unmarshal(result, &resp)
	if err != nil {
		return nil, err
	}
	if resp.State == "token-expire-at" || resp.Msg == "登录过期了，请重新登录！" {
		err = c.GetToken(ctx)
		if err != nil {
			return nil, err
		}
		return c.CpList(ctx, ticketId)
	}

	return resp, nil
}

func (c *ShowStartClient) OrderList(ctx context.Context, req *OrderListReq) (*OrderListResp, error) {
	path := "/nj/coupon/order_list"
	data, err := jsoniter.MarshalToString(req)
	if err != nil {
		return nil, err
	}
	//e = "{\"orderDetails\":[{\"goodsType\":1,\"skuType\":1,\"num\":\"1\",\"goodsId\":233792,\"skuId\":\"5f62b72525b041791b237911d064609a\",\"price\":188,\"goodsPhoto\":\"https://s2.showstart.com/img/2024/0703/15/30/00f9cd985a664df2b38797b2dcd78558_1200_1600_3167694.0x0.png\",\"dyPOIType\":2,\"goodsName\":\"SFNT「City Walk 都市漫游」巡回演唱会 -深圳站\"}],\"commonPerfomerIds\":[5773864],\"areaCode\":\"86_CN\",\"telephone\":\"15813830747\",\"addressId\":\"\",\"teamId\":\"\",\"couponId\":\"\",\"checkCode\":\"\",\"source\":0,\"discount\":0,\"sessionId\":3249212,\"freight\":0,\"amountPayable\":\"188.00\",\"totalAmount\":\"188.00\",\"partner\":\"\",\"orderSource\":1,\"videoId\":\"\",\"payVideotype\":\"\",\"st_flpv\":\"kgUBxCv3t7EN1wWr6Yeo\",\"sign\":\"f65f2dbcd7387b2376bcefe13754856b\",\"trackPath\":\"\"}", n = "341sW411imd3cgg7"

	result, err := c.Post(ctx, path, data)
	if err != nil {
		return nil, err
	}

	var resp *OrderListResp
	err = jsoniter.Unmarshal(result, &resp)
	if err != nil {
		return nil, err
	}

	if resp.State == "token-expire-at" || resp.Msg == "登录过期了，请重新登录！" {
		err = c.GetToken(ctx)
		if err != nil {
			return nil, err
		}
		return c.OrderList(ctx, req)
	}

	return resp, nil
}

func (c *ShowStartClient) Order(ctx context.Context, req *OrderReq) (*OrderResp, error) {
	req.StFlpv = c.StFlpv
	req.Sign = c.Sign

	path := "/nj/order/order"
	data, err := jsoniter.MarshalToString(req)
	if err != nil {
		return nil, err
	}
	result, err := c.Post(ctx, path, data)
	if err != nil {
		return nil, err
	}

	// 将 Order 的返回值打印到Debug日志中
	log.Logger.Debug("Order:", zap.String("result", string(result)))

	var resp *OrderResp
	err = jsoniter.Unmarshal(result, &resp)
	if err != nil {
		return nil, err
	}

	if resp.State == "token-expire-at" || resp.Msg == "登录过期了，请重新登录！" {
		err = c.GetToken(ctx)
		if err != nil {
			return nil, err
		}
		return c.Order(ctx, req)
	}

	if resp.Success || resp.State == "1" {
		return resp, nil
	}

	return nil, errors.New(resp.Msg)
}

func (c *ShowStartClient) CoreOrder(ctx context.Context, coreOrderKey string) (*OrderCoreResp, error) {
	path := "/nj/order/coreOrder"

	body := fmt.Sprintf(`{"coreOrderKey":"%s","st_flpv":"%s","sign":"%s","trackPath":""}`, coreOrderKey, c.StFlpv, c.Sign)

	result, err := c.Post(ctx, path, body)
	if err != nil {
		return nil, err
	}

	var resp *OrderCoreResp
	err = jsoniter.Unmarshal(result, &resp)
	if err != nil {
		return nil, err
	}

	if resp.State == "token-expire-at" || resp.Msg == "登录过期了，请重新登录！" {
		err = c.GetToken(ctx)
		if err != nil {
			return nil, err
		}
		return c.CoreOrder(ctx, coreOrderKey)
	}

	if resultStr, ok := resp.Result.(string); ok {
		if resultStr == "pending" {
			time.Sleep(200 * time.Millisecond)
			return c.CoreOrder(ctx, coreOrderKey)
		}
	}

	return resp, nil
}

func (c *ShowStartClient) GetOrderResult(ctx context.Context, orderJobKey string) (*GetOrderResultResp, error) {
	path := "/nj/order/getOrderResult"

	body := fmt.Sprintf(`{"orderJobKey":"%s","st_flpv":"%s","sign":"%s","trackPath":""}`, orderJobKey, c.StFlpv, c.Sign)

	result, err := c.Post(ctx, path, body)
	if err != nil {
		return nil, err
	}

	// 将 GetOrderResult 的返回值打印到Debug日志中
	log.Logger.Debug("GetOrderResult:", zap.String("result", string(result)))

	// 修复返回值为 pending 时的解析问题
	type CommonResp struct {
		Success bool        `json:"success"`
		Result  interface{} `json:"result"`
	}

	var commonResp CommonResp
	err = jsoniter.Unmarshal(result, &commonResp)
	if err != nil {
		return nil, err
	}

	if commonResp.Success && commonResp.Result == "pending" {
		// 如果 success 为 true 且 result 为 "pending"，则不解析 GetOrderResultResp 中的 result 键
		return &GetOrderResultResp{
			ShowStartCommonResp: &ShowStartCommonResp{
				Success: commonResp.Success,
			},
		}, nil
	}

	// 否则，继续解析为 GetOrderResultResp
	var resp GetOrderResultResp
	err = jsoniter.Unmarshal(result, &resp)
	if err != nil {
		return nil, err
	}

	if resp.State == "token-expire-at" || resp.Msg == "登录过期了，请重新登录！" {
		err = c.GetToken(ctx)
		if err != nil {
			return nil, err
		}
		return c.GetOrderResult(ctx, orderJobKey)
	}

	return &resp, nil
}
