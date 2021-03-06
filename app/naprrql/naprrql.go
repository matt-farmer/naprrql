package main

// Alternative naprr server offering graphql data feeds and KV store backend
// to try and elimnate issues with tcp, memory and disk access
// on Win32 that hamper all other deployments so far.

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"

	"github.com/matt-farmer/naprrql"
)

var ingest = flag.Bool("ingest", false, "Loads data from results file. Exisitng data is overwritten.")
var report = flag.Bool("report", false, "Creates .csv reports. Existing reports are overwritten")

func main() {

	flag.Parse()

	// shutdown handler
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		closeDB()
		os.Exit(1)
	}()

	// ingest results data, rebuild reports, and exit to save memory
	if *ingest {
		ingestData()
		startWebServer(true)
		writeReports()
		// shut down
		closeDB()
		os.Exit(1)
	}

	// create the csv reports
	if *report {
		// launch web-server
		startWebServer(true)
		// if requested regenerate reports
		if *report {
			writeReports()
		}
		// shut down
		closeDB()
		os.Exit(1)
	}

	// otherwise just start the webserver
	startWebServer(false)

	// wait for shutdown
	for {
		runtime.Gosched()
	}

}

//
// iterate & load any r/r data files provided
//
func ingestData() {
	// ingest the data
	log.Println("invoking data ingest...")
	clearDBWorkingDirectory()
	resultsFiles := parseResultsFileDirectory()
	for _, resultsFile := range resultsFiles {
		naprrql.IngestResultsFile(resultsFile)
	}
}

//
// launch the webserver
//
func startWebServer(silent bool) {
	go naprrql.RunQLServer()
	if !silent {
		fmt.Printf("\n\nBrowse to follwing locations:\n")
		fmt.Printf("\n\thttp://localhost:1329/ui\n\n for qa report user interface\n")
		fmt.Printf("\n\thttp://localhost:1329/sifql\n\n for data explorer\n\n")
	}

}

//
// create .csv reports
//
func writeReports() {
	clearReportsDirectory()
	log.Println("generating reports...")
	naprrql.GenerateReports()
	log.Println("reports generated...")
}

//
// look for results data files
//
func parseResultsFileDirectory() []string {

	files := make([]string, 0)

	zipFiles, _ := filepath.Glob("./in/*.zip")
	xmlFiles, _ := filepath.Glob("./in/*.xml")

	files = append(files, zipFiles...)
	files = append(files, xmlFiles...)
	if len(files) == 0 {
		log.Fatalln("No results data *.xml.zip or *.xml files found in input folder /in.")
	}

	return files

}

//
// ensure clean shutdown of data store
//
func closeDB() {
	log.Println("Closing datastore...")
	naprrql.GetDB().Close()
	log.Println("Datastore closed.")
}

//
// remove working files of datastore
//
func clearDBWorkingDirectory() {

	// remove existing logs and recreate the directory
	err := os.RemoveAll("kvs")
	if err != nil {
		log.Println("Error trying to reset datastore working directory: ", err)
	}
	createDBWorkingDirectory()
}

//
// remove reports working directory
//
func clearReportsDirectory() {
	// remove existing logs and recreate the directory
	err := os.RemoveAll("out")
	if err != nil {
		log.Println("Error trying to reset reports directory: ", err)
	}
	createReportsDirectory()

}

//
// create folder for .csv reports
//
func createReportsDirectory() {
	err := os.Mkdir("out", os.ModePerm)
	if !os.IsExist(err) && err != nil {
		log.Fatalln("Error trying to create reports directory: ", err)
	}

}

//
// create folder for datastore
//
func createDBWorkingDirectory() {
	err := os.Mkdir("kvs", os.ModePerm)
	if !os.IsExist(err) && err != nil {
		log.Fatalln("Error trying to create datastore working directory: ", err)
	}

}
