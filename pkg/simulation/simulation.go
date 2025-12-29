package Simulation

import (
	Map "AI30_-_BlackFriday/pkg/map"
	"fmt"
	"sync"
	"time"
)

type Simulation struct {
	Env       Environment
	syncChans sync.Map
	tickCount int
}

func NewSimulation(agentCount int, ticDuration int, mapData *Map.Map, deltaTime float64, searchRadius float64, mapName string, shoppingListsPath string) (simu *Simulation) {

	simu = &Simulation{
		Env: *NewEnvironment(mapData, ticDuration, deltaTime, searchRadius, mapName, shoppingListsPath),
	}
	return simu
}

func (s *Simulation) Agents() []Agent {
	return s.Env.Agents
}
func (s *Simulation) AddClient(agtId string, aggressiveness float64) error {

	_, ok := s.syncChans.Load(AgentID(agtId))
	if ok {
		return fmt.Errorf("Agent with id %s was already loaded", agtId)
	}
	syncChan := make(chan int)
	s.syncChans.Store(AgentID(agtId), syncChan)
	s.Env.AddClient(agtId, aggressiveness, syncChan)
	return nil
}
func (s *Simulation) AddGuard(agtId string) error {
	_, ok := s.syncChans.Load(AgentID(agtId))
	if ok {
		return fmt.Errorf("Agent with id %s was already loaded", agtId)
	}

	syncChan := make(chan int)
	s.syncChans.Store(AgentID(agtId), syncChan)
	_, err := s.Env.AddGuard(agtId, syncChan)
	if err != nil {
		return err
	}
	return nil
}
func (s *Simulation) SetTicDuration(value int) {
	s.Env.ticDuration = value
}
func (s *Simulation) Run() {
	//avoid dependency from env to simulation
	s.Env.Start(func(agtId AgentID) {
		s.syncChans.Delete(agtId)
	})
	go func() {
		for {
			s.tickCount++
			s.Env.currentTick = s.tickCount
			time.Sleep(time.Duration(s.Env.ticDuration) * time.Millisecond)
		}
	}()
	for _, agt := range s.Env.Agents {
		go agt.Start()

		go func(agt Agent) {
			step := 0
			for {
				s.Env.pauseWg.Wait()
				step++
				c, ok := s.syncChans.Load(agt.ID())
				if !ok {
					fmt.Printf("No sync channel found for agent %s\n", agt.ID())
					return
				}
				c.(chan int) <- step
				time.Sleep(1 * time.Millisecond * time.Duration(s.Env.ticDuration))
				<-c.(chan int)
			}
		}(agt)
	}

}

func (s *Simulation) TogglePause() {
	s.Env.isPaused = !s.Env.isPaused
	if s.Env.isPaused {
		s.Env.pauseWg.Add(1)
	} else {
		s.Env.pauseWg.Done()
	}
}

func (s *Simulation) Stop() {
	// export ?
	s.Env.isStopped = true
	s.Env.cancel()
	s.Env.stopWg.Wait()

}
