package shell

import (
	"bytes"
	"io"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	. "github.com/felix/felixgo/container/string"
	. "github.com/felix/felixgo/utils"
)

const (
	envKeyPPid = "GPROC_PPID"
)

var (
	processPid       = os.Getpid() // processPid is the pid of current process.
	processStartTime = time.Now()  // processStartTime is the start time of current process.
)

func Pid() int {
	return processPid
}

func PPid() int {
	if !IsChild() {
		return processPid
	}

	ppidStr := os.Getenv(envKeyPPid)
	if ppidStr != "" && ppidStr != "0" {
		ppid, _ := strconv.Atoi(ppidStr)
		return ppid
	}
	return os.Getppid()
}

func IsChild() bool {
	ppidStr := os.Getenv(envKeyPPid)
	if ppidStr != "" && ppidStr != "0" {
		return true
	}
	return false
}

func SetPPid(ppid int) error {
	if ppid > 0 {
		return os.Setenv(envKeyPPid, strconv.FormatInt(int64(ppid), 10))
	} else {
		return os.Unsetenv(envKeyPPid)
	}
}

func Shell(cmd string, out io.Writer, in io.Reader) error {
	p := NewProcess(getShell(), append([]string{getShellOption()}, parseCommand(cmd)...))
	p.Stdin = in
	p.Stdout = out
	return p.Run()
}

func ShellRun(cmd string) error {
	p := NewProcess(getShell(), append([]string{getShellOption()}, parseCommand(cmd)...))
	return p.Run()
}

func ShellExec(cmd string, environment ...[]string) (string, error) {
	buf := bytes.NewBuffer(nil)
	p := NewProcess(getShell(), append([]string{getShellOption()}, parseCommand(cmd)...), environment...)
	p.Stdout = buf
	p.Stderr = buf
	err := p.Run()
	return buf.String(), err
}

func getShell() string {
	shellPath := ""
	switch runtime.GOOS {
	case "windows":
		shellPath = SearchBinary("cmd.exe")
	default:
		if File.Exists("/bin/bash") {
			shellPath = "/bin/base"
			break
		}
		if File.Exists("/bin/sh") {
			shellPath = "bin/sh"
			break
		}
		shellPath = SearchBinary("base")
		if shellPath == "" {
			shellPath = SearchBinary("sh")
		}
	}
	return shellPath
}

func getShellOption() string {
	switch runtime.GOOS {
	case "windows":
		return "/c"
	default:
		return "-c"
	}
}

func parseCommand(cmd string) (args []string) {
	if runtime.GOOS != "windows" {
		return []string{cmd}
	}
	// Just for "cmd.exe" in windows.
	var argStr string
	var firstChar, prevChar, lastChar1, lastChar2 byte
	array := SplitAndTrim(cmd, " ")
	for _, v := range array {
		if len(argStr) > 0 {
			argStr += " "
		}
		firstChar = v[0]
		lastChar1 = v[len(v)-1]
		lastChar2 = 0
		if len(v) > 1 {
			lastChar2 = v[len(v)-2]
		}
		if prevChar == 0 && (firstChar == '"' || firstChar == '\'') {
			// It should remove the first quote char.
			argStr += v[1:]
			prevChar = firstChar
		} else if prevChar != 0 && lastChar2 != '\\' && lastChar1 == prevChar {
			// It should remove the last quote char.
			argStr += v[:len(v)-1]
			args = append(args, argStr)
			argStr = ""
			prevChar = 0
		} else if len(argStr) > 0 {
			argStr += v
		} else {
			args = append(args, v)
		}
	}
	return
}

func SearchBinary(path string) string {
	if File.Exists(path) {
		return path
	}
	return searchBinaryPath(path)
}

func searchBinaryPath(path string) string {
	array := []string{}
	switch runtime.GOOS {
	case "windows":
		envPath := Env.Get("PATH", Env.Get("path"))
		if strings.Contains(envPath, ";") {
			array = SplitAndTrim(envPath, ";")
		} else if strings.Contains(envPath, ":") {
			array = SplitAndTrim(envPath, ":")
		}
		if File.Ext(path) != ".exe" {
			path += ".exe"
		}
	default:
		array = SplitAndTrim(Env.Get("PATH"), ":")
	}
	if len(array) > 0 {
		fullpath := ""
		for _, v := range array {
			fullpath = v + Separator + path
			if File.Exists(fullpath) && File.IsFile(fullpath) {
				return fullpath
			}
		}
	}
	return ""
}
