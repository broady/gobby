package main

import (
	"bytes"
	"fmt"
	irc "github.com/fluffle/goirc/client"
	"regexp"
	"strings"
	"sync"
)

var plusOneMap = make(map[string]int)
var plusOneLock sync.Mutex
var plusOneRe = regexp.MustCompile(`^[^ :]+: \+1$`)

func PlusOne(conn *irc.Conn, channel string) {
	AutoJoin(conn, channel)
	conn.AddHandler("PRIVMSG", func(c *irc.Conn, l *irc.Line) {
		if l.Args[0] != channel {
			return
		}
		if l.Args[1] == ".scores" {
			plusOneScores(c, channel)
			return
		}
		if !plusOneRe.MatchString(l.Args[1]) {
			return
		}
		nick := strings.Split(l.Args[1], ":", 2)[0]
		nick = Canonical(nick)
		if len(nick) == 0 {
			return
		}
		plusOneCount(nick)
	})
}

func plusOneCount(nick string) {
	plusOneLock.Lock()
	defer plusOneLock.Unlock()
	plusOneMap[nick]++
}

func plusOneScores(c *irc.Conn, channel string) {
	plusOneLock.Lock()
	defer plusOneLock.Unlock()
	var buf bytes.Buffer
	for nick, score := range plusOneMap {
		if buf.Len() > 0 {
			fmt.Fprint(&buf, ", ")
		}
		fmt.Fprintf(&buf, "%s %d", nick, score)
	}
	c.Privmsg(channel, buf.String())
}
