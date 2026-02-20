package snmp_huawei

import (
	"log"
	"time"

	"github.com/gosnmp/gosnmp"

	"snmp-onu-monitor/models"
	repository "snmp-onu-monitor/repository/huawei"
)

func CollectDataOnuTxRxPower(device models.Device) {

	// ✅ Huawei OLT only
	if device.DeviceVendor != "HUAWEI" || device.DeviceCategory != "OLT" {
		return
	}

	snmpClient := &gosnmp.GoSNMP{
		Target:    device.IPAddress,
		Port:      161,
		Community: device.SNMPCommunity,
		Version:   gosnmp.Version2c,
		Timeout:   5 * time.Second,
		Retries:   1,
		MaxOids:   gosnmp.MaxOids,
	}

	if err := snmpClient.Connect(); err != nil {
		log.Println("❌ SNMP connect failed:", err)
		return
	}
	defer snmpClient.Conn.Close()

	// -------- Load OIDs --------
	oids := VendorOIDs[device.DeviceVendor][device.DeviceCategory][device.DeviceType]

	onuRxOID := oids["onuRxPower"]
	onuTxOID := oids["onuTxPower"]

	log.Println("📡 Huawei ONU TX/RX Power collection started")

	// Walk RX power (RX always exists)
	err := snmpClient.Walk(onuRxOID, func(pdu gosnmp.SnmpPDU) error {

		ifIndex := HuaweiExtractTxOnuIfIndex(pdu.Name)
		if ifIndex == "" {
			return nil
		}

		rxPower := HuaweiDecodeRxPower(pdu)

		// TX Power (GET)
		txResp, err := snmpClient.Get([]string{onuTxOID + "." + ifIndex})
		if err != nil || len(txResp.Variables) == 0 {
			return nil
		}

		txPower := HuaweiDecodeTxPower(txResp.Variables[0])

		// Skip invalid ONUs
		if rxPower == 0 && txPower == 0 {
			return nil
		}

		repository.InsertOnuTxRxPower(
			device.DeviceID,
			ifIndex, // 4194312192.3
			rxPower, // float dBm
			txPower, // float dBm
		)

		return nil
	})

	if err != nil {
		log.Println("❌ SNMP walk failed:", err)
	}
}
