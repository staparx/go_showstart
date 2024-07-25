package client

type ShowStartCommonResp struct {
	State   string `json:"state"`
	Success bool   `json:"success"`
	Msg     string `json:"msg"`
	Status  int    `json:"status"`
	TraceID string `json:"traceId"`
}

type GetTokenResp struct {
	*ShowStartCommonResp
	Result struct {
		AccessToken struct {
			AccessToken string `json:"access_token"`
			Expire      int    `json:"expire"`
		} `json:"accessToken"`
		IDToken struct {
			IDToken string `json:"id_token"`
			Expire  int    `json:"expire"`
		} `json:"idToken"`
	} `json:"result"`
}

type ActivityDetailResp struct {
	*ShowStartCommonResp
	Result struct {
		ActivityID        int           `json:"activityId"`
		IsLogin           int           `json:"isLogin"`
		ActivityName      string        `json:"activityName"`
		Price             string        `json:"price"`
		ActivityLevel     int           `json:"activityLevel"`
		ActivityLevelName string        `json:"activityLevelName"`
		ShowTime          string        `json:"showTime"`
		ShowTimeType      int           `json:"showTimeType"`
		Avatar            string        `json:"avatar"`
		Album             []string      `json:"album"`
		Document          string        `json:"document"`
		SellIdentity      int           `json:"sellIdentity"`
		ActivityTag       string        `json:"activityTag"`
		Tags              string        `json:"tags"`
		Banner            []interface{} `json:"banner"`
		Music             []interface{} `json:"music"`
		ShowLetter        bool          `json:"showLetter"`
		WhetherWantTo     bool          `json:"whetherWantTo"`
		WantToNum         int           `json:"wantToNum"`
		Site              struct {
			ID        int     `json:"id"`
			Name      string  `json:"name"`
			Avatar    string  `json:"avatar"`
			Address   string  `json:"address"`
			Photo     string  `json:"photo"`
			Contact   string  `json:"contact"`
			Longitude float64 `json:"longitude"`
			Latitude  float64 `json:"latitude"`
			CityName  string  `json:"cityName"`
			UserType  int     `json:"userType"`
		} `json:"site"`
		Host []struct {
			ID               int    `json:"id"`
			Name             string `json:"name"`
			Avatar           string `json:"avatar"`
			UserType         int    `json:"userType"`
			ActivityRoleType int    `json:"activityRoleType"`
			IsCollect        int    `json:"isCollect"`
		} `json:"host"`
		UserInfos        []interface{} `json:"userInfos"`
		SessionUserInfos []struct {
			Title     string `json:"title"`
			SessionID int    `json:"sessionId"`
			Selected  int    `json:"selected"`
			IsEnd     int    `json:"isEnd"`
			UserInfos []struct {
				ID               int    `json:"id"`
				Name             string `json:"name"`
				Avatar           string `json:"avatar"`
				UserType         int    `json:"userType"`
				ActivityRoleType int    `json:"activityRoleType"`
				RoleType         int    `json:"roleType"`
				IsCollect        int    `json:"isCollect"`
			} `json:"userInfos"`
		} `json:"sessionUserInfos"`
		URL           string `json:"url"`
		Title         string `json:"title"`
		ShowStartTime int64  `json:"showStartTime"`
		ShowEndTime   int64  `json:"showEndTime"`
		GoodsList     []struct {
			GoodsID          int    `json:"goodsId"`
			GoodsName        string `json:"goodsName"`
			GoodsPoster      string `json:"goodsPoster"`
			BindName         string `json:"bindName"`
			Price            string `json:"price"`
			BuyGroupType     int    `json:"buyGroupType"`
			DeliveryTimeType int    `json:"deliveryTimeType"`
		} `json:"goodsList"`
		GoodsNum                 int           `json:"goodsNum"`
		RealName                 int           `json:"realName"`
		RealNameValidType        int           `json:"realNameValidType"`
		ShareActivityIds         []interface{} `json:"shareActivityIds"`
		SellTerminal             int           `json:"sellTerminal"`
		IsShowCollection         int           `json:"isShowCollection"`
		ServiceTemplateContent   string        `json:"serviceTemplateContent"`
		ServiceTemplateEnContent string        `json:"serviceTemplateEnContent"`
		BeginTimeConfirmed       int           `json:"beginTimeConfirmed"`
		Notices                  []interface{} `json:"notices"`
		Advertising              []interface{} `json:"advertising"`
		SquadStatus              int           `json:"squadStatus"`
		ServiceTemplates         []string      `json:"serviceTemplates"`
		StylesTopic              []struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"stylesTopic"`
		Coupons        []interface{} `json:"coupons"`
		Labels         []interface{} `json:"labels"`
		DouyinStatus   int           `json:"douyinStatus"`
		IsPreAuth      int           `json:"isPreAuth"`
		PreAuthHelpURL string        `json:"preAuthHelpUrl"`
		PreAuthTips    string        `json:"preAuthTips"`
		OpenStrategy   int           `json:"openStrategy"`
	} `json:"result"`
}

type ActivityTicketListResp struct {
	*ShowStartCommonResp
	Result []*ActivityTicket `json:"result"`
}

type ActivityTicket struct {
	SessionName                 string        `json:"sessionName"`
	SessionID                   int           `json:"sessionId"`
	IsConfirmedStartTime        int           `json:"isConfirmedStartTime"`
	CommonPerformerDocumentType string        `json:"commonPerformerDocumentType"`
	IsSupportTransform          int           `json:"isSupportTransform"`
	TicketList                  []*TicketInfo `json:"ticketList"`
	TicketPriceList             []struct {
		Price      string        `json:"price"`
		TicketList []*TicketInfo `json:"ticketList"`
	} `json:"ticketPriceList"`
}

type TicketInfo struct {
	TicketID                  string        `json:"ticketId"`
	TicketType                string        `json:"ticketType"`
	SellingPrice              string        `json:"sellingPrice"`
	IsUndetermined            int           `json:"isUndetermined"`
	CostPrice                 string        `json:"costPrice"`
	TicketNum                 int           `json:"ticketNum"`
	ValidateType              int           `json:"validateType"`
	Time                      string        `json:"time"`
	Instruction               string        `json:"instruction"`
	Countdown                 int           `json:"countdown"`
	RemainTicket              int           `json:"remainTicket"`
	SaleStatus                int           `json:"saleStatus"`
	ActivityID                int           `json:"activityId"`
	GoodType                  int           `json:"goodType"`
	Telephone                 string        `json:"telephone"`
	AreaCode                  string        `json:"areaCode"`
	LimitBuyNum               int           `json:"limitBuyNum"`
	CanBuyNum                 int           `json:"canBuyNum"`
	CityName                  string        `json:"cityName"`
	UnPayOrderNum             int           `json:"unPayOrderNum"`
	Type                      int           `json:"type"`
	PickupAddress             string        `json:"pickupAddress"`
	EntityMailInstruction     string        `json:"entityMailInstruction"`
	EntityPickupInstruction   string        `json:"entityPickupInstruction"`
	BuyType                   int           `json:"buyType"`
	CanAddGoods               int           `json:"canAddGoods"`
	TicketRecordStatus        int           `json:"ticketRecordStatus"`
	StartSellNoticeStatus     int           `json:"startSellNoticeStatus"`
	ShowRuleTip               bool          `json:"showRuleTip"`
	StartTime                 int64         `json:"startTime"`
	SessionID                 int           `json:"sessionId"`
	ShowTime                  string        `json:"showTime"`
	MemberNum                 int           `json:"memberNum"`
	Labels                    []interface{} `json:"labels"`
	EndTime                   string        `json:"endTime"`
	TransformNum              int           `json:"transformNum"`
	IsPreAuth                 int           `json:"isPreAuth"`
	ConfirmPreOrderDetailTips string        `json:"confirmPreOrderDetailTips"`
	GroupID                   int           `json:"groupId"`
	IsCommonPerformerList     int           `json:"isCommonPerformerList"`
}

type ConfirmResp struct {
	*ShowStartCommonResp
	Result struct {
		OrderInfoVo struct {
			Title         string `json:"title"`
			Poster        string `json:"poster"`
			SiteName      string `json:"siteName"`
			CityName      string `json:"cityName"`
			SessionID     int    `json:"sessionId"`
			ActivityID    int    `json:"activityId"`
			ShowTime      string `json:"showTime"`
			AreaCode      string `json:"areaCode"`
			Telephone     string `json:"telephone"`
			TicketPriceVo struct {
				TicketID                string  `json:"ticketId"`
				TicketName              string  `json:"ticketName"`
				Price                   float64 `json:"price"`
				TicketType              int     `json:"ticketType"`
				LimitBuyNum             int     `json:"limitBuyNum"`
				CanBuyNum               int     `json:"canBuyNum"`
				Instruction             string  `json:"instruction"`
				EntityMailInstruction   string  `json:"entityMailInstruction"`
				EntityPickupInstruction string  `json:"entityPickupInstruction"`
				PickupAddress           string  `json:"pickupAddress"`
				RemainTicket            int     `json:"remainTicket"`
				TransformNum            int     `json:"transformNum"`
				DyPOIType               int     `json:"dyPOIType"`
			} `json:"ticketPriceVo"`
			RealName                    int    `json:"realName"`
			ValidateType                int    `json:"validateType"`
			BuyType                     int    `json:"buyType"`
			MemberNum                   int    `json:"memberNum"`
			SellAreaType                int    `json:"sellAreaType"`
			DouyinStatus                int    `json:"douyinStatus"`
			CommonPerformerDocumentType string `json:"commonPerformerDocumentType"`
			IsSupportTransform          int    `json:"isSupportTransform"`
			IsCommonPerformerList       int    `json:"isCommonPerformerList"`
		} `json:"orderInfoVo"`
		ActivityTips []string `json:"activityTips"`
	} `json:"result"`
}

type CpListResp struct {
	*ShowStartCommonResp
	Result []struct {
		ID                 int    `json:"id"`
		UserID             int    `json:"userId"`
		CanBuy             int    `json:"canBuy"`
		Name               string `json:"name"`
		DocumentType       int    `json:"documentType"`
		DocumentTypeStr    string `json:"documentTypeStr"`
		ShowDocumentNumber string `json:"showDocumentNumber"`
		IsSelf             int    `json:"isSelf"`
	} `json:"result"`
}

type OrderListResp struct {
	*ShowStartCommonResp
	Result struct {
		CouponList   []interface{} `json:"couponList"`
		HelpURL      string        `json:"helpURL"`
		CanUseNum    int           `json:"canUseNum"`
		CanNotUseNum int           `json:"canNotUseNum"`
	} `json:"result"`
}

type OrderResp struct {
	*ShowStartCommonResp
	//DebugMsgInner struct {
	//	Db []struct {
	//		DbCost float64 `json:"dbCost"`
	//		SQL    string  `json:"sql"`
	//	} `json:"db"`
	//	Docheck [][][]int `json:"docheck"`
	//	HTTP    []struct {
	//		Body string `json:"body"`
	//		Cmd  string `json:"cmd"`
	//		Res  struct {
	//			Body   string `json:"body"`
	//			Status int    `json:"status"`
	//			Header struct {
	//				ContentLength int    `json:"Content-Length"`
	//				ContentType   string `json:"Content-Type"`
	//			} `json:"header"`
	//			Truncated bool `json:"truncated"`
	//		} `json:"res"`
	//		Cost int `json:"cost"`
	//	} `json:"http"`
	//	Redis []struct {
	//		Key  []string `json:"key"`
	//		Cmd  string   `json:"cmd"`
	//		Cost float64  `json:"cost"`
	//	} `json:"redis"`
	//} `json:"debugMsgInner"`
	Result struct {
		OrderJobKey  string  `json:"orderJobKey"`
		Sleep        float64 `json:"sleep"`
		CoreOrderKey string  `json:"coreOrderKey"`
		SleepExipre  float64 `json:"sleepExipre"`
	} `json:"result"`
}

type OrderCoreResp struct {
	*ShowStartCommonResp
	Result interface{} `json:"result"`
	//Result struct {
	//	OrderJobKey string `json:"orderJobKey"`
	//} `json:"result"`
}

type GetOrderResultResp struct {
	*ShowStartCommonResp
	//DebugMsgInner struct {
	//	Redis []struct {
	//		Key  string  `json:"key"`
	//		Cmd  string  `json:"cmd"`
	//		Cost float64 `json:"cost"`
	//	} `json:"redis"`
	//} `json:"debugMsgInner"`
	Result struct {
		OrderSn string  `json:"orderSn"`
		Cost    float64 `json:"cost"`
		OrderID string  `json:"orderId"`
	} `json:"result"`
}
