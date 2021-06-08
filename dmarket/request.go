package dmarket

import (
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func (k Keys)request(body, path, method string) []byte {
	timestamp := strconv.Itoa(int(time.Now().UTC().Unix()))
	unsigned := method + path + body + timestamp
	signature, _ := k.Sign(unsigned)

	client := &http.Client{}
	req, _ := http.NewRequest(method, getRootUrl()+path, ioutil.NopCloser(strings.NewReader(body)))
	req.Header.Set("X-Sign-Date", timestamp)
	req.Header.Set("X-Request-Sign", "dmar ed25519 "+signature)
	req.Header.Set("X-Api-Key", k.Public)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	res, _ := client.Do(req)

	defer res.Body.Close()
	buf, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}
	return buf
}

