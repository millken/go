package main
import (
	"expvar"
	"flag"
	"github.com/qiniu/log"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"runtime/pprof"
	"strings"
	"time"
)

var VERSION string = "2.2.6"
var gitVersion string
var serverId string
var serverIP string
var serverGroups []string

var timeStarted = time.Now()
var qCounter = expvar.NewInt("qCounter")

var (
	flagconfig      = flag.String("config", "./", "directory of zone files")
	flagcheckconfig = flag.Bool("checkconfig", false, "check configuration and exit")
	flagidentifier  = flag.String("identifier", "", "identifier (hostname, pop name or similar)")
	flaginter       = flag.String("interface", "*", "set the listener address")
	flagport        = flag.String("port", "53", "default port number")
	flaghttp        = flag.String("http", ":8053", "http listen address (:8053)")
	flaglog         = flag.Bool("log", false, "be more verbose")

	cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
	memprofile = flag.String("memprofile", "", "write memory profile to this file")
)

func init() {
	if len(gitVersion) > 0 {
		VERSION = VERSION + "/" + gitVersion
	}

	log.SetOutputLevel(log.Ldebug)
}

func main() {
	flag.Parse()

	if len(*flagidentifier) > 0 {
		ids := strings.Split(*flagidentifier, ",")
		serverId = ids[0]
		if len(ids) > 1 {
			serverGroups = ids[1:]
		}
	}

	configFileName := filepath.Clean(*flagconfig + "/ydsync.conf")
	log.Debugf("config file '%s'\n", configFileName)

	if *flagcheckconfig {

		err := configReader(configFileName)
		if err != nil {
			log.Println("Errors reading config", err)
			os.Exit(2)
		}

		return
	}

	log.Printf("Starting mkdns %s\n", VERSION)

	if *cpuprofile != "" {
		prof, err := os.Create(*cpuprofile)
		if err != nil {
			panic(err.Error())
		}

		pprof.StartCPUProfile(prof)
		defer func() {
			log.Println("closing file")
			prof.Close()
		}()
		defer func() {
			log.Println("stopping profile")
			pprof.StopCPUProfile()
		}()
	}

	go configWatcher(configFileName)

	if *flaginter == "*" {
		addrs, _ := net.InterfaceAddrs()
		ips := make([]string, 0)
		for _, addr := range addrs {
			ip, _, err := net.ParseCIDR(addr.String())
			if err != nil {
				continue
			}
			if !(ip.IsLoopback() || ip.IsGlobalUnicast()) {
				continue
			}
			ips = append(ips, ip.String())
		}
		*flaginter = strings.Join(ips, ",")
	}

	inter := getInterfaces()

	for _, host := range inter {
		go listenAndServe(host)
	}

	terminate := make(chan os.Signal)
	signal.Notify(terminate, os.Interrupt)

	<-terminate
	log.Printf("geodns: signal received, stopping")

	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.WriteHeapProfile(f)
		f.Close()
	}

	//os.Exit(0)

}
