package repository

import (
	"log"

	"snmp-onu-monitor/config"
)

func InsertOnuStatus(
	deviceID int64,
	ifIndex int,
	onuStatus int,
) {

	_, err := config.DB.Exec(`
		INSERT INTO olt_onu_status (
			device_id,
			if_index,
			onu_status
		)
		VALUES ($1, $2, $3)
	`,
		deviceID,
		ifIndex,
		onuStatus,
	)

	if err != nil {
		log.Println("❌ Insert failed:", err)
	}
}
