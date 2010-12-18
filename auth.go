package main

import (
	"crypto/md5"
	"fmt"
	irc "github.com/nf/goirc/client"
	"strings"
)

func Auth(conn *irc.Conn, channel, secret string) {
	conn.AddHandler("353", func(c *irc.Conn, l *irc.Line) {
		if l.Args[2] != channel {
			return
		}
		for _, nick := range strings.Split(l.Args[3], " ", -1) {
			if !strings.HasPrefix(nick, "@") {
				continue
			}
			nick = nick[1:]
			hash := makeHash(channel, c.Me.Nick, nick, secret)
			c.Privmsg(nick, ".authop "+channel+" "+hash)
		}
	})
	conn.AddHandler("PRIVMSG", func(c *irc.Conn, l *irc.Line) {
		if strings.HasPrefix(l.Nick, "#") {
			return
		}
		p := strings.Split(l.Args[1], " ", -1)
		if len(p) != 3 || p[0] != ".authop" || p[1] != channel {
			return
		}
		if p[2] != makeHash(channel, l.Nick, c.Me.Nick, secret) {
			return
		}
		c.Mode(channel, l.Nick+" +o")
	})
}

func makeHash(channel, requester, requestee, secret string) string {
	h := md5.New()
	fmt.Fprintf(h, "%s-%s-%s-%s", channel, requester, requestee, secret)
	return fmt.Sprintf("%x", h.Sum())
}
