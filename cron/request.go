package cron

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"washwise/util"
)

var c = http.Client{}

func header() http.Header {
	return http.Header{
		"accept":          []string{"*/*"},
		"accept-language": []string{"zh-CN,zh;q=0.9"},
		"channel":         []string{"wechat"},
		"content-type":    []string{"application/x-www-form-urlencoded"},
		"sec-fetch-dest":  []string{"empty"},
		"sec-fetch-mode":  []string{"cors"},
		"sec-fetch-site":  []string{"cross-site"},
		"uid":             []string{"undefined"},
		"xweb_xhr":        []string{"1"},
		"User-Agent": []string{
			"BUPTGateWay Mozilla/5.0 (Windows NT 10.0; Win64; x64) ",
			"AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36",
		},
	}
}

func doPost[G any](ctx context.Context, url string, bodyData any) (*G, error) {
	body, err := util.UrlEncode(bodyData)
	if err != nil {
		return nil, err
	}

	bodyReader := strings.NewReader(body)

	req, err := http.NewRequestWithContext(ctx, "POST", url, bodyReader)
	if err != nil {
		return nil, err
	}
	req.Header = header()

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var respBody CommonResp[*G]
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return nil, err
	}

	if respBody.Code != 0 {
		return nil, errors.New(respBody.Msg)
	}

	return respBody.Data, nil
}

func GetMachineTypes(ctx context.Context, shopId string) (*GetMachineTypesResp, error) {
	return doPost[GetMachineTypesResp](ctx, "https://userapi.qiekj.com/machineModel/nearByList", GetMachineTypesReq{shopId})
}

func GetMachines(ctx context.Context, shopId string, machineTypeId string, pageSize, page int) (*GetMachinesResp, error) {
	return doPost[GetMachinesResp](ctx, "https://userapi.qiekj.com/machineModel/near/machines", GetMachinesReq{
		ShopId:        shopId,
		MachineTypeId: machineTypeId,
		PageSize:      pageSize,
		Page:          page,
	})
}

func GetMachineDetail(ctx context.Context, goodsId int64) (*GetMachineDetailResp, error) {
	return doPost[GetMachineDetailResp](ctx, "https://userapi.qiekj.com/goods/normal/details", GetMachineDetailReq{
		GoodsId: goodsId,
	})
}
