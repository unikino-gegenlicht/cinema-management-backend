/*
 * Copyright (c) 2024.  Jan Eike Suchard. All rights reserved
 * SPDX-License-Identifier: MIT
 */

package configurationTypes

import "fmt"

type DatabaseConfiguration struct {
	// Host contains the hostname or ip address the postgres database lies on
	Host string `toml:"host" comment:"the host on which the postgres database lies on"`

	// Port contains the TCP port used to access the postgres database on the
	// Host
	Port uint16 `toml:"port" comment:"the port the postgres database listens on for connections"`

	// Username contains the username used to access the postgres database
	Username string `toml:"username" comment:"database user"`

	// Password contains the password for the Username used to access the
	// postgres database
	Password string `toml:"password" comment:"database password"`
}

func (dc DatabaseConfiguration) ToDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/cinema_management",
		dc.Username, dc.Password, dc.Host, dc.Port)
}
