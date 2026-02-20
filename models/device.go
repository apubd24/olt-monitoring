package models

type Device struct {
	DeviceID       int64
	CustomerID     int
	CustomerName   string
	DeviceName     string
	DeviceVendor   string
	DeviceCategory string
	DeviceType     string
	IPAddress      string
	SNMPCommunity  string
	SNMPVersion    string
	IsActive       bool
}
