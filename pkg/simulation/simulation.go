package Simulation

import (
	"sync"
	"time"
)

type Simulation struct {
	NClients    int
	env Environment
	Speed       *float32
	agents      []Agent
	syncChans   sync.Map
}

func NewSimulation(agentCount int,width int, height int) (simu *Simulation) {
	simu = &Simulation{}
	simu.agents = make([]Agent, 0, agentCount)

	simu.env = *NewEnvironment(width,height)
	//voir quand initialiser les agents
	return simu
}

func (s *Simulation) Run() {
	s.env.Start()
	for i := 0; i < s.NClients; i++ {
		go s.env.Clients[i].Start()
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
