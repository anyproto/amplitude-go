// Package amplitude provides access to the Amplitude API
package amplitude

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"time"
)

const eventEndpoint = "https://api2.amplitude.com/2/httpapi"
const identifyEndpoint = "https://api.amplitude.com/identify"

// Client manages the communication to the Amplitude API
type Client struct {
	eventEndpoint    string
	identifyEndpoint string
	key              string
	client           *http.Client
}

type Event struct {
	Adid               string                 `json:"adid,omitempty"`
	AppVersion         string                 `json:"app_version,omitempty"`
	Carrier            string                 `json:"carrier,omitempty"`
	City               string                 `json:"city,omitempty"`
	Country            string                 `json:"country,omitempty"`
	DeviceBrand        string                 `json:"device_brand,omitempty"`
	DeviceId           string                 `json:"device_id,omitempty"`
	DeviceManufacturer string                 `json:"device_manufacturer,omitempty"`
	DeviceModel        string                 `json:"device_model,omitempty"`
	DeviceType         string                 `json:"device_type,omitempty"`
	Dma                string                 `json:"dma,omitempty"`
	EventId            int                    `json:"event_id,omitempty"`
	EventProperties    map[string]interface{} `json:"event_properties,omitempty"`
	EventType          string                 `json:"event_type,omitempty"`
	Groups             map[string]interface{} `json:"groups,omitempty"`
	Ifda               string                 `json:"ifda,omitempty"`
	InsertId           string                 `json:"insert_id,omitempty"`
	Ip                 string                 `json:"ip,omitempty"`
	Language           string                 `json:"language,omitempty"`
	LocationLat        string                 `json:"location_lat,omitempty"`
	LocationLng        string                 `json:"location_lng,omitempty"`
	OsName             string                 `json:"os_name,omitempty"`
	OsVersion          string                 `json:"os_version,omitempty"`
	Paying             string                 `json:"paying,omitempty"`
	Platform           string                 `json:"platform,omitempty"`
	Price              float64                `json:"price,omitempty"`
	ProductId          string                 `json:"productId,omitempty"`
	Quantity           int                    `json:"quantity,omitempty"`
	Region             string                 `json:"region,omitempty"`
	Revenue            float64                `json:"revenue,omitempty"`
	RevenueType        string                 `json:"revenueType,omitempty"`
	SessionId          int64                  `json:"session_id,omitempty"`
	StartVersion       string                 `json:"start_version,omitempty"`
	Time               int64                  `json:"time,omitempty"`
	UserId             string                 `json:"user_id,omitempty"`
	UserProperties     map[string]interface{} `json:"user_properties,omitempty"`
}

type EventRequest struct {
	ApiKey string  `json:"api_key,omitempty"`
	Events []Event `json:"events,omitempty"`
}

type Identify struct {
	AppVersion         string                 `json:"app_version,omitempty"`
	Carrier            string                 `json:"carrier,omitempty"`
	City               string                 `json:"city,omitempty"`
	Country            string                 `json:"country,omitempty"`
	DeviceBrand        string                 `json:"device_brand,omitempty"`
	DeviceId           string                 `json:"device_id,omitempty"`
	DeviceManufacturer string                 `json:"device_manufacturer,omitempty"`
	DeviceModel        string                 `json:"device_model,omitempty"`
	DeviceType         string                 `json:"device_type,omitempty"`
	Dma                string                 `json:"dma,omitempty"`
	Groups             map[string]interface{} `json:"groups,omitempty"`
	Language           string                 `json:"language,omitempty"`
	OsName             string                 `json:"os_name,omitempty"`
	OsVersion          string                 `json:"os_version,omitempty"`
	Paying             string                 `json:"paying,omitempty"`
	Platform           string                 `json:"platform,omitempty"`
	Region             string                 `json:"region,omitempty"`
	StartVersion       string                 `json:"start_version,omitempty"`
	UserId             string                 `json:"user_id,omitempty"`
	UserProperties     map[string]interface{} `json:"user_properties,omitempty"`
}

// New client with API key
func New(key string) *Client {
	return &Client{
		eventEndpoint:    eventEndpoint,
		identifyEndpoint: identifyEndpoint,
		key:              key,
		client:           new(http.Client),
	}
}

func (c *Client) SetClient(client *http.Client) {
	c.client = client
}

func (c *Client) Events(events []Event) error {
	req := EventRequest{
		ApiKey: c.key,
		Events: events,
	}
	evJson, err := json.Marshal(req)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var buf bytes.Buffer
	g := gzip.NewWriter(&buf)
	if _, err := g.Write(evJson); err != nil {
		return err
	}
	if err = g.Close(); err != nil {
		return err
	}

	r, err := http.NewRequestWithContext(ctx, "POST", eventEndpoint, &buf)
	if err != nil {
		return err
	}

	r.Header.Set("content-type", "application/json")
	r.Header.Set("content-encoding", "gzip")
	resp, err := c.client.Do(r)
	if err == nil {
		return resp.Body.Close()
	}

	return err
}

func (c *Client) Event(msg Event) error {
	return c.Events([]Event{msg})
}

func (c *Client) Identify(msg Identify) error {
	msgJson, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	values := url.Values{}
	values.Add("api_key", c.key)
	values.Add("identification", string(msgJson))

	resp, err := c.client.PostForm(identifyEndpoint, values)
	if err == nil {
		resp.Body.Close()
	}
	return err
}
