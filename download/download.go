package download

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type Artifact struct {
	GroupId         string
	Id              string
	Version         string
	Extension       string
	Classifier      string
	IsSnapshot      bool
	SnapshotVersion string
	RepositoryUrl   string
	RepoUser        string
	RepoPassword    string
	Downloader      func(string, string, string) (*http.Response, error)
}

type Metadata struct {
	LatestVersion  string `xml:"versioning>latest"`
	ReleaseVersion string `xml:"versioning>release"`
	Timestamp      string `xml:"versioning>snapshot>timestamp"`
	BuildNumber    string `xml:"versioning>snapshot>buildNumber"`
}

func Download(name, dest, repo, filename, extension, user, pwd string) (string, error) {
	a, err := ParseName(name)
	if err != nil {
		return "", err
	}

	a.Downloader = httpGetCustom
	a.RepositoryUrl = repo
	a.RepoUser = user
	a.RepoPassword = pwd
	a.Extension = extension

	url, err := ArtifactUrl(a)
	if err != nil {
		return "", err
	}

	resp, err := a.Downloader(url, user, pwd)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if filename == "" {
		filename = FileName(a)
	}

	filepath := dest + "/" + filename

	out, err := os.Create(filepath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", err
	}

	return filepath, nil
}

func DownloadMavenMetadata(name, repo, user, pwd string) (*Metadata, error) {
	a, err := ParseName(name)
	if err != nil {
		return nil, err
	}

	a.Downloader = httpGetCustom
	a.RepositoryUrl = repo
	a.RepoUser = user
	a.RepoPassword = pwd

	return fetchMavenMetadata(a)
}

func httpGetCustom(url, user, pwd string) (*http.Response, error) {
	if user != "" && pwd != "" {
		client := &http.Client{}
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}
		req.SetBasicAuth(user, pwd)
		return client.Do(req)
	}

	return http.Get(url)
}
func ParseName(n string) (Artifact, error) {
	parts := strings.Split(n, ":")
	artifact := Artifact{}
	len := len(parts)
	if len >= 3 {
		artifact.GroupId = parts[0]
		artifact.Id = parts[1]
		artifact.Version = parts[len-1]

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

	return a.RepositoryUrl + artifactPath(a), nil
}

func fetchMavenMetadata(a Artifact) (*Metadata, error) {
	return downloadMavenMetadata(a, metadataPath(a))
}

func LatestSnapshotVersion(a Artifact) (string, error) {
	m, err := downloadMavenMetadata(a, snapshotMetadataPath(a))

	if err != nil {
		return "", nil
	}

	return fmt.Sprintf("%s-%s", m.Timestamp, m.BuildNumber), nil
}

func downloadMavenMetadata(a Artifact, metadataUrl string) (*Metadata, error) {
	resp, err := a.Downloader(metadataUrl, a.RepoUser, a.RepoPassword)
	if err != nil {
		return nil, err
	} else if resp.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("unable to fetch maven Metadata from %s Http statusCode: %d", metadataUrl, resp.StatusCode))
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	m := Metadata{}
	err = xml.Unmarshal(body, &m)

	if err != nil {
		return nil, nil
	}

	return &m, err
}

func artifactPath(a Artifact) string {
	return groupPath(a) + "/" + FileName(a)
}

func groupPath(a Artifact) string {
	parts := append(strings.Split(a.GroupId, "."), a.Id)
	if a.IsSnapshot {
		return strings.Join(append(parts, a.Version+"-SNAPSHOT"), "/")
	} else {
		return strings.Join(append(parts, a.Version), "/")
	}
}

func metadataPath(a Artifact) string {
	parts := append(strings.Split(a.GroupId, "."), a.Id)
	return a.RepositoryUrl + strings.Join(append(parts), "/") + "/maven-metadata.xml"
}

func snapshotMetadataPath(a Artifact) string {
	return a.RepositoryUrl + groupPath(a) + "/maven-metadata.xml"
}
