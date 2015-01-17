package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"time"
)

func calculateBandwidth(interfaceName string, statisticsName string, delay int, output chan int) {
	var rawBytes []byte
	var bytes int
	var previousBytes int
	var err error

	for {
		rawBytes, err = ioutil.ReadFile("/sys/class/net/" + interfaceName + "/statistics/" + statisticsName)

		if err != nil {
			log.Println("Error reading interface ", interfaceName, statisticsName, ":", err)
			continue
		}

		bytes, err = strconv.Atoi(strings.Trim(string(rawBytes), "\n"))

		if err != nil {
			log.Println("Error parsing", statisticsName, err)
			continue
		}

		if previousBytes != 0 {
			output <- (bytes - previousBytes) / 1000 / delay
		}

		previousBytes = bytes

		time.Sleep(time.Duration(delay) * time.Second)
	}
}

func main() {
	interfaceName, delay := parseArgs()

	log.Println("Monitor interface", interfaceName)
	log.Println("Monitor delay", delay)

	downloadSpeeds := make(chan int)
	uploadSpeeds := make(chan int)

	go calculateBandwidth(interfaceName, "rx_bytes", delay, downloadSpeeds)
	go calculateBandwidth(interfaceName, "tx_bytes", delay, uploadSpeeds)

	for {
		select {
		case downloadSpeed := <-downloadSpeeds:
			fmt.Println("D:", downloadSpeed, "KB/s")
		case uploadSpeed := <-uploadSpeeds:
			fmt.Println("U:", uploadSpeed, "KB/s")
		}
	}
}

func parseArgs() (string, int) {
	interfaceName := flag.String("interface", "eno1", "The interface to monitor")
	delay := flag.Int("delay", 1, "The monitor delay")
	flag.Parse()
	return *interfaceName, *delay
}
