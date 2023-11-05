package raspi

import (
	"bufio"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func getCpuTemp() float64 {
	cmd := exec.Command("vcgencmd", "measure_temp")
	out, err := cmd.Output()

	if err != nil {
		log.Println("Error executing vcgencmd measure_temp. Error:", err.Error())
		return 0
	}

	output := string(out)
	output = strings.Replace(output, "temp=", "", -1)
	output = strings.Replace(output, "'C", "", -1)

	outputFloat, err := strconv.ParseFloat(output, 64)

	if err != nil {
		log.Println("Failed to convert vcgencmd measure_temp output to Float. Error:", err.Error())
		return 0
	}

	return outputFloat
}

func getCpuUsage() []float64 {

	cpuUsages := []float64{}
	prevIdleTime, prevTotalTime := []float64{}, []float64{}

	for i := 0; i < 2; i++ {
		statFile, err := os.OpenFile("/proc/stat", os.O_RDONLY, 0444)
		if err != nil {
			log.Println("Failed to open /proc/stat, Error: ", err.Error())
			return cpuUsages
		}

		defer statFile.Close()

		scanner := bufio.NewScanner(statFile)
		scanner.Split(bufio.ScanLines)

		cpu := 0
		for scanner.Scan() {
			line := scanner.Text()
			if !strings.HasPrefix(line, "cpu") {
				continue
			}
			line = line[5:]
			splittedLine := strings.Fields(line)
			idleTime, _ := strconv.ParseInt(splittedLine[3], 10, 64)
			totalTime := int64(0)
			for _, s := range splittedLine {
				u, _ := strconv.ParseInt(s, 10, 64)
				totalTime += u
			}
			if i > 0 {
				deltaIdleTime := float64(idleTime) - prevIdleTime[cpu]
				deltaTotalTime := float64(totalTime) - prevTotalTime[cpu]
				averageIdleTime := float64(deltaIdleTime) / float64(deltaTotalTime)
				cpuUsage := 1.0 - averageIdleTime
				cpuUsages = append(cpuUsages, cpuUsage*100.0)
			}

			prevIdleTime = append(prevIdleTime, float64(idleTime))
			prevTotalTime = append(prevTotalTime, float64(totalTime))
			cpu = cpu + 1

		}
		time.Sleep(100 * time.Millisecond)
	}

	return cpuUsages
}
