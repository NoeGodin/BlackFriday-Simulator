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

type StealRecord struct {
	Tick      int
	StealerID string
	VictimID  string
	ItemName  string
}

type StealTracker struct {
	mapName      string
	simulationID string
	records      []StealRecord
	mutex        sync.Mutex
	lastExported int
}

func NewStealTracker(mapName string, id string) *StealTracker {
	tracker := &StealTracker{
		mapName:      mapName,
		simulationID: id,
		records:      make([]StealRecord, 0),
	}

	go func() {
		ticker := time.NewTicker(constants.SALES_EXPORT_INTERVAL)
		defer ticker.Stop()

		for range ticker.C {
			if len(tracker.records) > 0 {
				fmt.Printf("### AUTO-EXPORT: %d vols (interval: %v) ###\n", len(tracker.records), constants.SALES_EXPORT_INTERVAL)
				tracker.ExportToCSV()
			}
		}
	}()

	return tracker
}

func (st *StealTracker) RecordSteal(stealerID, victimID, itemName string, tick int) {
	st.mutex.Lock()
	defer st.mutex.Unlock()

	record := StealRecord{
		Tick:      tick,
		StealerID: stealerID,
		VictimID:  victimID,
		ItemName:  itemName,
	}
	st.records = append(st.records, record)

	fmt.Printf("*** TRACKER: vol enregistré: %s a volé %s à %s au tick %d (records: %d) ***\n",
		stealerID, itemName, victimID, tick, len(st.records))
}

func (st *StealTracker) ExportToCSV() error {
	st.mutex.Lock()
	defer st.mutex.Unlock()

	if st.lastExported >= len(st.records) {
		return nil
	}

	statsDir := "stats"
	if err := os.MkdirAll(statsDir, 0755); err != nil {
		return fmt.Errorf("error creating folder stats: %v", err)
	}

	filename := filepath.Join(statsDir, "steal_tracker.csv")

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

	if !fileExists {
		header := []string{"simulation_id", "map_name", "tick", "stealer_id", "victim_id", "item_name"}
		if err := writer.Write(header); err != nil {
			return fmt.Errorf("error writing columns: %v", err)
		}
	}

	newRecords := st.records[st.lastExported:]
	for _, record := range newRecords {
		row := []string{
			st.simulationID,
			st.mapName,
			strconv.Itoa(record.Tick),
			record.StealerID,
			record.VictimID,
			record.ItemName,
		}
		if err := writer.Write(row); err != nil {
			return fmt.Errorf("error writing a log : %v", err)
		}
	}

	st.lastExported = len(st.records)
	fmt.Printf("### EXPORT FINISHED: %d new steal data (total: %d) ###\n",
		len(newRecords), st.lastExported)
	return nil
}

func (st *StealTracker) GetMapName() string {
	return st.mapName
}

func (st *StealTracker) GetSimulationID() string {
	return st.simulationID
}

func (st *StealTracker) GetRecordsCount() int {
	st.mutex.Lock()
	defer st.mutex.Unlock()
	return len(st.records)
}
