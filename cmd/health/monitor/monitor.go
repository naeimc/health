package monitor

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/naeimc/health/internal/network"
	"github.com/naeimc/health/log"
	"github.com/naeimc/health/monitor"
)

func Monitor(arguments []string) {

	var address string

	fs := flag.NewFlagSet("monitor", flag.ExitOnError)
	fs.StringVar(&address, "address", "127.0.0.1:80", "the address to listen on")
	fs.StringVar(&address, "a", "127.0.0.1:80", "the address to listen on")
	if err := fs.Parse(arguments); err != nil {
		panic(err)
	}

	logger := log.NewPrintLogger()
	go logger.Run()

	queue := make(chan network.Packet, 256)

	listener, err := network.NewTCPListener(address, queue)
	if err != nil {
		panic(err)
	}

	monitor := monitor.NewMonitor()
	monitor.Log = logger.Log

	server := network.NewServer(queue)
	server.Handle = monitor.Handle
	server.Log = logger.Log
	server.Listener(listener)

	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGINT)
		signal.Notify(signals, syscall.SIGTERM)
		logger.Log("received os signal: " + (<-signals).String())
		server.Shutdown()
	}()

	if err := server.Serve(); err != nil {
		panic(err)
	}

	logger.Shutdown()
	logger.Wait()

}
