package enums

import (
	"kw_tool/util/constants"
)

type Platform uint

func ToPlatform(s string) (p Platform) {
	switch s {
	case Google.String():
		return Google
	case Amazon.String():
		return Amazon
	}
	return
}

const (
	Google Platform = iota
	Amazon
)

func (u Platform) URI() string {
	return [...]string{"http://suggestqueries.google.com/complete/search"}[u]
}

func (u Platform) String() string {
	return constants.Platforms[u]
}

func (u Platform) Options() map[string]string {
	d := []map[string]string{
		{
			"output": fireFox.String(),
			"hl":     en.String(),
		},
	}
	return d[u]
}

type Output uint

const (
	fireFox Output = iota
	chrome
)

func (o Output) String() string {
	return [...]string{"firefox", "chrome"}[o]
}

type Lang uint

const (
	en Lang = iota
)

func (l Lang) String() string {
	return [...]string{"en"}[l]
}
