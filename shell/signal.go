package shell

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type SignalHandler func(sig os.Signal)

var (
	signalHandlerMap  = make(map[os.Signal][]SignalHandler)
	shutdownSignalMap = map[os.Signal]struct{}{
		syscall.SIGINT:  {},
		syscall.SIGQUIT: {},
		syscall.SIGKILL: {},
		syscall.SIGTERM: {},
		syscall.SIGABRT: {},
	}
)

func init() {
	for sig, _ := range shutdownSignalMap {
		signalHandlerMap[sig] = make([]SignalHandler, 0)
	}
}

func AddSigHandler(handler SignalHandler, signals ...os.Signal) {
	for _, sig := range signals {
		signalHandlerMap[sig] = append(signalHandlerMap[sig], handler)
	}
}

func AddShutdownHandler(handlers ...SignalHandler) {
	for _, handler := range handlers {
		for sig, _ := range shutdownSignalMap {
			signalHandlerMap[sig] = append(signalHandlerMap[sig], handler)
		}
	}
}

func SignalWait() {
	// 获取所有设置的信号
	signals := make([]os.Signal, 0)
	for sig, _ := range signalHandlerMap {
		signals = append(signals, sig)
	}
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, signals...)
	var sig os.Signal
	for {
		var wg sync.WaitGroup
		sig = <-sigChan
		if handlers, ok := signalHandlerMap[sig]; ok {
			for _, handler := range handlers {
				// 利用协程执行信号函数
				wg.Add(1)
				go func(handler SignalHandler, sig os.Signal) {
					defer wg.Done()
					handler(sig)
				}(handler, sig)
			}
		}
		if _, ok := shutdownSignalMap[sig]; ok {
			wg.Wait()
			return
		}
	}
}
