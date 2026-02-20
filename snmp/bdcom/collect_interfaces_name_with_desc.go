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

func CollectDeviceInterfaces(device models.Device) {

	// ✅ Run ONLY for BDCOM OLT
	if device.DeviceVendor != "BDCOM" || device.DeviceCategory != "OLT" {
		return
	}

	snmpClient := &gosnmp.GoSNMP{
		Target:    device.IPAddress,
		Port:      161,
		Community: device.SNMPCommunity,
		Version:   gosnmp.Version2c,
		Timeout:   3 * time.Second,
		Retries:   1,
		MaxOids:   gosnmp.MaxOids,
	}

	if err := snmpClient.Connect(); err != nil {
		log.Println("SNMP connect failed:", err)
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

	ifDescrOID := oids["ifDescr"]
	ifAliasOID := oids["ifAlias"]

	// -------- Walk ifDescr --------
	err := snmpClient.Walk(ifDescrOID, func(pdu gosnmp.SnmpPDU) error {

		// Extract index from OID
		indexStr := strings.TrimPrefix(pdu.Name, ifDescrOID+".")
		ifIndex, err := strconv.Atoi(indexStr)
		if err != nil {
			return nil
		}

		ifDescr := strings.TrimSpace(string(pdu.Value.([]byte)))
		if ifDescr == "" {
			return nil
		}

		// ifAlias (single GET, safe)
		ifAlias := snmpGetString(
			snmpClient,
			ifAliasOID+"."+indexStr,
		)

		// -------- Insert / Update DB --------
		repository.UpsertDeviceInterface(
			device.DeviceID,
			ifDescr,
			ifIndex,
			ifAlias,
		)

		return nil
	})

	if err != nil {
		log.Println("SNMP walk failed:", err)
	}
}
