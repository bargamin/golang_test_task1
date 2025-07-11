package application

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"

	"golang_test_task1/domain"
)

const (
	checkUrlTimeout = time.Second * 10

	//Possible number of parallel http calls to avoid a reject from firewall, DDOS protector. etc.
	maxHttpConnections = 20
)

type UrlScrapper struct {
	httpClient *http.Client
}

func NewUrlScrapper(httpClient *http.Client) *UrlScrapper {
	return &UrlScrapper{
		httpClient: httpClient,
	}
}

func (s *UrlScrapper) GetInfoByURL(ctx context.Context, rawURL string) (*domain.WebsiteInfo, error) {
	bodyBytes, err := s.getContentByURL(ctx, rawURL)
	if err != nil {
		return nil, err
	}

	return s.analyzeContent(ctx, string(bodyBytes), rawURL)
}

func (s *UrlScrapper) analyzeContent(
	ctx context.Context,
	body string,
	url string,
) (*domain.WebsiteInfo, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("parsing HTML: %w", err)
	}

	result := domain.WebsiteInfo{
		Url: url,
	}

	result.HTMLVersion = s.detectHTMLVersion(body)
	result.Title = strings.TrimSpace(doc.Find("title").First().Text())
	result.HeadingsCounts = s.countHeadings(doc)

	links := doc.Find("a[href]")
	result.InternalLinks, result.ExternalLinks, result.InaccessibleLinks = s.analyzeLinks(ctx, url, links)

	result.IsExistLoginForm = s.hasLoginForm(doc)

	return &result, nil
}

func (s *UrlScrapper) getContentByURL(ctx context.Context, u string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(ctx, checkUrlTimeout)
	defer cancel()

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetching %s: %w", u, err)
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("%+v\n", err)
		}
	}()

	if resp.StatusCode >= 400 {
		return nil, &HTTPError{Status: resp.StatusCode, Description: resp.Status}
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading body: %w", err)
	}
	return bodyBytes, nil
}

func (*UrlScrapper) detectHTMLVersion(raw string) string {
	raw = strings.ToLower(raw)
	switch {
	case strings.Contains(raw, "<!doctype html>"):
		return "HTML5"
	case strings.Contains(raw, "-//w3c//dtd html 4.01"):
		return "HTML 4.01"
	case strings.Contains(raw, "-//w3c//dtd xhtml"):
		return "XHTML"
	default:
		return "Unknown"
	}
}

func (*UrlScrapper) countHeadings(doc *goquery.Document) map[string]int {
	h := make(map[string]int, 6)
	for i := 1; i <= 6; i++ {
		tag := fmt.Sprintf("h%d", i)
		h[tag] = doc.Find(tag).Length()
	}

	return h
}

func (s *UrlScrapper) isSameOrigin(a, b *url.URL) bool {
	return a.Host == b.Host
}

func (s *UrlScrapper) hasLoginForm(doc *goquery.Document) bool {
	found := false

	doc.Find("form").EachWithBreak(func(_ int, f *goquery.Selection) bool {
		if f.Find("input[type='password']").Length() > 0 {
			found = true

			return false
		}

		return true
	})

	return found
}

func (s *UrlScrapper) analyzeLinks(ctx context.Context, baseRawUrl string, sel *goquery.Selection) (int, int, int) {
	var internal, external, inactive int

	checkLinkResult := make(chan bool)
	sem := make(chan struct{}, maxHttpConnections)
	wgCounter := &sync.WaitGroup{}
	wgCounter.Add(1)

	go func() {
		defer wgCounter.Done()

		for res := range checkLinkResult {
			if res {
				inactive++
			}
		}
	}()

	wgCheckers := &sync.WaitGroup{}
	baseUrl, _ := url.Parse(baseRawUrl)

	sel.Each(func(_ int, selection *goquery.Selection) {
		href, _ := selection.Attr("href")
		link, err := baseUrl.Parse(href)
		if err != nil {
			return
		}

		if s.isSameOrigin(baseUrl, link) {
			internal++
		} else {
			external++
		}

		sem <- struct{}{} //wait until semaphore is unblocked

		wgCheckers.Add(1)
		go s.isLinkUnavailable(ctx, link, checkLinkResult, sem, wgCheckers)
	})

	wgCheckers.Wait()
	close(checkLinkResult)
	wgCounter.Wait()

	return internal, external, inactive
}

func (s *UrlScrapper) isLinkUnavailable(
	ctx context.Context,
	link *url.URL,
	results chan<- bool,
	sem <-chan struct{},
	wg *sync.WaitGroup,
) {
	defer func() {
		<-sem
		wg.Done()
	}()

	ctx, cancel := context.WithTimeout(ctx, checkUrlTimeout)
	defer cancel()

	//makes query parameters URL encoded
	link.RawQuery = link.Query().Encode()

	req, _ := http.NewRequestWithContext(ctx, http.MethodHead, link.String(), nil)
	resp, err := s.httpClient.Do(req)
	if err != nil || resp.StatusCode >= 400 {
		log.Printf("URL is unavailable: %s : %v HTTP response: %v", link, err, resp)
		results <- true

		return
	}

	results <- false
}
