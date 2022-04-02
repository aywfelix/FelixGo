package shell

import (
	"os"

	. "github.com/aywfelix/felixgo/container/map"
)

type ProcessManager struct {
	processMap *IntAnyMap
}

func NewProcessManager() *ProcessManager {
	return &ProcessManager{
		processMap: NewIntAnyMap(),
	}
}

func (m *ProcessManager) NewProcess(path string, args []string, environment []string) *Process {
	process := NewProcess(path, args, environment)
	process.Manager = m
	return process
}

func (m *ProcessManager) Get(pid int) *Process {
	if v := m.processMap.Get(pid); v != nil {
		process, _ := v.(*Process)
		return process
	}
	return nil
}

func (m *ProcessManager) Add(pid int) {
	if m.processMap.Get(pid) == nil {
		if process, err := os.FindProcess(pid); err == nil {
			p := NewProcess("", nil, nil)
			p.Process = process
			m.processMap.Set(pid, p)
		}
	}
}

func (m *ProcessManager) Remove(pid int) {
	m.processMap.Remove(pid)
}

func (m *ProcessManager) Processes() []*Process {
	processes := make([]*Process, 0)
	m.processMap.RLockFunc(func(m map[int]interface{}) {
		for _, v := range m {
			processes = append(processes, v.(*Process))
		}
	})
	return processes
}

func (m *ProcessManager) Pids() []int {
	return m.processMap.Keys()
}

func (m *ProcessManager) WaitAll() {
	for _, p := range m.Processes() {
		p.Wait()
	}
}

func (m *ProcessManager) KillAll() error {
	for _, p := range m.Processes() {
		if err := p.Kill(); err != nil {
			return err
		}
	}
	return nil
}
