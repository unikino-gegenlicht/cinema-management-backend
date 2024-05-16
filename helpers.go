package main

import (
	"runtime/debug"
)

var Commit = func() string {
	commitHash := "unkown"
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range info.Settings {
			if setting.Key == "vcs.revision" {
				commitHash = setting.Value
			}
		}

		for _, setting := range info.Settings {
			if setting.Key == "vcs.modified" {
				if setting.Value == "true" {
					commitHash += " (modified)"
				}
			}
		}
	}

	return commitHash
}()

var CommitTime = func() string {
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range info.Settings {
			if setting.Key == "vcs.time" {
				return setting.Value
			}
		}
	}

	return "unkown"
}()
