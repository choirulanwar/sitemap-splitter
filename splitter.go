package sitemapsplitter

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

// URL represents a single URL entry in the sitemap
type URL struct {
	XMLName    xml.Name `xml:"url"`
	Loc        string   `xml:"loc"`
	LastMod    string   `xml:"lastmod,omitempty"`
	ChangeFreq string   `xml:"changefreq,omitempty"`
	Priority   string   `xml:"priority,omitempty"`
}

// URLSet represents the root element of a sitemap
type URLSet struct {
	XMLName xml.Name `xml:"urlset"`
	XMLNS   string   `xml:"xmlns,attr"`
	XHTML   string   `xml:"xmlns:xhtml,attr"`
	URLs    []URL    `xml:"url"`
}

// SitemapIndex represents the root element of a sitemap index
type SitemapIndex struct {
	XMLName  xml.Name  `xml:"sitemapindex"`
	XMLNS    string    `xml:"xmlns,attr"`
	Sitemaps []Sitemap `xml:"sitemap"`
}

// Sitemap represents a single sitemap entry in the index
type Sitemap struct {
	XMLName xml.Name `xml:"sitemap"`
	Loc     string   `xml:"loc"`
	LastMod string   `xml:"lastmod"`
}

// SitemapSplitter handles splitting large sitemaps into smaller ones
type SitemapSplitter struct {
	path  string // Absolute or relative path to sitemap file
	limit int    // Maximum number of URLs per sitemap file
}

// NewSitemapSplitter creates a new SitemapSplitter instance
func NewSitemapSplitter(path string, limit int) (*SitemapSplitter, error) {
	if path == "" {
		return nil, fmt.Errorf("sitemap path is required")
	}
	if limit <= 0 {
		return nil, fmt.Errorf("limit must be greater than 0")
	}

	return &SitemapSplitter{
		path:  path,
		limit: limit,
	}, nil
}

// Split reads the sitemap and splits it into multiple files
func (s *SitemapSplitter) Split() error {
	// Read and parse the original sitemap
	data, err := ioutil.ReadFile(s.path)
	if err != nil {
		return fmt.Errorf("error reading sitemap file: %v", err)
	}

	var urlset URLSet
	if err := xml.Unmarshal(data, &urlset); err != nil {
		return fmt.Errorf("error parsing XML: %v", err)
	}

	if len(urlset.URLs) == 0 {
		return fmt.Errorf("no URLs found in sitemap")
	}

	// Get directory and filename from path
	dir := filepath.Dir(s.path)
	filename := filepath.Base(s.path)
	baseFilename := filename[:len(filename)-len(filepath.Ext(filename))]

	var sitemapFiles []struct {
		BaseURL     string
		Name        string
		LastModDate string
	}

	// Split URLs into chunks
	for i := 0; i*s.limit < len(urlset.URLs); i++ {
		start := i * s.limit
		end := (i + 1) * s.limit
		if end > len(urlset.URLs) {
			end = len(urlset.URLs)
		}

		chunk := urlset.URLs[start:end]
		if len(chunk) == 0 {
			break
		}

		// Create new URLSet for this chunk
		newURLSet := URLSet{
			XMLNS: "http://www.sitemaps.org/schemas/sitemap/0.9",
			XHTML: "http://www.w3.org/1999/xhtml",
			URLs:  chunk,
		}

		// Generate sitemap name
		sitemapName := fmt.Sprintf("%s-%d.xml", baseFilename, i+1)

		// Get base URL from the last URL in chunk
		lastURL := chunk[len(chunk)-1]
		parsedURL, err := url.Parse(lastURL.Loc)
		if err != nil {
			return fmt.Errorf("error parsing URL: %v", err)
		}
		baseURL := fmt.Sprintf("%s://%s/", parsedURL.Scheme, parsedURL.Host)

		// Get last modification date
		lastMod := lastURL.LastMod
		if lastMod == "" {
			lastMod = time.Now().Format(time.RFC3339)
		}

		sitemapFiles = append(sitemapFiles, struct {
			BaseURL     string
			Name        string
			LastModDate string
		}{
			BaseURL:     baseURL,
			Name:        sitemapName,
			LastModDate: lastMod,
		})

		// Write sitemap file
		xmlData, err := xml.MarshalIndent(newURLSet, "", "  ")
		if err != nil {
			return fmt.Errorf("error marshaling XML: %v", err)
		}

		xmlHeader := []byte(xml.Header)
		xmlData = append(xmlHeader, xmlData...)

		outputPath := filepath.Join(dir, sitemapName)
		if err := os.WriteFile(outputPath, xmlData, 0644); err != nil {
			return fmt.Errorf("error writing sitemap file: %v", err)
		}
	}

	// Create sitemap index
	sitemapIndex := SitemapIndex{
		XMLNS: "http://www.sitemaps.org/schemas/sitemap/0.9",
	}

	for _, file := range sitemapFiles {
		sitemapIndex.Sitemaps = append(sitemapIndex.Sitemaps, Sitemap{
			Loc:     file.BaseURL + file.Name,
			LastMod: file.LastModDate,
		})
	}

	// Write sitemap index
	xmlData, err := xml.MarshalIndent(sitemapIndex, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling XML: %v", err)
	}

	xmlHeader := []byte(xml.Header)
	xmlData = append(xmlHeader, xmlData...)

	indexPath := filepath.Join(dir, "sitemap-index.xml")
	if err := os.WriteFile(indexPath, xmlData, 0644); err != nil {
		return fmt.Errorf("error writing sitemap index: %v", err)
	}

	return nil
}
