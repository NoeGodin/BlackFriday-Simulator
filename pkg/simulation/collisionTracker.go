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

type CollisionRecord struct {
	Tick            int
	CollisionCount  int
	TotalCollisions int
}

type CollisionTracker struct {
	mapName         string
	simulationID    string
	records         []CollisionRecord
	mutex           sync.Mutex
	lastExported    int // last export index
	totalCollisions int // cumulative collision count
}

func NewCollisionTracker(mapName string) *CollisionTracker {
	tracker := &CollisionTracker{
		mapName:         mapName,
		simulationID:    generateSimulationID(),
		records:         make([]CollisionRecord, 0),
		totalCollisions: 0,
	}

	// Auto-export based on config interval
	go func() {
		ticker := time.NewTicker(constants.SALES_EXPORT_INTERVAL)
		defer ticker.Stop()

		for range ticker.C {
			if len(tracker.records) > 0 {
				fmt.Printf("### AUTO-EXPORT: %d collision records (interval: %v) ###\n", len(tracker.records), constants.SALES_EXPORT_INTERVAL)
				tracker.ExportToCSV()
			}
		}
	}()

	return tracker
}

func (ct *CollisionTracker) RecordCollision(collisionCount int, tick int) {
	ct.mutex.Lock()
	defer ct.mutex.Unlock()

	ct.totalCollisions += collisionCount

	record := CollisionRecord{
		Tick:            tick,
		CollisionCount:  collisionCount,
		TotalCollisions: ct.totalCollisions,
	}
	ct.records = append(ct.records, record)

	fmt.Printf("*** COLLISION TRACKER: %d collisions enregistrées au tick %d, Total: %d (records: %d) ***\n",
		collisionCount, tick, ct.totalCollisions, len(ct.records))
}

func (ct *CollisionTracker) ExportToCSV() error {
	ct.mutex.Lock()
	defer ct.mutex.Unlock()

	// No new collisions
	if ct.lastExported >= len(ct.records) {
		return nil
	}

	statsDir := "stats"
	if err := os.MkdirAll(statsDir, 0755); err != nil {
		return fmt.Errorf("error creating folder stats: %v", err)
	}

	filename := filepath.Join(statsDir, "collision_tracker.csv")

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
		header := []string{"simulation_id", "map_name", "tick", "collision_count", "total_collisions"}
		if err := writer.Write(header); err != nil {
			return fmt.Errorf("error writing columns: %v", err)
		}
	}

	// Only writing new records
	newRecords := ct.records[ct.lastExported:]
	for _, record := range newRecords {
		row := []string{
			ct.simulationID,
			ct.mapName,
			strconv.Itoa(record.Tick),
			strconv.Itoa(record.CollisionCount),
			strconv.Itoa(record.TotalCollisions),
		}
		if err := writer.Write(row); err != nil {
			return fmt.Errorf("error writing a log: %v", err)
		}
	}

	// new index
	ct.lastExported = len(ct.records)
	fmt.Printf("### EXPORT FINISHED: %d new collision records (total: %d) ###\n",
		len(newRecords), ct.lastExported)
	return nil
}

func (ct *CollisionTracker) GetMapName() string {
	return ct.mapName
}

func (ct *CollisionTracker) GetSimulationID() string {
	return ct.simulationID
}

func (ct *CollisionTracker) GetRecordsCount() int {
	ct.mutex.Lock()
	defer ct.mutex.Unlock()
	return len(ct.records)
}

func (ct *CollisionTracker) GetTotalCollisions() int {
	ct.mutex.Lock()
	defer ct.mutex.Unlock()
	return ct.totalCollisions
}
