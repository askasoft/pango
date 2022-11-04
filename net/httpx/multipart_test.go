package httpx

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"testing"
)

const testdir = "_testdir"

func TestMultipartWriteFile(t *testing.T) {
	defer os.RemoveAll(testdir)
	os.MkdirAll(testdir, os.FileMode(0777))

	files := []string{}
	for i := 0; i < 2; i++ {
		f := path.Join(testdir, fmt.Sprintf("%d.txt", i))
		os.WriteFile(f, []byte(f), os.FileMode(0666))
		files = append(files, f)
	}

	buf := &bytes.Buffer{}
	mw := NewMultipartWriter(buf)
	ct := mw.FormDataContentType()
	for _, file := range files {
		err := mw.WriteFile("files", file)
		if err != nil {
			t.Fatalf("Failed to write file %s", file)
			return
		}
	}
	mw.Close()

	url := "https://panda-demo.azurewebsites.net/files/uploads"
	res, err := http.Post(url, ct, buf)
	if err != nil {
		t.Fatalf("Failed to post file to %s - %v", url, err)
		return
	}
	defer res.Body.Close()

	io.Copy(os.Stdout, res.Body)
	fmt.Println()

	if res.StatusCode != http.StatusOK {
		t.Errorf("Invalid response status code: %d", res.StatusCode)
	}
}

func TestMultipartWriteFilePipe(t *testing.T) {
	defer os.RemoveAll(testdir)
	os.MkdirAll(testdir, os.FileMode(0777))

	files := []string{}
	for i := 0; i < 2; i++ {
		f := path.Join(testdir, fmt.Sprintf("%d.txt", i))
		os.WriteFile(f, []byte(f), os.FileMode(0666))
		files = append(files, f)
	}

	var errw error

	pr, pw := io.Pipe()
	mw := NewMultipartWriter(pw)
	ct := mw.FormDataContentType()
	go func() {
		defer pw.Close()
		defer mw.Close()
		for _, file := range files {
			errw = mw.WriteFile("files", file)
			if errw != nil {
				return
			}
		}
	}()

	url := "https://panda-demo.azurewebsites.net/files/uploads"
	res, err := http.Post(url, ct, pr)
	if err != nil {
		t.Fatalf("Failed to post file to %s - %v", url, err)
		return
	}
	defer res.Body.Close()

	io.Copy(os.Stdout, res.Body)
	fmt.Println()

	if res.StatusCode != http.StatusOK {
		t.Errorf("Invalid response status code: %d", res.StatusCode)
	}

	if errw != nil {
		t.Fatalf("Failed to write files %v", errw)
	}
}
