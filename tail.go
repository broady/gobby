package main

import (
	"bufio"
	irc "github.com/fluffle/goirc/client"
	"log"
	"os"
	"time"
)

func Tail(conn *irc.Conn, channel, filename string, fn func(string) bool) {
	AutoJoin(conn, channel)
	f, err := os.Open(filename, 0, 0)
	if err != nil {
		log.Println(err)
		return
	}
	_, err = f.Seek(0, 2)
	if err != nil {
		log.Println(err)
		return
	}
	b := bufio.NewReader(f)
	go func() {
		for {
			l, err := b.ReadBytes('\n')
			if err != nil {
				if err == os.EOF {
					time.Sleep(1e9)
					b = bufio.NewReader(f)
					continue
				}
				log.Println(err)
				break
			}
			s := string(l[:len(l)-1])
			if fn(s) {
				conn.Privmsg(channel, s)
			}
		}
	}()
}
