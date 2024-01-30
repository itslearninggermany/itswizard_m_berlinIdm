package itswizard_m_berlinIdm

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func createUrl(query map[string]string, method string, dev bool) (*url.URL, error) {
	u := new(url.URL)
	var err error

	if dev {
		u, err = url.Parse(apiEndpointDev + method)
		if err != nil {
			return nil, err
		}
	} else {
		u, err = url.Parse(apiEndpoint + method)
		if err != nil {
			return nil, err
		}
	}

	u.Scheme = "https"
	q := u.Query()
	for k, v := range query {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()
	return u, nil
}

func (p *BlusdConnection) callAPI(t time.Time, sid string, service string, usetime bool) (out []byte, err error) {
	return
}

func (p *BlusdConnection) callAPIPerson(anzahlAnVergangenenMinuten int, sid string, service string, usetime bool) (out []byte, err error) {
	q := make(map[string]string)

	if anzahlAnVergangenenMinuten > 0 {
		loc, _ := time.LoadLocation("Europe/Berlin")
		// handle err
		time.Local = loc // -> this is setting the global timezone
		t := time.Now()

		t = t.Add(time.Duration(-anzahlAnVergangenenMinuten) * time.Minute)
		fmt.Println(t)
		q["dtLetzteAenderung"] = parseDate(t)
	} else {
		q["dtLetzteAenderung"] = "2021-01-01 00:01:00.000"
	}

	if sid != "" {
		q["schuleUID"] = sid
	}

	q["einschluss"] = "bezugspersonen"

	url, err := createUrl(q, service, p.dev)
	p.url = url
	if err != nil {
		return
	}
	log.Println(p.url.String())
	req, _ := http.NewRequest("GET", p.url.String(), nil)
	req.Header.Add("Authorization", p.authentificationHeader)
	req.Header.Add("Accept", "application/json")
	resp, err := p.client.Do(req)
	resp_body, _ := ioutil.ReadAll(resp.Body)

	if strings.Contains(string(resp_body), "listError") {
		var errorResp blusdError
		err = json.Unmarshal(resp_body, &errorResp)
		if err != nil {
			return nil, errors.New(err.Error() + " " + string(resp_body))
		}
		var errorString string
		for _, k := range errorResp.ListError {
			errorString = errorString + " --- " + k.ErrorText + " : " + k.ErrorSource + " : " + " : " + k.Reference
		}
		err = resp.Body.Close()
		if err != nil {
			return
		}
		return nil, errors.New(errorString)
	}

	out = resp_body
	err = resp.Body.Close()
	if err != nil {
		return
	}
	return
}

func parseDate(t time.Time) string {
	var date string
	date = strconv.Itoa(t.Year()) + "-"

	m := int(t.Month())

	if m < 10 {
		date = date + "0" + strconv.Itoa(m) + "-"
	} else {
		date = date + strconv.Itoa(m) + "-"
	}

	if t.Day() < 10 {
		date = date + "0" + strconv.Itoa(t.Day()) + " "
	} else {
		date = date + strconv.Itoa(t.Day()) + " "
	}

	h := int(t.Hour())
	if h < 10 {
		date = date + "0" + strconv.Itoa(h) + ":"
	} else {
		date = date + strconv.Itoa(h) + ":"
	}

	mi := int(t.Minute())
	if mi < 10 {
		date = date + "0" + strconv.Itoa(mi) + ":"
	} else {
		date = date + strconv.Itoa(mi) + ":"
	}

	date = date + "00.111"

	return date
}
