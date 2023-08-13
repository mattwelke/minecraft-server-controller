package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

func getStatus(serviceName string) (string, error) {
	cmd := exec.Command("systemctl", "is-active", serviceName)
	outputBytes, _ := cmd.CombinedOutput()
	output := strings.TrimSpace(string(outputBytes))
	if output != "inactive" && output != "active" {
		return "", errors.New("failed to get status (error encountered or unknown status returned)")
	}
	return output, nil
}

func waitForStatus(serviceName, desiredStatus string) error {
	for {
		status, err := getStatus(serviceName)
		if err != nil {
			return fmt.Errorf("failed to wait for status %q for server %q: %w", desiredStatus, serviceName, err)
		}
		if status == desiredStatus {
			return nil
		}
		time.Sleep(1 * time.Second)
	}
}

func main() {
	correctPassword := os.Getenv("SHUTDOWN_PASSWORD")
	if correctPassword == "" {
		log.Fatal("SHUTDOWN_PASSWORD environment variable is not set")
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			fmt.Fprint(w, form("", ""))
			return
		}

		providedPassword := r.FormValue("password")
		if providedPassword != correctPassword {
			fmt.Fprint(w, form("Incorrect password!", ""))
			return
		}

		serviceName := r.FormValue("server")
		if serviceName != "minecraft" && serviceName != "minecraft-private" {
			fmt.Fprint(w, form("Invalid server selected!", ""))
			return
		}

		currentStatus, err := getStatus(serviceName)
		if err != nil {
			logAndDisplayError(w, fmt.Sprintf("Failed to check current status of server %q: %%s", serviceName), err)
			return
		}

		action := r.FormValue("action")
		if (action == "stop" && currentStatus == "inactive") || (action == "start" && currentStatus == "active") {
			message := fmt.Sprintf("Server %q is already %q, no action needed.", serviceName, currentStatus)
			fmt.Fprint(w, form(message, serviceName))
			return
		}

		cmd := exec.Command("sudo", "systemctl", action, serviceName)
		_, err = cmd.CombinedOutput()
		if err != nil {
			logAndDisplayError(w, fmt.Sprintf("Failed to execute command for server %q: %%s", serviceName), err)
			return
		}

		desiredStatus := "active"
		if action == "stop" {
			desiredStatus = "inactive"
		}
		err = waitForStatus(serviceName, desiredStatus)
		if err != nil {
			logAndDisplayError(w, fmt.Sprintf("Failed to wait for desired status for server %q: %%s", serviceName), err)
			return
		}

		message := fmt.Sprintf("Server %q status has been changed to %q successfully.", serviceName, desiredStatus)
		fmt.Fprint(w, form(message, serviceName))
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func form(message, serverName string) string {
	return fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head>
			<title>Minecraft Server Control</title>
		</head>
		<body>
			<h1>Minecraft Server Controller</h1>
			<p>Enter the password, select the server to work with, then either start or stop it. Save the password to your password manager for faster server changes in the future.</p>
			<form method="post">
				<label for="password">Password:</label>
				<input type="password" id="password" name="password" required>
				<label for="server">Server:</label>
				<select id="server" name="server">
					<option value="minecraft">Minecraft</option>
					<option value="minecraft-private">Minecraft Private</option>
				</select>
				<button name="action" value="start">Start Server</button>
				<button name="action" value="stop">Stop Server</button>
			</form>
			<p>%s</p>
		</body>
		</html>`, message)
}

func logAndDisplayError(w http.ResponseWriter, format string, err error) {
	message := fmt.Sprintf(format, err)
	log.Println(message)
	fmt.Fprint(w, form(message, ""))
}
