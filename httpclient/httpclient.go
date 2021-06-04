package httpclient

import (
	"bytes"
	"crypto/tls"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/wneessen/sotbot/version"
	"io"
	"net/http"
	"net/http/cookiejar"
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
		Timeout:   10 * time.Second,
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

	if serverResp.StatusCode != 200 {
		l.Errorf("HTTP request failed: %v", serverResp.Status)
		return []byte{}, fmt.Errorf("%v", serverResp.StatusCode)
	}

	// Read the response body
	var respBody bytes.Buffer
	_, err = io.Copy(&respBody, serverResp.Body)
	if err != nil {
		return []byte{}, err
	}

	return respBody.Bytes(), nil
}

// Set package specific HTTP header
func setReqHeader(h *http.Request) {
	h.Header.Set("User-Agent", fmt.Sprintf("SoT Discord Bot v%v", version.Version))
	//h.Header.Set("Accept", "application/json")
}
