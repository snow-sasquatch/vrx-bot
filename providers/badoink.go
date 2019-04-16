package providers

import (
	"bytes"
	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
	"net/http"
)

const badoinkURL = "https://badoinkvr.com"

type Badoink struct {
	c *http.Client
}

func NewBadoinkProvider(c *http.Client) (p *Badoink) {
	p = &Badoink{c}
	return p
}

func (p *Badoink) Content() {
	req, err := http.NewRequest(http.MethodGet, badoinkURL, nil)
	if err != nil {
		log.Warn("Creating Request Failed: %v", err)
	}
	accessCookie := &http.Cookie{Name: "legal_age", Value: "true"}
	req.AddCookie(accessCookie)
	res, err := p.c.Do(req)
	defer res.Body.Close()
	if err != nil {
		log.Warn("Request Badoink failed: %v", err)
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(res.Body)
	s := buf.String()
	print(s)
}
