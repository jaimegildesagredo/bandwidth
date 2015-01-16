package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)

const (
	DELAY = 1 * time.Second
)

func main() {
	interfaceName := getInterfaceName()

	fmt.Println("Interface", interfaceName)

	go func() {
		var rawRxBytes []byte
		var rxBytes int
		var oldRxBytes int
		var err error

		for {
			rawRxBytes, err = ioutil.ReadFile("/sys/class/net/" + interfaceName + "/statistics/rx_bytes")

			if err != nil {
				fmt.Println("Error", err)
				continue
			}

			rxBytes, err = strconv.Atoi(strings.Trim(string(rawRxBytes), "\n"))

			if err != nil {
				fmt.Println("Error", err)
				continue
			}

			if oldRxBytes != 0 {
				fmt.Println("D:", (rxBytes-oldRxBytes)/1000/(int(DELAY/time.Second)), "KB/s")
			}

			oldRxBytes = rxBytes

			time.Sleep(DELAY)
		}
	}()

	go func() {
		var rawTxBytes []byte
		var txBytes int
		var oldTxBytes int
		var err error

		for {
			rawTxBytes, err = ioutil.ReadFile("/sys/class/net/" + interfaceName + "/statistics/tx_bytes")

			if err != nil {
				fmt.Println("Error", err)
				continue
			}

			txBytes, err = strconv.Atoi(strings.Trim(string(rawTxBytes), "\n"))

			if err != nil {
				fmt.Println("Error", err)
				continue
			}

			if oldTxBytes != 0 {
				fmt.Println("U:", (txBytes-oldTxBytes)/1000/(int(DELAY/time.Second)), "KB/s")
			}

			oldTxBytes = txBytes

			time.Sleep(DELAY)
		}
	}()

	for {
		time.Sleep(1 * time.Minute)
	}
}

func getInterfaceName() string {
	value := flag.String("interface", "eno1", "The interface name to monitor")
	flag.Parse()
	return *value
}
