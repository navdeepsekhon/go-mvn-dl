package download

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestFetchMavenMetadata(t *testing.T) {
	a := Artifact{
		GroupId:       "ch.qos.logback",
		Id:            "logback-classic",
		Version:       "1.2.11",
		RepositoryUrl: "http://greatestRepo.com/",
		Downloader:    mockHttpGetMetadata,
	}
	metadata, err := fetchMavenMetadata(a)
	if err != nil {
		t.Errorf("Error fetching maven-Metadata.xml.  Error: %s", err)
	}
	if got, want := metadata.LatestVersion, "2.11.0"; got != want {
		t.Errorf("Wrong latest version. Expected: %s, got: %s", want, got)
	}
	if got, want := metadata.ReleaseVersion, "2.10.0"; got != want {
		t.Errorf("Wrong release version. Expected: %s, got: %s", want, got)
	}
}

func TestMetadataUrl(t *testing.T) {
	a := Artifact{
		GroupId:       "ch.qos.logback",
		Id:            "logback-classic",
		Version:       "1.2.11",
		RepositoryUrl: "http://greatestRepo.com/",
	}

	if got, want := metadataPath(a), "http://greatestRepo.com/ch/qos/logback/logback-classic/maven-metadata.xml"; got != want {
		t.Errorf("Wrong Metadata url. Expected: %s, got: %s", want, got)
	}
}

func TestSnapshotMetadataUrl(t *testing.T) {
	a := Artifact{
		GroupId:       "ch.qos.logback",
		Id:            "logback-classic",
		Version:       "1.2.11",
		IsSnapshot:    true,
		RepositoryUrl: "http://greatestRepo.com/",
	}

	if got, want := snapshotMetadataPath(a), "http://greatestRepo.com/ch/qos/logback/logback-classic/1.2.11-SNAPSHOT/maven-metadata.xml"; got != want {
		t.Errorf("Wrong Metadata url. Expected: %s, got: %s", want, got)
	}

}

func mockHttpGetMetadata(url, user, pwd string) (*http.Response, error) {
	body := []byte("<Metadata><versioning><latest>2.11.0</latest><release>2.10.0</release></versioning></Metadata>")
	resp := &http.Response{
		StatusCode: 200,
		Body:       ioutil.NopCloser(bytes.NewReader(body)),
	}

	return resp, nil
}
