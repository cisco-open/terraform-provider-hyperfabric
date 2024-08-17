// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// SPDX-License-Identifier: MPL-2.0

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
		// 	Name:  "HYPERFABRIC-Request-Signature",
		// 	Value: signature,
		// })
		// req.AddCookie(&http.Cookie{
		// 	Name:  "HYPERFABRIC-Certificate-Algorithm",
		// 	Value: "v1.0",
		// })

		// // Actual certificate fingerprint/thumbprint generation is not required
		// // Simply setting cookie to fingerprint is sufficient for cert-based requests.
		// req.AddCookie(&http.Cookie{
		// 	Name:  "HYPERFABRIC-Certificate-Fingerprint",
		// 	Value: "fingerprint",
		// })
		// if client.appUserName != "" {
		// 	req.AddCookie(&http.Cookie{
		// 		Name:  "HYPERFABRIC-Certificate-ID",
		// 		Value: fmt.Sprintf("uni/userext/appuser-%s/usercert-%s", client.appUserName, client.adminCert),
		// 	})
		// } else {
		// 	req.AddCookie(&http.Cookie{
		// 		Name:  "HYPERFABRIC-Certificate-ID",
		// 		Value: fmt.Sprintf("uni/userext/user-%s/usercert-%s", client.username, client.adminCert),
		// 	})
		// }
		log.Printf("[DEBUG] finished signature creation")
		return req, nil
	} else {
		// return req, fmt.Errorf("one of token or private_key is required")
		return req, fmt.Errorf("An Hyperfabric API Bearer Token is required. Set the `HYPERFABRIC_TOKEN` environment variable or set the `token` attribute under the provider configuration")
	}
}
