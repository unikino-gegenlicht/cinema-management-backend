/*
 * Copyright (c) 2023.  Jan Eike Suchard. All rights reserved
 * SPDX-License-Identifier: MIT
 */

package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"regexp"

	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/rs/zerolog/log"

	backendErrors "github.com/unikino-gegenlicht/cinema-management-backend/errors"
	configurationTypes "github.com/unikino-gegenlicht/cinema-management-backend/types/configuration"
)

type userInfo struct {
	FullName string        `bson:"fullName"`
	Username string        `bson:"username"`
	Groups   []interface{} `bson:"groups"`
}

// OpenIDConnectJWTAuthentication uses the Access Token present in the request
// headers to authenticate and check the authorization of the user making a
// call to the backend.
// To check a request for authorization, the middleware validates the access
// token via the JWKS uri and checks the information contained in the access
// token.
// Furthermore, it will also check that the access token contains the correct
// scopes to allow access to the backend.
// To allow the individual access control to some routes, the middleware
// attaches all scopes found to the request context.
// This way, the routes may filter the scopes further, if needed.
func OpenIDConnectJWTAuthentication(config configurationTypes.OpenIDConnectConfiguration, apiErrors map[string]backendErrors.APIError) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		// create the handlerFunction in which all the code of the handler will
		// be contained in
		handlerFunction := func(w http.ResponseWriter, r *http.Request) {
			// start the handler by extracting the value from the Authorization
			// header
			headerValues, headerSet := r.Header["Authorization"]
			// now check if the header is even present
			if !headerSet {
				apiErrors["NO_AUTHENTICATION"].SendError(w)
				return
			}
			// now check that there is only one Authorization header present
			if len(headerValues) != 1 {
				apiErrors["MULTIPLE_AUTHORIZATIONS"].SendError(w)
				return
			}
			// since the id token is transmitted as bearer token, get the token
			// without "Bearer" prepended
			replaceRegex := regexp.MustCompile(`(?i)^bearer\s*`)
			rawAccessToken := headerValues[0]
			rawAccessToken = replaceRegex.ReplaceAllString(rawAccessToken, "")

			// now try to get the jwks from the identity provider
			res, err := http.Get(*config.JWKSEndpointUri)
			if err != nil {
				log.Error().Err(err).Msg("unable to pull jwks from open id provider")
				apiErrors["JWKS_NOT_LOADABLE"].SendError(w)
				return
			}
			rawJWKS, err := io.ReadAll(res.Body)
			// now parse the response body into the jwks
			jwks, err := jwk.Parse(rawJWKS)
			if err != nil {
				log.Error().Err(err).Msg("unable to parse jwks from open id provider")
				apiErrors["JWKS_PARSE_FAILED"].SendError(w)
				w.WriteHeader(500)
				return
			}

			// since the id token and the jwks have been loaded, parse the raw
			// id token into the jwt used for checking the authentication
			// information
			options := []jwt.ParseOption{
				jwt.WithKeySet(jwks),
				jwt.WithIssuer(*config.Issuer),
				jwt.WithAudience(config.ClientID),
			}
			accessToken, err := jwt.ParseString(rawAccessToken, options...)
			if err != nil {
				switch err {
				case jwt.ErrInvalidAudience():
					log.Warn().Str("err", "'audience' incorrect").Msg("backend accessed with invalid token")
					apiErrors["INCORRECT_AUDIENCE"].SendError(w)
					return
				case jwt.ErrInvalidIssuedAt():
					log.Warn().Str("err", "'issued at' incorrect").Msg("backend accessed with invalid token")
					apiErrors["ISSUED_AT_INVALID"].SendError(w)
					return
				case jwt.ErrInvalidIssuer():
					log.Warn().Str("err", "'issuer' incorrect").Msg("backend accessed with invalid token")
					apiErrors["INCORRECT_ISSUER"].SendError(w)
					return
				case jwt.ErrRequiredClaim():
					log.Warn().Str("err", "'claims' incorrect").Msg("backend accessed with invalid token")
					apiErrors["MISSING_CLAIM"].SendError(w)
					return
				case jwt.ErrTokenExpired():
					log.Warn().Str("err", "token expired").Msg("backend accessed with invalid token")
					apiErrors["TOKEN_EXPIRED"].SendError(w)
					return
				case jwt.ErrTokenNotYetValid():
					log.Warn().Str("err", "token not valid yet").Msg("backend accessed with invalid token")
					apiErrors["TOKEN_NOT_ALIVE"].SendError(w)
					return
				default:
					log.Warn().Err(err).Msg("backend accessed with invalid token")
					apiErrors["TOKEN_NOT_ALIVE"].SendError(w)
					return
				}
			}

			// since all checks on the access token have passed, try to get the
			// information of the userinfo endpoint for the access token if the
			// userinfo endpoint was discovered
			if config.UserinfoEndpointUri == nil {
				log.Error().Msg("no userinfo endpoint discovered during init")
				apiErrors["USERINFO_ENDPOINT_DISCOVERY_FAIL"].SendError(w)
				return
			}

			// now check if the token already has been inspected one
			var info userInfo
			info = userInfo{}
			req, err := http.NewRequest("GET", *config.UserinfoEndpointUri, nil)
			if err != nil {
				log.Error().Err(err).Msg("unable to build request for userinfo endpoint")
				backendErrors.SendInternalError(err, w)
				return
			}
			req.Header.Set("Authorization", headerValues[0])

			// now execute the request
			httpClient := http.Client{}
			res, err = httpClient.Do(req)
			if err != nil {
				log.Error().Err(err).Msg("unable to request userinfo endpoint")
				backendErrors.SendInternalError(err, w)
				return
			}

			// now parse the response
			var userinfo map[string]interface{}
			err = json.NewDecoder(res.Body).Decode(&userinfo)
			if err != nil {
				log.Error().Err(err).Msg("unable to parse userinfo response")
				backendErrors.SendInternalError(err, w)
				return
			}

			// now check the userinfo response for the following fields:
			//   - sub [the subject the userinfo is for]
			//   - preferred_username [the username of the user]
			//   - name [the name of the user]
			//   - groups [the groups the user is a member of]
			subject, isSet := userinfo["sub"].(string)
			if !isSet {
				log.Error().Msg("userinfo response did not contain the subject the userinfo is issued for")
				w.WriteHeader(500)
				return
			}
			// now validate that the userinfo is issued for the same subject
			// as the jwt
			if accessToken.Subject() != subject {
				w.WriteHeader(401)
				return
			}

			// now get the username and real name of the user accessing
			// the backend
			username, isSet := userinfo["preferred_username"].(string)
			if !isSet {
				log.Warn().Msg("user not identifiable, disallowing access")
				w.WriteHeader(500)
				return
			}
			info.Username = username

			// now get the full name
			fullName, isSet := userinfo["name"].(string)
			if !isSet {
				log.Warn().Msg("user not identifiable, disallowing access")
				w.WriteHeader(500)
				return
			}
			info.FullName = fullName

			// now check the groups the user is a member of the userinfo
			// response and has the correct type
			groups, isSet := userinfo["groups"].([]interface{})
			if !isSet {
				log.Warn().Msg("no groups in userinfo response found with type []interface{}")
			}
			info.Groups = groups
			// todo write object into database

			// now add the collected values to the context
			ctx := context.WithValue(r.Context(), "username", info.Username)
			ctx = context.WithValue(ctx, "fullName", info.FullName)
			ctx = context.WithValue(ctx, "groups", info.Groups)
			ctx = context.WithValue(ctx, "subject", accessToken.Subject())

			// now serve the request with the updated context
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		// return the handler function as http.Handler
		return http.HandlerFunc(handlerFunction)
	}
}

// ExtractUserInfo allows the retrieval of the available user information that
// has been set by the [OpenIDConnectJWTAuthentication] middleware. If any of
// the expected user information fields is empty, it will return an errors
func ExtractUserInfo(r *http.Request) (username string, fullName string, subject string, groups []string, err error) {
	// get the context from the request
	ctx := r.Context()
	// now check if a username has been set
	username, usernameSet := ctx.Value("username").(string)
	if !usernameSet {
		return "", "", "", nil, errors.New("no username available with type string")
	}
	// now check if the full name has been set
	fullName, fullNameSet := ctx.Value("fullName").(string)
	if !fullNameSet {
		return "", "", "", nil, errors.New("no full name available with type string")
	}
	// now check if the subject of the user has been set
	subject, subjectSet := ctx.Value("subject").(string)
	if !subjectSet {
		return "", "", "", nil, errors.New("no subject available with type string")
	}
	// now check if the groups have been set
	groups, groupsSet := ctx.Value("groups").([]string)
	if !groupsSet {
		return "", "", "", nil, errors.New("no groups available with type []string")
	}
	// now return the result
	return username, fullName, subject, groups, nil
}
