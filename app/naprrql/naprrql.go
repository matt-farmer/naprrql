package main

// Alternative naprr server offering graphql data feeds and KV store backend
// to try and elimnate issues with tcp, memory and disk access
// on Win32 that hamper all other deployments so far.

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"

	"github.com/matt-farmer/naprrql"
)

var ingest = flag.Bool("ingest", false, "Loads data from results file. Exisitng data is overwritten.")

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

	// ingest results data, and exit to save memory
	if *ingest {
		clearDBWorkingDirectory()
		resultsFiles := parseResultsFileDirectory()
		for _, resultsFile := range resultsFiles {
			naprrql.IngestResultsFile(resultsFile)
		}
		closeDB()
		os.Exit(1)
	} else {
		// launch sif-ql service
		go naprrql.RunQLServer()
		// launch ui service
	}

	// wait for shutdown
	for {
		runtime.Gosched()
	}

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

func createDBWorkingDirectory() {
	err := os.Mkdir("kvs", os.ModePerm)
	if !os.IsExist(err) && err != nil {
		log.Fatalln("Error trying to create datastore working directory: ", err)
	}

}
