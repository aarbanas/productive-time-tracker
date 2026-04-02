package main

import (
	"fmt"
	"os"

	"github.com/aarbanas/productive-time-tracker/api"
	"github.com/aarbanas/productive-time-tracker/appinit"
	"github.com/aarbanas/productive-time-tracker/utilities"
)

func main() {
	token, orgID, err := appinit.LoadCredentials()
	if err != nil {
		fmt.Printf("Failed to initialize client: %v\n", err)
		os.Exit(1)
	}

	client := api.NewClient(token, orgID)

	totalAbsenceMinutes, err := utilities.AbsenceMinutes(client)
	if err != nil {
		fmt.Printf("Failed to calculate absence minutes: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Total absence hours: %d\n", totalAbsenceMinutes/60)

	totalTimeEntriesMinutes, err := utilities.TimeEntriesMinutes(client)
	if err != nil {
		fmt.Printf("Failed to calculate time entries minutes: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Total time entries hours: %d\n", totalTimeEntriesMinutes/60)

	totalMinutes := totalAbsenceMinutes + totalTimeEntriesMinutes
	requiredMinutes := utilities.RequiredWorkingMinutesPreviousMonth()
	fmt.Printf("Required hours to track for previous month: %d\n", requiredMinutes/60)
	fmt.Printf("Total hours tracked: %d\n", totalMinutes/60)

	if totalMinutes < requiredMinutes {
		fmt.Printf("You are %d minutes behind schedule.\n", requiredMinutes-totalMinutes)
	} else {
		fmt.Println("Great job! You are on track!")
	}
}
