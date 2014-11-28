// goGetMyIP project main.go
package main

import (
	"fmt"
	"net"
	"os"
)

type NetCard struct {
	name string
	ipv4 []net.Addr
}

const (
	LOG_ULTRA_LESS = 0
	LOG_LESS       = 1
	LOG_ALL        = 2
)

func usage() {
	fmt.Println("Usages about this tool: \n\t$ goGetMyIP      # print out IPv4 addr only\n\t$ goGetMyIP all  # print out all interfaces and address\n\t$ goGetMyIP less # print the interfaces that conatains IP address\n\n")
}

func PrintOut(interface_out []net.Interface, log_level int) {
	if log_level == LOG_ALL {
		fmt.Println("Here're the all of your interfaces:\n")
	} else if log_level == LOG_LESS {
		fmt.Println("Here're the interfaces that conatian IP address:\n")
	}

	for _, iface_out := range interface_out {
		if log_level != LOG_ULTRA_LESS {
			fmt.Printf("Interface[%d]: %s\n", iface_out.Index, iface_out.Name)
			fmt.Printf("Hardware Address: %v", iface_out.HardwareAddr)
			if len(iface_out.HardwareAddr) == 0 {
				fmt.Println("(nil)")
			} else {
				fmt.Println()
			}
		}

		addrs, err := iface_out.Addrs()
		if err != nil {
			fmt.Println(err)
			continue
		}
		for _, addr := range addrs {

			ipnet, ok := addr.(*net.IPNet)
			if !ok {
				continue
			}
			v4 := ipnet.IP.To4()
			if (log_level != LOG_ALL) && (v4 == nil || v4[0] == 127) { // loopback address
				continue
			}

			fmt.Println(addr.String())
		}

		if log_level != LOG_ULTRA_LESS {
			// print out delimiter
			fmt.Println()
		}

	}
}

func main() {
	var log_level int
	var interface_out []net.Interface

	// hint: os.Args[0] is the command name
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "help":
			usage()
			return
		case "all":
			log_level = LOG_ALL
		case "less":
			log_level = LOG_LESS
		}
	} else {
		log_level = LOG_ULTRA_LESS
	}

	interfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for _, iface := range interfaces {
		addrs, err := iface.Addrs()
		if err != nil {
			fmt.Println(err)
			continue
		}

		if (log_level != LOG_ALL) && (len(addrs) == 0) {
			continue
		}

		isOut := false
		for _, addr := range addrs {

			ipnet, ok := addr.(*net.IPNet)
			if !ok {
				continue
			}
			v4 := ipnet.IP.To4()
			if (log_level != LOG_ALL) && (v4 == nil || v4[0] == 127) { // loopback address
				continue
			}

			isOut = true

		}
		if isOut {
			interface_out = append(interface_out, iface)
		}
	}

	PrintOut(interface_out, log_level)

	if len(interfaces) == 0 {
		fmt.Println("The interfaces didn't conatian any IP address")
	}
}
