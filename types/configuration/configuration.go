/*
 * Copyright (c) 2023.  Jan Eike Suchard. All rights reserved
 * SPDX-License-Identifier: MIT
 */

package configurationTypes

type Configuration struct {
	RedisUri      string                     `toml:"redis-uri" comment:"the uri used to connect to the redis database used as cache"`
	Database      DatabaseConfiguration      `toml:"database" comment:"the configuration for the database storing the data"`
	OpenIdConnect OpenIDConnectConfiguration `toml:"openid-connect" comment:"the configuration for authenticating users accessing the backend"`
}
