package main

import (
	irc "github.com/fluffle/goirc/client"
	"regexp"
	"xml"
	"http"
)

var (
	youtubeRe = regexp.MustCompile(`http.*youtube[^\?]+\?([^ #$]+)`)
)

type Video struct {
	XMLName		xml.Name "entry"
	Title		string
}

func YouTuber(conn *irc.Conn, channel string) {
	AutoJoin(conn, channel)
	conn.AddHandler("PRIVMSG", func(c *irc.Conn, l *irc.Line) {
		if l.Args[0] != channel {
			return
		}
		m := l.Args[1]
		url := youtubeRe.FindStringSubmatch(m)
		if url == nil {
			return
		}
		q, _ := http.ParseQuery(url[1])
		key, found := q["v"]
		if !found {
			return
		}
		printMetadataForKey(c, channel, key[0])
	})
}

func printMetadataForKey(c *irc.Conn, channel string, key string) {
	r, _, _ := http.Get("http://gdata.youtube.com/feeds/api/videos/" + key)
	if r.StatusCode != http.StatusOK {
		return
	}
	var m Video = Video{xml.Name{}, ""}
	xml.Unmarshal(r.Body, &m)
	c.Privmsg(channel, "  -> " + m.Title)
}
