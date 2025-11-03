package cron

type CommonResp[G any] struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data G      `json:"data,omitempty"`
	T    int64  `json:"t"`
}

type GetMachineTypesReq struct {
	ShopId string `url:"shopId"`
}

type GetMachineTypesResp struct {
	Items []struct {
		MachineTypeId   string `json:"machineTypeId"`
		MachineTypeName string `json:"machineTypeName"`
	} `json:"items"`
}

type GetMachinesReq struct {
	ShopId        string `url:"shopId"`
	MachineTypeId string `url:"machineTypeId"`
	PageSize      int    `url:"pageSize"`
	Page          int    `url:"page"`
}

type GetMachinesResp struct {
	Items []struct {
		Id     string `json:"id"`
		Type   int    `json:"type"`
		Name   string `json:"name"`
		Status int    `json:"status"`
	} `json:"items"`
	GoodsPage int `json:"goodsPage"`
}

type GetMachineDetailReq struct {
	GoodsId int64 `url:"goodsId"`
}

type GetMachineDetailResp struct {
	GoodsId             int64   `json:"goodsId"`
	Name                string  `json:"name"`
	DeviceId            int64   `json:"deviceId"`
	RemainTime          int     `json:"remainTime"`
	MachineId           string  `json:"machineId"`
	ShopId              string  `json:"shopId"`
	DeviceErrorCode     int     `json:"deviceErrorCode"`
	DeviceErrorMsg      *string `json:"deviceErrorMsg"`
	AnnouncementContent *string `json:"announcementContent"`
}
