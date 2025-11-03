package model

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
	ShopId      string `gorm:"index"`
}
