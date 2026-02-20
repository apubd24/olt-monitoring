package snmp_bdcom

import (
	"fmt"
	"log"
	"time"

	"github.com/gosnmp/gosnmp"

	"snmp-onu-monitor/models"
)

const (
	ifDescrOID   = ".1.3.6.1.2.1.2.2.1.2"
	onuMacOID    = ".1.3.6.1.4.1.3320.101.10.1.1.76"
	onuStatusOID = ".1.3.6.1.4.1.3320.101.11.1.1.6"
)

func CollectOnuStatus(device models.Device) {

	// ✅ RUN ONLY FOR BDCOM EPON
	if device.DeviceVendor != "BDCOM" || device.DeviceType != "EPON" {
		log.Printf(
			"[SKIP] Device=%s Vendor=%s Type=%s (Not BDCOM EPON)\n",
			device.DeviceName,
			device.DeviceVendor,
			device.DeviceType,
		)
		return
	}

	log.Printf(
		"[START] Collecting ONU status for BDCOM EPON device: %s (%s)\n",
		device.DeviceName,
		device.IPAddress,
	)

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

	sysName := getSysName(snmp)
	log.Println("Connected to:", sysName)

	// -----------------------------
	// STEP 1: ifIndex → ifDescr
	// -----------------------------
	ifDescrByIndex := make(map[string]string)

	_ = snmp.Walk(ifDescrOID, func(pdu gosnmp.SnmpPDU) error {
		index := extractIndex(ifDescrOID, pdu.Name)
		ifDescrByIndex[index] = fmt.Sprintf("%s", pdu.Value)

		log.Printf("[IF-DESCR] index=%s value=%s\n", index, ifDescrByIndex[index])
		return nil
	})

	log.Printf("Total ifDescr mapped: %d\n", len(ifDescrByIndex))

	// -----------------------------
	// STEP 2: ONU Index → MAC
	// -----------------------------
	onuIndexToMac := make(map[string]string)

	_ = snmp.Walk(onuMacOID, func(pdu gosnmp.SnmpPDU) error {
		index := extractIndex(onuMacOID, pdu.Name)
		mac := parseMac(pdu.Value)
		onuIndexToMac[index] = mac

		log.Printf("[ONU-MAC] index=%s mac=%s\n", index, mac)
		return nil
	})

	log.Printf("Total ONU MAC mapped: %d\n", len(onuIndexToMac))

	// -----------------------------
	// STEP 3: ONU STATUS (FINAL)
	// -----------------------------
	log.Println("========== ONU STATUS (BDCOM EPON) ==========")

	err := snmp.Walk(onuStatusOID, func(pdu gosnmp.SnmpPDU) error {

		index := extractIndex(onuStatusOID, pdu.Name)

		mac, ok := onuIndexToMac[index]
		if !ok {
			log.Printf("[SKIP] index=%s MAC not found\n", index)
			return nil
		}

		ifDescr, ok := ifDescrByIndex[index]
		if !ok {
			log.Printf("[SKIP] index=%s ifDescr not found\n", index)
			return nil
		}

		status := parseOnuStatus(pdu.Value)

		log.Printf(
			"[ONU]\n"+
				" Device   : %s\n"+
				" Index    : %s\n"+
				" ifDescr  : %s\n"+
				" MAC      : %s\n"+
				" Status   : %d\n",
			device.DeviceName,
			index,
			ifDescr,
			mac,
			status,
		)

		return nil
	})

	if err != nil {
		log.Println("ONU status walk failed:", err)
	}
}
