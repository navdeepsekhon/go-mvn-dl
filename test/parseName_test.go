package test

import (
	"testing"
	"go-mvn-dl/download"
)

func TestParseName(t *testing.T) {
	a, err := download.ParseName("org.apache.commons:commons-lang3:3.4")
	if err != nil { t.Error(err) }
	verifyArtifact(a, "org.apache.commons", "commons-lang3", "3.4", "", false, t)
}

func TestParseNameWithExtension(t *testing.T) {
	a, err := download.ParseName("org.apache.commons:commons-lang3:war:3.4")
	if err != nil { t.Error(err) }
	verifyArtifact(a, "org.apache.commons", "commons-lang3", "3.4", "war", false, t)
}

func TestParseNameWithSnapshot(t *testing.T) {
	a, err := download.ParseName("org.apache.commons:commons-lang3:3.4-SNAPSHOT")
	if err != nil { t.Error(err) }
	verifyArtifact(a, "org.apache.commons", "commons-lang3", "3.4", "", true, t)
}

func verifyArtifact(a download.Artifact, gId, id, v, e string, ss bool, t *testing.T) {
	if a.GroupId != gId {
		t.Errorf("Incorrect groupId. Expected: %s, Got:%s", gId, a.GroupId)
	}

	if a.Id != id {
		t.Errorf("Incorrect Id. Expected: %s, Got:%s", id, a.Id)
	}

	if a.Version != v {
		t.Errorf("Incorrect version. Expected: %s, Got:%s", v, a.Version)
	}

	if a.Extension != e {
		t.Errorf("Incorrect Extension. Expected: %s, Got:%s", e, a.Extension)
	}

	if a.IsSnapshot != ss {
		t.Errorf("Incorrect IsSnapshot. Expected: %t, Got:%t", ss, a.IsSnapshot)
	}
}