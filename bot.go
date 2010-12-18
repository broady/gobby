package main

import (
	"flag"
	irc "github.com/nf/goirc/client"
	"log"
	"os"
	"strings"
)

var (
	nick    = flag.String("nick", "", "irc nick")
	server  = flag.String("server", "", "server host:port")
	pass    = flag.String("pass", "", "server password")
	channel = flag.String("chan", "", "channel")
	secret  = flag.String("secret", "", "secret")
)

func authLog(c *irc.Conn) {
	Tail(c, "#log", "/var/log/auth.log", func(s string) bool {
		return !strings.Contains(s, "cron:session")
	})
}

func main() {
	flag.Parse()
	if *nick == "" || *server == "" || *channel == "" {
		flag.Usage()
		os.Exit(1)
	}
	c := irc.New(*nick, *nick, *nick)
	c.SSL = true
	authLog(c)
	PlusOne(c, *channel)
	if *secret != "" {
		Auth(c, *channel, *secret)
	}
	if err := c.Connect(*server, *pass); err != nil {
		log.Exit(err)
	}
	select {
	}
}
