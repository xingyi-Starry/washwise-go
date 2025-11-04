package servicev1

import (
	"strconv"
	"time"
	"washwise/model"
	"washwise/util"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

// api/v1/getLaundryMachines?LaundryID=1
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
		if machine.Code == model.MachineCodeInUse {
			remainTime = int(max(machine.LastUseTime+40*60-time.Now().Unix(), 0))
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

// /api/v1/getMachineDetail?MachineID=1
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
