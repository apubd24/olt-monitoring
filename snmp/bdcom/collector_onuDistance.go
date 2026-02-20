package snmp_bdcom

import (
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/gosnmp/gosnmp"

	"snmp-onu-monitor/models"
	repository "snmp-onu-monitor/repository/bdcom"
)

func CollectDataOnuDistance(device models.Device) {
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

	onuDistanceOID := oids["onuDistance"]

	log.Println("📡 Walking ONU Distance OID:", onuDistanceOID)

	// -------- Walk ONU distance directly --------
	err := snmpClient.Walk(onuDistanceOID, func(pdu gosnmp.SnmpPDU) error {

		// Extract ifIndex
		parts := strings.Split(pdu.Name, ".")
		ifIndexStr := parts[len(parts)-1]

		ifIndex, err := strconv.Atoi(ifIndexStr)
		if err != nil {
			return nil
		}

		onuDistance := gosnmp.ToBigInt(pdu.Value).Int64()
		if onuDistance <= 0 {
			return nil
		}

		// ✅ Always INSERT (no update, no conflict)
		repository.InsertOnuDistance(
			device.DeviceID,
			ifIndex,
			int(onuDistance),
		)

		return nil
	})

	if err != nil {
		log.Println("SNMP walk failed:", err)
	}
}
