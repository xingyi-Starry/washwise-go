package model

const (
	MachineStatusAvailable = 0
	MachineStatusOffline   = 1
	MachineStatusInUse     = 2
)

type Machine struct {
	Id          int64 `gorm:"primaryKey"`
	Name        string
	Status      int
	LastUseTime int64
	Msg         string
	ShopId      string `gorm:"index"`
}
