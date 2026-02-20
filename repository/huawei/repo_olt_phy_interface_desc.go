package repository

import (
	"log"
	"snmp-onu-monitor/config"
)

func UpsertDeviceInterface(
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

	_, err := config.DB.Exec(query, deviceID, ifDescr, ifIndex, ifAlias)
	if err != nil {
		log.Println("❌ Insert failed:", err)
	}
}
