// Accept requests to combine PDF documents located on Amazon S3 into
// a single document, upload back to S3, and notify a callback URL.
// TODO fail fast if S3 connection is not usable.
// TODO need a way to enable auth in responses
package main

import (
	"flag"
	"launchpad.net/goamz/aws"
	"os"
	"pdfcombiner/server"
)

var (
	port       string
	serverMode bool
	bucket     string
	employerId int
)

func init() {
	flag.BoolVar(&serverMode, "server", false, "run in server mode")
	flag.StringVar(&port, "port", "8080", "port to listen on for server mode")
	flag.StringVar(&bucket, "bucket", "", "bucket name to use in standalone mode")
	flag.IntVar(&employerId, "employer", 0, "id number of the employer to combine for in standalone mode")
	flag.Parse()
	flag.Usage = func() {
		println("Usage: pdfcombiner [OPTS] [FILE...]")
		flag.PrintDefaults()
		os.Exit(1)
	}
}

func main() {
	verifyAws()
	switch {
	case serverMode:
		server.ListenOn(port)
	default:
		combineSynchronously()
	}
}

// Combine the requested files and return the status to standard out.
func combineSynchronously() {
	if flag.NArg() < 1 {
		println("Cannot start in standalone mode with no files to combine.")
		flag.Usage()
	}
	pdfFiles := flag.Args()
	println(pdfFiles)
}

// Verify AWS credentials are set in the environment:
// - AWS_ACCESS_KEY_ID
// - AWS_SECRET_ACCESS_KEY
func verifyAws() {
	_, err := aws.EnvAuth()
	if err != nil {
		panic(err)
	}
}
