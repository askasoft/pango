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
	ds.Username = ""
	ds.Password = ""
	ds.Port = 25
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
		err := ds.directSends(domain, addrs, mail)
		if err != nil {
			return fmt.Errorf("Failed to direct send email to %s: %v", domain, err)
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

		ds.Host = mx.Host
		ds.Port = 25

		if ds.Host[len(ds.Host)-1] == '.' {
			ds.Host = ds.Host[:len(ds.Host)-1]
		}
		ds.tlsConfig().ServerName = ds.Host

		err = ds.Dial()
		if err == nil {
			err = ds.Send(mail)
			ds.Close()
			break
		}
	}

	return err
}
