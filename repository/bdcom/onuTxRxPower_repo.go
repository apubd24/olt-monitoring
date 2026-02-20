package repository

import (
	"log"

	"snmp-onu-monitor/config"
)

func InsertOnuTxRxPower(
	deviceID int64,
	onuIndex int,
	onuRxPower float64,
	onuTxPower float64,
) {

	_, err := config.DB.Exec(`
		INSERT INTO onu_tx_rx_powers (
			device_id, onu_ifindex, onuRx_power, onuTx_power
		)
		VALUES ($1,$2,$3,$4)
	`,
		deviceID,
		onuIndex,
		onuRxPower,
		onuTxPower,
	)

	if err != nil {
		log.Println("❌ Insert failed:", err)
	}
}
