package main

import (
	"fmt"
	"github.com/go-co-op/gocron"
	"github.com/tseech/bronco-build-tracker/internal/settings"
	"github.com/tseech/bronco-build-tracker/internal/state"
	"github.com/tseech/bronco-build-tracker/internal/tracker"
	"time"
)

func main() {
	settings := settings.ReadSettings()

	// Either run this once and quit or run a cron job
	if settings.RunOnce {
		run(settings)
	} else {
		s := gocron.NewScheduler(time.UTC)
		s.Every(settings.Interval).Do(
			func() {
				run(settings)
				_, nextTime := s.NextRun()
				fmt.Println("Next run at " + nextTime.Local().Format(time.Kitchen))
			})
		s.StartBlocking()
	}
}

// Main execution that can be run once or multiple times
func run(settings settings.Settings) {
	fmt.Println("Checking pizza tracker...")
	tracker.CheckStatus(settings)
	fmt.Println("Checking backdoor tracker...")
	tracker.CheckBackdoorTracker(settings)
	fmt.Println("Current state:")
	state.PrintState()
}
