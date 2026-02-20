package main

import (
	"log"
	"sync"
	"time"

	"snmp-onu-monitor/config"
	"snmp-onu-monitor/models"
	repository "snmp-onu-monitor/repository/bdcom"
	snmp_huawei "snmp-onu-monitor/snmp/huawei"
)

const (
	// Polling intervals
	POLLING_INTERVAL_HOUR = 6 * time.Hour
	POWER_INTERVAL        = 5 * time.Minute

	// Max concurrent OLT polling
	MAX_WORKERS = 5
)

func main() {
	log.Println("ONU Monitoring Collector Started")

	config.InitDB()
	defer config.DB.Close()

	go hourScheduler()
	go powerScheduler()

	// Block forever
	select {}
}

// ===============================
// DISTANCE SCHEDULER (6 hours)
// ===============================
func hourScheduler() {
	ticker := time.NewTicker(POLLING_INTERVAL_HOUR)
	defer ticker.Stop()

	for {
		log.Println("[DISTANCE] Polling started")
		runDistanceJob()
		log.Println("[DISTANCE] Polling finished")

		<-ticker.C
	}
}

func runDistanceJob() {
	devices := repository.GetActiveDevices()

	sem := make(chan struct{}, MAX_WORKERS)
	var wg sync.WaitGroup

	for _, d := range devices {
		wg.Add(1)
		sem <- struct{}{}

		go func(device models.Device) {
			defer wg.Done()
			defer func() { <-sem }()

			// snmp_bdcom.CollectDataOnuDistance(device)
			snmp_huawei.CollectDataOnuDistance(device)

		}(d)
	}

	wg.Wait()
}

// ===============================
// RX/TX POWER SCHEDULER (10 min)
// ===============================
func powerScheduler() {
	ticker := time.NewTicker(POWER_INTERVAL)
	defer ticker.Stop()

	for {
		log.Println("[POWER] Polling started")
		runPowerJob()
		log.Println("[POWER] Polling finished")

		<-ticker.C
	}
}

func runPowerJob() {
	devices := repository.GetActiveDevices()

	sem := make(chan struct{}, MAX_WORKERS)
	var wg sync.WaitGroup

	for _, d := range devices {
		wg.Add(1)
		sem <- struct{}{}

		go func(device models.Device) {
			defer wg.Done()
			defer func() { <-sem }()

			// snmp.CollectOnuStatus(device)
			// snmp_bdcom.CollectDataOnuStatus(device)
			// snmp_bdcom.CollectDeviceInterfaces(device)
			// snmp_huawei.CollectDeviceInterfaces(device)
			// snmp_huawei.CollectHuaweiONUInterfaces(device)
			// snmp_bdcom.CollectDataOnuTxRxPower(device)
			snmp_huawei.CollectDataOnuTxRxPower(device)
		}(d)
	}

	wg.Wait()
}
