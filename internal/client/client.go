// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Jeffail/gabs/v2"
)

// Default timeout for NGINX in Hyperfabric service is 90 Seconds.
// Allow the client to set a shorter or longer time depending on their
// environment
const DefaultReqTimeoutVal int = 100
const DefaultBackoffMinDelay int = 4
const DefaultBackoffMaxDelay int = 60
const DefaultBackoffDelayFactor float64 = 3
const DefaultCandidate string = "default"

// Client is the main entry point
type Client struct {
	baseURL            *url.URL
	pathURL            string
	httpClient         *http.Client
	apiToken           string
	privatekey         string
	adminCert          string
	insecure           bool
	reqTimeoutSet      bool
	reqTimeoutVal      uint32
	proxyUrl           string
	proxyCreds         string
	preserveBaseUrlRef bool
	skipLoggingPayload bool
	maxRetries         int
	backoffMinDelay    int
	backoffMaxDelay    int
	backoffDelayFactor float64
	autoCommit         bool
	candidate          string
	// lockRequest        sync.Mutex
	changedFabrics    map[string]string
	lockChangedFabric sync.Mutex
}

// singleton implementation of a client
var clientImpl *Client

type Option func(*Client)

func Insecure(insecure bool) Option {
	return func(client *Client) {
		client.insecure = insecure
	}
}

func PathURL(pathUrl string) Option {
	return func(sm *Client) {
		sm.pathURL = pathUrl
	}
}

func APIToken(apiToken string) Option {
	return func(client *Client) {
		client.apiToken = apiToken
	}
}

func PrivateKey(privatekey string) Option {
	return func(client *Client) {
		client.privatekey = privatekey
	}
}

func AdminCert(adminCert string) Option {
	return func(client *Client) {
		client.adminCert = adminCert
	}
}

func ProxyUrl(pUrl string) Option {
	return func(client *Client) {
		client.proxyUrl = pUrl
	}
}

func ProxyCreds(pcreds string) Option {
	return func(client *Client) {
		client.proxyCreds = pcreds
	}
}

func MaxRetries(maxRetries int) Option {
	return func(client *Client) {
		client.maxRetries = maxRetries
	}
}

func BackoffMinDelay(backoffMinDelay int) Option {
	return func(client *Client) {
		client.backoffMinDelay = backoffMinDelay
	}
}

func BackoffMaxDelay(backoffMaxDelay int) Option {
	return func(client *Client) {
		client.backoffMaxDelay = backoffMaxDelay
	}
}

func BackoffDelayFactor(backoffDelayFactor float64) Option {
	return func(client *Client) {
		client.backoffDelayFactor = backoffDelayFactor
	}
}

func AutoCommit(autoCommit bool) Option {
	return func(client *Client) {
		client.autoCommit = autoCommit
	}
}

func Candidate(candidate string) Option {
	return func(client *Client) {
		client.candidate = candidate
	}
}

func (c *Client) AddChangedFabric(fabricId string) {
	c.lockChangedFabric.Lock()
	if c.changedFabrics != nil {
		c.changedFabrics[fabricId] = ""
	} else {
		c.changedFabrics = map[string]string{fabricId: ""}
	}
	c.lockChangedFabric.Unlock()
}

func (c *Client) DoAutoCommit() {
	if clientImpl != nil {
		clientImpl.lockChangedFabric.Lock()
		if len(clientImpl.changedFabrics) > 0 && clientImpl.autoCommit {
			log.Printf("[DEBUG] Start of the auto-committing process due to auto_commit (true) and change detected.")
			candidate := DefaultCandidate
			if clientImpl.candidate != "" {
				candidate = clientImpl.candidate
			}
			marshalPayload, err := json.Marshal(map[string]string{"comments": "Terraform Auto-Commit"})
			if err != nil {
				log.Printf("[DEBUG] Marshalling of candidate JSON payload failed. Err: %s. Please report this issue to the provider developers.", err)
				return
			}
			jsonPayload, err := gabs.ParseJSON(marshalPayload)
			if err != nil {
				log.Printf("[DEBUG] Construction of candidate JSON payload failed. Err: %s. Please report this issue to the provider developers.", err)
				return
			}
			for fabricId := range clientImpl.changedFabrics {
				log.Printf("[TRACE] Auto-committing for fabric %s.", fabricId)
				_, errorDiag := clientImpl.DoRestRequest(fmt.Sprintf("/api/v1/fabrics/%s/candidates/%s", fabricId, candidate), "POST", jsonPayload)
				if errorDiag != nil {
					log.Printf("[DEBUG] Error when commiting Fabric %s, %s. %s.", fabricId, errorDiag.Summary, errorDiag.Detail)
				}
			}

			log.Printf("[DEBUG] End of the auto-committing process.")
		}
		clientImpl.lockChangedFabric.Unlock()
	}
}

// HttpClient option: allows for caller to set 'httpClient' with 'Transport'.
// When this option is set 'client.proxyUrl' option is ignored.
func HttpClient(httpcl *http.Client) Option {
	return func(client *Client) {
		client.httpClient = httpcl
	}
}

func SkipLoggingPayload(skipLoggingPayload bool) Option {
	return func(client *Client) {
		client.skipLoggingPayload = skipLoggingPayload
	}
}

func PreserveBaseUrlRef(preserveBaseUrlRef bool) Option {
	return func(client *Client) {
		client.preserveBaseUrlRef = preserveBaseUrlRef
	}
}

func ReqTimeout(timeout uint32) Option {
	return func(client *Client) {
		client.reqTimeoutSet = true
		client.reqTimeoutVal = timeout
	}
}

func initClient(clientUrl, apiToken string, options ...Option) *Client {
	var transport *http.Transport
	baseUrl, err := url.Parse(clientUrl)
	if err != nil {
		// cannot move forward if url is undefined
		log.Fatal(err)
	}
	client := &Client{
		baseURL:  baseUrl,
		apiToken: apiToken,
	}

	for _, option := range options {
		option(client)
	}

	if client.httpClient == nil {
		transport = client.useInsecureHTTPClient(client.insecure)
		if client.proxyUrl != "" {
			transport = client.configProxy(transport)
		}
		client.httpClient = &http.Client{
			Transport: transport,
		}
	}

	var timeout time.Duration
	if client.reqTimeoutSet {
		timeout = time.Second * time.Duration(client.reqTimeoutVal)
	} else {
		timeout = time.Second * time.Duration(DefaultReqTimeoutVal)
	}

	client.httpClient.Timeout = timeout
	return client
}

// GetClient returns a singleton
func GetClient(clientUrl, apiToken string, options ...Option) *Client {
	if clientImpl == nil {
		clientImpl = initClient(clientUrl, apiToken, options...)
	} else {
		// making sure it is the same client
		baseUrl, err := url.Parse(clientUrl)
		if err != nil {
			// cannot move forward if url is undefined
			log.Fatal(err)
		}
		if baseUrl != clientImpl.baseURL {
			clientImpl = initClient(clientUrl, apiToken, options...)
		}
	}
	return clientImpl
}

// NewClient returns a new Instance of the client - allowing for simultaneous connections to the Hyperfabric service
func NewClient(clientUrl, apiToken string, options ...Option) *Client {
	// making sure it is the same client
	_, err := url.Parse(clientUrl)
	if err != nil {
		// cannot move forward if url is undefined
		log.Fatal(err)
	}

	// initClient always returns a new struct, so always create a new pointer to allow for
	// multiple object instances
	newClientImpl := initClient(clientUrl, apiToken, options...)

	return newClientImpl
}

func (c *Client) configProxy(transport *http.Transport) *http.Transport {
	log.Printf("[DEBUG]: Using Proxy Server: %s ", c.proxyUrl)
	pUrl, err := url.Parse(c.proxyUrl)
	if err != nil {
		log.Fatal(err)
	}
	transport.Proxy = http.ProxyURL(pUrl)

	if c.proxyCreds != "" {
		basicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(c.proxyCreds))
		transport.ProxyConnectHeader = http.Header{}
		transport.ProxyConnectHeader.Add("Proxy-Authorization", basicAuth)
	}

	return transport
}

func (c *Client) useInsecureHTTPClient(insecure bool) *http.Transport {
	transport := http.DefaultTransport.(*http.Transport)

	transport.TLSClientConfig = &tls.Config{
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
		},
		PreferServerCipherSuites: true,
		InsecureSkipVerify:       insecure,
		MinVersion:               tls.VersionTLS11,
		MaxVersion:               tls.VersionTLS13,
	}

	return transport
}

func (c *Client) MakeRestRequest(method string, path string, payloadContainer *gabs.Container, payloadByteArray []byte, authenticated bool) (*http.Request, error) {

	pathURL, err := url.Parse(path)
	if err != nil {
		return nil, err
	}

	url, err := url.Parse(c.baseURL.String())
	if err != nil {
		return nil, err
	}

	if c.preserveBaseUrlRef {
		// Default is false for preserveBaseUrlRef - matching original behavior to strip out baseURL
		urlStr := url.String() + pathURL.String()
		url, err = url.Parse(urlStr)
		if err != nil {
			return nil, err
		}
	} else {
		// Original behavior to strip down baseURL
		url = url.ResolveReference(pathURL)
	}

	var req *http.Request
	log.Printf("[DEBUG] baseURL: %s, pathURL: %s, url: %s", c.baseURL.String(), pathURL.String(), url.String())
	if method == "GET" || method == "DELETE" {
		req, err = http.NewRequest(method, url.String(), nil)
	} else if payloadContainer != nil {
		req, err = http.NewRequest(method, url.String(), bytes.NewBuffer((payloadContainer.Bytes())))
	} else {
		req, err = http.NewRequest(method, url.String(), bytes.NewBuffer(payloadByteArray))
	}
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	// req.Header.Set("Cisco-FA-ShowEmpty", "true")

	if c.skipLoggingPayload {
		log.Printf("HTTP request %s %s", method, path)
	} else {
		log.Printf("HTTP request %s %s %v", method, path, req)
	}

	if authenticated {
		req, err = c.InjectAuthenticationHeader(req, "")
		if err != nil {
			return req, err
		}
	}

	if !c.skipLoggingPayload {
		log.Printf("HTTP request after auth header injection %s %s %v", method, path, req)
	}

	return req, nil
}

func StrtoInt(s string, startIndex int, bitSize int) (int64, error) {
	return strconv.ParseInt(s, startIndex, bitSize)
}

func (c *Client) Do(req *http.Request) (*gabs.Container, *http.Response, error) {
	log.Printf("[DEBUG] Begining Do method %s", req.URL.String())

	// retain the request body across multiple attempts
	var body []byte
	if req.Body != nil && c.maxRetries != 0 {
		body, _ = io.ReadAll(req.Body)
	}

	for attempts := 0; ; attempts++ {
		log.Printf("[TRACE] HTTP Request Method and URL: %s %s", req.Method, req.URL.String())
		if c.maxRetries != 0 {
			req.Body = io.NopCloser(bytes.NewBuffer(body))
		}
		if !c.skipLoggingPayload {
			log.Printf("[TRACE] HTTP Request Body: %v", req.Body)
		}

		resp, err := c.httpClient.Do(req)
		if err != nil {
			if strings.Contains(err.Error(), " tls: ") {
				log.Printf("[ERROR] HTTP Connection failed due to TLS Error: %+v", err)
				return nil, nil, fmt.Errorf("failed to connect due to a TLS error. Verify that you are connecting to the correct Hyperfabric service.\nError message: %+v", err)
			} else {
				if ok := c.backoff(attempts); !ok {
					log.Printf("[ERROR] HTTP Connection error occured: %+v", err)
					log.Printf("[DEBUG] Exit from Do method")
					return nil, nil, fmt.Errorf("failed to connect to the Hyperfabric service. Verify that you are connecting to the correct Hyperfabric service.\nError message: %+v", err)
				} else {
					log.Printf("[ERROR] HTTP Connection failed: %s, retries: %v", err, attempts)
					continue
				}
			}
		}

		if !c.skipLoggingPayload {
			log.Printf("[TRACE] HTTP Response: %d %s %v", resp.StatusCode, resp.Status, resp)
		} else {
			log.Printf("[TRACE] HTTP Response: %d %s", resp.StatusCode, resp.Status)
		}

		bodyBytes, _ := io.ReadAll(resp.Body)
		bodyStr := string(bodyBytes)
		resp.Body.Close()

		if !c.skipLoggingPayload {
			log.Printf("[DEBUG] HTTP response unique string %s %s %s", req.Method, req.URL.String(), bodyStr)
		}
		if resp.StatusCode == 204 {
			return nil, resp, nil
		} else if resp.StatusCode == 200 || resp.StatusCode == 201 {
			obj, err := gabs.ParseJSON(bodyBytes)
			if err != nil {
				log.Printf("[ERROR] Error occured while JSON parsing response with status code 200 and 201: %+v", err)

				// If nginx is too busy or the page is not found, nginx might respond with an HTML doc instead of a JSON Response.
				// In those cases, parse the HTML response for the message and return that to the user
				// htmlErr := c.checkHtmlResp(bodyStr)
				// log.Printf("[ERROR] Error occured while json parsing: %s", htmlErr.Error())
				log.Printf("[DEBUG] Exit from Do method")
				return nil, resp, fmt.Errorf("failed to parse JSON response with status code 200 and 201 from: %s. Verify that you are connecting to the correct Hyperfabric service.\nHTTP response status: %s\nMessage: %s", req.URL.String(), resp.Status, bodyBytes)
			}
			log.Printf("[DEBUG] Exit from Do method")
			return obj, resp, nil
		} else {
			var unrecoverableError bool
			obj, err := gabs.ParseJSON(bodyBytes)
			if err != nil {
				unrecoverableError = true
			} else {
				restError := NewRestError(obj.Data().(map[string]interface{}))
				if restError.ErrCode != "ERR_CODE_SERVICE_UNAVAILABLE" && restError.ErrCode != "ERR_CODE_TOO_MANY_REQUESTS" {
					unrecoverableError = true
				}
			}
			if ok := c.backoff(attempts); unrecoverableError || !ok {
				if err != nil {
					log.Printf("[ERROR] Error occured while json parsing: %+v", err)

					// If nginx is too busy or the page is not found, nginx might respond with an HTML doc instead of a JSON Response.
					// In those cases, parse the HTML response for the message and return that to the user
					// htmlErr := c.checkHtmlResp(bodyStr)
					// log.Printf("[ERROR] Error occured while json parsing: %s", htmlErr.Error())
					log.Printf("[DEBUG] Exit from Do method")
					return nil, resp, fmt.Errorf("failed to parse JSON response from: %s. Verify that you are connecting to the correct Hyperfabric service.\nHTTP response status: %s\nMessage: %s", req.URL.String(), resp.Status, bodyBytes)
				}
				log.Printf("[DEBUG] Exit from Do method")
				return obj, resp, nil
			} else {
				log.Printf("[ERROR] HTTP Request failed: StatusCode %v, Retries: %v", resp.StatusCode, attempts)
				continue
			}
		}
	}
}

func (c *Client) backoff(attempts int) bool {
	log.Printf("[DEBUG] Begining backoff method: attempts %v on %v", attempts, c.maxRetries)
	if attempts >= c.maxRetries {
		log.Printf("[DEBUG] Exit from backoff method with return value false")
		return false
	}

	minDelay := time.Duration(DefaultBackoffMinDelay) * time.Second
	if c.backoffMinDelay != 0 {
		minDelay = time.Duration(c.backoffMinDelay) * time.Second
	}

	maxDelay := time.Duration(DefaultBackoffMaxDelay) * time.Second
	if c.backoffMaxDelay != 0 {
		maxDelay = time.Duration(c.backoffMaxDelay) * time.Second
	}

	factor := DefaultBackoffDelayFactor
	if c.backoffDelayFactor != 0 {
		factor = c.backoffDelayFactor
	}

	min := float64(minDelay)
	backoff := min * math.Pow(factor, float64(attempts))
	if backoff > float64(maxDelay) {
		backoff = float64(maxDelay)
	}
	backoff = (rand.Float64()/2+0.5)*(backoff-min) + min
	backoffDuration := time.Duration(backoff)
	log.Printf("[TRACE] Starting sleeping for %v", backoffDuration.Round(time.Second))
	time.Sleep(backoffDuration)
	log.Printf("[DEBUG] Exit from backoff method with return value true")
	return true
}

type RestError struct {
	Causes     []string
	Critical   bool
	ErrCode    string // int32
	Field      string
	Message    string
	Notes      string
	Status     float64
	TrackingId string
	Value      string
}

type DiagError struct {
	Summary string
	Detail  string
}

func getDiagError(summary string, detail string) *DiagError {
	return &DiagError{
		Summary: summary,
		Detail:  detail,
	}
}

func (e RestError) ToString() string {
	return fmt.Sprintf(
		"Status: %v Error Code: %s Message: %s Field: %s Value: %s Critical: %v Notes: %s",
		e.Status,
		e.ErrCode,
		e.Message,
		e.Field,
		e.Value,
		e.Critical,
		e.Notes,
	)
}

// func getEmptyRestError() RestError {
// 	return RestError{
// 		Causes:     make([]string, 0),
// 		Critical:   false,
// 		ErrCode:    ""
// 		Field:      ""
// 		Message:    basetypes.NewStringNull(),
// 		Notes:      basetypes.NewStringNull(),
//		Status:
// 		TrackingId: basetypes.NewStringNull(),
// 		Value:      basetypes.NewStringNull(),
// 	}
// }

func NewRestError(data map[string]interface{}) RestError {
	var restError RestError
	restError.ErrCode = "NO_ERROR_CODE"
	for attributeName, attributeValue := range data {
		if attributeName == "causes" && attributeValue != nil {
			newCauses := make([]string, 0)
			for _, cause := range attributeValue.([]interface{}) {
				newCauses = append(newCauses, fmt.Sprint(cause))
			}
			restError.Causes = newCauses
		} else if attributeName == "critical" && attributeValue != nil {
			restError.Critical = attributeValue.(bool)
		} else if attributeName == "errCode" && attributeValue != nil {
			restError.ErrCode = attributeValue.(string)
		} else if attributeName == "field" && attributeValue != nil {
			restError.Field = attributeValue.(string)
		} else if attributeName == "message" && attributeValue != nil {
			restError.Message = attributeValue.(string)
		} else if attributeName == "notes" && attributeValue != nil {
			restError.Notes = attributeValue.(string)
		} else if attributeName == "status" && attributeValue != nil {
			restError.Status = attributeValue.(float64)
		} else if attributeName == "trackingId" && attributeValue != nil {
			restError.TrackingId = attributeValue.(string)
		} else if attributeName == "value" && attributeValue != nil {
			restError.Value = attributeValue.(string)
		}
	}
	return restError
}

func (c *Client) DoRestRequest(path, method string, payload *gabs.Container) (*gabs.Container, *DiagError) {
	restRequest, err := c.MakeRestRequest(method, path, payload, nil, true)
	if err != nil {
		errString := fmt.Sprintf("Error: %s. Please report this issue to the provider developers.", err)
		if strings.HasPrefix(err.Error(), "An Hyperfabric API Bearer Token is required.") {
			errString = fmt.Sprintf("Error: %s.", err)
		}
		diagError := getDiagError(
			fmt.Sprintf("Creation of %s REST request to %s failed", strings.ToUpper(method), path),
			errString,
		)
		return nil, diagError
	}
	// c.lockRequest.Lock()
	container, restResponse, err := c.Do(restRequest)
	// c.lockRequest.Unlock()

	if restResponse != nil && container.Data() != nil && (restResponse.StatusCode != 200 && restResponse.StatusCode != 204) {
		restError := NewRestError(container.Data().(map[string]interface{}))

		// Need error codes for:  Cannot create object, Cannot delete object
		log.Printf("[DEBUG] The %s REST request to %s failed with HTTP Status Code %d, %s", strings.ToUpper(method), path, restResponse.StatusCode, restError.ToString())

		if restResponse.StatusCode == 404 && (strings.ToLower(method) == "get" || strings.ToLower(method) == "delete") {
			return nil, nil
		} else {
			diagError := getDiagError(
				fmt.Sprintf("The %s REST request to %s failed with HTTP Status Code %d", strings.ToUpper(method), path, restResponse.StatusCode),
				fmt.Sprintf("%s, err: %v. Please report this issue to the provider developers.", restError.ToString(), err),
			)
			return nil, diagError
		}
	} else if err != nil {
		if restResponse == nil || !(restResponse.StatusCode == 404 && (strings.ToLower(method) == "get" || strings.ToLower(method) == "delete")) {
			diagError := getDiagError(
				fmt.Sprintf("The %s REST request to %s failed with HTTP Status Code %d", strings.ToUpper(method), path, restResponse.StatusCode),
				fmt.Sprintf("Err: %s. Please report this issue to the provider developers.", err),
			)
			return nil, diagError
		}
		return nil, nil
	}

	return container, nil
}
