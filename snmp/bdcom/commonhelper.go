package snmp_bdcom

import (
	"fmt"
	"strings"

	"github.com/gosnmp/gosnmp"
)

func extractIndex(baseOID, fullOID string) string {
	return strings.TrimPrefix(fullOID, baseOID+".")
}

func parseMac(value interface{}) string {
	b, ok := value.([]byte)
	if !ok {
		return ""
	}

	mac := make([]string, 0)
	for _, v := range b {
		mac = append(mac, fmt.Sprintf("%02X", v))
	}
	return strings.Join(mac, ":")
}

func parseOnuStatus(value interface{}) int {
	switch v := value.(type) {
	case int:
		return v
	case uint:
		return int(v)
	case int64:
		return int(v)
	default:
		return 0
	}
}

func getSysName(snmp *gosnmp.GoSNMP) string {
	pkt, err := snmp.Get([]string{".1.3.6.1.2.1.1.5.0"})
	if err != nil || len(pkt.Variables) == 0 {
		return ""
	}
	return fmt.Sprintf("%s", pkt.Variables[0].Value)
}
