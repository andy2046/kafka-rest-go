package kafka

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/pkg/errors"
)

type (
	// Topic data
	Topic struct {
		Name       string          `json:"name"`
		Configs    json.RawMessage `json:"configs"`
		Partitions []Partition     `json:"partitions"`
	}

	// Topics data
	Topics struct {
		Kafka *Kafka
		List  []Topic
	}

	// TopicNames data
	TopicNames []string

	// ProducerMessage is the wrapper for the Topic / Partition data
	ProducerMessage struct {
		KeySchema   string `json:"key_schema,omitempty"`
		KeySchemaID int    `json:"key_schema_id,omitempty"`
		//either value schema or value schema id must be provided for avro messages
		ValueSchema   string           `json:"value_schema,omitempty"`
		ValueSchemaID int              `json:"value_schema_id,omitempty"`
		Records       []ProducerRecord `json:"records"`
	}

	// ProducerRecord is an individual message for Topic / Partition
	ProducerRecord struct {
		Key       json.RawMessage `json:"key,omitempty"`
		Value     json.RawMessage `json:"value"`
		Partition int             `json:"partition,omitempty"`
	}

	// ProducerResponse is the Topic / Partition response
	ProducerResponse struct {
		KeySchemaID   int               `json:"key_schema_id"`
		ValueSchemaID int               `json:"value_schema_id"`
		Offsets       []ProducerOffsets `json:"offsets"`
	}

	// ProducerOffsets are the resulting offsets for Topic / Partition
	ProducerOffsets struct {
		Partition int    `json:"partition"`
		Offset    int64  `json:"offset"`
		ErrorCode int64  `json:"error_code"`
		Error     string `json:"error"`
	}
)

// Topics lists all topics.
func (ts *Topics) Topics() ([]Topic, error) {
	ns, err := ts.Names()
	if err != nil {
		return nil, err
	}

	for _, n := range ns {
		t, err := ts.Topic(n)
		if err != nil {
			return ts.List, err
		}
		ts.List = append(ts.List, t)
	}

	return ts.List, nil
}

// Names lists all topic names.
func (ts *Topics) Names() (TopicNames, error) {
	client := ts.Kafka.HTTPClient()
	url, err := URLJoin(ts.Kafka.URL, "topics")
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", ts.Kafka.Accept)

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer closeBody(res)

	err = validateStatusCode(res)
	if err != nil {
		return nil, err
	}

	var tn TopicNames

	err = json.NewDecoder(res.Body).Decode(&tn)
	if err != nil && err != io.EOF {
		return nil, err
	}

	return tn, nil
}

// Topic returns the Topic with provided topicName.
func (ts *Topics) Topic(topicName string) (Topic, error) {
	client := ts.Kafka.HTTPClient()
	url, err := URLJoin(ts.Kafka.URL, "topics", topicName)
	if err != nil {
		return Topic{}, err
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return Topic{}, err
	}
	req.Header.Set("Accept", ts.Kafka.Accept)

	res, err := client.Do(req)
	if err != nil {
		return Topic{}, err
	}
	defer closeBody(res)

	err = validateStatusCode(res)
	if err != nil {
		return Topic{}, err
	}

	t := Topic{}

	err = json.NewDecoder(res.Body).Decode(&t)
	if err != nil && err != io.EOF {
		return Topic{}, err
	}

	return t, nil
}

// Produce post message to the Topic with provided topicName.
func (ts *Topics) Produce(topicName string, message *ProducerMessage) (*ProducerResponse, error) {
	if ts.Kafka.Format == Avro && message.ValueSchema == "" && message.ValueSchemaID == 0 {
		return nil, fmt.Errorf("Must provide a value schema or value schema id for Avro format")
	}

	client := ts.Kafka.HTTPClient()
	url, err := URLJoin(ts.Kafka.URL, "topics", topicName)
	if err != nil {
		return nil, err
	}

	b := &bytes.Buffer{}
	json.NewEncoder(b).Encode(message)
	req, err := http.NewRequest("POST", url, b)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", ts.Kafka.Accept)
	req.Header.Set("Content-Type", ts.Kafka.ContentType)

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer closeBody(res)

	err = validateStatusCode(res)
	if err != nil {
		return nil, err
	}

	pr := &ProducerResponse{}

	err = json.NewDecoder(res.Body).Decode(pr)
	if err != nil && err != io.EOF {
		return nil, err
	}

	cause, hasError := errors.New("Error: produce messages to topic "+topicName), false
	for _, offset := range pr.Offsets {
		if offset.ErrorCode != 0 {
			hasError = true
			cause = errors.Wrap(cause, offset.Error)
		}
	}

	if hasError {
		return pr, cause
	}

	return pr, nil
}

// NewPartitions returns a Partitions instance.
func (ts *Topics) NewPartitions(t ...*Topic) *Partitions {
	ps := Partitions{
		Kafka: ts.Kafka,
		Topic: nil,
		List:  nil,
	}

	if len(t) > 0 {
		ps.Topic = t[0]
	}

	return &ps
}

// Stringer returns formated string.
func Stringer(v interface{}) (string, error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}
