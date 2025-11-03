package model

import (
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// CreateMachine 创建新机器
func CreateMachine(machine *Machine) error {
	return db.Create(machine).Error
}

// GetMachineByID 根据ID获取机器
func GetMachineByID(id int64) (*Machine, error) {
	var machine Machine
	err := db.First(&machine, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &machine, nil
}

// GetMachinesByShopID 根据商店ID获取所有机器
func GetMachinesByShopID(shopId string) ([]Machine, error) {
	var machines []Machine
	err := db.Where("shop_id = ?", shopId).Find(&machines).Error
	return machines, err
}

// GetAllMachines 获取所有机器
func GetAllMachines() ([]Machine, error) {
	var machines []Machine
	err := db.Find(&machines).Error
	return machines, err
}

// UpdateMachine 更新机器信息
func UpdateMachine(machine *Machine) error {
	return db.Save(machine).Error
}

// UpdateMachineFields 更新机器的指定字段
func UpdateMachineFields(id int64, updates map[string]interface{}) error {
	return db.Model(&Machine{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteMachine 删除机器
func DeleteMachine(id int64) error {
	return db.Delete(&Machine{}, id).Error
}

// DeleteMachinesByShopID 删除指定商店的所有机器
func DeleteMachinesByShopID(shopId string) error {
	return db.Where("shop_id = ?", shopId).Delete(&Machine{}).Error
}

// UpsertMachine 插入或更新机器（基于主键ID）
func UpsertMachine(machine *Machine) error {
	return db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"name", "status", "last_use_time", "msg", "shop_id"}),
	}).Create(machine).Error
}

// UpsertMachines 批量插入或更新机器
func UpsertMachines(machines []Machine) error {
	if len(machines) == 0 {
		return nil
	}
	return db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"name", "status", "last_use_time", "msg", "shop_id"}),
	}).Create(&machines).Error
}

// GetMachinesByStatus 根据状态获取机器列表
func GetMachinesByStatus(status int) ([]Machine, error) {
	var machines []Machine
	err := db.Where("status = ?", status).Find(&machines).Error
	return machines, err
}

// CountMachinesByShopID 统计指定商店的机器数量
func CountMachinesByShopID(shopId string) (int64, error) {
	var count int64
	err := db.Model(&Machine{}).Where("shop_id = ?", shopId).Count(&count).Error
	return count, err
}

// CountMachinesByStatus 统计指定状态的机器数量
func CountMachinesByStatus(status int) (int64, error) {
	var count int64
	err := db.Model(&Machine{}).Where("status = ?", status).Count(&count).Error
	return count, err
}
