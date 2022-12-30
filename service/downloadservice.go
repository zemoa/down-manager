package service

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func GetLinkDetails(url string) (filename string, rangeSupported bool, length int, error error) {
	response, err := http.Head(url)
	if err != nil || (response.StatusCode != 200 && response.StatusCode != 206) {
		return url, false, 0, fmt.Errorf("error while creating request : %v", err)
	}

	headers := response.Header
	conLen := headers.Get("Content-Length")
	filename_ := getFileName(url, &headers)
	cl, err := strconv.Atoi(conLen)
	if err != nil {
		return filename_, false, 0, fmt.Errorf("error Parsing content length : %v", err)
	}

	if headers.Get("Accept-Ranges") == "bytes" {
		return filename_, true, cl, nil
	}

	return filename_, false, cl, nil
}

func DownloadFile(url string, filename string, path string) error {
	filePath := path + "/" + filename
	out, err := os.Create(filePath + ".tmp")
	if err != nil {
		out.Close()
		return err
	}
	resp, err := http.Get(url)
	if err != nil {
		out.Close()
		return err
	}
	defer resp.Body.Close()
	counter := &WriteCounter{}
	if _, err = io.Copy(out, io.TeeReader(resp.Body, counter)); err != nil {
		out.Close()
		return err
	}
	out.Close()
	if err = os.Rename(filePath+".tmp", filePath); err != nil {
		return err
	}
	return nil
}

func getFileName(url string, headers *http.Header) string {
	filename := headers.Get("Content-Disposition")
	if filename == "" {
		splittedUrl := strings.Split(url, "/")
		filename = splittedUrl[len(splittedUrl)-1]
	}
	return filename
}

type WriteCounter struct {
	Total uint64
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += uint64(n)
	wc.PrintProgress()
	return n, nil
}

func (wc WriteCounter) PrintProgress() {
	// Clear the line by using a character return to go back to the start and remove
	// the remaining characters by filling it with spaces
	fmt.Printf("\r%s", strings.Repeat(" ", 35))

	// Return again and print current status of download
	// We use the humanize package to print the bytes in a meaningful way (e.g. 10 MB)
	fmt.Printf("\rDownloading... %d complete", wc.Total)
}
