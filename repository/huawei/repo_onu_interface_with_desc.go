package repository

import (
	"log"
	"snmp-onu-monitor/config"
)

func UpsertonuInterface(
	deviceID int64,
	ifDescr string,
	ifIndex string,
	ifAlias string,
) {

	query := `
	INSERT INTO device_interface_descs
	(device_id, if_descr, if_index, if_alias)
	VALUES ($1, $2, $3, $4)
	ON CONFLICT (device_id, if_index)
	DO UPDATE SET
		if_descr = EXCLUDED.if_descr,
		if_alias = EXCLUDED.if_alias,
		updated_at = now();
	`

	_, err := config.DB.Exec(query,
		deviceID,
		ifDescr,
		ifIndex,
		ifAlias,
	)
	if err != nil {
		log.Println("❌ Insert failed:", err)
	}
}

func GetPonIfDescr(deviceID int64, ifIndex string) (string, error) {

	var ifDescr string

	query := `
	SELECT if_descr
	FROM device_interface_descs
	WHERE device_id = $1
	  AND if_index = $2
	LIMIT 1
	`

	err := config.DB.QueryRow(query, deviceID, ifIndex).Scan(&ifDescr)
	return ifDescr, err
}
