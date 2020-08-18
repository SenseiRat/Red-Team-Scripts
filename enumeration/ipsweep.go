package main

import (
    "os"
    "fmt"
    "flag"
    "os/exec"
    "strconv"
    "strings"
)

var ipaddr string
const chan_limit int = 1000

func main() {
    flag.StringVar(&ipaddr, "i", "", "IP Address to enumerate")
    flag.Parse()

    var chans [chan_limit]chan bool
    for i := 0; i < chan_limit; i++ {
        chans[i] = make(chan bool)
    }


    octets := splitAddress(ipaddr)              // Split the address into octets
    oct := convertStr(octets)                   // Convert strings to int
    oct1, oct2, oct3, oct4 := calcRanges(oct)   // Generate the ranges of octets
    loopRanges(oct1, oct2, oct3, oct4, chans)       // Loop through the ranges
}

func splitAddress(ipaddr string) []string {
    octets := strings.Split(ipaddr, ".")
    return octets
}

func convertStr(octets []string) [4]int {
    var oct [4]int
    for i := 0; i < len(octets); i++ {
        octet, err := strconv.Atoi(octets[i])
        if err == nil {
            oct[i] = octet
        } else {
            fmt.Println("Error converting octets.")
            os.Exit(1)
        }
    }
    return oct
}

func calcRanges(oct [4]int) (map[string]int, map[string]int, map[string]int, map[string]int) {
    oct1 := make(map[string]int)
    oct2 := make(map[string]int)
    oct3 := make(map[string]int)
    oct4 := make(map[string]int)

    if oct[0] == 192 && oct[1] == 168 {
        oct1["min"], oct1["max"] = 192, 192
        oct2["min"], oct2["max"] = 168, 168
        oct3["min"], oct3["max"] = 0, 255
        oct4["min"], oct4["max"] = 0, 255
    } else if oct[0] == 172 && oct[1] >= 16 && oct[1] <= 31 {
        oct1["min"], oct1["max"] = 172, 172
        oct2["min"], oct2["max"] = 16, 31
        oct3["min"], oct3["max"] = 0, 255
        oct4["min"], oct4["max"] = 0, 255
    } else if oct[0] == 10 {
        oct1["min"], oct1["max"] = 10, 10
        oct2["min"], oct2["max"] = 0, 255
        oct3["min"], oct3["max"] = 0, 255
        oct4["min"], oct4["max"] = 0, 255
    } else {
        fmt.Println("IP address is not in a private range.")
        os.Exit(1)
    }
    return oct1, oct2, oct3, oct4
}

func loopRanges(oct1 map[string]int, oct2 map[string]int, oct3 map[string]int, oct4 map[string]int, chans [1000]chan bool) {
    for i := 0; i < len(chans); i++ {
        for i1 := oct1["min"]; i1 <= oct1["max"]; i1++ {
            for i2 := oct2["min"]; i2 <= oct2["max"]; i2++ {
                for i3 := oct3["min"]; i3 <= oct3["max"]; i3++ {
                    for i4 := oct4["min"]; i4 <= oct4["max"]; i4++ {
                        address := fmt.Sprintf("%d.%d.%d.%d", i1, i2, i3, i4)

                        go ping(address, chans[i])
                    }
                }
            }
        }
        ping_check := <- chans[i]

        if ping_check == true {
            fmt.Println(address)
        }
    }
}

func ping(address string, channel chan bool) {
    out, _ := exec.Command("ping", address, "-c 1").Output()
    if strings.Contains(string(out), "64 bytes") {
        channel <- true
    } else {
        channel <- false
    }
}