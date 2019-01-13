package lib

import (
	"encoding/json"
	"io/ioutil"
	"testing"
)

var (
	invalidEmail string
	invalidPass  string
	validEmail   string
	validPass    string
)

type ClientConfig struct {
	InvalidEmail string `json:"invalidEmail"`
	InvalidPass  string `json:"invalidPass"`
	ValidEmail   string `json:"validEmail"`
	ValidPass    string `json:"validPass"`
}

func init() {
	raw, err := ioutil.ReadFile("./test_data.json")
	if err != nil {
		panic(err)
	}

	cc := &ClientConfig{}
	err = json.Unmarshal(raw, cc)
	if err != nil {
		panic(err)
	}

	invalidEmail = cc.InvalidEmail
	invalidPass = cc.InvalidPass
	validEmail = cc.ValidEmail
	validPass = cc.ValidPass
}

func TestLogin(t *testing.T) {
	c1 := NewClient()
	err := c1.Login(invalidEmail, validPass)
	if err == nil {
		t.Error("failed Login test(invalid email pattern)")
	}

	c2 := NewClient()
	err = c2.Login(validEmail, invalidPass)
	if err == nil {
		t.Error("failed Login test(invalid pass pattern)")
	}

	c3 := NewClient()
	err = c3.Login(validEmail, validPass)
	if err != nil {
		t.Error("failed Login test(valid email and valid pass pattern)")
	}

}

func TestAuth(t *testing.T) {
	c := NewClient()
	err := c.Login(validEmail, validPass)
	if err != nil {
		t.Error("failed")
	}

	token, partialKey, err := c.Auth1()
	if token == "" || partialKey == "" || err != nil {
		t.Error("failed Auth1 test")
	}

	err = c.Auth2(token, partialKey)
	if err != nil {
		t.Error("failed Auth2 test")
	}

}

func TestGetTimeFreeM3U8(t *testing.T) {
	var (
		invalidStationId = "XXXXstation"
		validStationId   = "MBS"
		invalidFt        = "99999999999999"
		validFt          = "20190108050000"
		invalidTo        = "99999999999999"
		validTo          = "20190108060000"
	)

	// invalid station id
	c1, token1, err := authenticated()
	if err != nil {
		t.Error("failed authenticated")
	}
	m3u8, err := c1.GetTimeFreeM3U8(invalidStationId, validFt, validTo, token1)
	if err == nil || m3u8 != "" {
		t.Error("failed GetTimeFreeM3U8(invalid station id pattern)")
	}

	// invalid ft
	c2, token2, err := authenticated()
	if err != nil {
		t.Error("failed authenticated")
	}
	m3u8, err = c2.GetTimeFreeM3U8(validStationId, invalidFt, validTo, token2)
	if err == nil || m3u8 != "" {
		t.Error("failed GetTimeFreeM3U8(invalid ft pattern)")
	}

	// invalid to
	c3, token3, err := authenticated()
	if err != nil {
		t.Error("failed authenticated")
	}
	m3u8, err = c3.GetTimeFreeM3U8(validStationId, validFt, invalidTo, token3)
	if err == nil || m3u8 != "" {
		t.Error("failed GetTimeFreeM3U8(invalid to pattern)")
	}

	//success
	c4, token4, err := authenticated()
	if err != nil {
		t.Error("failed authenticated")
	}
	m3u8, err = c4.GetTimeFreeM3U8(validStationId, validFt, validTo, token4)
	if err != nil || m3u8 == "" {
		t.Error("failed GetTimeFreeM3U8(success pattern)")
	}

}

func authenticated() (*Client, string, error) {
	c := NewClient()
	err := c.Login(validEmail, validPass)
	if err != nil {
		return nil, "", err
	}

	token, partialKey, err := c.Auth1()
	if err != nil {
		return nil, "", err
	}

	err = c.Auth2(token, partialKey)
	if err != nil {
		return nil, "", err
	}

	return c, token, nil
}
