package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"time"
	"trafficGenerator/src/utils"
	"trafficGenerator/src/version"
)

var hosts []string

func pickOneHost() string {
	rand.Seed(time.Now().UnixNano())
	var index = rand.Intn(len(hosts))

	return hosts[index]
}

func genUdpData() []byte {
	M := 1 * 1024

	data := make([]int, M)
	for i := 0; i < M; i++ {
		data[i] = i
	}

	// size:11963323 bytes
	db, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("err")
	}

	log.Printf("size:%v", len(db))
	return db
}

func main() {
	var argName = flag.String("name", "default", "node name")
	// var argRate = flag.Int("r", 100, "pkts per second")
	// var argPkts = flag.Int("k", 0, "# of pkts to send")
	// var argTime = flag.String("t", "1s", "duration")
	var argDstPort = flag.String("p", "12345", "remote port")
	var argLog = flag.String("log", "", "log file dir")

	flag.Parse()
	if *argLog != "" {
		log.SetOutput(utils.GetLogFile(*argLog, fmt.Sprintf("%s-sender.log", *argName)))
	}

	for i := 1; i <= 16; i++ {
		hosts = append(hosts, fmt.Sprintf("10.0.0.%v", i))
	}
	log.Print(hosts)

	version.PrintVersion()

	udpdata := genUdpData()
	// TODO 记录流量的速度
	for {

		host := pickOneHost()
		t1 := time.Now().UnixNano()
		// host := "10.0.0.1"

		log.Printf("Start traffic. host:%v port:%v", host, *argDstPort)
		addr, err := net.ResolveUDPAddr("udp", host+":12345")
		if err != nil {
			fmt.Println("Can't resolve address: ", err)
			os.Exit(1)
		}
		conn, err := net.DialUDP("udp", nil, addr)
		if err != nil {
			fmt.Println("Can't dial: ", err)
			os.Exit(1)
		}
		for i := 0; i < 256; i++ {
			_, err = conn.Write(udpdata)
			if err != nil {
				fmt.Println("failed:", err)
				os.Exit(1)
			}
		}

		log.Printf("data send over,size:%v", len(udpdata))
		// resp := make([]byte, 4)
		// _, err = conn.Read(resp)
		// if err != nil {
		// 	fmt.Println("failed to read UDP msg because of ", err)
		// 	os.Exit(1)
		// }

		// t := binary.BigEndian.Uint32(resp)
		// log.Println(time.Unix(int64(t), 0).String())
		conn.Close()
		time.Sleep(1 * time.Second)

		t2 := time.Now().UnixNano()
		log.Printf("totoal time:%vms", (t2-t1)/1e6)
	}
}
