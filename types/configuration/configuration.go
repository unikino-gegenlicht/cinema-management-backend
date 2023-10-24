/*
 * Copyright (c) 2023.  Jan Eike Suchard. All rights reserved
 * SPDX-License-Identifier: MIT
 */

package configurationTypes

type Configuration struct {
	MongoDbUri    string                     `toml:"mongoDB-uri" comment:"the uri used to connect to the mongodb in which all data is stored"`
	RedisUri      string                     `toml:"redis-uri" comment:"the uri used to connect to the redis database used as cache"`
	OpenIdConnect OpenIDConnectConfiguration `toml:"openid-connect" comment:"the configuration for authenticating users accessing the backend"`
}
