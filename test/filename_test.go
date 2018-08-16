package test

import (
	"testing"
	"go-mvn-dl/download"
)

func TestFilename(t *testing.T) {
	f := download.FileName(download.Artifact{GroupId: "org.apache.commons", Id:"commons-lang3", Version: "3.4"})
	expected := "commons-lang3-3.4.jar"
	if f != expected {
		t.Errorf("Incorrect filename. Expected:%s Got:%s", expected, f)
	}
}

func TestWarFilename(t *testing.T){
	f := download.FileName(download.Artifact{GroupId: "org.apache.commons", Id:"commons-lang3", Version: "3.4", Extension: "war"})
	expected := "commons-lang3-3.4.war"
	if f != expected {
		t.Errorf("Incorrect filename. Expected:%s Got:%s", expected, f)
	}
}

func TestClassifierFilename(t *testing.T){
	f := download.FileName(download.Artifact{GroupId: "org.apache.commons", Id:"commons-lang3", Version: "3.4", Classifier: "test"})
	expected := "commons-lang3-3.4-test.jar"
	if f != expected {
		t.Errorf("Incorrect filename. Expected:%s Got:%s", expected, f)
	}
}

func TestSnapshotVersionFilename(t *testing.T){
	f := download.FileName(download.Artifact{GroupId: "org.apache.commons", Id:"commons-lang3", Version: "3.4", IsSnapshot: true, SnapshotVersion:"123"})
	expected := "commons-lang3-3.4-123.jar"
	if f != expected {
		t.Errorf("Incorrect filename. Expected:%s Got:%s", expected, f)
	}
}

func TestSnapshotFilename(t *testing.T){
	f := download.FileName(download.Artifact{GroupId: "org.apache.commons", Id:"commons-lang3", Version: "3.4", IsSnapshot: true})
	expected := "commons-lang3-3.4-SNAPSHOT.jar"
	if f != expected {
		t.Errorf("Incorrect filename. Expected:%s Got:%s", expected, f)
	}
}