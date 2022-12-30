package service

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"strconv"
	"strings"
)

func GetLinkDetails(url string) (filename string, rangeSupported bool, length int, error error) {
	request, err := http.NewRequest("HEAD", url, strings.NewReader(""))
	if err != nil {
		return url, false, 0, fmt.Errorf("error while creating request : %v", err)
	}

	statusCode_, headers, _, err := doAPICall(request)
	if err != nil || (statusCode_ != 200 && statusCode_ != 206) {
		return url, false, 0, fmt.Errorf("error calling url : %v", err)
	}

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

func getFileName(url string, headers *http.Header) string {
	filename := headers.Get("Content-Disposition")
	if filename == "" {
		splittedUrl := strings.Split(url, "/")
		filename = splittedUrl[len(splittedUrl)-1]
	}
	return filename
}

// doAPICall will do the api call and return statuscode,headers,data,error respectively
func doAPICall(request *http.Request) (statusCode int, header http.Header, body []byte, error error) {

	client := http.Client{
		Timeout: 0,
	}

	response, err := client.Do(request)
	if err != nil {
		return 0, http.Header{}, []byte{}, fmt.Errorf("error while doing request : %v", err)
	}
	defer response.Body.Close()

	data, err := httputil.DumpResponse(response, true)
	if err != nil {
		return 0, http.Header{}, []byte{}, fmt.Errorf("error while reading response body : %v", err)
	}

	return response.StatusCode, response.Header, data, nil

}
