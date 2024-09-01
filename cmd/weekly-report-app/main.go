package main

import (
	"log"
	"wbapp/internal/weekly-report-app"
)

func main() {
	err := weekly_report_app.Run()
	if err != nil {
		log.Fatal(err.Error())
	}
}
