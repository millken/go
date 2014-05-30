package  main
import (
	"time"
	"fmt"
	"os"
	"os/signal"	
)

func init() {
	
}
var chans = make(map[string]chan bool)

func main() {
	ips := []string{"1.1.1.1","2.2.2.2","3.3.3.3"}
	
	for i := 0; i < 3; i++ {
	
	for _, ip := range ips {
   		chans[ip] = make(chan bool)
		//fmt.Printf("ip : %s\n", ip)
		go w1(i, ip)
		
		
		
	}

	time.Sleep(3 * time.Second)
	for ip,_ := range chans {
		chans[ip] <- true
	}
}
	for i,_ := range chans {
		close(chans[i])
	}
	terminate := make(chan os.Signal)
	signal.Notify(terminate, os.Interrupt)

	<-terminate	
}

func w1(i int,ip string) {
	go w2(i, ip)
}

func w2(i int, ip string) {
	for {
		select{
			case <-time.After(time.Second * 1):
			fmt.Printf("w2 ip %d: %s\n", i, ip)
		case <-chans[ip]:
			return
		}
	}
}