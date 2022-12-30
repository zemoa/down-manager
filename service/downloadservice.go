package service

import (
	"fmt"
	"net/http"
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

func DowloadFile(url string) {
	http.Get(url)
}

func getFileName(url string, headers *http.Header) string {
	filename := headers.Get("Content-Disposition")
	if filename == "" {
		splittedUrl := strings.Split(url, "/")
		filename = splittedUrl[len(splittedUrl)-1]
	}
	return filename
}
