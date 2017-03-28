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

func TestUnmarshalForecastData(t *testing.T) {
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
                "forecast": [
                {
                    "conditionDay": "多云",
                    "conditionIdDay": "1",
                    "conditionIdNight": "31",
                    "conditionNight": "多云",
                    "predictDate": "2016-09-01",
                    "tempDay": "27",
                    "tempNight": "18",
                    "updatetime": "2016-09-01 09:07:08",
                    "windDirDay": "西北风",
                    "windDirNight": "西北风",
                    "windLevelDay": "3",
                    "windLevelNight": "2"
                },
                {
                    "conditionDay": "多云",
                    "conditionIdDay": "1",
                    "conditionIdNight": "31",
                    "conditionNight": "多云",
                    "predictDate": "2016-09-02",
                    "tempDay": "27",
                    "tempNight": "20",
                    "updatetime": "2016-09-01 09:07:08",
                    "windDirDay": "北风",
                    "windDirNight": "东南风",
                    "windLevelDay": "2",
                    "windLevelNight": "1"
                },
                {
                    "conditionDay": "多云",
                    "conditionIdDay": "1",
                    "conditionIdNight": "31",
                    "conditionNight": "多云",
                    "predictDate": "2016-09-03",
                    "tempDay": "28",
                    "tempNight": "20",
                    "updatetime": "2016-09-01 09:07:08",
                    "windDirDay": "南风",
                    "windDirNight": "东风",
                    "windLevelDay": "2",
                    "windLevelNight": "2"
                },
                {
                    "conditionDay": "局部阵雨",
                    "conditionIdDay": "3",
                    "conditionIdNight": "31",
                    "conditionNight": "多云",
                    "predictDate": "2016-09-04",
                    "tempDay": "28",
                    "tempNight": "18",
                    "updatetime": "2016-09-01 09:07:08",
                    "windDirDay": "南风",
                    "windDirNight": "西风",
                    "windLevelDay": "2",
                    "windLevelNight": "2"
                },
                {
                    "conditionDay": "多云",
                    "conditionIdDay": "1",
                    "conditionIdNight": "31",
                    "conditionNight": "少云",
                    "predictDate": "2016-09-05",
                    "tempDay": "30",
                    "tempNight": "21",
                    "updatetime": "2016-09-01 09:07:08",
                    "windDirDay": "西北风",
                    "windDirNight": "西北风",
                    "windLevelDay": "2",
                    "windLevelNight": "2"
                },
                {
                    "conditionDay": "多云",
                    "conditionIdDay": "1",
                    "conditionIdNight": "31",
                    "conditionNight": "多云",
                    "predictDate": "2016-09-06",
                    "tempDay": "30",
                    "tempNight": "20",
                    "updatetime": "2016-09-01 09:07:08",
                    "windDirDay": "西北风",
                    "windDirNight": "西风",
                    "windLevelDay": "3",
                    "windLevelNight": "2"
                }
                ]
            },
            "msg": "success",
            "rc": {
                "c": 0,
                "p": "success"
            }
        }`)

	ret, err := unmarshalForecastData(data)
	if err != nil {
		t.Error(err)
	}

	expCity := City{
		ID:       284609,
		Name:     "东城区",
		County:   "中国",
		Province: "北京市",
	}

	if ret.City != expCity {
		t.Error("parsed city not correct")
	}

	if len(ret.Forecasts) != 6 {
		t.Error("forecasts count not match 6")
	}

	now := time.Now()
	loc := now.Location()

	expDate := time.Date(2016, time.September, 1, 0, 0, 0, 0, loc)

	expDayCondition := Condition{
		ID:            1,
		Name:          "多云",
		Humidity:      0,
		Temp:          27,
		UpdateTime:    time.Date(2016, time.September, 1, 9, 7, 8, 0, loc),
		WindDirection: "西北风",
		WindLevel:     3,
	}

	expNightCondition := Condition{
		ID:            31,
		Name:          "多云",
		Humidity:      0,
		Temp:          18,
		UpdateTime:    time.Date(2016, time.September, 1, 9, 7, 8, 0, loc),
		WindDirection: "西北风",
		WindLevel:     2,
	}

	first := ret.Forecasts[0]
	if first.PredictDate != expDate {
		t.Error("parsed predict date not correct")
	}

	if *first.DayCondition != expDayCondition {
		t.Error("parsed day condition not correct")
	}

	if *first.NightCondition != expNightCondition {
		t.Error("parsed night condition not correct")
	}
}
