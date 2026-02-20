package repository

import (
	"log"

	"snmp-onu-monitor/config"
)

func parseOnuStatus(
	customerName, deviceName, vendor, deviceType, deviceIP, sysName string,
	ifDescr string,
	status int,
) {

	_, err := config.DB.Exec(`
		INSERT INTO onu_status (
			customer_name, device_name, device_vendor, device_type,
			device_ip, sys_name, if_descr,
			onu_status
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
	`,
		customerName,
		deviceName,
		vendor,
		deviceType,
		deviceIP,
		sysName,
		ifDescr,
		status,
	)

	if err != nil {
		log.Println("❌ Insert failed:", err)
	}
}
