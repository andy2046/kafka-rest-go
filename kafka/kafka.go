// Package kafka provides a thin wrapper around the REST API,
// providing a more convenient interface for accessing cluster metadata and producing and consuming data.
package kafka

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/pkg/errors"
)

type (
	// Kafka represents a Kafka REST API.
	Kafka struct {
		URL         string
		Timeout     time.Duration
		Accept      string
		ContentType string
		Format      Format
		Offset      Offset
		Version     Version
	}

	kafkaInterface interface {
		// TODO
	}

	// Broker data
	Broker struct {
		Brokers []int `json:"brokers"`
	}

	// Format is one of json, binary or avro
	Format string

	// Offset is either earliest or latest
	Offset string

	// Version is the API version
	Version string

	// ErrorMessage for API response
	ErrorMessage struct {
		ErrorCode int    `json:"error_code,omitempty"`
		Message   string `json:"message,omitempty"`
	}
)

const (
	// JSON formated consumer
	JSON = Format("json")
	// Binary formated consumer
	Binary = Format("binary")
	// Avro formated consumer
	Avro = Format("avro")

	// Earliest is the oldest offset for API v2
	Earliest = Offset("earliest")
	// Latest is the newest offset for API v2
	Latest = Offset("latest")

	// Smallest is the oldest offset for API v1
	Smallest = Offset("smallest")
	// Largest is the newest offset for API v1
	Largest = Offset("largest")

	// V2 is API v2
	V2 = Version("v2")

	// V1 is API v1
	V1 = Version("v1")
)

// Defaults for Kafka
var Defaults = Kafka{
	URL:         "http://localhost:8082",
	Timeout:     60 * time.Second,
	Accept:      "application/vnd.kafka+json, application/json",
	ContentType: "application/vnd.kafka+json",
	Format:      Binary,
	Offset:      Largest,
	Version:     V1,
}

// HTTPClient creates a new http.Client with timeout.
func (k *Kafka) HTTPClient() *http.Client {
	var netClient = &http.Client{
		Timeout: k.Timeout,
	}
	return netClient
}

// SetOption takes one or more option function and applies them in order to Kafka.
func (k *Kafka) SetOption(options ...func(*Kafka) error) error {
	for _, opt := range options {
		if err := opt(k); err != nil {
			return err
		}
	}
	return nil
}

// EarliestOffset set Offset to Earliest
func EarliestOffset(k *Kafka) error {
	k.Offset = Earliest
	return nil
}

// LatestOffset set Offset to Latest
func LatestOffset(k *Kafka) error {
	k.Offset = Latest
	return nil
}

// SmallestOffset set Offset to Earliest
func SmallestOffset(k *Kafka) error {
	k.Offset = Smallest
	return nil
}

// LargestOffset set Offset to Latest
func LargestOffset(k *Kafka) error {
	k.Offset = Largest
	return nil
}

// BinaryFormat set Format to Binary
func BinaryFormat(k *Kafka) error {
	k.Format = Binary
	return nil
}

// V1Version set Version to V1
func V1Version(k *Kafka) error {
	k.Version = V1
	return nil
}

// V2Version set Version to V2
func V2Version(k *Kafka) error {
	k.Version = V2
	return nil
}

// JSONFormat set Format to JSON
func JSONFormat(k *Kafka) error {
	k.Format = JSON
	return nil
}

// AvroFormat set Format to Avro
func AvroFormat(k *Kafka) error {
	k.Format = Avro
	return nil
}

// SetTimeout applies Timeout to Kafka.
func SetTimeout(timeout time.Duration) func(*Kafka) error {
	return func(k *Kafka) error {
		k.Timeout = timeout
		return nil
	}
}

// SetURL applies URL to Kafka.
func SetURL(url string) func(*Kafka) error {
	return func(k *Kafka) error {
		k.URL = url
		return nil
	}
}

// SetAccept applies Accept to Kafka.
func SetAccept(accept string) func(*Kafka) error {
	return func(k *Kafka) error {
		k.Accept = accept
		return nil
	}
}

// SetContentType applies ContentType to Kafka.
func SetContentType(contentType string) func(*Kafka) error {
	return func(k *Kafka) error {
		k.ContentType = contentType
		return nil
	}
}

func applyDefaults(k *Kafka) {
	k.URL = Defaults.URL
	k.Timeout = Defaults.Timeout
	k.Accept = Defaults.Accept
	k.ContentType = Defaults.ContentType
	k.Format = Defaults.Format
	k.Offset = Defaults.Offset
	k.Version = Defaults.Version
}

func validateStatusCode(res *http.Response, expectedStatusCode ...int) error {
	code := http.StatusOK
	if len(expectedStatusCode) > 0 {
		code = expectedStatusCode[0]
	}

	if res.StatusCode != code {
		errMsg := &ErrorMessage{}
		err := json.NewDecoder(res.Body).Decode(errMsg)
		if err != nil && err != io.EOF {
			body, _ := ioutil.ReadAll(res.Body)
			return errors.Wrap(errors.New("API Error"), string(body))
		}
		return errors.Wrap(errors.New("API Error"),
			fmt.Sprintf("StatusCode %v %v ErrorCode %v %v", res.StatusCode, res.Status, errMsg.ErrorCode, errMsg.Message))
	}
	return nil
}

func getConsumerGroup(cs *Consumers, consumerGroup []string) (string, error) {
	switch {
	case len(consumerGroup) > 0:
		return consumerGroup[0], nil
	case cs.ConsumerGroup != "":
		return cs.ConsumerGroup, nil
	default:
		return "", errors.New("Error: empty consumerGroup")
	}
}

func getTopicName(t *Topic, topicName []string) (string, error) {
	switch {
	case len(topicName) > 0:
		return topicName[0], nil
	case t != nil:
		return t.Name, nil
	default:
		return "", errors.New("Error: empty topicName")
	}
}

func doRequest(method, url string, client *http.Client, b *bytes.Buffer, requestHooker func(*http.Request), responseHooker func(*http.Response) error) error {
	req, err := http.NewRequest(method, url, b)
	if err != nil {
		return err
	}

	requestHooker(req)

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer closeBody(res)

	err = responseHooker(res)
	if err != nil {
		return err
	}

	return nil
}

// URLJoin joins url with path and return the whole url string.
func URLJoin(urlstr string, pathstrs ...string) (string, error) {
	u, err := url.Parse(urlstr)
	if err != nil {
		return "", err
	}
	str := strings.Join(pathstrs, "/")
	u.Path = path.Join(u.Path, str)
	return u.String(), nil
}

// New returns a Kafka instance with default setting.
func New(options ...func(*Kafka) error) (*Kafka, error) {
	var k Kafka
	applyDefaults(&k)
	err := k.SetOption(options...)
	if err != nil {
		return nil, err
	}
	return &k, nil
}

// NewTopics returns a Topics instance.
func (k *Kafka) NewTopics() *Topics {
	return &Topics{
		Kafka: k,
		List:  nil,
	}
}

// NewConsumers returns a Consumers instance.
func (k *Kafka) NewConsumers(consumerGroup ...string) *Consumers {
	cs := Consumers{
		Kafka: k,
	}

	if len(consumerGroup) > 0 {
		cs.ConsumerGroup = consumerGroup[0]
	}

	return &cs
}

// Broker returns the brokers.
func (k *Kafka) Broker() (*Broker, error) {
	client := k.HTTPClient()
	url, err := URLJoin(k.URL, "brokers")
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", k.Accept)

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer closeBody(res)

	err = validateStatusCode(res)
	if err != nil {
		return nil, err
	}

	b := &Broker{}

	err = json.NewDecoder(res.Body).Decode(b)
	if err != nil && err != io.EOF {
		return nil, err
	}

	return b, nil
}

func closeBody(res *http.Response) {
	// Drain and close the body to let the Transport reuse the connection
	io.Copy(ioutil.Discard, res.Body)
	res.Body.Close()
}
