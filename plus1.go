package main

import (
	"bytes"
	"fmt"
	irc "github.com/fluffle/goirc/client"
	"regexp"
	"strings"
	"sync"
)

var (
	plusOneMap  = make(map[string]int)
	plusOneLock sync.Mutex
	plusOneRe   = regexp.MustCompile(`^[^ :]+: [+\-]1$`)
)

func PlusOne(conn *irc.Conn, channel string) {
	AutoJoin(conn, channel)
	conn.AddHandler("PRIVMSG", func(c *irc.Conn, l *irc.Line) {
		if l.Args[0] != channel {
			return
		}
		m := l.Args[1]
		if m == ".scores" {
			plusOneScores(c, channel)
			return
		}
		if !plusOneRe.MatchString(m) {
			return
		}
		nick := strings.Split(m, ":", 2)[0]
		if c.GetNick(nick) == nil {
			return
		}
		nick = Canonical(nick)
		if len(nick) == 0 {
			return
		}
		n := 1
		if m[len(m)-2] == '-' {
			n = -1
		}
		plusOneCount(nick, n)
	})
}

func plusOneCount(nick string, n int) {
	plusOneLock.Lock()
	plusOneMap[nick] += n
	plusOneLock.Unlock()
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
