package lib

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"

	gq "github.com/PuerkitoBio/goquery"
)

type Client struct {
	c *http.Client
}

func NewClient() *Client {
	jar, err := cookiejar.New(nil)
	if err != nil {
		panic(err)
	}

	return &Client{
		c: &http.Client{
			Jar: jar,
		},
	}
}

/* radiko login */
func (c *Client) Login(email, pass string) error {
	const loginUrl = "https://radiko.jp/ap/member/login/login"

	params := url.Values{}
	params.Add("mail", email)
	params.Add("pass", pass)
	req, err := http.NewRequest("POST", loginUrl, strings.NewReader(params.Encode()))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.c.Do(req)
	if err != nil {
		return err
	}

	doc, err := gq.NewDocumentFromResponse(resp)
	if err != nil {
		return err
	}
	if doc.Find("#member .login-area > .caution").Size() != 0 {
		return errors.New(fmt.Sprintf("invalid email(%s) or password(%s)", email, pass))
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New("failed to login")
	}

	return nil

}

/* get auth token */
func (c *Client) Auth1() (string, string, error) {
	const authUrl = "https://radiko.jp/v2/api/auth1"
	req, err := http.NewRequest("GET", authUrl, nil)
	if err != nil {
		return "", "", err
	}

	req.Header.Set("User-Agent", "curl/7.56.1")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("pragma", "no-cache")
	req.Header.Set("x-radiko-app", "pc_html5")
	req.Header.Set("x-radiko-app-version", "0.0.1")
	req.Header.Set("x-radiko-device", "pc")
	req.Header.Set("x-radiko-user", "dummy_user")
	resp, err := c.c.Do(req)
	if err != nil {
		return "", "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", "", errors.New("failed to get auth token")
	}

	token := resp.Header.Get("X-Radiko-AuthToken")
	keyLength, _ := strconv.Atoi(resp.Header.Get("X-Radiko-KeyLength"))
	keyOffset, _ := strconv.Atoi(resp.Header.Get("X-Radiko-KeyOffset"))

	// TODO: 正規表現でsrcからちゃんと取得するようにする
	key := "bcd151073c03b352e1ef2fd66c32209da9ca0afa"
	cnvKey := key[keyOffset : keyOffset+keyLength]
	partialKey := base64.StdEncoding.EncodeToString([]byte(cnvKey))

	return token, partialKey, nil
}

/* enable auto token  */
func (c *Client) Auth2(token, partialKey string) error {
	const auth2Url = "https://radiko.jp/v2/api/auth2"
	req, err := http.NewRequest("GET", auth2Url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("X-Radiko-AuthToken", token)
	req.Header.Set("x-radiko-device", "pc")
	req.Header.Set("x-radiko-partialkey", partialKey)
	req.Header.Set("x-radiko-user", "dummy_user")
	resp, err := c.c.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New("failed to enable auto token")
	}

	return nil
}

/* get TimeFreeM3U8 */
func (c *Client) GetTimeFreeM3U8(stationId, ft, to, token string) (string, error) {
	const M3U8url = "https://radiko.jp/v2/api/ts/playlist.m3u8?station_id=%s&ft=%s&to=%s"

	req, err := http.NewRequest("POST", fmt.Sprintf(M3U8url, stationId, ft, to), nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("pragma", "no-cache")
	req.Header.Set("X-Radiko-AuthToken", token)
	resp, err := c.c.Do(req)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", errors.New("faled to get TimeFreeM3U8")
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return strings.Split(string(b), "\n")[3], nil
}
