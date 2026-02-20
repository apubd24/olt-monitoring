package snmp_huawei

import (
	"log"
	"strings"
	"time"

	"github.com/gosnmp/gosnmp"

	"snmp-onu-monitor/models"
	repository "snmp-onu-monitor/repository/huawei"
)

func CollectHuaweiONUInterfaces(device models.Device) {

	// ✅ Only Huawei OLT
	if device.DeviceVendor != "HUAWEI" || device.DeviceCategory != "OLT" {
		return
	}

	snmp := &gosnmp.GoSNMP{
		Target:    device.IPAddress,
		Port:      161,
		Community: device.SNMPCommunity,
		Version:   gosnmp.Version2c,
		Timeout:   5 * time.Second,
		Retries:   1,
		MaxOids:   gosnmp.MaxOids,
	}

	if err := snmp.Connect(); err != nil {
		log.Println("SNMP connect failed:", err)
		return
	}
	defer snmp.Conn.Close()

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

	onuAliasOID := oids["onu_ifAlias"]

	log.Println("📡 Collecting Huawei ONU descriptions:", device.IPAddress)

	err := snmp.Walk(onuAliasOID, func(pdu gosnmp.SnmpPDU) error {

		// Example OID:
		// .1.3.6.1...9.4194312192.3
		indexPart := strings.TrimPrefix(pdu.Name, onuAliasOID+".")

		parts := strings.Split(indexPart, ".")
		if len(parts) != 2 {
			return nil
		}

		ponIfIndex := parts[0] // 4194312192
		onuNo := parts[1]      // 3

		ifAlias := strings.TrimSpace(string(pdu.Value.([]byte)))
		if ifAlias == "" || ifAlias == "ONT_NO_DESCRIPTION" {
			return nil
		}

		// 🔎 Find physical PON interface name
		ponIfDescr, err := repository.GetPonIfDescr(
			device.DeviceID,
			ponIfIndex,
		)
		if err != nil || ponIfDescr == "" {
			return nil
		}

		// Build ONU interface name
		// GPON 0/1/0:3
		onuIfDescr := ponIfDescr + ":" + onuNo
		onuIfIndex := ponIfIndex + "." + onuNo

		repository.UpsertonuInterface(
			device.DeviceID,
			onuIfDescr,
			onuIfIndex,
			ifAlias,
		)

		return nil
	})

	if err != nil {
		log.Println("SNMP walk failed:", err)
	}
}
