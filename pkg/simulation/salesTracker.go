package Simulation

import (
	"AI30_-_BlackFriday/pkg/constants"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

type SaleRecord struct {
	Tick        int
	Amount      float64
	TotalProfit float64
}

type SalesTracker struct {
	mapName      string
	simulationID string
	records      []SaleRecord
	mutex        sync.Mutex
	lastExported int // last export index
}

func NewSalesTracker(mapName string) *SalesTracker {
	tracker := &SalesTracker{
		mapName:      mapName,
		simulationID: generateSimulationID(),
		records:      make([]SaleRecord, 0),
	}

	// Auto-export based on config interval
	go func() {
		ticker := time.NewTicker(constants.SALES_EXPORT_INTERVAL)
		defer ticker.Stop()

		for range ticker.C {
			if len(tracker.records) > 0 {
				fmt.Printf("### AUTO-EXPORT: %d sells (interval: %v) ###\n", len(tracker.records), constants.SALES_EXPORT_INTERVAL)
				tracker.ExportToCSV()
			}
		}
	}()

	return tracker
}

func generateSimulationID() string {
	return fmt.Sprintf("sim_%d", time.Now().UnixNano())
}

func (st *SalesTracker) RecordSale(saleAmount float64, totalProfit float64, tick int) {
	st.mutex.Lock()
	defer st.mutex.Unlock()

	record := SaleRecord{
		Tick:        tick,
		Amount:      saleAmount,
		TotalProfit: totalProfit,
	}
	st.records = append(st.records, record)

	fmt.Printf("*** TRACKER: Vente enregistrée: %.2f€ au tick %d, Total: %.2f€ (records: %d) ***\n",
		saleAmount, tick, totalProfit, len(st.records))
}

func (st *SalesTracker) ExportToCSV() error {
	st.mutex.Lock()
	defer st.mutex.Unlock()

	// No new sells
	if st.lastExported >= len(st.records) {
		return nil
	}

	statsDir := "stats"
	if err := os.MkdirAll(statsDir, 0755); err != nil {
		return fmt.Errorf("error creating folkder stats: %v", err)
	}

	filename := filepath.Join(statsDir, "sales_tracker.csv")

	fileExists := true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		fileExists = false
	}

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error opening csv file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// write column name if file dont exist
	if !fileExists {
		header := []string{"simulation_id", "map_name", "tick", "montant_vente", "profit_total"}
		if err := writer.Write(header); err != nil {
			return fmt.Errorf("error writing columns: %v", err)
		}
	}

	// Only writing new records
	newRecords := st.records[st.lastExported:]
	for _, record := range newRecords {
		row := []string{
			st.simulationID,
			st.mapName,
			strconv.Itoa(record.Tick),
			strconv.FormatFloat(record.Amount, 'f', 2, 64),
			strconv.FormatFloat(record.TotalProfit, 'f', 2, 64),
		}
		if err := writer.Write(row); err != nil {
			return fmt.Errorf("error writing a log : %v", err)
		}
	}

	// new index
	st.lastExported = len(st.records)
	fmt.Printf("### EXPORT FINISHED: %d new sells data (total: %d) ###\n",
		len(newRecords), st.lastExported)
	return nil
}

func (st *SalesTracker) GetMapName() string {
	return st.mapName
}

func (st *SalesTracker) GetSimulationID() string {
	return st.simulationID
}

func (st *SalesTracker) GetRecordsCount() int {
	st.mutex.Lock()
	defer st.mutex.Unlock()
	return len(st.records)
}
