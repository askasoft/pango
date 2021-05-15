package email

import (
	"fmt"
	"net"
)

// DirectSender direct send email
type DirectSender struct {
	Sender
}

// DirectSend lookup MX record and direct send to recipter's mail server
func (ds *DirectSender) DirectSend(ms ...*Email) error {
	for i, m := range ms {
		if err := ds.directSend(m); err != nil {
			return fmt.Errorf("Failed to direct send email %d: %v", i+1, err)
		}
	}

	return nil
}

func (ds *DirectSender) directSend(mail *Email) error {
	rm := mail.GetRcptsByDomain()
	for domain, addrs := range rm {
		if err := ds.directSends(domain, addrs, mail); err != nil {
			return err
		}
	}

	return nil
}

func (ds *DirectSender) directSends(domain string, addrs []string, mail *Email) error {
	mxrs, err := net.LookupMX(domain)
	if err != nil {
		return err
	}

	for _, mx := range mxrs {
		if mx.Host == "" {
			continue
		}

		host := mx.Host
		port := 25

		if host[len(host)-1] == '.' {
			host = host[:len(host)-1]
		}
		ds.tlsConfig(host).ServerName = host

		err = ds.dial(host, port)
		if err == nil {
			err = ds.send(mail)
			ds.Close()
			break
		}
	}

	return err
}
