package tracker

import (
	"encoding/json"
	"fmt"
	"github.com/tseech/bronco-build-tracker/internal/notification"
	"github.com/tseech/bronco-build-tracker/internal/settings"
	"github.com/tseech/bronco-build-tracker/internal/state"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const ReservationUrl = "https://www.authagent.ford.com/api/secure-purchase/gep/USA/reservations/%s?lang=en_us&includeBSLData=false"

type AuthResponse struct {
	AccessToken string `json:"access_token"`
}

type Reservation struct {
	Status  string             `json:"status"`
	Entries []ReservationEntry `json:"entries"`
}

type ReservationEntry struct {
	BuiltWeek         string `json:"builtWeek"`
	DisplayStatus     string `json:"displayStatus"`
	VehicleStatusCode string `json:"vehicleStatusCode"`
}

func CheckBackdoorTracker(settings settings.Settings) {
	currentState := state.GetState()

	if settings.ReservationId == "" || settings.ReservationId == "-" || settings.RefreshToken == "" || settings.RefreshToken == "-" {
		currentState.BackdoorTrackerStatus = "Not configured - Reservation ID and Refresh Token required"
		state.SaveState(currentState)
		return
	}

	// Get the access token
	accessToken, err := getAccessToken(settings)
	if err != nil {
		currentState.BackdoorTrackerStatus = err.Error()
		state.SaveState(currentState)
		return
	}

	// Get the reservation info
	reservation, err := getReservation(settings, accessToken)
	if err != nil {
		currentState.BackdoorTrackerStatus = err.Error()
		state.SaveState(currentState)
		return
	}

	// Check for changes and update the currentState
	changed := false
	if len(reservation.Entries) > 0 {
		changed = update(&reservation.Status, &currentState.Status) || changed
		changed = update(&reservation.Entries[0].BuiltWeek, &currentState.BuiltWeek) || changed
		changed = update(&reservation.Entries[0].DisplayStatus, &currentState.DisplayStatus) || changed
		changed = update(&reservation.Entries[0].VehicleStatusCode, &currentState.VehicleStatusCode) || changed
	}

	// Report changes if there are any
	if changed {
		notification.Notify(getReservationUrl(settings), fmt.Sprintf("Backdoor tracker change detected: %s", reservation), settings)
	}

	// Save the currentState and mark as successful
	currentState.BackdoorTrackerStatus = "Success"
	state.SaveState(currentState)
}

// Gets an access token to access the tracker data using the refresh token
func getAccessToken(settings settings.Settings) (string, error) {
	values := url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {settings.RefreshToken},
	}

	request, err := http.NewRequest("POST", "https://api.mps.ford.com/api/oauth2/v1/token", strings.NewReader(values.Encode()))
	if err != nil {
		return "", err
	}
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Add("application-id", "e7ec653e-113e-41e9-83ab-7e54719e2977")
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	var authResponse AuthResponse
	json.Unmarshal(body, &authResponse)
	return authResponse.AccessToken, nil
}

// Gets the reservation data using the access token for authentication
func getReservation(settings settings.Settings, accessToken string) (Reservation, error) {
	var reservation Reservation
	request, err := http.NewRequest("GET", getReservationUrl(settings), nil)
	if err != nil {
		return reservation, err
	}
	request.Header.Add("X-Identity-Authorization", fmt.Sprintf("Bearer %s", accessToken))
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return reservation, err
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return reservation, err
	}

	json.Unmarshal(body, &reservation)
	return reservation, nil
}

func getReservationUrl(settings settings.Settings) string {
	return fmt.Sprintf(ReservationUrl, settings.ReservationId)
}

// Updates a destination string from a source string and returns true if the destination has changed
func update(src *string, dest *string) bool {
	if *src == *dest || *src == "" {
		return false
	}
	*dest = *src
	return true
}
