package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

func pingList(list []string) []string {
	openHostList = []string{}
	fmt.Println("\nhost discovering...")
	for _, host := range list {
		//ch <- true
		//wg.Add(1)
		fmt.Printf("%s is scanning\n", host)
		ping(host)
	}
	//wg.Wait()
	return openHostList
}

func ping(host string) {
	conn, err := net.Dial("ip:icmp", host)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	defer conn.Close()

	var msg [512]byte
	msg[0] = 8
	msg[1] = 0
	msg[2] = 0
	msg[3] = 0
	msg[4] = 0
	msg[5] = 13
	msg[6] = 0
	msg[7] = 37
	msg[8] = 99
	len := 9
	check := checkSum(msg[0:len])
	msg[2] = byte(check >> 8)
	msg[3] = byte(check & 0xff)
	//fmt.Println(msg[0:len])
	for i := 0; i < 2; i++ {
		_, err = conn.Write(msg[0:len])
		if err != nil {
			//fmt.Println(err.Error())
			continue
		}

		conn.SetReadDeadline((time.Now().Add(time.Millisecond * 400)))
		_, err := conn.Read(msg[0:])
		if err != nil {
			//fmt.Println(err.Error())
			continue
		}

		//fmt.Println(msg[0 : 20+len])
		//fmt.Println("Got response")
		if msg[20+5] == 13 && msg[20+7] == 37 && msg[20+8] == 99 {
			fmt.Printf("%s open\n", ljust(host, 21))
			openHostList = append(openHostList, host)
			//<-ch
			//wg.Done()
			return
		}
	}
	//<-ch
	//wg.Done()
}

func checkSum(msg []byte) uint16 {
	sum := 0

	len := len(msg)
	for i := 0; i < len-1; i += 2 {
		sum += int(msg[i])*256 + int(msg[i+1])
	}
	if len%2 == 1 {
		sum += int(msg[len-1]) * 256
	}

	sum = (sum >> 16) + (sum & 0xffff)
	sum += (sum >> 16)
	var answer uint16 = uint16(^sum)
	return answer
}
