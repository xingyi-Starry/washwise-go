package model

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

// CountUsagesByMachineIDAndTimeRange 统计指定机器在指定时间范围的使用次数
func CountUsagesByMachineIDAndTimeRange(machineId, startTime, endTime int64) (int64, error) {
	var count int64
	err := db.Model(&Usage{}).
		Where("machine_id = ? AND start_time >= ? AND start_time <= ?", machineId, startTime, endTime).
		Count(&count).Error
	return count, err
}
