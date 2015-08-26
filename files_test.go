package main

import "os"
import "testing"
import "io/ioutil"

func TestUnzip(t *testing.T) {
	dataDir := "testdata"
	archive := dataDir + "/test.zip"
	testDir := dataDir + "/some-dir"
	testNestedDir := dataDir + "/some-dir/nested"
	testFile := testDir + "/x"
	testFileWithContent := testNestedDir + "/y"
	testContent := `y
`

	os.RemoveAll(testDir)
	_, err := os.Open(testDir)
	if err == nil {
		t.Errorf("didn't expect %v to exist yet.", testDir)
		return
	}

	unzip(archive, dataDir)

	files := []string{testDir, testNestedDir, testFile, testFileWithContent}
	for _, f := range files {
		_, err = os.Open(f)
		if err != nil {
			t.Errorf("expected %v to exist", f)
		}
	}

	content, err := ioutil.ReadFile(testFileWithContent)
	if err != nil || testContent != string(content) {
		t.Errorf("expected to find [%v] in %v. Got: [%v]", testContent, testFileWithContent, string(content))
	}
}
