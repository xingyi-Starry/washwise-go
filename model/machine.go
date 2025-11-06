package model

import (
	"gorm.io/gorm/clause"
)

const (
	MachineCodeAvailable = 0
	MachineCodeOffline   = 1
	MachineCodeInUse     = 2
)

type Machine struct {
	Id          int64 `gorm:"primaryKey"`
	Name        string
	Code        int
	LastUseTime int64
	Msg         string
	AvgUseTime  int64  // 平均使用时间，单位秒
	ShopId      string `gorm:"index"`
	Type        string
}

// GetMachinesByShopID 根据商店ID获取所有机器
func GetMachinesByShopID(shopId string) ([]Machine, error) {
	var machines []Machine
	err := db.Where("shop_id = ?", shopId).Find(&machines).Error
	return machines, err
}

// GetMachineByID 根据机器ID获取单个机器
func GetMachineByID(machineId int64) (*Machine, error) {
	var machine Machine
	err := db.Where("id = ?", machineId).First(&machine).Error
	return &machine, err
}

// GetAllMachines 获取所有机器
func GetAllMachines() ([]*Machine, error) {
	var machines []*Machine
	err := db.Find(&machines).Error
	return machines, err
}

// UpdateMachine 更新机器信息
func UpdateMachine(machine *Machine) error {
	return db.Save(machine).Error
}

// InsertMachinesIfNotExists 批量插入机器（如果不存在）
func InsertMachinesIfNotExists(machines []Machine) error {
	if len(machines) == 0 {
		return nil
	}
	return db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoNothing: true,
	}).Create(&machines).Error
}
