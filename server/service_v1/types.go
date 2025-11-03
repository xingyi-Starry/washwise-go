package servicev1

type MachineInfo struct {
	Name       string `json:"name"`
	DeviceCode int    `json:"deviceCode"`
	DeviceMsg  string `json:"deviceMsg"`
	RemainTime int    `json:"remainTime"` // 剩余时间，单位：分钟
	ErrorCount int    `json:"errorCount"`
}

type GetLaundryMachinesResp struct {
	Data map[string]MachineInfo `json:"洗衣机"`
}
