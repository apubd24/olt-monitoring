package repository

import (
	"log"

	"snmp-onu-monitor/config"
	"snmp-onu-monitor/models"
)

func GetActiveDevices() []models.Device {

	rows, err := config.DB.Query(`
		SELECT device_id, customer_id, customer_name,
		       device_name, device_vendor, device_category, device_type,
		       ip_address, snmp_community, snmp_version, is_active
		FROM devices
		WHERE is_active = true
	`)
	if err != nil {
		log.Fatal(err)
	}

	var devices []models.Device

	for rows.Next() {
		var d models.Device
		rows.Scan(
			&d.DeviceID,
			&d.CustomerID,
			&d.CustomerName,
			&d.DeviceName,
			&d.DeviceVendor,
			&d.DeviceCategory,
			&d.DeviceType,
			&d.IPAddress,
			&d.SNMPCommunity,
			&d.SNMPVersion,
			&d.IsActive,
		)
		devices = append(devices, d)
	}

	return devices
}
