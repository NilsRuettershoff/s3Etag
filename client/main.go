package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/CologneBroadcastingCenter/s3Etag"
)

func main() {
	var path string
	flag.StringVar(&path, "path", "", "path to file for checksum")
	flag.Parse()
	etag, err := s3Etag.CalculateLocalETag(path, 5)
	if err != nil {
		log.Printf("err: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(etag)
}
