package repository

import (
	"log"

	"snmp-onu-monitor/config"
)

func InsertOnuDistance(
	deviceID int64,
	ifIndex int,
	onuDistance int,
) {

	_, err := config.DB.Exec(`
		INSERT INTO olt_onu_distance (
			device_id,
			if_index,
			onu_distance
		)
		VALUES ($1, $2, $3)
	`,
		deviceID,
		ifIndex,
		onuDistance,
	)

	if err != nil {
		log.Println("❌ Insert failed:", err)
	}
}
