package Simulation

import (
	Map "AI30_-_BlackFriday/pkg/map"
	"fmt"
	"sync"
	"time"
)

type Simulation struct {
	NClients  int
	Env       Environment
	Speed     *float32
	agents    []Agent
	syncChans sync.Map
}

func NewSimulation(agentCount int, mapData *Map.Map) (simu *Simulation) {

	simu = &Simulation{agents: make([]Agent, agentCount), Env: *NewEnvironment(mapData)}
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
	for i := range s.NClients {
		go s.Env.Clients[i].Start()
	}

	for _, agt := range s.agents {
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
				time.Sleep(1 * time.Millisecond * 100)
				<-c.(chan int)
			}
		}(agt)
	}

}
