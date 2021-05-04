package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"strings"
	"trafficGenerator/src/utils"
)

var limitChan = make(chan bool, 100)

func process(conn *net.UDPConn) {
	data := make([]byte, 1024)
	n, remoteAddr, err := conn.ReadFromUDP(data)
	if err != nil {
		fmt.Println("failed to read UDP msg because of ", err.Error())
		return
	}
	log.Printf("read form %v length:%v", remoteAddr, n)

	// daytime := time.Now().Unix()
	// b := make([]byte, 4)
	// binary.BigEndian.PutUint32(b, uint32(daytime))
	// conn.WriteToUDP(b, remoteAddr)
	<-limitChan
}

func getHostIp() string {
	addrs, err := net.InterfaceAddrs()

	if err != nil {
		log.Fatalf("get net addr error\n")
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			ip := ipnet.IP
			if ip.To4() != nil && strings.Index(ip.String(), "10.0") != -1 {
				fmt.Println(ip.String())
				return ip.String()
			}

		}
	}
	return "127.0.0.1"
}

func StartUdp() {
	host := getHostIp()
	if host == "" {
		log.Fatalf("can not resolve cif\n")
	}
	log.Println("ip addr:", host)

	addr, err := net.ResolveUDPAddr("udp", host+":12345")
	if err != nil {
		log.Fatalf("Can't resolve address")
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatalf("Error listening:%v", err)
	}

	defer conn.Close()
	for {
		limitChan <- true
		go process(conn)
	}

}

func main() {
	var argName = flag.String("name", "h.log", "host name")
	var argLog = flag.String("log", "/home/caoyuhua/go/src/trafficGenerator/src", "log dir")
	log.Printf("Recv begin host:%v \n", argName)

	flag.Parse()
	if *argLog != "" {
		log.SetOutput(utils.GetLogFile(*argLog, fmt.Sprintf("%s-recv.log", *argName)))
	}

	StartUdp()

	log.Printf("Recv end, ip:%v", getHostIp())
}
