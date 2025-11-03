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
			remainTime = int(machine.LastUseTime + 40*60 - time.Now().Unix())
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
// TODO
