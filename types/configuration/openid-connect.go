/*
 * Copyright (c) 2023.  Jan Eike Suchard. All rights reserved
 * SPDX-License-Identifier: MIT
 */

package configurationTypes

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

var errDiscoveryDisabled = errors.New("discovery is disabled in the configuration")
var errNoDiscoveryUriSet = errors.New("no discovery uri set in configuration")
var errEmptyDiscoveryUri = errors.New("empty discovery uri set in configuration")
var errNoIssuerDiscovered = errors.New("discovery result did not contain an issuer in field 'iss'")
var errNoJWKSUriDiscovered = errors.New("discovery result did not contain an jwks_uri in field 'jwks_uri'")
var errNoUserinfoUriDiscovered = errors.New("discovery result did not contain an userinfo endpoint in field 'userinfo_endpoint'")
var errFieldWrongType = errors.New("field has invalid type")

type OpenIDConnectConfiguration struct {
	// ClientID contains the client id that is used to get access and id tokens
	// during the login process in the fronted
	ClientID string `toml:"client-id" comment:"The client id used for the open id connect provider"`

	// UseDiscovery indicates if the openid provider metadata should be used
	// for discovering the following endpoints of the configured openid
	// provider:
	//   - issuer
	//   - userinfo_endpoint
	//   - jwks_uri
	UseDiscovery bool `toml:"use-discovery" comment:"Use the openid provider metadata returned by the discoveryEndpoint to discover other needed endpoints"`

	// DiscoveryEndpointUri contains the URI pointing to the discovery endpoint
	// of the identity provider used to provide the authentication in the
	// frontend
	DiscoveryEndpointUri *string `toml:"discovery-endpoint" comment:"The URI of the discovery endpoint that should be used for discovering the other endpoints"`

	// JWKSEndpointUri contains the URI pointing to the JSON Web Key Set that
	// the identity provider uses for signing and encrypting the JSON Web Tokens
	// issued for authenticating users
	JWKSEndpointUri *string `toml:"jwks-endpoint" comment:"The URI for the JWKS endpoint, that returns the keys used for signing and possibly encrypting the JSON Web Tokens used as access control mechanism"`

	// UserinfoEndpointUri contains the URI pointing to the Userinfo endpoint
	// which allows the backend to get the groups a user is in when accessing
	// the backend. Furthermore, it is used to validate the access token used
	// in the Authorization header
	UserinfoEndpointUri *string `toml:"userinfo-endpoint" comment:"The URI for the userinfo endpoint, that is used to get additional information about the user accessing the backend"`

	// Issuer contains the identification of the openid provider issuing the
	// JSON Web Tokens. It is used for verifying the origin of a JSON Web Token
	Issuer *string `toml:"issuer" comment:"The identification of the openid provider that issues the JWTs used to access the backend"`
}

// Discover executes the discovery of the open id connect endpoints needed
// for the backend authentication and authorization
func (c *OpenIDConnectConfiguration) Discover() error {
	// check if the discovery is enabled
	if !c.UseDiscovery {
		return errDiscoveryDisabled
	}

	// now check if a discovery endpoint uri has been set
	if c.DiscoveryEndpointUri == nil {
		return errNoDiscoveryUriSet
	}

	// since the discovery uri has been set, check if it contains anything
	if strings.TrimSpace(*c.DiscoveryEndpointUri) == "" {
		return errEmptyDiscoveryUri
	}

	// now get the discovery data from the discovery endpoint
	discoveryResponse, err := http.Get(*c.DiscoveryEndpointUri)
	if err != nil {
		return fmt.Errorf("unable to execute discovery request: %w", err)
	}

	// now parse the discovery response into a map
	var discoveryResult map[string]interface{}
	err = json.NewDecoder(discoveryResponse.Body).Decode(&discoveryResult)
	if err != nil {
		return fmt.Errorf("unable to parse discovery result: %w", err)
	}

	// now check the discovery result for the following fields:
	//   - issuer [denotes the issuer of the JWTs]
	//   - jwks_uri [denotes the endpoint under which the JWKS is available]
	//   - userinfo_endpoint [denotes the endpoint for the userinfo retrieval]
	rawIssuer, issuerSet := discoveryResult["issuer"]
	if !issuerSet {
		return errNoIssuerDiscovered
	}
	// now check if the issuer is a string
	issuer, isString := rawIssuer.(string)
	if !isString {
		return fmt.Errorf("unable to get issuer: %w, expected string", errFieldWrongType)
	}
	// since the issuer has been extracted, set it in the configuration
	c.Issuer = &issuer

	rawJWKSEndpoint, jwksEndpointSet := discoveryResult["jwks_uri"]
	if !jwksEndpointSet {
		return errNoJWKSUriDiscovered
	}
	// now check if the issuer is a string
	jwksUri, isString := rawJWKSEndpoint.(string)
	if !isString {
		return fmt.Errorf("unable to get jwks uri: %w, expected string", errFieldWrongType)
	}
	// since the issuer has been extracted, set it in the configuration
	c.JWKSEndpointUri = &jwksUri

	rawUserinfoEndpoint, userinfoEndpointSet := discoveryResult["userinfo_endpoint"]
	if !userinfoEndpointSet {
		return errNoUserinfoUriDiscovered
	}
	// now check if the issuer is a string
	userinfoUri, isString := rawUserinfoEndpoint.(string)
	if !isString {
		return fmt.Errorf("unable to get userinfo uri: %w, expected string", errFieldWrongType)
	}
	// since the issuer has been extracted, set it in the configuration
	c.UserinfoEndpointUri = &userinfoUri
	return nil
}
