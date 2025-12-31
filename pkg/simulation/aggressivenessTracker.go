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

type AggressivenessRecord struct {
	Tick           int
	Aggressiveness float64
}

type AggressivenessTracker struct {
	mapName      string
	simulationID string
	records      []AggressivenessRecord
	mutex        sync.Mutex
	lastExported int // last export index
}

func NewAggressivenessTracker(mapName string, id string) *AggressivenessTracker {
	tracker := &AggressivenessTracker{
		mapName:      mapName,
		simulationID: id,
		records:      make([]AggressivenessRecord, 0),
	}

	// Auto-export based on config interval
	go func() {
		ticker := time.NewTicker(constants.SALES_EXPORT_INTERVAL)
		defer ticker.Stop()

		for range ticker.C {
			if len(tracker.records) > 0 {
				fmt.Printf("### AUTO-EXPORT: %d agressivité (interval: %v) ###\n", len(tracker.records), constants.SALES_EXPORT_INTERVAL)
				tracker.ExportToCSV()
			}
		}
	}()

	return tracker
}

func (st *AggressivenessTracker) RecordAggressiveness(aggressiveness float64, tick int) {
	st.mutex.Lock()
	defer st.mutex.Unlock()

	record := AggressivenessRecord{
		Tick:           tick,
		Aggressiveness: aggressiveness,
	}
	st.records = append(st.records, record)

	fmt.Printf("*** TRACKER: agressivité enregistrée: %.2f au tick %d (records: %d) ***\n",
		aggressiveness, tick, len(st.records))
}

func (st *AggressivenessTracker) ExportToCSV() error {
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

	filename := filepath.Join(statsDir, "aggressiveness_tracker.csv")

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
		header := []string{"simulation_id", "map_name", "tick", "aggressiveness"}
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
			strconv.FormatFloat(record.Aggressiveness, 'f', 2, 64),
		}
		if err := writer.Write(row); err != nil {
			return fmt.Errorf("error writing a log : %v", err)
		}
	}

	// new index
	st.lastExported = len(st.records)
	fmt.Printf("### EXPORT FINISHED: %d new aggressiveness data (total: %d) ###\n",
		len(newRecords), st.lastExported)
	return nil
}

func (st *AggressivenessTracker) GetMapName() string {
	return st.mapName
}

func (st *AggressivenessTracker) GetSimulationID() string {
	return st.simulationID
}

func (st *AggressivenessTracker) GetRecordsCount() int {
	st.mutex.Lock()
	defer st.mutex.Unlock()
	return len(st.records)
}
