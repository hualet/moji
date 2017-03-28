package moji

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const (
	mojiHost         = "http://aliv1.data.moji.com"
	briefForcast     = "/whapi/json/aliweather/briefforecast6days"
	briefCurrent     = "/whapi/json/aliweather/briefcondition"
	token            = "0f9d7e535dfbfad15b8fd2a84fee3e36"
	appCodeEnv       = "MOJI_APP_CODE"
	updateTimeLayout = "2006-01-02 15:04:00"
)

// City represents a city in moji service.
type City struct {
	ID       int64  `json:"cityId"`
	Name     string `json:"name"`
	County   string `json:"counname"`
	Province string `json:"pname"`
}

// Condition is the current condition of the target city.
type Condition struct {
	ID            int       `json:"icon,string"`
	Name          string    `json:"condition"`
	Humidity      int       `json:"humidity,string"`
	Temp          int       `json:"temp,string"`
	UpdateTime    time.Time `json:"updatetime"`
	WindDirection string    `json:"windDir"`
	WindLevel     int       `json:"windLevel,string"`
}

// Forecast stands for one forecast records of the future days.
type Forecast struct {
	PredictDate    time.Time `json:"predictDate"`
	DayCondition   *Condition
	NightCondition *Condition
}

// ConditionData wraps the response of the condition request.
type ConditionData struct {
	City      City      `json:"city"`
	Condition Condition `json:"condition"`
}

// ForecastData wraps the response of the forecast request.
type ForecastData struct {
	City      City       `json:"city"`
	Forecasts []Forecast `json:"forecast"`
}

// Client do all the requests against moji service for you.
type Client struct {
	appCode    string
	token      string
	httpClient *http.Client
}

// NewClient creates a new instance of Client
func NewClient() (*Client, error) {
	m := &Client{}

	m.appCode = os.Getenv(appCodeEnv)
	if m.appCode == "" {
		return nil, ErrAppCodeEnv
	}

	m.token = token
	m.httpClient = &http.Client{
		Timeout: time.Second * 5,
	}

	return m, nil
}

// ConditionByLatLong fetch the current condition of target city by latitude and longitude.
func (m *Client) ConditionByLatLong(lat, long string) (*ConditionData, error) {
	// fetch current condition
	form := m.generatePostData(lat, long)

	req, _ := http.NewRequest("POST", mojiHost+briefCurrent, strings.NewReader(form.Encode()))
	resp, err := m.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	ret, err := unmarshalConditionData(content)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

// ForecastByLatLong fetch the forecasts data of the target city by latitude and longitude.
func (m *Client) ForecastByLatLong(lat, long string) (*ForecastData, error) {
	// fetch current condition
	form := m.generatePostData(lat, long)

	req, _ := http.NewRequest("POST", mojiHost+briefCurrent, strings.NewReader(form.Encode()))
	resp, err := m.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	ret, err := unmarshalForecastData(content)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

// UnmarshalJSON unmarshals Forecast from json data.
func (f *Forecast) UnmarshalJSON(data []byte) error {
	conv := &struct {
		*Forecast
		ConditionDay       string    `json:"conditionDay"`
		ConditionNight     string    `json:"conditionNight"`
		ConditionIDDay     int       `json:"conditionIdDay"`
		ConditionIDNight   int       `json:"conditionIdNight"`
		UpdateTime         time.Time `json:"updatetime"`
		TempDay            int       `json:"tempDay"`
		TempNight          int       `json:"tempNight"`
		WindDirectionDay   string    `json:"windDirDay"`
		WindDirectionNight string    `jsong:"windDirNight"`
		WindLevelDay       int       `json:"windLevelDay"`
		WindLevelNight     int       `json:"windLevelNight"`
	}{}
	err := json.Unmarshal(data, conv)

	if err != nil {
		return err
	}

	f.DayCondition = &Condition{
		ID:            conv.ConditionIDDay,
		Name:          conv.ConditionDay,
		Humidity:      0,
		Temp:          conv.TempDay,
		UpdateTime:    conv.UpdateTime,
		WindDirection: conv.WindDirectionDay,
		WindLevel:     conv.WindLevelDay,
	}

	f.NightCondition = &Condition{
		ID:            conv.ConditionIDNight,
		Name:          conv.ConditionNight,
		Humidity:      0,
		Temp:          conv.TempNight,
		UpdateTime:    conv.UpdateTime,
		WindDirection: conv.WindDirectionNight,
		WindLevel:     conv.WindLevelNight,
	}

	return nil
}

// UnmarshalJSON unmarshals Condition from json data.
func (c *Condition) UnmarshalJSON(data []byte) error {
	type Alias Condition

	aux := &struct {
		UpdateTime string `json:"updatetime"`
		*Alias
	}{
		Alias: (*Alias)(c),
	}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}

	tm := time.Now()
	updateTime, err := time.ParseInLocation(updateTimeLayout, aux.UpdateTime, tm.Location())
	if err != nil {
		return err
	}

	c.UpdateTime = updateTime

	return nil
}

func (m *Client) generatePostData(lat, long string) url.Values {
	ret := url.Values{}

	ret.Add("lat", lat)
	ret.Add("lon", long)
	ret.Add("token", m.token)

	return ret
}

func unmarshalForecastData(content []byte) (*ForecastData, error) {
	ret := &ForecastData{}

	val := struct {
		*ForecastData
		Code int `json:"code"`
	}{ret, 0}

	err := json.Unmarshal(content, &val)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func unmarshalConditionData(content []byte) (*ConditionData, error) {
	ret := &ConditionData{}

	val := struct {
		Condition *ConditionData `json:"data"`
		Code      int            `json:"code"`
	}{
		Condition: ret,
	}

	err := json.Unmarshal(content, &val)
	if err != nil {
		return nil, err
	}

	return ret, nil
}