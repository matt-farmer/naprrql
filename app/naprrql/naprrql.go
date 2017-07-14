package main

// Alternative naprr server offering graphql data feeds and KV store backend
// to try and elimnate issues with tcp, memory and disk access
// on Win32 that hamper all other deployments so far.

import (
	"flag"
	"runtime"

	"github.com/matt-farmer/naprrql"
)

var ingest = flag.Bool("ingest", false, "forces re-load of data from results file.")

// var webonly = flag.Bool("webonly", false, "just launch web data explorer")

func main() {

	flag.Parse()

	// ingest results data
	if *ingest {
		naprrql.IngestResultsFile("master_nap.xml.zip")
	}

	// launch sif-ql service
	// wait group?
	go naprrql.RunQLServer()

	// launch ui service

	runtime.Goexit()

}
