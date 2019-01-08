package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"
)

func main() {
	/* radiko login */
	email := ""
	pass := ""
	loginUrl := "https://radiko.jp/ap/member/login/login"

	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatal(err)
	}
	client := &http.Client{Jar: jar}

	params := url.Values{}
	params.Add("mail", email)
	params.Add("pass", pass)
	loginReq, err := http.NewRequest("POST", loginUrl, strings.NewReader(params.Encode()))
	if err != nil {
		log.Fatal(err)
	}
	loginReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	loginResp, err := client.Do(loginReq)
	if err != nil {
		log.Fatal(err)
	}
	if loginResp.StatusCode != http.StatusOK {
		log.Fatal("failed to login")
	}

	params.Del("mail")
	params.Del("pass")

	/* get auth token */
	authUrl := "https://radiko.jp/v2/api/auth1"
	authReq, err := http.NewRequest("GET", authUrl, nil)
	if err != nil {
		log.Fatal(err)
	}

	authReq.Header.Set("User-Agent", "curl/7.56.1")
	authReq.Header.Set("Accept", "*/*")
	authReq.Header.Set("pragma", "no-cache")
	authReq.Header.Set("x-radiko-app", "pc_html5")
	authReq.Header.Set("x-radiko-app-version", "0.0.1")
	authReq.Header.Set("x-radiko-device", "pc")
	authReq.Header.Set("x-radiko-user", "dummy_user")
	authResp, err := client.Do(authReq)
	if err != nil {
		log.Fatal(err)
	}
	if authResp.StatusCode != http.StatusOK {
		log.Fatal("failed to get auth token")
	}

	token := authResp.Header.Get("X-Radiko-AuthToken")
	keyLength, _ := strconv.Atoi(authResp.Header.Get("X-Radiko-KeyLength"))
	keyOffset, _ := strconv.Atoi(authResp.Header.Get("X-Radiko-KeyOffset"))

	// TODO: 正規表現でsrcからちゃんと取得するようにする
	key := "bcd151073c03b352e1ef2fd66c32209da9ca0afa"
	cnvKey := key[keyOffset : keyOffset+keyLength]
	partialKey := base64.StdEncoding.EncodeToString([]byte(cnvKey))

	/* enable auto token  */
	auth2Url := "https://radiko.jp/v2/api/auth2"
	auth2Req, err := http.NewRequest("GET", auth2Url, nil)
	if err != nil {
		log.Fatal(err)
	}

	auth2Req.Header.Set("X-Radiko-AuthToken", token)
	auth2Req.Header.Set("x-radiko-device", "pc")
	auth2Req.Header.Set("x-radiko-partialkey", partialKey)
	auth2Req.Header.Set("x-radiko-user", "dummy_user")
	auth2Resp, err := client.Do(auth2Req)
	if err != nil {
		log.Fatal(err)
	}
	if auth2Resp.StatusCode != http.StatusOK {
		log.Fatal("failed to enable auto token")
	}

	/* get TimeFreeM3U8 */
	//params.Add("station_id", "MBS")
	//params.Add("l", "15")
	//params.Add("ft", "20190108050000")
	//params.Add("to", "20190108060000")
	//M3U8url := "https://radiko.jp/v2/api/ts/playlist.m3u8"
	//reqM3U8, err := http.NewRequest("POST", M3U8url, strings.NewReader(params.Encode()))
	M3U8url := "https://radiko.jp/v2/api/ts/playlist.m3u8?station_id=MBS&l=15&ft=20190108050000&to=20190108060000"
	reqM3U8, err := http.NewRequest("POST", M3U8url, nil)
	if err != nil {
		log.Fatal(err)
	}

	reqM3U8.Header.Set("pragma", "no-cache")
	reqM3U8.Header.Set("X-Radiko-AuthToken", token)
	respM3U8, err := client.Do(reqM3U8)
	if err != nil {
		log.Fatal(err)
	}
	if respM3U8.StatusCode != http.StatusOK {
		log.Fatal("faled to get TimeFreeM3U8")
	}

	b, err := ioutil.ReadAll(respM3U8.Body)
	if err != nil {
		log.Fatal()
	}

	fmt.Println("---------m3u8---------")
	fmt.Println(string(b))
	fmt.Println("---------m3u8---------")

}
