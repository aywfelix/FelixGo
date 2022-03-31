package shell

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

type Process struct {
	exec.Cmd
	Manager *ProcessManager
	PPid    int
}

func NewProcess(path string, args []string, environment ...[]string) *Process {
	env := os.Environ()
	if len(environment) > 0 {
		for k, v := range environment[0] {
			env[k] = v
		}
	}
	process := &Process{
		Manager: NewProcessManager(),
		PPid:    os.Getpid(),
		Cmd: exec.Cmd{
			Args:       []string{path},
			Path:       path,
			Stdin:      os.Stdin,
			Stdout:     os.Stdout,
			Stderr:     os.Stderr,
			Env:        env,
			ExtraFiles: make([]*os.File, 0),
		},
	}
	process.Dir, _ = os.Getwd()
	if len(args) > 0 {
		start := 0
		if strings.EqualFold(path, args[0]) {
			start = 1
		}
		process.Args = append(process.Args, args[start:]...)
	}
	return process
}

func (p *Process) Pid() int {
	if p.Process != nil {
		return p.Process.Pid
	}
	return 0
}

func (p *Process) Kill() error {
	if err := p.Process.Kill(); err == nil {
		if p.Manager != nil {
			p.Manager.processMap.Remove(p.Pid())
		}
		if runtime.GOOS != "windows" {
			err = p.Process.Release()
		}
		_, err = p.Process.Wait()
		if err != nil {
			return err
		}
		return nil
	} else {
		return err
	}
}

func (p *Process) Start() (int, error) {
	if p.Process != nil {
		return p.Pid(), nil
	}
	p.Env = append(p.Env, fmt.Sprintf("%s=%d", envKeyPPid, p.PPid))
	if err := p.Cmd.Start(); err == nil {
		if p.Manager != nil {
			p.Manager.processMap.Set(p.Process.Pid, p)
		}
		return p.Process.Pid, nil
	} else {
		return 0, err
	}
}

func (p *Process) Run() error {
	if _, err := p.Start(); err == nil {
		return p.Wait()
	} else {
		return err
	}
}

func (p *Process) Release() error {
	return p.Process.Release()
}

func (p *Process) Signal(sig os.Signal) error {
	return p.Process.Signal(sig)
}
