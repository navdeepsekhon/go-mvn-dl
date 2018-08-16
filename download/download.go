package download

import (
	"fmt"
	"strings"
	"net/http"
	"io/ioutil"
	"encoding/xml"
	"io"
	"os"
	"errors"
)

type Artifact struct {
	GroupId	string
	Id	string
	Version	string
	Extension	string
	Classifier	string
	IsSnapshot bool
	SnapshotVersion	string
	RepositoryUrl	string
	Downloader	func (string) (*http.Response, error)
}

type metadata struct {
	Timestamp	string	`xml:"versioning>snapshot>timestamp"`
	BuildNumber	string	`xml:"versioning>snapshot>buildNumber"`
}

func Download(name, dest, repo, filename string) (string, error) {
	a, err := ParseName(name)
	if err != nil { return "", err}

	a.Downloader = http.Get
	a.RepositoryUrl = repo

	url, err := ArtifactUrl(a)
	if err != nil { return "", err}

	resp, err := http.Get(url)
	if err != nil { return "", err}
	defer resp.Body.Close()

	if filename == "" {
		filename = FileName(a)
	}

	filepath := dest + "/" + filename

	out, err := os.Create(filepath)
	if err != nil { return "", err}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil { return "", err}

	return filepath, nil
}

func ParseName(n string) (Artifact, error) {
	parts := strings.Split(n, ":")
	artifact := Artifact{}
	len := len(parts)
	if len >= 3 {
		artifact.GroupId = parts[0]
		artifact.Id = parts[1]
		artifact.Version = parts[len - 1]

		if len > 3 {
			artifact.Extension = parts[2]
		}
		if len > 4 {
			artifact.Classifier = parts[3]
		}

		if strings.HasSuffix(artifact.Version, "-SNAPSHOT") {
			artifact.IsSnapshot = true
			artifact.Version = strings.Trim(artifact.Version, "-SNAPSHOT")
		}

		return artifact, nil
	}

	return artifact, errors.New("invalid package name. Try groupId:artifactId:version")
}

func FileName(a Artifact) string {
	ext := "jar"
	if a.Extension != "" {
		ext = a.Extension
	}

	v := a.Version

	if a.IsSnapshot {
		if a.SnapshotVersion != "" {
			v += "-" + a.SnapshotVersion
		} else {
			v += "-SNAPSHOT"
		}
	}

	if a.Classifier != "" {
		return fmt.Sprintf("%s-%s-%s.%s", a.Id, v, a.Classifier, ext)
	} else {
		return fmt.Sprintf("%s-%s.%s", a.Id, v, ext)
	}
}

func ArtifactUrl(a Artifact) (string, error) {
	if a.RepositoryUrl == "" {
		a.RepositoryUrl = "https://repo1.maven.org/maven2/"
	}

	if a.IsSnapshot {
		var err error
		a.SnapshotVersion, err = LatestSnapshotVersion(a)
		if err != nil {
			return "", err
		}
	}

	return a.RepositoryUrl + ArtifactPath(a), nil
}

func LatestSnapshotVersion(a Artifact) (string, error) {
	metadataUrl := a.RepositoryUrl + GroupPath(a) + "/maven-metadata.xml"
	resp, err := a.Downloader(metadataUrl)
	if err != nil {
		return "", err
	} else if resp.StatusCode != 200 {
		return "", errors.New(fmt.Sprintf("unable to fetch maven metadata from %s Http statusCode: %d", metadataUrl, resp.StatusCode))
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	m := metadata{}
	err = xml.Unmarshal(body, &m)

	if err != nil {
		return "", nil
	}

	return fmt.Sprintf("%s-%s", m.Timestamp, m.BuildNumber), nil
}

func ArtifactPath(a Artifact) string {
	return GroupPath(a) + "/" + FileName(a)
}

func GroupPath(a Artifact) string {
	parts := append(strings.Split(a.GroupId, "."), a.Id)
	if a.IsSnapshot {
		return strings.Join(append(parts, a.Version + "-SNAPSHOT"), "/")
	} else {
		return strings.Join(append(parts, a.Version), "/")
	}
}