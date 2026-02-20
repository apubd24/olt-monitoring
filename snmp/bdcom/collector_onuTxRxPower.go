package snmp_bdcom

import (
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/gosnmp/gosnmp"

	"snmp-onu-monitor/models"
	repository "snmp-onu-monitor/repository/bdcom"
	// repository "snmp-onu-monitor/repository/bdcom/"
)

func CollectDataOnuTxRxPower(device models.Device) {

	// ✅ Run ONLY for BDCOM OLT
	if device.DeviceVendor != "BDCOM" || device.DeviceCategory != "OLT" {
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
		log.Println("SNMP connect failed:", err)
		return
	}
	defer snmpClient.Conn.Close()

	// -------- Vendor OIDs --------
	vendorMap, ok := VendorOIDs[device.DeviceVendor]
	if !ok {
		log.Println("Unknown vendor")
		return
	}

	categoryMap, ok := vendorMap[device.DeviceCategory]
	if !ok {
		log.Println("Unknown category")
		return
	}

	oids, ok := categoryMap[device.DeviceType]
	if !ok {
		log.Println("Unknown device type")
		return
	}

	onuRxOID := oids["onuRxPower"]
	onuTxOID := oids["onuTxPower"]

	log.Println("📡 Walking ONU RX Power OID:", onuRxOID)

	// ===============================
	// WALK RX POWER (PRIMARY INDEX SOURCE)
	// ===============================
	err := snmpClient.Walk(onuRxOID, func(pdu gosnmp.SnmpPDU) error {

		// Extract ifIndex from OID
		parts := strings.Split(pdu.Name, ".")
		ifIndexStr := parts[len(parts)-1]

		ifIndex, err := strconv.Atoi(ifIndexStr)
		if err != nil {
			return nil
		}

		// RX Power
		rxPower := RxSnmpGetFloatdBm(snmpClient, onuRxOID+"."+ifIndexStr)

		// TX Power (single GET, safe)
		txPower := TxSnmpGetFloatdBm(snmpClient, onuTxOID+"."+ifIndexStr)

		// Skip invalid ONU
		if rxPower == 0 && txPower == 0 {
			return nil
		}

		// INSERT (historical data)
		repository.InsertOnuTxRxPower(
			device.DeviceID,
			ifIndex,
			rxPower,
			txPower,
		)

		return nil
	})

	if err != nil {
		log.Println("SNMP walk failed:", err)
	}
}
