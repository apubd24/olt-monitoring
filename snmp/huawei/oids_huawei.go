package snmp_huawei

var VendorOIDs = map[string]map[string]map[string]map[string]string{

	// ================= HUAWEI =================
	"HUAWEI": {

		// -------- OLT CATEGORY --------
		"OLT": {

			"GPON": {
				"ifDescr":     ".1.3.6.1.2.1.31.1.1.1.1",
				"ifAlias":     ".1.3.6.1.2.1.31.1.1.1.18",
				"onu_ifAlias": ".1.3.6.1.4.1.2011.6.128.1.1.2.43.1.9",
				"onuDistance": "1.3.6.1.4.1.2011.6.128.1.1.2.46.1.20",
				"onuRxPower":  ".1.3.6.1.4.1.2011.6.128.1.1.2.51.1.4", //ONT RX Power: (OLT->ONT)
				"onuTxPower":  ".1.3.6.1.4.1.2011.6.128.1.1.2.51.1.6", //OLT RX Power: (ONT<-OLT)
			},

			"EPON": {
				"ifDescr":     ".1.3.6.1.2.1.31.1.1.1.1",
				"ifAlias":     ".1.3.6.1.2.1.31.1.1.1.18",
				"onu_ifAlias": ".1.3.6.1.4.1.2011.6.128.1.1.2.53.1.9",
				"onuRxPower":  ".1.3.6.1.4.1.2011.6.128.1.1.2.104.1.5",
				"onuTxPower":  ".1.3.6.1.4.1.2011.6.128.1.1.2.51.1.6",
			},
		},
	},
}
