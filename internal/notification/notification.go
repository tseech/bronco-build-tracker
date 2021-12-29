package notification

import (
	"fmt"
	"github.com/tseech/bronco-build-tracker/internal/settings"
	"io/ioutil"
	"net/http"
	"net/url"
)

// Sends a notification to a device based on the settings that are populated
// If the configuration allows it, multiple notifications may be sent
func Notify(srcUrl string, message string, settings settings.Settings) {
	// If Pushover info is provided send the notification via Pushover
	if settings.PushoverToken != "-" && settings.PushoverUser != "-" && settings.PushoverToken != "" && settings.PushoverUser != "" {
		pushoverNotify(srcUrl, message, settings)
	}
	// If Textbelt info is provided, send the notification via Textbelt
	if settings.PhoneNumber != "-" && settings.TextbeltKey != "-" && settings.PhoneNumber != "" && settings.TextbeltKey != "" {
		textbeltNotify(srcUrl, message, settings)
	}
}

// Sends a notification via Pushover.net
func pushoverNotify(srcUrl string, message string, settings settings.Settings) {
	values := url.Values{
		"token":     {settings.PushoverToken},
		"user":      {settings.PushoverUser},
		"message":   {message},
		"url":       {srcUrl},
		"url_title": {"Link"},
	}
	response, _ := http.PostForm("https://api.pushover.net/1/messages.json", values)
	body, _ := ioutil.ReadAll(response.Body)
	fmt.Println(string(body))
}

// Sends a notification via textbelt.com
func textbeltNotify(srcUrl string, message string, settings settings.Settings) {
	values := url.Values{
		"phone":   {settings.PhoneNumber},
		"message": {message},
		"key":     {settings.TextbeltKey},
	}

	response, _ := http.PostForm("https://textbelt.com/text", values)
	body, _ := ioutil.ReadAll(response.Body)
	fmt.Println(string(body))
}
