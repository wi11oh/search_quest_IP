package main

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

func main() {
	loop()

	for {
		fmt.Println("\n無入力Enterで終了、R入力で再検索:")
		var input string
		fmt.Scanln(&input)

		if input == "R" || input == "r" {
			loop()
		} else {
			fmt.Println("終了します。")
			break
		}
	}
}

func loop() {
	fmt.Println("\n===== QUEST IP SEARCH =====")
	exec.Command("chcp", "65001").Run()

	ip := "192.168.0."

	var wg1 sync.WaitGroup
	var reachableIPs []string
	totalIPs := 254

	fmt.Println("\nローカルIPを探索しています:")
	for i := 1; i <= totalIPs; i++ {
		wg1.Add(1)
		go func(i int) {
			scanningIp := ip + strconv.Itoa(i)

			defer wg1.Done()

			if ping(scanningIp) {
				reachableIPs = append(reachableIPs, scanningIp)
			}

			fmt.Printf("\rScanned: %-15s%s", scanningIp, createProgressBar(i, totalIPs))
		}(i)
		time.Sleep(10 * time.Millisecond)
	}

	wg1.Wait()

	// 小細工すぎる
	fmt.Println("\rScanned: 192.168.0.255  #################### 100.00%")

	var wg2 sync.WaitGroup
	var dm sync.Mutex
	devices := map[string]string{}

	for i := 0; i <= len(reachableIPs)-1; i++ {
		wg2.Add(1)
		arpingIp := reachableIPs[i]
		go func(arpingIp string) {
			defer wg2.Done()

			macAddr := getMacAddr(arpingIp)

			dm.Lock()
			devices[macAddr] = arpingIp
			dm.Unlock()

		}(arpingIp)
	}

	wg2.Wait()

	fmt.Println("\n検出されたデバイス:")
	for k, v := range devices {
		fmt.Printf("%-15s | %s\n", v, k)
	}

	fmt.Println("\nQUESTっぽい奴:")
	quests := searchOculus(devices)
	if quests == nil {
		fmt.Println("なし")
	} else {
		for k, v := range quests {
			fmt.Printf("%-15s | %s\n", v, k)
		}
	}
}

func ping(ip string) bool {
	err := exec.Command("ping", "-n", "1", "-w", "500", ip).Run()
	return err == nil
}

func createProgressBar(current, total int) string {
	const barLength = 20
	progress := float64(current) / float64(total)
	numBars := int(progress * barLength)

	var progressBar string
	for i := 0; i < barLength; i++ {
		if i < numBars {
			progressBar += "#"
		} else {
			progressBar += "-"
		}
	}

	return fmt.Sprintf("%s %.2f%%", progressBar, progress*100)
}

func getMacAddr(ip string) (macAddress string) {
	output, err := exec.Command("arp", "-a", ip).Output()

	if err != nil {
		macAddress = ""
	} else {
		matches := regexp.MustCompile(`([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2})`).FindStringSubmatch(string(output))
		if len(matches) > 0 {
			macAddress = matches[0]
		} else {
			macAddress = "このPC"
		}
	}

	return
}

func searchOculus(devices map[string]string) map[string]string {
	prefixes := []string{"00-01-61", "80-f3-ef", "88-25-08", "94-f9-29", "b4-17-a8", "c0-dd-8a", "cc-a1-74"}
	oculuses := make(map[string]string)

	for key, value := range devices {
		for _, prefix := range prefixes {
			if strings.HasPrefix(key, prefix) {
				oculuses[key] = value
			}
		}
	}

	if len(oculuses) == 0 {
		oculuses = nil
	}

	return oculuses
}
