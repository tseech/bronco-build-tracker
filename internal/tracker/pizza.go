package tracker

import (
	"crypto/sha256"
	"fmt"
	"github.com/tseech/bronco-build-tracker/internal/notification"
	"github.com/tseech/bronco-build-tracker/internal/settings"
	"github.com/tseech/bronco-build-tracker/internal/state"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
)

const FordStatusUrl = "https://shop.ford.com/vehicleordertracking/status/?vin=%s&orderId=%s"
const FordStatusUrlNoOrder = "https://shop.ford.com/vehicleordertracking/status/?vin=%s"
const StatusMatchString = "<li id=\"step-%s\" aria-live=\"polite\" aria-describedby=\"aria-%s\" class=\"brand-tracker-step brand-tracker-step__%s\""

// Checks the status of the ford tracker
func CheckStatus(settings settings.Settings) {
	// Get previous currentState so only changes are reported
	currentState := state.GetState()
	currentState.LastRunTime = time.Now().Local().Format(time.RFC822)
	// Try to read the ford status site
	statusSite, url, success := getStatusSite(settings.VIN, settings.OrderNumber)

	// If the site cannot be read, return
	if !success {
		fmt.Println("Status site is not available")
		currentState.PizzaTrackerStatus = "Status site is not available"
		state.SaveState(currentState)
		return
	}

	// Look for the status of each step
	changed := processStatus("order_confirmed", &currentState.OrderConfirmed, statusSite)
	changed = processStatus("in_production", &currentState.InProduction, statusSite) || changed
	changed = processStatus("built", &currentState.Built, statusSite) || changed
	changed = processStatus("shipped", &currentState.Shipped, statusSite) || changed
	changed = processStatus("delivery", &currentState.Delivered, statusSite) || changed

	// Report if and of the status values have changed
	if changed {
		fmt.Println("Status change found!")
		notification.Notify(url, "Status change found!", settings)
	}

	// Get the window sticker link from the site and hash the PDF
	wsURL := getWindowStickerLink(statusSite)
	wsHash := getResponseSha256Sum(wsURL)

	// If the hash has not been seen, notify and record the new hash
	if !contains(currentState.WindowStickers, wsHash) {
		fmt.Println("New window sticker found!")
		currentState.WindowStickers = append(currentState.WindowStickers, wsHash)
		notification.Notify(wsURL, "New window sticker found!", settings)
	}

	// Save the currentState for the next run
	currentState.PizzaTrackerStatus = "Success"
	state.SaveState(currentState)
}

// Gets the status site and tries to tell if the page is good or if the server is not
// providing good data. It will try multiple URLs because Ford seems to have servers that
// behave differently, and you cannot control which one you hit.
func getStatusSite(vin string, orderNumber string) (string, string, bool) {
	// Try the first URL
	url := fmt.Sprintf(FordStatusUrlNoOrder, vin)
	body := getResponseBody(url)
	bodyString := string(body)
	matchFound, _ := regexp.MatchString("vehicle-status", bodyString)
	if matchFound {
		return bodyString, url, true
	}

	// Try the second URL
	url = fmt.Sprintf(FordStatusUrl, vin, orderNumber)
	body = getResponseBody(url)
	bodyString = string(body)
	matchFound, _ = regexp.MatchString("vehicle-status", bodyString)
	if matchFound {
		return bodyString, url, true
	}

	// Failure
	return "", "", false
}

// Checks a step status and updates the state data
func processStatus(step string, stateField *string, statusSite string) bool {
	inProgress, confirmed := checkStepStatus(step, statusSite)
	return updateStatus(stateField, inProgress, confirmed)
}

// Trys to detect of a step is in progress, confirmed, or neither on the status site
func checkStepStatus(step string, statusSite string) (bool, bool) {
	inProgress, _ := regexp.MatchString(fmt.Sprintf(StatusMatchString, step, step, "inprogress"), statusSite)
	confirmed, _ := regexp.MatchString(fmt.Sprintf(StatusMatchString, step, step, "completed"), statusSite)
	return inProgress, confirmed
}

// Trys to find the window sticker link from the status page and return it
func getWindowStickerLink(statusSite string) string {
	const WsPrefix = "<a id=\"windowSticker\" class=\" brand-button-linknu__button-light  \" href=\""
	start := strings.Index(statusSite, WsPrefix)
	sub := statusSite[start+len(WsPrefix):]
	end := strings.Index(sub, "\"")
	sub = sub[:end]
	return sub
}

// Gets the SHA256 sum from the response to a request to a URL
func getResponseSha256Sum(url string) string {
	body := getResponseBody(url)
	sumBytes := sha256.Sum256(body)
	sumString := fmt.Sprintf("%x", sumBytes)
	return sumString
}

// Updates a status field in the state object based on status booleans
func updateStatus(field *string, inProgress bool, complete bool) bool {
	status := state.GetStatus(inProgress, complete)
	if *field == status {
		return false
	}
	if *field == "" {
		*field = status
		return true
	}
	if *field == state.InProgress && status == state.Completed {
		*field = status
		return true
	}
	return false
}

// Returns the response body from a GET request to a URL
func getResponseBody(url string) []byte {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return body
}

// Checks if a string value is in a string array
func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
