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
