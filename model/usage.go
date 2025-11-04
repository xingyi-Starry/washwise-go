package model

import (
	"errors"

	"gorm.io/gorm"
)

type Usage struct {
	Id        int64 `gorm:"primaryKey,autoIncrement"`
	MachineId int64 `gorm:"index"`
	StartTime int64 `gorm:"index"`
	EndTime   int64
}

// CreateUsage 创建新使用记录
func CreateUsage(usage *Usage) error {
	return db.Create(usage).Error
}

// GetUsageByID 根据ID获取使用记录
func GetUsageByID(id int64) (*Usage, error) {
	var usage Usage
	err := db.First(&usage, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &usage, nil
}

// GetUsagesByMachineID 根据机器ID获取所有使用记录
func GetUsagesByMachineID(machineId int64) ([]Usage, error) {
	var usages []Usage
	err := db.Where("machine_id = ?", machineId).Order("start_time DESC").Find(&usages).Error
	return usages, err
}

// GetUsagesByTimeRange 根据时间范围获取使用记录
func GetUsagesByTimeRange(startTime, endTime int64) ([]Usage, error) {
	var usages []Usage
	err := db.Where("start_time >= ? AND start_time <= ?", startTime, endTime).
		Order("start_time DESC").
		Find(&usages).Error
	return usages, err
}

// GetUsagesByMachineIDAndTimeRange 根据机器ID和时间范围获取使用记录
func GetUsagesByMachineIDAndTimeRange(machineId, startTime, endTime int64) ([]Usage, error) {
	var usages []Usage
	err := db.Where("machine_id = ? AND start_time >= ? AND start_time <= ?", machineId, startTime, endTime).
		Order("start_time DESC").
		Find(&usages).Error
	return usages, err
}

// GetAllUsages 获取所有使用记录
func GetAllUsages() ([]Usage, error) {
	var usages []Usage
	err := db.Order("start_time DESC").Find(&usages).Error
	return usages, err
}

// UpdateUsage 更新使用记录
func UpdateUsage(usage *Usage) error {
	return db.Save(usage).Error
}

// UpdateUsageFields 更新使用记录的指定字段
func UpdateUsageFields(id int64, updates map[string]interface{}) error {
	return db.Model(&Usage{}).Where("id = ?", id).Updates(updates).Error
}

// UpdateUsageEndTime 更新使用记录的结束时间
func UpdateUsageEndTime(id int64, endTime int64) error {
	return db.Model(&Usage{}).Where("id = ?", id).Update("end_time", endTime).Error
}

// DeleteUsage 删除使用记录
func DeleteUsage(id int64) error {
	return db.Delete(&Usage{}, id).Error
}

// GetOngoingUsage 获取机器当前正在进行的使用记录（end_time = 0）
func GetOngoingUsage(machineId int64) (*Usage, error) {
	var usage Usage
	err := db.Where("machine_id = ? AND end_time = 0", machineId).
		Order("start_time DESC").
		First(&usage).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &usage, nil
}

// CountUsagesByMachineID 统计指定机器的使用次数
func CountUsagesByMachineID(machineId int64) (int64, error) {
	var count int64
	err := db.Model(&Usage{}).Where("machine_id = ?", machineId).Count(&count).Error
	return count, err
}

// CountUsagesByTimeRange 统计指定时间范围的使用次数
func CountUsagesByTimeRange(startTime, endTime int64) (int64, error) {
	var count int64
	err := db.Model(&Usage{}).
		Where("start_time >= ? AND start_time <= ?", startTime, endTime).
		Count(&count).Error
	return count, err
}

// CountUsagesByMachineIDAndTimeRange 统计指定机器在指定时间范围的使用次数
func CountUsagesByMachineIDAndTimeRange(machineId, startTime, endTime int64) (int64, error) {
	var count int64
	err := db.Model(&Usage{}).
		Where("machine_id = ? AND start_time >= ? AND start_time <= ?", machineId, startTime, endTime).
		Count(&count).Error
	return count, err
}

// GetLatestUsageByMachineID 获取机器最新的使用记录
func GetLatestUsageByMachineID(machineId int64) (*Usage, error) {
	var usage Usage
	err := db.Where("machine_id = ?", machineId).
		Order("start_time DESC").
		First(&usage).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &usage, nil
}
