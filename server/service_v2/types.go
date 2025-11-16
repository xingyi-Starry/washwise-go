package servicev2

type CommonResp[G any] struct {
	Code int    `json:"code"`
	Msg  string `json:"msg,omitempty"`
	Data G      `json:"data,omitempty"`
}

type GetShopsResp struct {
	Items []*GetShopsRespItem `json:"items"`
}

type GetShopsRespItem struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type GetMachinesReq struct {
	ShopId string `query:"shopId"`
}

type GetMachinesResp struct {
	Items []*GetMachinesRespItem `json:"items"`
}

type GetMachinesRespItem struct {
	Id         int64  `json:"id"`
	Name       string `json:"name"`
	Type       string `json:"type"`
	Msg        string `json:"msg"`
	Status     int    `json:"status"`
	UsageCount int    `json:"usageCount"`
	RemainTime int64  `json:"remainTime"`
}

type MachineDetailResp struct {
	Id         int64  `json:"id"`
	Name       string `json:"name"`
	Type       string `json:"type"`
	Msg        string `json:"msg"`
	Status     int    `json:"status"`
	RemainTime int64  `json:"remainTime"`

	AvgUseTime  int64          `json:"avgUseTime"`  // 平均使用时间，即预计使用时间，单位秒
	LastUseTime int64          `json:"lastUseTime"` // 上个人开始使用时间
	History     map[string]int `json:"history"`     // date -> usage count
}
