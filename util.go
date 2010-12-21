package main

import (
	irc "github.com/nf/goirc/client"
	"regexp"
	"sync"
)

var autoJoinChannels = make(map[string]bool)
var autoJoinLock sync.Mutex

func AutoJoin(conn *irc.Conn, channel string) {
	autoJoinLock.Lock()
	defer autoJoinLock.Unlock()
	if autoJoinChannels[channel] {
		return
	}
	autoJoinChannels[channel] = true
	conn.AddHandler("CONNECTED", func(c *irc.Conn, l *irc.Line) {
		conn.Join(channel)
	})
}

var canonRe = regexp.MustCompile("[^a-zA-Z]+")
const canonMaxLen = 10

func Canonical(nick string) string {
	nick = canonRe.ReplaceAllString(nick, "")
	if len(nick) > canonMaxLen {
		nick = nick[:canonMaxLen]
	}
	return nick
}
