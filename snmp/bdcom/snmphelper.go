package snmp_bdcom

import (
	"fmt"

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

// snmpGetInt - Get integer value from SNMP OID
func snmpGetInt(client *gosnmp.GoSNMP, oid string) int {
	result, err := client.Get([]string{oid})
	if err != nil {
		return 0
	}

	if len(result.Variables) == 0 {
		return 0
	}

	switch v := result.Variables[0].Value.(type) {
	case int:
		return v
	case int64:
		return int(v)
	case uint:
		return int(v)
	case uint64:
		return int(v)
	case int32:
		return int(v)
	case uint32:
		return int(v)
	default:
		return 0
	}
}

func snmpToString(v interface{}) string {
	switch t := v.(type) {
	case string:
		return t
	case []byte:
		return string(t)
	default:
		return fmt.Sprintf("%v", t)
	}
}

func parseSNMPInt(v interface{}) int {
	switch val := v.(type) {
	case int:
		return val
	case int64:
		return int(val)
	case uint:
		return int(val)
	case uint64:
		return int(val)
	default:
		return 0
	}
}

// ifDescr
func getIfDescr(snmp *gosnmp.GoSNMP, index string) string {
	oid := CommonOIDs["ifDescr"] + "." + index
	result, err := snmp.Get([]string{oid})
	if err != nil || len(result.Variables) == 0 {
		return "N/A"
	}
	return snmpToString(result.Variables[0].Value)
}

func snmpValueToString(value interface{}) string {
	switch v := value.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	case int, int64, uint, uint64:
		return fmt.Sprintf("%d", v)
	case nil:
		return ""
	default:
		return fmt.Sprintf("%v", v)
	}
}

// For the Rx Power
func RxSnmpGetFloatdBm(snmp *gosnmp.GoSNMP, oid string) float64 {
	result, err := snmp.Get([]string{oid})
	if err != nil || len(result.Variables) == 0 {
		return 0
	}

	switch v := result.Variables[0].Value.(type) {
	case int:
		return float64(v) / 10
	case int64:
		return float64(v) / 10
	case uint:
		return float64(v) / 10
	case uint64:
		return float64(v) / 10
	default:
		return 0
	}
}

// For the Tx Power
func TxSnmpGetFloatdBm(snmp *gosnmp.GoSNMP, oid string) float64 {
	result, err := snmp.Get([]string{oid})
	if err != nil || len(result.Variables) == 0 {
		return 0
	}

	switch v := result.Variables[0].Value.(type) {
	case int:
		return float64(v) / 10
	case int64:
		return float64(v) / 10
	case uint:
		return float64(v) / 10
	case uint64:
		return float64(v) / 10
	default:
		return 0
	}
}

func snmpValueToFloat(val interface{}) float64 {
	switch v := val.(type) {
	case int:
		return float64(v)
	case int64:
		return float64(v)
	case uint:
		return float64(v)
	case uint64:
		return float64(v)
	default:
		return 0
	}
}
