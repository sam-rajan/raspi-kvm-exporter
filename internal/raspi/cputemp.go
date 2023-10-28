package raspi

import (
	"errors"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

func getCpuTemp() (error, float64) {
	cmd := exec.Command("vcgencmd", "measure_temp")
	out, err := cmd.Output()

	if err != nil {
		log.Println("Error executing vcgencmd measure_temp. Error:", err.Error())
		return errors.New("Failed to fetch temp"), 0
	}

	output := string(out)
	output = strings.Replace(output, "temp=", "", -1)
	output = strings.Replace(output, "'C", "", -1)

	outputFloat, err := strconv.ParseFloat(output, 64)

	if err != nil {
		log.Println("Failed to convert vcgencmd measure_temp output to Float. Error:", err.Error())
		return errors.New("Failed to fetch temp"), 0
	}

	return nil, outputFloat
}
