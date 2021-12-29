package state

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

const StateFile = "state.json"

const (
	InProgress = "InProgress"
	Completed  = "Completed"
)

const DefaultWindowStickerHash = "21979c27f520821587157e7dd3af9af3872998d527834f141cd7dc0535aab5b5"

type State struct {
	LastRunTime           string   `json:"lastRunTime"`
	PizzaTrackerStatus    string   `json:"pizzaTrackerStatus"`
	BackdoorTrackerStatus string   `json:"backdoorTrackerStatus"`
	OrderConfirmed        string   `json:"orderConfirmed"`
	InProduction          string   `json:"inProduction"`
	Built                 string   `json:"built"`
	Shipped               string   `json:"shipped"`
	Delivered             string   `json:"delivered"`
	WindowStickers        []string `json:"windowStickers"`
	Status                string   `json:"status"`
	DisplayStatus         string   `json:"displayStatus"`
	VehicleStatusCode     string   `json:"vehicleStatusCode"`
	BuiltWeek             string   `json:"builtWeek"`
}

// Gets the saved state for the system
func GetState() State {
	var state State
	jsonData, err := ioutil.ReadFile(getFullPath(StateFile))
	if err == nil {
		err = json.Unmarshal(jsonData, &state)
	}
	if state.WindowStickers == nil || len(state.WindowStickers) == 0 {
		state.WindowStickers = make([]string, 1)
		state.WindowStickers[0] = DefaultWindowStickerHash
	}
	return state
}

// Saves the state out to file
func SaveState(state State) {
	jsonOutput, _ := json.MarshalIndent(state, "", "  ")
	_ = ioutil.WriteFile(getFullPath(StateFile), jsonOutput, 0644)
}

func PrintState() {
	state := GetState()
	jsonOutput, _ := json.MarshalIndent(state, "", "  ")
	fmt.Println(string(jsonOutput))
}

// Gets status value string from status booleans
func GetStatus(inProgress bool, complete bool) string {
	if inProgress && !complete {
		return InProgress
	}
	if complete {
		return Completed
	}
	return ""
}

func getFullPath(file string) string {
	e, _ := os.Executable()
	basePath := path.Dir(e)
	return path.Join(basePath, file)
}
