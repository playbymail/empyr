// Copyright (c) 2023-2025 Michael D Henderson. All rights reserved.

package main

import (
	"github.com/emersion/go-imap/v2"
	"github.com/emersion/go-imap/v2/imapclient"
	"github.com/playbymail/empyr/pkg/dotenv"
	"log"
	"net"
	"os"
	"time"
)

func main() {
	started := time.Now()

	log.SetFlags(log.LstdFlags | log.LUTC)

	if err := dotenv.Load("EMPYR"); err != nil {
		log.Fatalf("main: %+v\n", err)
	}
	imapHost := os.Getenv("EMPYR_IMAP_HOST")
	imapPort := os.Getenv("EMPYR_IMAP_PORT")
	imapAccount := os.Getenv("EMPYR_IMAP_ACCOUNT")
	if imapAccount == "" {
		log.Fatalf("%q is not set\n", "EMPYR_IMAP_ACCOUNT")
	}
	imapSecret := os.Getenv("EMPYR_IMAP_SECRET")
	if imapSecret == "" {
		log.Fatalf("%q is not set\n", "EMPYR_IMAP_SECRET")
	}

	run(true, imapHost, imapPort, imapAccount, imapSecret)
	log.Printf("completed in %v\n", time.Now().Sub(started))
}

func run(verbose bool, host, port, account, secret string) {
	if port == "" {
		port = "993"
	}
	imapHost := net.JoinHostPort(host, port)
	log.Printf("connecting to %q %q %q\n", imapHost, account, secret)

	c, err := imapclient.DialTLS(imapHost, nil)
	if err != nil {
		log.Fatalf("failed to dial IMAP server: %v", err)
	}
	defer c.Close()

	if err := c.Login(account, secret).Wait(); err != nil {
		log.Printf("failed to login: %v", err)
		return
	}

	mailboxes, err := c.List("", "%", nil).Collect()
	if err != nil {
		log.Printf("failed to list mailboxes: %v", err)
		return
	}
	if verbose {
		log.Printf("Found %v mailboxes", len(mailboxes))
		for _, mbox := range mailboxes {
			log.Printf(" - %v", mbox.Mailbox)
		}
	}

	for _, mbox := range []string{"INBOX", "INBOX.Received"} {
		selectedMbox, err := c.Select(mbox, nil).Wait()
		if err != nil {
			log.Printf("mbox: failed to select %q: %v\n", mbox, err)
			continue
		}
		if verbose {
			log.Printf("%s: contains %d messages", mbox, selectedMbox.NumMessages)
		}

		var id uint32
		for id = 1; id <= selectedMbox.NumMessages; id++ {
			seqSet := imap.SeqSetNum(id)
			fetchOptions := &imap.FetchOptions{Envelope: true, Flags: true}
			messages, err := c.Fetch(seqSet, fetchOptions).Collect()
			if err != nil {
				log.Printf("%s: msg %d: failed to fetch: %v\n", mbox, id, err)
				continue
			}
			if verbose {
				log.Printf("%s: msg %d:      to %q\n", mbox, id, messages[0].Envelope.To)
				log.Printf("%s: msg %d: subject %q\n", mbox, id, messages[0].Envelope.Subject)
				log.Printf("%s: msg %d:    from %q\n", mbox, id, messages[0].Envelope.From)
				log.Printf("%s: msg %d:  sender %q\n", mbox, id, messages[0].Envelope.Sender)
				for _, flag := range messages[0].Flags {
					log.Printf("%s: msg %d:    flag %q\n", mbox, id, flag)
				}
			}
		}
	}

	if err := c.Logout().Wait(); err != nil {
		log.Printf("failed to logout: %v", err)
	}
}
