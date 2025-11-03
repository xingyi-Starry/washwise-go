package model

const (
	MachineStatusAvailable = 0
	MachineStatusOffline   = 1
	MachineStatusInUse     = 2
)

type Machine struct {
	Id             int64 `gorm:"primaryKey"`
	Name           string
	Status         int
	LastUpdateTime int64
	RemainTime     int64
	ShopId         int64
}
