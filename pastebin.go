package pastebin

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

func string2int(s string) (int, error) {
	r := strings.Replace(s, `"`, ``, -1)

	q, err := strconv.ParseInt(r, 10, 64)
	if err != nil {
		return 0, err
	}

	return int(q), nil
}

func string2time(s string) (*time.Time, error) {
	r := strings.Replace(s, `"`, ``, -1)

	q, err := strconv.ParseInt(r, 10, 64)
	if err != nil {
		return nil, err
	}

	if q == 0 {
		return nil, nil
	}

	t := time.Unix(q, 0)
	return &t, nil
}

type Paste struct {
	ScrapeURL string     `json:"scrape_url,omitempty"`
	FullURL   string     `json:"full_url,omitempty"`
	Date      *time.Time `json:"date,omitempty"`
	Key       string     `json:"key,omitempty"`
	Size      int        `json:"size,omitempty"`

	Expire *time.Time `json:"expire,omitempty"`

	Title  string `json:"title,omitempty"`
	Syntax string `json:"syntax,omitempty"`
	User   string `json:"user,omitempty"`
	Hits   int    `json:"hits,omitempty"`

	Raw string
}

func (t *Paste) UnmarshalJSON(s []byte) (err error) {
	v := map[string]interface{}{}
	if err := json.Unmarshal(s, &v); err != nil {
		return err
	}

	t.ScrapeURL = v["scrape_url"].(string)
	t.FullURL = v["full_url"].(string)
	t.Key = v["key"].(string)
	t.Title = v["title"].(string)
	t.User = v["user"].(string)
	t.Syntax = v["syntax"].(string)

	if v, err := string2time(v["date"].(string)); err == nil {
		t.Date = v
	}

	if v, err := string2int(v["size"].(string)); err == nil {
		t.Size = v
	}

	if v, err := string2time(v["expire"].(string)); err == nil {
		t.Expire = v
	}

	if v, ok := v["hits"]; !ok {
	} else if v, err := string2int(v.(string)); err == nil {
		t.Hits = v
	}

	return
}

type PastebinClient struct {
	*http.Client
	BaseURL *url.URL
}

func New(u *url.URL) *PastebinClient {
	return &PastebinClient{
		Client:  http.DefaultClient,
		BaseURL: u,
	}
}

/*
type Scrape struct {
	PasteChan chan Paste
	ErrorChan chan error
}

func (pc *PastebinClient) Scrape(*context.Context) (*Scrape, error) {

	scrape := &Scrape{
		PasteChan: make(chan Paste, 100),
		ErrorChan: make(chan error, 100),
	}

	pasteChan

	go func() {
		select {
		case paste := <-pasteChan:
			raw, err := pc.GetRaw(paste.Key)
			if err != nil {
				scrape.ErrorChan <- err
				continue
			}

			defer raw.Close()

			b, err := ioutil.ReadAll(raw)
			fmt.Printf("%#v\n", string(b))

			paste.Raw = string(b)
			scrape.PasteChan <- paste
		}
	}()

	return scrape, nil
}
*/

func (pc *PastebinClient) Recent(size int) ([]Paste, error) {
	req, err := pc.NewRequest("GET", "/api_scraping.php")
	if err != nil {
		return nil, err
	}

	resp, err := pc.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Paste returned status code %d", resp.StatusCode)
	}

	pastes := []Paste{}

	err = json.NewDecoder(io.TeeReader(resp.Body, os.Stdout)).Decode(&pastes)
	if err != nil {
		return nil, err
	}

	return pastes, nil
}

func (pc *PastebinClient) GetRaw(key string) (io.ReadCloser, error) {
	req, err := pc.NewRequest("GET", fmt.Sprintf("/api_scrape_item.php?i=%s", key))
	if err != nil {
		return nil, err
	}

	resp, err := pc.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Paste returned status code %d", resp.StatusCode)
	}

	return resp.Body, nil
}

func (pc *PastebinClient) NewRequest(method, urlStr string) (*http.Request, error) {
	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	u := pc.BaseURL.ResolveReference(rel)

	req, err := http.NewRequest(method, u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "text/json; charset=UTF-8")
	req.Header.Add("Accept", "text/json")
	return req, nil
}
