package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/aarbanas/productive-time-tracker/api"
	"github.com/aarbanas/productive-time-tracker/appinit"
	"github.com/aarbanas/productive-time-tracker/utilities"
)

const specURL = "https://developer.productive.io/reference/download_spec"
const specPath = "api/productive.yaml"

func main() {
	currentMonth := flag.Bool("c", false, "Show current month")
	downloadSpec := flag.Bool("download-spec", false, "Download the latest Productive OpenAPI spec and save it to "+specPath)
	flag.Parse()

	if *downloadSpec {
		if err := fetchSpec(); err != nil {
			fmt.Printf("Failed to download spec: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Spec saved to %s\n", specPath)
		return
	}

	token, orgID, err := appinit.LoadCredentials()
	if err != nil {
		fmt.Printf("Failed to initialize client: %v\n", err)
		os.Exit(1)
	}

	client, err := api.NewProductiveClient(token, orgID)
	if err != nil {
		fmt.Printf("Failed to create client: %v\n", err)
		os.Exit(1)
	}

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

func fetchSpec() error {
	resp, err := http.Get(specURL) //#nosec G107 -- URL is a hardcoded constant
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return os.WriteFile(specPath, data, 0600)
}