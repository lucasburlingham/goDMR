package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"

	"gopkg.in/ini.v1"
)

// EngineStatus represents the DMR engine configuration and status
type EngineStatus struct {
	Callsign  string            `json:"callsign"`
	DMRID     int               `json:"dmr_id"`
	Frequency float64           `json:"frequency"`
	Timeslot  int               `json:"timeslot"`
	ColorCode int               `json:"color_code"`
	Services  map[string]string `json:"services"`
}

// Path to config.ini
const configPath = "config.ini"

// Load configuration from config.ini
func readConfig() EngineStatus {
	cfg, err := ini.Load(configPath)
	if err != nil {
		log.Println("Failed to load config.ini, using defaults")
		return EngineStatus{
			Callsign:  "N0CALL",
			DMRID:     1234567,
			Frequency: 438.800,
			Timeslot:  2,
			ColorCode: 1,
			Services:  map[string]string{"dmr": "stopped"},
		}
	}

	section := cfg.Section("dmr")
	return EngineStatus{
		Callsign:  section.Key("callsign").MustString("N0CALL"),
		DMRID:     section.Key("dmr_id").MustInt(1234567),
		Frequency: section.Key("frequency").MustFloat64(438.800),
		Timeslot:  section.Key("timeslot").MustInt(1),
		ColorCode: section.Key("color_code").MustInt(1),
		Services:  map[string]string{"dmr": "stopped"},
	}
}

// writeConfig saves EngineStatus to config.ini
func writeConfig(cfg EngineStatus) error {
	file, err := ini.Load(configPath)
	if err != nil {
		file = ini.Empty()
	}

	section := file.Section("dmr")
	section.Key("callsign").SetValue(cfg.Callsign)
	section.Key("dmr_id").SetValue(fmt.Sprintf("%d", cfg.DMRID))
	section.Key("frequency").SetValue(fmt.Sprintf("%.3f", cfg.Frequency))
	section.Key("timeslot").SetValue(fmt.Sprintf("%d", cfg.Timeslot))
	section.Key("color_code").SetValue(fmt.Sprintf("%d", cfg.ColorCode))

	// Log the config write
	log.Println("DMR configuration saved to config.ini:", cfg)

	return file.SaveTo(configPath)

}

// apiStatus handles GET /api/status
func apiStatus(w http.ResponseWriter, r *http.Request) {
	status := readConfig()

	// Check if DMR engine binary is running
	cmd := exec.Command("pgrep", "-f", "mmdvm") // replace "mmdvm" with your binary
	if err := cmd.Run(); err == nil {
		status.Services["dmr"] = "running"
	}

	jsonData, _ := json.Marshal(status)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
	// Log the status request
	log.Println("DMR status requested:", status)
}

// apiConfig handles POST /api/config
func apiConfig(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var cfg EngineStatus
	if err := json.NewDecoder(r.Body).Decode(&cfg); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := writeConfig(cfg); err != nil {
		http.Error(w, "Failed to save config", http.StatusInternalServerError)
		return
	}

	// Restart DMR engine after config change
	// _ = exec.Command("systemctl", "restart", "mmdvm.service").Run()

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Pragma", "no-cache")
	// Access-Control-Allow-Origin for local UI access
	w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:8080")
	w.Write([]byte(`{"result":"ok"}`))

	// Log the config change
	log.Println("DMR configuration updated and engine restarted with the following settings:", cfg)
}

// resetConfig handles POST /api/reset
func resetConfig(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	defaultCfg := EngineStatus{
		Callsign:  "N0CALL",
		DMRID:     1234567,
		Frequency: 438.800,
		Timeslot:  2,
		ColorCode: 1,
		Services:  map[string]string{"dmr": "stopped"},
	}

	if err := writeConfig(defaultCfg); err != nil {
		http.Error(w, "Failed to reset config", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"result":"config reset to defaults"}`))

	// Log the config reset
	log.Println("DMR configuration reset to defaults:", defaultCfg)
}

func backupConfig(w http.ResponseWriter, r *http.Request) {
	// redirect user to download the config.ini file
	http.ServeFile(w, r, configPath)
	log.Println("DMR configuration backup requested")
}

func main() {
	http.HandleFunc("/api/status", apiStatus)
	http.HandleFunc("/api/config", apiConfig)
	http.HandleFunc("/api/reset", resetConfig)
	http.HandleFunc("/api/backup", backupConfig)

	log.Println("Starting DMR API server on http://127.0.0.1:8080")
	if err := http.ListenAndServe("127.0.0.1:8080", nil); err != nil {
		log.Fatal(err)
	}
}
