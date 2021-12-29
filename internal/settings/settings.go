package settings

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/jessevdk/go-flags"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
)

const SettingsFile = "settings.json"

type Settings struct {
	OrderNumber   string `json:"orderNumber" short:"o" long:"order-number" required:"false" description:"Order number (required for pizza tracker)"`
	VIN           string `json:"vin" short:"v" long:"vin" required:"false" description:"Vehicle VIN (required for pizza tracker)"`
	ReservationId string `json:"reservationId" short:"r" long:"reservation" required:"false" description:"Reservation ID (required for backdoor tracker)"`
	RefreshToken  string `json:"refreshToken" long:"refreshToken" required:"false" description:"Refresh token (required for backdoor tracker)"`
	Interval      string `json:"interval" short:"i" long:"interval" required:"false" description:"Interval between checks"`
	PhoneNumber   string `json:"phoneNumber" short:"p" long:"phone" required:"false" description:"Phone number to text (required for textbelt notification)"`
	TextbeltKey   string `json:"textbeltKey" long:"textbelt-key" required:"false" description:"Textbelt API key (required for textbelt notification)"`
	PushoverToken string `json:"pushoverToken" long:"pushover-token" required:"false" description:"Pushover token (required for pushover notification)"`
	PushoverUser  string `json:"pushoverUser" long:"pushover-user" required:"false" description:"Pushover user (required for pushover notification)"`
	RunOnce       bool   `json:"-" long:"once" required:"false" description:"Flag to run check once and stop."`
	Quiet         bool   `json:"-" long:"quiet" short:"q" required:"false" description:"Quiet run - don't prompt for any settings"`
}

// Reads setting from file and command line arguments
// Command line arguments override file settings and will be written to the settings file
func ReadSettings() Settings {
	var settingsFromFile Settings

	// Read settings from file
	jsonData, err := ioutil.ReadFile(getFullPath(SettingsFile))
	if err == nil {
		err = json.Unmarshal(jsonData, &settingsFromFile)
	}

	// Read settings from args
	var settingsFromArgs Settings
	_, err = flags.Parse(&settingsFromArgs)
	if err != nil {
		log.Fatal(err)
	}

	// Override settings with args
	setStringIfEmpty(&settingsFromArgs.VIN, &settingsFromFile.VIN, "", "VIN", settingsFromArgs.Quiet)
	setStringIfEmpty(&settingsFromArgs.OrderNumber, &settingsFromFile.OrderNumber, "", "Order Number", settingsFromArgs.Quiet)
	setStringIfEmpty(&settingsFromArgs.ReservationId, &settingsFromFile.ReservationId, "-", "Reservation Id", settingsFromArgs.Quiet)
	setStringIfEmpty(&settingsFromArgs.RefreshToken, &settingsFromFile.RefreshToken, "-", "RefreshToken", settingsFromArgs.Quiet)
	setStringIfEmpty(&settingsFromArgs.PhoneNumber, &settingsFromFile.PhoneNumber, "-", "Phone Number", settingsFromArgs.Quiet)
	setStringIfEmpty(&settingsFromArgs.TextbeltKey, &settingsFromFile.TextbeltKey, "-", "Textbelt Key", settingsFromArgs.Quiet)
	setStringIfEmpty(&settingsFromArgs.Interval, &settingsFromFile.Interval, "60m", "Check Interval", settingsFromArgs.Quiet)
	setStringIfEmpty(&settingsFromArgs.PushoverToken, &settingsFromFile.PushoverToken, "-", "Pushover token", settingsFromArgs.Quiet)
	setStringIfEmpty(&settingsFromArgs.PushoverUser, &settingsFromFile.PushoverUser, "-", "Pushover user", settingsFromArgs.Quiet)
	settingsFromFile.Quiet = settingsFromArgs.Quiet
	settingsFromFile.RunOnce = settingsFromArgs.RunOnce

	// Write the settings out to file
	jsonOutput, _ := json.MarshalIndent(settingsFromFile, "", "  ")
	_ = ioutil.WriteFile(getFullPath(SettingsFile), jsonOutput, 0644)

	// Echo settings to terminal so the all settings can be seen
	fmt.Println("---- Settings ----")
	fmt.Println(string(jsonOutput))

	return settingsFromFile
}

// Copies a setting from a source to a destination, including a default value if the source is empty
// Will also prompt the user to change a setting and will terminate if a setting is empty
func setStringIfEmpty(src *string, dest *string, defaultValue string, fieldName string, quiet bool) {
	// Copy src to dest of src is not empty
	if len(*src) > 0 {
		*dest = *src
	}

	// If dest is still empty, set to default value
	if len(*dest) == 0 {
		*dest = defaultValue
	}

	if !quiet {
		reader := bufio.NewReader(os.Stdin)
		changeValue := true
		// Ask the user if they would like to change the setting value
		if len(*dest) != 0 {
			fmt.Println(fieldName + " = " + *dest)
			fmt.Print("Would you like to change the value? (y/n): ")
			out, _ := reader.ReadString('\n')
			out = strings.ReplaceAll(out, "\n", "")
			changeValue = strings.EqualFold(out, "y")
		}
		// If the user wants to change the setting, get the value from the command line
		if changeValue {
			fmt.Print(fieldName + ": ")
			out, _ := reader.ReadString('\n')
			out = strings.ReplaceAll(out, "\n", "")
			if len(out) != 0 {
				*dest = out
			}
		}
	}

	// If dest is still empty throw an error
	if len(*dest) == 0 {
		log.Fatal("Required value missing")
	}
}

func getFullPath(file string) string {
	e, _ := os.Executable()
	basePath := path.Dir(e)
	return path.Join(basePath, file)
}
