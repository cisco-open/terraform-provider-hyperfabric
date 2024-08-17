// Copyright 2024 Cisco Systems, Inc. and its affiliates
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func (client *Client) InjectAuthenticationHeader(req *http.Request, path string) (*http.Request, error) {
	log.Printf("[DEBUG] Begin Injection")
	// client.l.Lock()
	// defer client.l.Unlock()
	if client.apiToken != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", client.apiToken))
		return req, nil
	} else if client.privatekey != "" {
		var bodyStr string
		if req.Method != "GET" {
			buffer, _ := ioutil.ReadAll(req.Body)
			rdr2 := ioutil.NopCloser(bytes.NewBuffer(buffer))

			req.Body = rdr2
			bodyStr = string(buffer)
		}
		contentStr := ""
		if bodyStr != "{}" {
			contentStr = fmt.Sprintf("%s%s%s", req.Method, path, bodyStr)
		} else {
			contentStr = fmt.Sprintf("%s%s", req.Method, path)

		}
		log.Printf("[DEBUG] Content %s", contentStr)
		// content := []byte(contentStr)

		// signature, err := createSignature(content, client.privatekey)
		// log.Printf("[DEBUG] Signature %s", signature)
		// if err != nil {
		// 	return req, err
		// }
		// req.AddCookie(&http.Cookie{
		// 	Name:  "APIC-Request-Signature",
		// 	Value: signature,
		// })
		// req.AddCookie(&http.Cookie{
		// 	Name:  "APIC-Certificate-Algorithm",
		// 	Value: "v1.0",
		// })

		// // Actual certificate fingerprint/thumbprint generation is not required
		// // Simply setting cookie to fingerprint is sufficient for cert-based requests.
		// req.AddCookie(&http.Cookie{
		// 	Name:  "APIC-Certificate-Fingerprint",
		// 	Value: "fingerprint",
		// })
		// if client.appUserName != "" {
		// 	req.AddCookie(&http.Cookie{
		// 		Name:  "APIC-Certificate-DN",
		// 		Value: fmt.Sprintf("uni/userext/appuser-%s/usercert-%s", client.appUserName, client.adminCert),
		// 	})
		// } else {
		// 	req.AddCookie(&http.Cookie{
		// 		Name:  "APIC-Certificate-DN",
		// 		Value: fmt.Sprintf("uni/userext/user-%s/usercert-%s", client.username, client.adminCert),
		// 	})
		// }
		log.Printf("[DEBUG] finished signature creation")
		return req, nil
	} else {

		return req, fmt.Errorf("one of token or private_key is required")
	}
}
