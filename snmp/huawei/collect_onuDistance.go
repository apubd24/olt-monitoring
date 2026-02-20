package snmp_huawei

import (
	"log"
	"time"

	"github.com/gosnmp/gosnmp"

	"snmp-onu-monitor/models"
	repository "snmp-onu-monitor/repository/huawei"
)

func CollectDataOnuDistance(device models.Device) {

	// ✅ Run ONLY for HUAWEI OLT
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

	// -------- Load Vendor OIDs --------
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

	err := snmpClient.Walk(onuDistanceOID, func(pdu gosnmp.SnmpPDU) error {

		// ✅ FIXED if_index extraction
		ifIndex := HuaweiExtractTxOnuIfIndex(pdu.Name)
		if ifIndex == "" {
			return nil
		}

		onuDistance := gosnmp.ToBigInt(pdu.Value).Int64()
		if onuDistance <= 0 {
			return nil
		}

		repository.InsertOnuDistance(
			device.DeviceID,
			ifIndex, // 4194312192.3
			int(onuDistance),
		)

		return nil
	})

	if err != nil {
		log.Println("❌ SNMP walk failed:", err)
	}
}
