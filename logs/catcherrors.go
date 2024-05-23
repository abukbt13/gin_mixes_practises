package logs

import (
	"fmt"
	"log"
	"os"
)

func LogToFile(msg string) {
	file, err := os.OpenFile("error.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening log file:", err)
		return
	}
	defer file.Close()

	log.SetOutput(file)
	log.Println(msg)
}
