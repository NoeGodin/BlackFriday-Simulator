package Simulation

import (
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

func NewSimulation(agentCount int, width int, height int) (simu *Simulation) {
	simu = &Simulation{}
	simu.agents = make([]Agent, 0, agentCount)

	simu.Env = *NewEnvironment(width, height)
	//voir quand initialiser les agents
	return simu
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
				c, _ := s.syncChans.Load(agt.ID())
				c.(chan int) <- step
				time.Sleep(1 * time.Millisecond)
				<-c.(chan int)
			}
		}(agt)
	}

}
