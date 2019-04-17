package badoink

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
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
	log.Infof("Badoink data folder is located at: %v", badoinkFolderPath)
	return p
}

func (p *Badoink) Content() {
	res := createRequest(badoinkURL, p)
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
					handleVideoLink(a.Val, p)
				}
			}
		}
	})
}

func handleVideoLink(l string, p *Badoink) {
	videoTitle := strings.Split(l, "/")[2]
	videoFolderPath := badoinkFolderPath + string(filepath.Separator) + videoTitle

	//If the data folder for a video release does not exist we create one and download the video assets
	if _, err := os.Stat(videoFolderPath); os.IsNotExist(err) {
		err = os.Mkdir(videoFolderPath, 0777)
		if err != nil {
			log.Warn("Couldn't create Badoink video directory: %v", err)
			return
		}
		downloadAssets(l, videoFolderPath, p)
	}
}

func downloadAssets(link, videoFolderPath string, p *Badoink) error {
	res := createRequest(badoinkURL+link, p)
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Warn("GoQuery Reader error: %v", err)
	}

	var wg sync.WaitGroup

	doc.Find(".gallery-item").Each(func(i int, s *goquery.Selection) {
		for _, n := range s.Nodes {
			for _, a := range n.Attr {
				if a.Key == "data-big-image" {
					wg.Add(1)
					go downloadImage(a.Val, videoFolderPath, i, &wg, p)
				}
			}
		}
	})
	wg.Wait()
	log.Info(fmt.Sprintf("downloaded picture pack: %s", videoFolderPath))
	return nil
}

func downloadImage(imageLink string, videoFolderPath string, pictureNum int, wg *sync.WaitGroup, p *Badoink) {
	defer wg.Done()
	response, e := p.c.Get(imageLink)
	if e != nil {
		log.Warn(fmt.Sprintf("Failed to Download image at %s: %v", imageLink, e))
	}
	defer response.Body.Close()

	file, err := os.Create(videoFolderPath + string(filepath.Separator) + strconv.Itoa(pictureNum) + ".jpg")
	if err != nil {
		log.Warn(fmt.Sprintf("Failed to create file for %s: %v", imageLink, err))
	}
	defer file.Close()

	// Use io.Copy to just dump the response body to the file. This supports huge files
	_, err = io.Copy(file, response.Body)
	if err != nil {
		log.Warn(fmt.Sprintf("Failed to save file from %s: %v", imageLink, err))
	}
}

func createRequest(URL string, p *Badoink) (res *http.Response) {
	req, err := http.NewRequest(http.MethodGet, URL, nil)
	if err != nil {
		log.Warn("Creating Request Failed: %v", err)
	}
	accessCookie := &http.Cookie{Name: "legal_age", Value: "true"}
	req.AddCookie(accessCookie)
	res, err = p.c.Do(req)
	if err != nil {
		log.Warn("Request Badoink failed: %v", err)
	}
	return res
}
