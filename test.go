package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	slack(querySpace())
}

func querySpace() int {
	cmd, _ := exec.Command("/bin/bash", "-c", "df -Ph / | tail -1 | awk '{print $5}'").CombinedOutput()
	space := strings.TrimRight(string(cmd), "\n")

	remaining, _ := strconv.Atoi(strings.TrimRight(space, "%"))
	available := 100 - remaining

	return available
}

func slack(space int) {

	client := http.Client{}

	slackMessage := `{"text":"` + strconv.Itoa(space) + `% of disk space is still available on the comms server."}`

	jsonStr := []byte(slackMessage)

	req, err := http.NewRequest("POST", os.Getenv("SLACK_URL"), bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Unable to reach the server.")
	}
	defer resp.Body.Close()
}
