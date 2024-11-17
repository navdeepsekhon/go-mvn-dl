package main

import (
	"go-mvn-dl/download"
	"log"
)

func main() {
	metadata, err := download.DownloadMavenMetadata(
		"ch.qos.logback:logback-classic:1",
		"https://repo1.maven.org/maven2/",
		"",
		"")
	if err != nil {
		log.Printf("Error downloading maven-metadata.xml.  Error: %s", err)
	}
	log.Printf("releaseVersion=%s, latestVersion=%s", metadata.ReleaseVersion, metadata.LatestVersion)

	file, err := download.Download(
		"ch.qos.logback:logback-classic:1.2.11",
		".",
		"https://repo1.maven.org/maven2/",
		"",
		"",
		"",
		"")

	if err != nil {
		log.Printf("Error downlaoding.  Error: %s", err)
	}
	log.Printf("file=%s", file)
}
