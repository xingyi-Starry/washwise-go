package servicev1

import (
	"strconv"
	"time"
	"washwise/model"
	"washwise/util"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

// @Summary 获取洗衣机列表
// @Description 获取洗衣机列表
// @Tags v1
// @Param LaundryID query string true "店铺ID"
// @Produce json
// @Success 200 {object} GetLaundryMachinesResp
// @Router /api/v1/getLaundryMachines [get]
func GetLaundryMachines(c *fiber.Ctx) error {
	shopId := c.Query("LaundryID")
	if shopId == "" {
		return util.BadRequest(c, "LaundryID is required")
	}
	machines, err := model.GetMachinesByShopID(shopId)
	if err != nil {
		logrus.WithError(err).Error("db error")
		return util.Internal(c)
	}

	// build resp
	resp := GetLaundryMachinesResp{make(map[string]MachineInfo)}
	for _, machine := range machines {
		k := strconv.FormatInt(machine.Id, 10)

		remainTime := 0
		predictUseTime := machine.AvgUseTime
		if predictUseTime == 0 {
			predictUseTime = 45 * 60 // 默认45分钟
		}
		if machine.Code == model.MachineCodeInUse {
			remainTime = int(max(machine.LastUseTime+predictUseTime-time.Now().Unix(), 0))
		}
		resp.Data[k] = MachineInfo{
			Name:       machine.Name,
			DeviceCode: machine.Code,
			DeviceMsg:  machine.Msg,
			RemainTime: remainTime,
			ErrorCount: 0,
		}
	}

	return c.JSON(resp)
}

// @Summary 获取洗衣机详情
// @Description 获取洗衣机详情
// @Tags v1
// @Param MachineID query string true "洗衣机ID"
// @Produce json
// @Success 200 {object} GetMachineDetailResp
// @Router /api/v1/getMachineDetail [get]
func GetMachineDetail(c *fiber.Ctx) error {
	machineIdStr := c.Query("MachineID")
	machineId, err := strconv.ParseInt(machineIdStr, 10, 64)
	if err != nil {
		return util.BadRequest(c, "MachineID is required")
	}

	now := time.Now()
	resp := make(GetMachineDetailResp)

	// 分7个区间依次查询并构建响应
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

		resp[dateStr] = int(count)
	}

	return c.JSON(resp)
}
