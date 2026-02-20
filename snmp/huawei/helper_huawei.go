package snmp_huawei

import (
	"strings"

	"github.com/gosnmp/gosnmp"
)

// onu Interface Description
func snmpGetString(snmp *gosnmp.GoSNMP, oid string) string {

	result, err := snmp.Get([]string{oid})
	if err != nil || len(result.Variables) == 0 {
		return ""
	}

	if b, ok := result.Variables[0].Value.([]byte); ok {
		return string(b)
	}

	return ""
}

// onuDistance

func HuaweieXtractOnuIfIndex(fullOID string) string {
	parts := strings.Split(fullOID, ".")
	if len(parts) < 2 {
		return ""
	}

	return parts[len(parts)-2] + "." + parts[len(parts)-1]
}

func HuaweiExtractTxOnuIfIndex(fullOID string) string {
	parts := strings.Split(fullOID, ".")
	if len(parts) < 2 {
		return ""
	}
	return parts[len(parts)-2] + "." + parts[len(parts)-1]
}

// fHuawei returns INTEGER in 0.01 dBm

func HuaweiDecodeRxPower(pdu gosnmp.SnmpPDU) float64 {
	val := gosnmp.ToBigInt(pdu.Value).Int64()

	if val == 0 || val == 32767 || val == 2147483647 {
		return 0
	}
	return float64(val) / 100.0
}

/*
Huawei TX power decoder (GET result)
*/
func HuaweiDecodeTxPower(pdu gosnmp.SnmpPDU) float64 {
	val := gosnmp.ToBigInt(pdu.Value).Int64()

	if val == 0 || val == 32767 || val == 2147483647 {
		return 0
	}
	return float64(val) / 100.0
}
