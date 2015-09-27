package pad

import "os"
import "net/http"
import "testing"
import "io/ioutil"

func TestInstallNotReDownloading(t *testing.T) {
	var df bool
	td := func(r string) (*http.Response, error) {
		df = true
		return &http.Response{}, nil
	}

	local := "gerd"
	url := "http://hans/" + local

	os.Create(local)
	defer os.Remove(local)

	download(td, url)

	if df {
		t.Errorf("Expected download to shortcircuit as file exists.")
	}

}

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

	os.RemoveAll(testDir)
}
