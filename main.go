package main

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/xjasonlyu/tun2socks/constant"
	"github.com/xjasonlyu/tun2socks/engine"
	"github.com/xjasonlyu/tun2socks/log"

	flag "github.com/spf13/pflag"
)

var (
	device  string
	iface   string
	level   string
	proxy   string
	secret  string
	stats   string
	mtu     int
	version bool
)

func init() {
	flag.StringVarP(&device, "device", "d", "", "Use this device [driver://]name")
	flag.StringVarP(&iface, "interface", "i", "", "Use network INTERFACE (Darwin/Linux only)")
	flag.StringVarP(&proxy, "proxy", "p", "", "Use this proxy [protocol://]host[:port]")
	flag.StringVarP(&level, "loglevel", "l", "info", "Log level [debug|info|warn|error|silent]")
	flag.StringVar(&secret, "secret", "", "HTTP statistic server auth secret")
	flag.StringVar(&stats, "stats", "", "HTTP statistic server listen address")
	flag.IntVarP(&mtu, "mtu", "m", 0, "Maximum transmission unit")
	flag.BoolVarP(&version, "version", "v", false, "Show version information and quit")
	flag.Parse()
}

func main() {
	if version {
		fmt.Printf("%s %s\n%s/%s, %s, %s\n",
			constant.Name,
			constant.Version,
			runtime.GOOS,
			runtime.GOARCH,
			runtime.Version(),
			constant.BuildTime,
		)
		os.Exit(0)
	}

	options := []engine.Option{
		engine.WithDevice(device),
		engine.WithInterface(iface),
		engine.WithLogLevel(level),
		engine.WithMTU(mtu),
		engine.WithProxy(proxy),
		engine.WithStats(stats, secret),
	}

	eng := engine.New(options...)
	if err := eng.Start(); err != nil {
		log.Fatalf("Start engine error: %v", err)
	}
	defer eng.Stop()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
}
