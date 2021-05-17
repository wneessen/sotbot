package httpclient

import (
	"bufio"
	"crypto/tls"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"time"
)

func NewHttpClient() (*http.Client, error) {
	tlsConfig := &tls.Config{
		MaxVersion:    tls.VersionTLS13,
		MinVersion:    tls.VersionTLS12,
		Renegotiation: tls.RenegotiateFreelyAsClient,
	}
	httpTransport := &http.Transport{TLSClientConfig: tlsConfig}
	cookieJar, err := cookiejar.New(nil)
	if err != nil {
		return &http.Client{}, err
	}
	httpClient := &http.Client{
		Transport: httpTransport,
		Timeout:   5 * time.Second,
		Jar:       cookieJar,
	}

	return httpClient, nil
}

// Perform a GET request
func HttpReqGet(u string, hc *http.Client, rc string, ref string) ([]byte, error) {
	l := log.WithFields(log.Fields{
		"action": "httpclient.HttpReqGet",
	})

	// Create a HTTP request
	httpReq, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return []byte{}, fmt.Errorf("an error occured creating new HTTP GET request: %v", err)
	}

	// Set cookie
	if rc != "" {
		ratCookie := &http.Cookie{
			Name:  "rat",
			Value: rc,
		}
		httpReq.AddCookie(ratCookie)
	}

	// Set referer
	if ref != "" {
		httpReq.Header.Set("referer", ref)
	}
	if ref == "" {
		httpReq.Header.Set("referer", "https://www.seaofthieves.com/profile/achievements")
	}

	// Set HTTP header
	setReqHeader(httpReq)
	serverResp, err := hc.Do(httpReq)
	if err != nil {
		return []byte{}, err
	}
	defer func() {
		err := serverResp.Body.Close()
		if err != nil {
			l.Printf("error while closing response body: %v", err)
		}
	}()
	if !strings.HasPrefix(serverResp.Header.Get("Content-Type"), "application/json") {
		return []byte{}, fmt.Errorf("API returned non-JSON response")
	}

	// Read the response body
	var respBody []byte
	scanner := bufio.NewScanner(serverResp.Body)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)
	for scanner.Scan() {
		respBody = scanner.Bytes()
	}
	if err = scanner.Err(); err != nil {
		return []byte{}, err
	}

	return respBody, nil
}

// Set package specific HTTP header
func setReqHeader(h *http.Request) {
	h.Header.Set("User-Agent", "curl/7.76.1")
	h.Header.Set("Accept", "*/*")
}
