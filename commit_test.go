package fscommit

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func skipIfError(t *testing.T, err error) {
	if err == nil {
		return
	}
	t.Skip(err)
}

func failIfError(t *testing.T, err error) {
	if err == nil {
		return
	}
	t.Error(err)
}

func Test_Commit_Execute_All(t *testing.T) {
	// Set up a little test folder
	name, err := ioutil.TempDir(os.TempDir(), "testing")
	skipIfError(t, err)
	defer os.RemoveAll(name)
	skipIfError(t, ioutil.WriteFile(filepath.Join(name, "testfile"), []byte{0x13, 0x37, 0xBA, 0xBE, 0xCA, 0xFE, 0xDE, 0xAD}, 0644))
	skipIfError(t, ioutil.WriteFile(filepath.Join(name, "testfile2"), []byte{0x13, 0x37, 0xBA, 0xBE, 0xCA, 0xFE, 0xDE, 0xAD}, 0644))
	skipIfError(t, os.Mkdir(filepath.Join(name, "testfolder"), 0644))

	commit := Commit{
		Rename(filepath.Join(name, "testfile2"), filepath.Join(name, "testy")),
		RemoveAll(filepath.Join(name, "testfile")),
		RemoveAll(filepath.Join(name, "testfolder")),
	}
	result := commit.Execute()
	failIfError(t, result.CommitError)
	if len(result.FinalizeErrors) > 0 {
		t.Error("expected result.FinalizeErrors to be empty")
	}
}

func Test_Commit_Execute_Fail(t *testing.T) {
	// Set up a little test folder
	name, err := ioutil.TempDir(os.TempDir(), "testing")
	skipIfError(t, err)
	defer os.RemoveAll(name)
	skipIfError(t, ioutil.WriteFile(filepath.Join(name, "testfile"), []byte{0x13, 0x37, 0xBA, 0xBE, 0xCA, 0xFE, 0xDE, 0xAD}, 0644))
	skipIfError(t, ioutil.WriteFile(filepath.Join(name, "testfile2"), []byte{0x13, 0x37, 0xBA, 0xBE, 0xCA, 0xFE, 0xDE, 0xAD}, 0644))
	skipIfError(t, os.Mkdir(filepath.Join(name, "testfolder"), 0644))

	commit := Commit{
		Rename(filepath.Join(name, "testfile2"), filepath.Join(name, "testy")),
		RemoveAll(filepath.Join(name, "testfile")),
		RemoveAll(filepath.Join(name, "testfolder")),
		Rename(filepath.Join(name, "testfolder"), filepath.Join(name, "testy")),
	}
	result := commit.Execute()
	if result.CommitError == nil {
		t.Error("expected result.CommitError to not be nil")
	}
	if len(result.RevertErrors) > 0 {
		t.Error("expected result.RevertErrors to be empty")
	}
	if !os.IsNotExist(result.CommitError) {
		t.Error("expected result.CommitError to be a NotExist error")
	}
	if _, err := os.Stat(filepath.Join(name, "testfolder")); err != nil {
		t.Error("expected testfolder to still be accessible")
	}
	if _, err := os.Stat(filepath.Join(name, "testfile")); err != nil {
		t.Error("expected testfile to still be accessible")
	}
	if _, err := os.Stat(filepath.Join(name, "testfile2")); err != nil {
		t.Error("expected testfile2 to still be accessible")
	}
}

func Test_RemoveAll_Execute(t *testing.T) {
	// Set up a little test folder
	name, err := ioutil.TempDir(os.TempDir(), "testing")
	skipIfError(t, err)
	defer os.RemoveAll(name)

	failIfError(t, RemoveAll(name).Execute())
}

func Test_Rename_Execute(t *testing.T) {
	// Set up a little test folder
	name, err := ioutil.TempDir(os.TempDir(), "testing")
	skipIfError(t, err)
	defer os.RemoveAll(name)

	name2, err := tmpPath(os.TempDir(), "testing")
	skipIfError(t, err)
	failIfError(t, Rename(name, name2).Execute())
	defer os.RemoveAll(name2)
}
