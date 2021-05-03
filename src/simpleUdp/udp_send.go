package main

import (
	"flag"
	"fmt"
	"log"
	"trafficGenerator/src/utils"
	"trafficGenerator/src/version"
)

func main() {

	var argName = flag.String("name", "default", "node name")
	var argRate = flag.Int("r", 100, "pkts per second")
	var argPkts = flag.Int("k", 0, "# of pkts to send")
	var argTime = flag.String("t", "1s", "duration")
	var argDstPort = flag.Int("p", 12345, "remote port")
	var argLog = flag.String("log", "", "log file dir")

	flag.Parse()
	if *argLog != "" {
		log.SetOutput(utils.GetLogFile(*argLog, fmt.Sprintf("%s-sender.log", *argName)))
	}

	version.PrintVersion()

}
