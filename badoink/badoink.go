package badoink

import (
	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const badoinkURL = "https://badoinkvr.com"
const badoinkDataFolder = "badoink-data"
const badoinkFolderPath = "." + string(filepath.Separator) + badoinkDataFolder

type Badoink struct {
	c *http.Client
}

func NewProvider(c *http.Client) (p *Badoink) {
	p = &Badoink{c}
	//create data folder if it does not exist already
	if _, err := os.Stat(badoinkFolderPath); os.IsNotExist(err) {
		err = os.Mkdir(badoinkFolderPath, 0777)
		if err != nil {
			log.Warn("Couldn't create Badoink data directory: %v", err)
		}
	}
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
	if err != nil {
		log.Warn("Request Badoink failed: %v", err)
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Warn("GoQuery Reader error: %v", err)
	}

	//parse main doc to find new videos
	doc.Find(".video-card-image-container").Each(func(i int, s *goquery.Selection) {
		for _, n := range s.Nodes {
			for _, a := range n.Attr {
				if a.Key == "href" {
					handleVideoLink(a.Val)
				}
			}
		}
	})
}

func handleVideoLink(l string) {
	videoTitle := strings.Split(l, "/")[2]
	videoFolderPath := badoinkFolderPath + string(filepath.Separator) + videoTitle
	if _, err := os.Stat(videoFolderPath); os.IsNotExist(err) {
		err = os.Mkdir(videoFolderPath, 0777)
		//If the data folder for a video release does not exist we create one and download the video assets
		if err != nil {
			log.Warn("Couldn't create Badoink video directory: %v", err)
		}
		downloadAssets
	}
}

func downloadAssets(l, videoFolderPath string) error {

	return nil
}
