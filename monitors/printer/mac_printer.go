package printer

import (
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// на печати в pdf невидно

func MacPrinter() {
	for {
		checkPrintJobs()
		fmt.Println("checkPrinter")
		// Пауза между проверками
		time.Sleep(1 * time.Second)
	}
}

func checkPrintJobs() {
	cmd := exec.Command("lpstat", "-o")
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("Failed to execute lpstat command:", err)
		return
	}

	outputLines := strings.Split(string(out), "\n")
	for _, line := range outputLines {
		if strings.HasPrefix(line, "printer ") && strings.Contains(line, " ") {
			parts := strings.Fields(line)
			jobID := parts[1]
			jobUser := parts[2]
			jobTitle := parts[5]
			fmt.Printf("Job ID: %s, User: %s, Title: %s\n", jobID, jobUser, jobTitle)
		}
	}
}
