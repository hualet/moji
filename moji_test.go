package moji

import (
	"os"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	env := os.Getenv(appCodeEnv)
	os.Unsetenv(appCodeEnv)

	c, err := NewClient()
	if c != nil || err != ErrAppCodeEnv {
		t.Error("should report err failed to fetch environment variable.")
	}

	os.Setenv(appCodeEnv, env)
}

func TestConditionByLatLong(t *testing.T) {
	// c, err := NewClient()
	// if err != nil {
	// 	t.Error(err)
	// }

	// cond, err := c.ConditionByLatLong("30.6", "114.4")
}

func TestUnmarshalConditionData(t *testing.T) {
	data := []byte(
		`{
      "code": 0,
      "data": {
        "city": {
          "cityId": 284609,
          "counname": "中国",
          "name": "东城区",
          "pname": "北京市"
        },
        "condition": {
          "condition": "晴",
          "humidity": "26",
          "icon": "0",
          "temp": "25",
          "updatetime": "2016-09-01 10:00:00",
          "windDir": "北风",
          "windLevel": "3"
        }
      },
      "msg": "success",
      "rc": {
        "c": 0,
        "p": "success"
      }
    }`)

	ret, err := unmarshalConditionData(data)
	if err != nil {
		t.Error(err)
	}

	expCity := City{
		ID:       284609,
		Name:     "东城区",
		County:   "中国",
		Province: "北京市",
	}

	now := time.Now()
	expCondition := Condition{
		ID:            0,
		Name:          "晴",
		Humidity:      26,
		Temp:          25,
		UpdateTime:    time.Date(2016, time.September, 1, 10, 0, 0, 0, now.Location()),
		WindDirection: "北风",
		WindLevel:     3,
	}

	if ret.City != expCity {
		t.Error("parsed city not correct")
	}

	if ret.Condition != expCondition {
		t.Errorf("parsed condition not correct")
	}
}
