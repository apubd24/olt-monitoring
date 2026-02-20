package repository

import (
	"log"

	"snmp-onu-monitor/config"
)

func InsertOnuTxRxPower(
	deviceID int64,
	ifIndex string,
	RxPower float64,
	TxPower float64,
) {

	_, err := config.DB.Exec(`
		INSERT INTO huawei_onu_tx_rx_powers (
			device_id, onu_ifindex, onuRx_power, onuTx_power
		)
		VALUES ($1,$2,$3,$4)
	`,
		deviceID,
		ifIndex,
		RxPower,
		TxPower,
	)

	if err != nil {
		log.Println("❌ Insert failed:", err)
	}
}
