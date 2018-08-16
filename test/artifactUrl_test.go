package test

import (
	"testing"
	"go-mvn-dl/download"
	"net/http"
	"io/ioutil"
	"bytes"
	"strings"
)

func TestSnapshotArtifactUrl(t *testing.T) {
	a := download.Artifact {
		GroupId: "my.group",
		Id: "theFact",
		Version: "7",
		IsSnapshot: true,
		RepositoryUrl: "http://greatestRepo.com/",
		Downloader: mockHttpGet,
	}

	url, err := download.ArtifactUrl(a)
	
	if err != nil {
		t.Error(err)
	}

	expected := "http://greatestRepo.com/my/group/theFact/7-SNAPSHOT/theFact-7-123-456.jar"
	if url != expected {
		t.Errorf("Incorrect url. Expected: %s Got: %s", expected, url)
	}
}

func mockHttpGet(url, user, pwd string) (*http.Response, error) {
	body := []byte("<metadata><versioning><snapshot><timestamp>123</timestamp><buildNumber>456</buildNumber></snapshot></versioning></metadata>")
	resp := &http.Response {
		StatusCode: 200,
		Body: ioutil.NopCloser(bytes.NewReader(body)),
	}

	return resp, nil
}

func TestPackagingArtifactUrl(t *testing.T) {
	a := download.Artifact {
		GroupId: "org.apache.commons",
		Id: "commons-lang3",
		Version: "3.4",
		Extension: "pom",
	}

	url, err := download.ArtifactUrl(a)

	if err != nil {
		t.Error(err)
	}
	expected := "org/apache/commons/commons-lang3/3.4/commons-lang3-3.4.pom"
	if !strings.HasSuffix(url, expected) {
		t.Errorf("Incorrect url. Should end with: %s Got: %s", expected, url)
	}
}

func TestClassifierArtifactUrl(t *testing.T) {
	a := download.Artifact {
		GroupId: "org.apache.commons",
		Id: "commons-lang3",
		Version: "3.4",
		Classifier: "all",
	}

	url, err := download.ArtifactUrl(a)

	if err != nil {
		t.Error(err)
	}
	expected := "org/apache/commons/commons-lang3/3.4/commons-lang3-3.4-all.jar"
	if !strings.HasSuffix(url, expected) {
		t.Errorf("Incorrect url. Should end with: %s Got: %s", expected, url)
	}
}

func TestContainsDefaultRepository(t *testing.T) {
	a := download.Artifact {
		GroupId: "org.apache.commons",
		Id: "commons-lang3",
		Version: "3.4",
	}

	url, err := download.ArtifactUrl(a)

	if err != nil {
		t.Error(err)
	}
	expected := "https://repo1.maven.org/maven2/"
	if !strings.HasPrefix(url, expected) {
		t.Errorf("Incorrect url. Should start with: %s Got: %s", expected, url)
	}
}

func TestContainsCustomRepository(t *testing.T) {
	a := download.Artifact {
		GroupId: "org.apache.commons",
		Id: "commons-lang3",
		Version: "3.4",
		RepositoryUrl: "my.repo.com/",
	}

	url, err := download.ArtifactUrl(a)

	if err != nil {
		t.Error(err)
	}
	expected := "my.repo.com/"
	if !strings.HasPrefix(url, expected) {
		t.Errorf("Incorrect url. Should start with: %s Got: %s", expected, url)
	}
}