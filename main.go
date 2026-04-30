package main

import (
	"flag"
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

	currentMonth := flag.Bool("c", false, "Show current month")
	flag.Parse()

	client := api.NewClient(token, orgID)

	reports, err := utilities.ReportMinutes(client, *currentMonth)
	if err != nil {
		fmt.Printf("Failed generate report: %v\n", err)
	}

	if reports > 0 {
		fmt.Printf("Missing %.2f hours in selected period.\n", float64(reports)/60)
		return
	}

	if reports < 0 {
		fmt.Printf("You have some overtime %.2f.\n", float64(reports)/60*-1)
		return
	}

	fmt.Printf("You have correctly tracked your time")
}
