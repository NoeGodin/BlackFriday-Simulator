package Simulation

import (
	"AI30_-_BlackFriday/pkg/logger"
	Map "AI30_-_BlackFriday/pkg/map"
	"fmt"
	"sync"
	"time"
)

type Simulation struct {
	NClients  int
	Env       Environment
	Speed     float64
	agents    []Agent
	syncChans sync.Map
}

func NewSimulation(agentCount int, speed float64, mapData *Map.Map, deltaTime float64, searchRadius float64, mapName string) (simu *Simulation) {

	simu = &Simulation{agents: make([]Agent, agentCount), Env: *NewEnvironment(mapData, deltaTime, searchRadius, mapName), Speed: speed}
	return simu
}

func (s *Simulation) Agents() []Agent {
	return s.agents
}
func (s *Simulation) AddClient(agtId string) error {

	_, ok := s.syncChans.Load(AgentID(agtId))
	if ok {
		return fmt.Errorf("Agent with id %s was already loaded", agtId)
	}

	syncChan := make(chan int)
	s.syncChans.Store(AgentID(agtId), syncChan)
	agt := s.Env.AddClient(agtId, syncChan)
	s.agents = append(s.agents, agt)
	s.NClients++
	return nil
}
func (s *Simulation) Run() {
	s.Env.Start()
	go s.Env.exitRequest(s)

	for _, agt := range s.agents {
		go agt.Start()

		go func(agt Agent) {
			step := 0
			for {
				step++
				c, ok := s.syncChans.Load(agt.ID())
				if !ok {
					fmt.Printf("No sync channel found for agent %s\n", agt.ID())
					return
				}
				c.(chan int) <- step
				time.Sleep(1 * time.Millisecond * time.Duration(s.Speed))
				<-c.(chan int)
			}
		}(agt)
	}

}

func (simu *Simulation) RemoveAgent(agentID AgentID) {
	newAgents := simu.agents[:0]
	for _, a := range simu.agents {
		if a.ID() != agentID {
			newAgents = append(newAgents, a)
		}
	}
	simu.agents = newAgents
	simu.Env.Clients = removeAgentFromClients(agentID, simu.Env.Clients)
	simu.syncChans.Delete(agentID)
}

func (env *Environment) exitRequest(simu *Simulation) {
	for exitRequest := range env.exitChan {
		simu.RemoveAgent(exitRequest.Agt.ID())
		exitRequest.ResponseChannel <- true

		// Export data if nore more agents
		if len(simu.agents) == 0 {
			if err := env.ExportSalesData(); err != nil {
				logger.Errorf("Error during sells data export: %v", err)
			}
		}
	}
}
