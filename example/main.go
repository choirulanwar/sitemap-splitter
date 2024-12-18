package main

import (
	"fmt"
	"log"

	sitemapsplitter "github.com/choirulanwar/sitemap-splitter"
)

func main() {
	// Create a new SitemapSplitter instance
	// Parameters:
	// 1. Path to the sitemap file to be split
	// 2. Maximum number of URLs per file (e.g., 10 URLs per file)
	splitter, err := sitemapsplitter.NewSitemapSplitter("./example/sitemap.xml", 10)
	if err != nil {
		log.Fatalf("Error creating splitter: %v", err)
	}

	// Perform the splitting process
	if err := splitter.Split(); err != nil {
		log.Fatalf("Error splitting sitemap: %v", err)
	}

	fmt.Println("Sitemap successfully split!")
}
