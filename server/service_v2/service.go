package servicev2

import (
	"strconv"
	"time"
	"washwise/config"
	"washwise/model"
	"washwise/util"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

var shops = map[string]string{
	"202401041041470000069996565184": "沙河雁北洗衣房",
	"202401041044000000069996552384": "沙河雁南洗衣房",
	"202302071714530000012067133598": "海淀西土城校区",
}

// 暂无获取洗衣房名称的接口，先硬编码
func getShopName(shopId string) string {
	name, ok := shops[shopId]
	if !ok {
		return "未知洗衣房"
	}
	return name
}

// @Summary 获取店铺列表
// @Description 获取店铺列表
// @Tags v2
// @Produce json
// @Success 200 {object} GetShopsResp
// @Router /api/v2/shops [get]
func GetShops(c *fiber.Ctx) error {
	resp := &GetShopsResp{}
	for _, shopId := range config.Get().Shops {
		resp.Items = append(resp.Items, &GetShopsRespItem{
			Id:   shopId,
			Name: getShopName(shopId),
		})
	}
	return c.JSON(resp)
}

// @Summary 获取洗衣机列表
// @Description 获取洗衣机列表
// @Tags v2
// @Param shopId query string true "店铺ID"
// @Produce json
// @Success 200 {object} GetMachinesResp
// @Router /api/v2/machines [get]
func GetMachines(c *fiber.Ctx) error {
	req := &GetMachinesReq{}
	if err := c.QueryParser(req); err != nil {
		return util.BadRequest(c, err.Error())
	}

	machines, err := model.GetMachinesWithUsageCount(req.ShopId, time.Now().AddDate(0, 0, -7).Unix(), time.Now().Unix())
	if err != nil {
		logrus.WithError(err).Error("db error")
		return util.Internal(c)
	}

	resp := &GetMachinesResp{}

	for _, machine := range machines {

		remainTime := int64(0)
		if machine.Code == model.MachineCodeInUse {
			predictUseTime := machine.AvgUseTime
			if predictUseTime == 0 {
				predictUseTime = 45 * 60 // 默认45分钟
			}
			remainTime = max(machine.LastUseTime+predictUseTime-time.Now().Unix(), 0)
		}

		resp.Items = append(resp.Items, &GetMachinesRespItem{
			Id:         machine.Id,
			Name:       machine.Name,
			Type:       machine.Type,
			Msg:        machine.Msg,
			Status:     machine.Code,
			UsageCount: machine.UsageCount,
			RemainTime: remainTime,
		})
	}
	return c.JSON(resp)
}

// @Summary 获取洗衣机详情
// @Description 获取洗衣机详情
// @Tags v2
// @Param machineId path string true "洗衣机ID"
// @Produce json
// @Success 200 {object} MachineDetailResp
// @Router /api/v2/machine/{machineId} [get]
func GetMachine(c *fiber.Ctx) error {
	machineIdStr := c.Params("machineId")
	machineId, err := strconv.ParseInt(machineIdStr, 10, 64)
	if err != nil {
		return util.BadRequest(c, "machineId is required")
	}

	// 获取机器信息
	machine, err := model.GetMachineByID(machineId)
	if err != nil {
		logrus.WithError(err).Error("db error")
		return util.Internal(c)
	}

	// 计算剩余时间
	remainTime := int64(0)
	if machine.Code == model.MachineCodeInUse {
		predictUseTime := machine.AvgUseTime
		if predictUseTime == 0 {
			predictUseTime = 45 * 60 // 默认45分钟
		}
		remainTime = max(machine.LastUseTime+predictUseTime-time.Now().Unix(), 0)
	}

	// 获取近7天的使用历史
	now := time.Now()
	history := make(map[string]int)

	for i := 6; i >= 0; i-- {
		date := now.AddDate(0, 0, -i)
		dateStr := date.Format("2006-01-02")

		// 计算该天的起止时间
		startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
		endOfDay := time.Date(date.Year(), date.Month(), date.Day(), 23, 59, 59, 0, date.Location())

		// 查询该天的使用次数
		count, err := model.CountUsagesByMachineIDAndTimeRange(
			machineId,
			startOfDay.Unix(),
			endOfDay.Unix(),
		)
		if err != nil {
			logrus.WithError(err).Error("db error")
			return util.Internal(c)
		}

		history[dateStr] = int(count)
	}

	// 构建响应
	resp := &MachineDetailResp{
		Id:          machine.Id,
		Name:        machine.Name,
		Type:        machine.Type,
		Msg:         machine.Msg,
		Status:      machine.Code,
		RemainTime:  remainTime,
		AvgUseTime:  machine.AvgUseTime,
		LastUseTime: machine.LastUseTime,
		History:     history,
	}

	return c.JSON(resp)
}
