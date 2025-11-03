package cron

import (
	"context"
	"strconv"
	"sync"
	"time"
	"washwise/config"
	"washwise/model"

	log "github.com/sirupsen/logrus"
)

// TaskManager 任务管理器
type TaskManager struct {
	ctx    context.Context
	cancel context.CancelFunc

	// 内存存储
	machineTypes    map[string]*GetMachineTypesResp // shopId -> types
	machineTypesMux sync.RWMutex

	// Tickers
	typeTicker    *time.Ticker
	machineTicker *time.Ticker
	detailTicker  *time.Ticker
}

var tm *TaskManager

// InitTaskManager 初始化任务管理器
func InitTaskManager() *TaskManager {
	ctx, cancel := context.WithCancel(context.Background())
	tm = &TaskManager{
		ctx:          ctx,
		cancel:       cancel,
		machineTypes: make(map[string]*GetMachineTypesResp),
	}
	return tm
}

// GetTaskManager 获取任务管理器实例
func GetTaskManager() *TaskManager {
	return tm
}

// Start 启动所有定时任务
func (tm *TaskManager) Start() {
	cfg := config.Get()

	// 立即执行一次初始化
	log.Info("开始初始化数据...")
	tm.fetchMachineTypes()
	tm.fetchMachines()
	tm.fetchMachineDetails()
	log.Info("数据初始化完成")

	// 启动定时任务
	tm.typeTicker = time.NewTicker(config.GetMachineTypesInterval())
	tm.machineTicker = time.NewTicker(config.GetMachinesInterval())
	tm.detailTicker = time.NewTicker(config.GetMachineDetailsInterval())

	go tm.runMachineTypesTask()
	go tm.runMachinesTask()
	go tm.runMachineDetailsTask()

	log.WithFields(log.Fields{
		"machine_types_interval":   cfg.Cron.MachineTypesInterval,
		"machines_interval":        cfg.Cron.MachinesInterval,
		"machine_details_interval": cfg.Cron.MachineDetailsInterval,
	}).Info("定时任务已启动")
}

// Stop 停止所有定时任务
func (tm *TaskManager) Stop() {
	log.Info("正在停止定时任务...")
	tm.cancel()
	if tm.typeTicker != nil {
		tm.typeTicker.Stop()
	}
	if tm.machineTicker != nil {
		tm.machineTicker.Stop()
	}
	if tm.detailTicker != nil {
		tm.detailTicker.Stop()
	}
	log.Info("定时任务已停止")
}

// runMachineTypesTask 运行机器类型获取任务（大周期）
func (tm *TaskManager) runMachineTypesTask() {
	for {
		select {
		case <-tm.ctx.Done():
			return
		case <-tm.typeTicker.C:
			tm.fetchMachineTypes()
		}
	}
}

// runMachinesTask 运行机器列表获取任务（中等周期）
func (tm *TaskManager) runMachinesTask() {
	for {
		select {
		case <-tm.ctx.Done():
			return
		case <-tm.machineTicker.C:
			tm.fetchMachines()
		}
	}
}

// runMachineDetailsTask 运行机器详情获取任务（短周期）
func (tm *TaskManager) runMachineDetailsTask() {
	for {
		select {
		case <-tm.ctx.Done():
			return
		case <-tm.detailTicker.C:
			tm.fetchMachineDetails()
		}
	}
}

// fetchMachineTypes 获取所有商店的机器类型
func (tm *TaskManager) fetchMachineTypes() {
	cfg := config.Get()
	log.Info("开始获取机器类型...")

	for _, shopId := range cfg.Shops {
		resp, err := GetMachineTypes(tm.ctx, shopId)
		if err != nil {
			log.WithError(err).WithField("shopId", shopId).Error("获取机器类型失败")
			continue
		}

		tm.machineTypesMux.Lock()
		tm.machineTypes[shopId] = resp
		tm.machineTypesMux.Unlock()

		log.WithFields(log.Fields{
			"shopId": shopId,
			"count":  len(resp.Items),
		}).Info("获取机器类型成功")
	}
}

// fetchMachines 获取所有商店、所有类型的机器列表
func (tm *TaskManager) fetchMachines() {
	cfg := config.Get()
	log.Info("开始获取机器列表...")

	totalCount := 0

	for _, shopId := range cfg.Shops {
		// 获取该商店的机器类型
		tm.machineTypesMux.RLock()
		types, exists := tm.machineTypes[shopId]
		tm.machineTypesMux.RUnlock()

		if !exists || types == nil {
			log.WithField("shopId", shopId).Warn("未找到机器类型，跳过")
			continue
		}

		// 遍历所有机器类型
		for _, machineType := range types.Items {
			resp, err := GetMachines(tm.ctx, shopId, machineType.MachineTypeId, 1000, 1)
			if err != nil {
				log.WithError(err).WithFields(log.Fields{
					"shopId":        shopId,
					"machineTypeId": machineType.MachineTypeId,
				}).Error("获取机器列表失败")
				break
			}

			if len(resp.Items) == 0 {
				break
			}

			// 将机器数据持久化到数据库
			machines := make([]model.Machine, 0, len(resp.Items))
			for _, item := range resp.Items {
				id, _ := strconv.ParseInt(item.Id, 10, 64)
				machines = append(machines, model.Machine{
					Id:     id,
					Name:   item.Name,
					Status: model.MachineStatusOffline,
					ShopId: shopId,
				})
			}

			if err := model.UpsertMachines(machines); err != nil {
				log.WithError(err).Error("保存机器数据失败")
			} else {
				totalCount += len(machines)
			}
		}
	}

	log.WithField("count", totalCount).Info("获取机器列表完成")
}

// fetchMachineDetails 获取所有机器的详情
func (tm *TaskManager) fetchMachineDetails() {
	log.Info("开始获取机器详情...")

	// 从数据库获取所有机器
	machines, err := model.GetAllMachines()
	if err != nil {
		log.WithError(err).Error("从数据库获取机器列表失败")
		return
	}

	successCount := 0
	for _, machine := range machines {
		detail, err := GetMachineDetail(tm.ctx, machine.Id)
		if err != nil {
			log.WithError(err).WithField("machineId", machine.Id).Warn("获取机器详情失败")
			continue
		}

		// 更新机器信息
		machine.Name = detail.Name
		machine.ShopId = detail.ShopId
		machine.Msg = ""
		if detail.DeviceErrorMsg != nil {
			machine.Msg = *detail.DeviceErrorMsg
		}
		// 仅当机器状态从可用变为使用中时，更新最后使用时间
		if machine.Status == model.MachineStatusAvailable && detail.DeviceErrorCode == model.MachineStatusInUse {
			machine.LastUseTime = time.Now().Unix()
		}
		machine.Status = detail.DeviceErrorCode

		if err := model.UpdateMachine(&machine); err != nil {
			log.WithError(err).WithField("machineId", machine.Id).Warn("更新机器信息失败")
			continue
		}

		successCount++
	}

	log.WithFields(log.Fields{
		"total":   len(machines),
		"success": successCount,
		"fail":    len(machines) - successCount,
	}).Info("获取机器详情完成")
}

// GetMachineTypes 获取指定商店的机器类型（从内存）
func (tm *TaskManager) GetMachineTypesFromMemory(shopId string) *GetMachineTypesResp {
	tm.machineTypesMux.RLock()
	defer tm.machineTypesMux.RUnlock()
	return tm.machineTypes[shopId]
}
