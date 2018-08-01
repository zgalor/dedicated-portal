/*
Copyright (c) 2018 Red Hat, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package auth to handle authentication
package auth

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"

	"github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

// jwtCert on jwt key
type jwtCert struct {
	KID string `json:"kid,omitempty"`
	Kty string `json:"kty,omitempty"`
	Alg string `json:"alg,omitempty"`
	Use string `json:"use,omitempty"`
	N   string `json:"n,omitempty"`
	E   string `json:"e,omitempty"`
}

// jwtKeys a list of JwtCerts
type jwtKeys struct {
	Keys []jwtCert `json:"keys"`
}

// CheckToken checks token validity and
// returns the user subject
func CheckToken(w http.ResponseWriter, r *http.Request) (sub string, err error) {
	var claims jwt.MapClaims
	var ok bool

	// Get user from request
	user := r.Context().Value("user")
	if user == nil {
		err = fmt.Errorf("can't parse token")
		OnAuthError(w, r, err.Error())
		return
	}

	// Get token
	userToken := user.(*jwt.Token)

	// Check validity
	if !userToken.Valid {
		err = fmt.Errorf("token invalid")
		OnAuthError(w, r, err.Error())
		return
	}

	// Get token claims
	if claims, ok = userToken.Claims.(jwt.MapClaims); !ok {
		err = fmt.Errorf("can't parse token claims")
		OnAuthError(w, r, err.Error())
		return
	}

	// TODO:
	// check aud, access ...

	// Get the token subject (keycloack's uniqe ID of our user)
	sub, _ = claims["sub"].(string)

	return
}

// OnAuthError returns an error json struct
func OnAuthError(w http.ResponseWriter, r *http.Request, err string) {
	msg, _ := json.Marshal(map[string]string{"error": fmt.Sprint(err)})
	if msg == nil {
		msg = []byte("{\"error\":\"Unknown error while converting an error to json\"}")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)

	// Send response body
	_, e := w.Write(msg)
	if e != nil {
		glog.Errorf("Write to client: %s", e)
	}
}

// downloadPublicKeys download public keys from URL.
func downloadPublicKeys(url string) (keyMap map[string]*rsa.PublicKey, err error) {
	var body []byte
	var certs jwtKeys
	var res *http.Response
	var pemStr string

	// Init keyMap
	keyMap = map[string]*rsa.PublicKey{}

	// If no errors getting response from cert server:
	res, err = http.Get(url)
	if err != nil {
		return
	}

	// Try to read the response body.
	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}

	// Try to parse the response body.
	err = json.Unmarshal(body, &certs)
	if err != nil {
		return
	}

	// Convert cert list to map.
	for _, c := range certs.Keys {
		// Try to convert cert to string.
		pemStr, err = certToPEM(c)
		if err != nil {
			return
		}

		keyMap[c.KID], err = jwt.ParseRSAPublicKeyFromPEM([]byte(pemStr))
		if err != nil {
			return
		}
	}

	return
}

// certToPEM convert JWT object to PEM
func certToPEM(c jwtCert) (string, error) {
	var out bytes.Buffer

	// Check key type.
	if c.Kty != "RSA" {
		return "", fmt.Errorf("invalid key type: %s", c.Kty)
	}

	// Decode the base64 bytes for e and n.
	nb, err := base64.RawURLEncoding.DecodeString(c.N)
	if err != nil {
		return "", err
	}
	eb, err := base64.RawURLEncoding.DecodeString(c.E)
	if err != nil {
		return "", err
	}

	// Generate new public key
	pk := &rsa.PublicKey{
		N: new(big.Int).SetBytes(nb),
		E: int(new(big.Int).SetBytes(eb).Int64()),
	}

	der, err := x509.MarshalPKIXPublicKey(pk)
	if err != nil {
		return "", err
	}

	block := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: der,
	}

	// Output pem as string
	pem.Encode(&out, block)
	return out.String(), nil
}

// Router creates an authentication router middleware.
func Router(certURL string, router *mux.Router) (*negroni.Negroni, error) {
	// Try to read the JWT public key object file.
	keyMap, err := downloadPublicKeys(certURL)
	if err != nil {
		return nil, err
	}

	// Add the JWT Middleware
	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			// Try to get the token kid.
			kid, ok := token.Header["kid"]
			if !ok {
				return nil, fmt.Errorf("auth.Router: no kid in token")
			}

			// Try to get currect cert from certs map.
			key, ok := keyMap[kid.(string)]
			if !ok {
				return nil, fmt.Errorf("auth.Router: can't find key for key id [%v]", kid)
			}

			return key, nil
		},
		ErrorHandler:  OnAuthError,
		SigningMethod: jwt.SigningMethodRS256,
	})

	// Enable the access authentication:
	authRouter := negroni.New(
		negroni.HandlerFunc(jwtMiddleware.HandlerWithNext))
	authRouter.UseHandler(router)

	// Return router, error
	return authRouter, nil
}
