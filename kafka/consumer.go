package kafka

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

type (
	// ConsumerRequest is the metadata needed to create a consumer instance
	ConsumerRequest struct {
		Format     Format `json:"format"`
		Offset     Offset `json:"auto.offset.reset"`
		AutoCommit string `json:"auto.commit.enable"` // true or false
		Name       string `json:"name,omitempty"`
	}

	// ConsumerInstance data
	ConsumerInstance struct {
		ConsumerName string `json:"instance_id"`
		BaseURI      string `json:"base_uri"`
	}

	// Consumers data
	Consumers struct {
		Kafka         *Kafka
		ConsumerGroup string
		List          []ConsumerInstance
	}

	// ConsumerOffset are the offsets to commit
	ConsumerOffset struct {
		Partition int    `json:"partition"`
		Offset    int64  `json:"offset"`
		Topic     string `json:"topic"`
		Metadata  string `json:"metadata,omitempty"`
	}

	// ConsumerOffsets data
	ConsumerOffsets struct {
		Offsets []ConsumerOffset `json:"offsets"`
	}

	// ConsumerOffsetsPartitions are the partitions for consumer committed offsets
	ConsumerOffsetsPartitions struct {
		Partitions []ConsumerPartitions `json:"partitions"`
	}

	// ConsumerPartitions are the partitions for consumer
	ConsumerPartitions struct {
		Partition int    `json:"partition"`
		Topic     string `json:"topic"`
	}

	// TopicSubscription data, topic_pattern and topics fields are mutually exclusive
	TopicSubscription struct {
		Topics       *TopicsSubscription
		TopicPattern *TopicPatternSubscription
	}

	// TopicsSubscription with topics
	TopicsSubscription struct {
		Topics []string `json:"topics"`
	}

	// TopicPatternSubscription with topic pattern
	TopicPatternSubscription struct {
		TopicPattern string `json:"topic_pattern"`
	}

	// Message is a single Kafka message
	Message struct {
		Topic     string          `json:"topic"`
		Key       json.RawMessage `json:"key"`
		Value     json.RawMessage `json:"value"`
		Partition int             `json:"partition"`
		Offset    int64           `json:"offset"`
	}

	// Argument is the argument for both method Records and Messages
	Argument struct {
		Timeout       int
		TopicName     string
		MaxBytes      int
		ConsumerName  string
		ConsumerGroup string
	}
)

// NewConsumer creates a new consumer instance in the consumer group.
func (cs *Consumers) NewConsumer(consumerRequest *ConsumerRequest, consumerGroup ...string) (*ConsumerInstance, error) {
	cg, err := getConsumerGroup(cs, consumerGroup)
	if err != nil {
		return nil, err
	}

	url, err := URLJoin(cs.Kafka.URL, "consumers", cg)
	if err != nil {
		return nil, err
	}

	b := &bytes.Buffer{}
	json.NewEncoder(b).Encode(consumerRequest)

	requestHooker := func(req *http.Request) {
		req.Header.Set("Accept", cs.Kafka.Accept)
		req.Header.Set("Content-Type", cs.Kafka.ContentType)
	}

	ci := &ConsumerInstance{}

	responseHooker := func(res *http.Response) error {
		err = validateStatusCode(res)
		if err != nil {
			return err
		}

		err = json.NewDecoder(res.Body).Decode(ci)
		if err != nil && err != io.EOF {
			return err
		}
		return nil
	}

	err = doRequest("POST", url, cs.Kafka.HTTPClient(), b, requestHooker, responseHooker)
	if err != nil {
		return nil, err
	}

	return ci, nil
}

// DeleteConsumer destroy the consumer instance.
func (cs *Consumers) DeleteConsumer(consumerName string, consumerGroup ...string) error {
	cg, err := getConsumerGroup(cs, consumerGroup)
	if err != nil {
		return err
	}

	client := cs.Kafka.HTTPClient()
	url, err := URLJoin(cs.Kafka.URL, "consumers", cg, "instances", consumerName)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", cs.Kafka.Accept)
	req.Header.Set("Content-Type", cs.Kafka.ContentType)

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer closeBody(res)

	err = validateStatusCode(res, http.StatusNoContent)
	if err != nil {
		return err
	}

	return nil
}

// CommitOffsets commits a list of offsets for the consumer.
func (cs *Consumers) CommitOffsets(consumerOffsets *ConsumerOffsets, consumerName string, consumerGroup ...string) error {
	cg, err := getConsumerGroup(cs, consumerGroup)
	if err != nil {
		return err
	}

	client := cs.Kafka.HTTPClient()
	url, err := URLJoin(cs.Kafka.URL, "consumers", cg, "instances", consumerName, "offsets")
	if err != nil {
		return err
	}

	b := &bytes.Buffer{}
	json.NewEncoder(b).Encode(consumerOffsets)
	req, err := http.NewRequest("POST", url, b)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", cs.Kafka.Accept)
	req.Header.Set("Content-Type", cs.Kafka.ContentType)

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer closeBody(res)

	err = validateStatusCode(res)
	if err != nil {
		return err
	}

	return nil
}

// Offsets get the last committed offsets for the given partitions.
func (cs *Consumers) Offsets(consumerOffsetsPartitions *ConsumerOffsetsPartitions, consumerName string, consumerGroup ...string) (*ConsumerOffsets, error) {
	cg, err := getConsumerGroup(cs, consumerGroup)
	if err != nil {
		return nil, err
	}

	client := cs.Kafka.HTTPClient()
	url, err := URLJoin(cs.Kafka.URL, "consumers", cg, "instances", consumerName, "offsets")
	if err != nil {
		return nil, err
	}

	b := &bytes.Buffer{}
	json.NewEncoder(b).Encode(consumerOffsetsPartitions)
	req, err := http.NewRequest("GET", url, b)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", cs.Kafka.Accept)

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer closeBody(res)

	err = validateStatusCode(res)
	if err != nil {
		return nil, err
	}

	cresp := &ConsumerOffsets{}

	err = json.NewDecoder(res.Body).Decode(cresp)
	if err != nil && err != io.EOF {
		return nil, err
	}

	return cresp, nil
}

// Subscribe to the given list of topics or a topic pattern.
func (cs *Consumers) Subscribe(topicSubscription *TopicSubscription, useTopicPattern bool, consumerName string, consumerGroup ...string) error {
	cg, err := getConsumerGroup(cs, consumerGroup)
	if err != nil {
		return err
	}

	client := cs.Kafka.HTTPClient()
	url, err := URLJoin(cs.Kafka.URL, "consumers", cg, "instances", consumerName, "subscription")
	if err != nil {
		return err
	}

	b := &bytes.Buffer{}
	switch {
	case useTopicPattern:
		json.NewEncoder(b).Encode(topicSubscription.TopicPattern)
	case !useTopicPattern:
		json.NewEncoder(b).Encode(topicSubscription.Topics)
	}

	req, err := http.NewRequest("POST", url, b)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", cs.Kafka.Accept)
	req.Header.Set("Content-Type", cs.Kafka.ContentType)

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer closeBody(res)

	err = validateStatusCode(res, http.StatusNoContent)
	if err != nil {
		return err
	}

	return nil
}

// Subscriptions get the current subscribed list of topics.
func (cs *Consumers) Subscriptions(consumerName string, consumerGroup ...string) (*TopicsSubscription, error) {
	cg, err := getConsumerGroup(cs, consumerGroup)
	if err != nil {
		return nil, err
	}

	client := cs.Kafka.HTTPClient()
	url, err := URLJoin(cs.Kafka.URL, "consumers", cg, "instances", consumerName, "subscription")
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", cs.Kafka.Accept)

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer closeBody(res)

	err = validateStatusCode(res, http.StatusOK)
	if err != nil {
		return nil, err
	}

	tsub := &TopicsSubscription{}

	err = json.NewDecoder(res.Body).Decode(tsub)
	if err != nil && err != io.EOF {
		return nil, err
	}

	return tsub, nil
}

// Unsubscribe from topics currently subscribed.
func (cs *Consumers) Unsubscribe(consumerName string, consumerGroup ...string) error {
	cg, err := getConsumerGroup(cs, consumerGroup)
	if err != nil {
		return err
	}

	client := cs.Kafka.HTTPClient()
	url, err := URLJoin(cs.Kafka.URL, "consumers", cg, "instances", consumerName, "subscription")
	if err != nil {
		return err
	}

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", cs.Kafka.Accept)

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer closeBody(res)

	err = validateStatusCode(res, http.StatusNoContent)
	if err != nil {
		return err
	}

	return nil
}

// Assign manually assign a list of partitions to this consumer.
func (cs *Consumers) Assign(consumerOffsetsPartitions *ConsumerOffsetsPartitions, consumerName string, consumerGroup ...string) error {
	cg, err := getConsumerGroup(cs, consumerGroup)
	if err != nil {
		return err
	}

	client := cs.Kafka.HTTPClient()
	url, err := URLJoin(cs.Kafka.URL, "consumers", cg, "instances", consumerName, "assignments")
	if err != nil {
		return err
	}

	b := &bytes.Buffer{}
	json.NewEncoder(b).Encode(consumerOffsetsPartitions)

	req, err := http.NewRequest("POST", url, b)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", cs.Kafka.Accept)
	req.Header.Set("Content-Type", cs.Kafka.ContentType)

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer closeBody(res)

	err = validateStatusCode(res, http.StatusNoContent)
	if err != nil {
		return err
	}

	return nil
}

// Assignments get the list of partitions currently manually assigned to this consumer.
func (cs *Consumers) Assignments(consumerName string, consumerGroup ...string) (*ConsumerOffsetsPartitions, error) {
	cg, err := getConsumerGroup(cs, consumerGroup)
	if err != nil {
		return nil, err
	}

	client := cs.Kafka.HTTPClient()
	url, err := URLJoin(cs.Kafka.URL, "consumers", cg, "instances", consumerName, "assignments")
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", cs.Kafka.Accept)

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer closeBody(res)

	err = validateStatusCode(res)
	if err != nil {
		return nil, err
	}

	cofp := &ConsumerOffsetsPartitions{}

	err = json.NewDecoder(res.Body).Decode(cofp)
	if err != nil && err != io.EOF {
		return nil, err
	}

	return cofp, nil
}

// Seek overrides the fetch offsets that the consumer will use for the next set of records to fetch.
func (cs *Consumers) Seek(consumerOffsets *ConsumerOffsets, consumerName string, consumerGroup ...string) error {
	cg, err := getConsumerGroup(cs, consumerGroup)
	if err != nil {
		return err
	}

	client := cs.Kafka.HTTPClient()
	url, err := URLJoin(cs.Kafka.URL, "consumers", cg, "instances", consumerName, "positions")
	if err != nil {
		return err
	}

	b := &bytes.Buffer{}
	json.NewEncoder(b).Encode(consumerOffsets)

	req, err := http.NewRequest("POST", url, b)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", cs.Kafka.Accept)
	req.Header.Set("Content-Type", cs.Kafka.ContentType)

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer closeBody(res)

	err = validateStatusCode(res, http.StatusNoContent)
	if err != nil {
		return err
	}

	return nil
}

// SeekToBeginning seek to the first offset for each of the given partitions.
func (cs *Consumers) SeekToBeginning(consumerOffsetsPartitions *ConsumerOffsetsPartitions, consumerName string, consumerGroup ...string) error {
	cg, err := getConsumerGroup(cs, consumerGroup)
	if err != nil {
		return err
	}

	client := cs.Kafka.HTTPClient()
	url, err := URLJoin(cs.Kafka.URL, "consumers", cg, "instances", consumerName, "positions", "beginning")
	if err != nil {
		return err
	}

	b := &bytes.Buffer{}
	json.NewEncoder(b).Encode(consumerOffsetsPartitions)

	req, err := http.NewRequest("POST", url, b)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", cs.Kafka.Accept)
	req.Header.Set("Content-Type", cs.Kafka.ContentType)

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer closeBody(res)

	err = validateStatusCode(res, http.StatusNoContent)
	if err != nil {
		return err
	}

	return nil
}

// SeekToEnd seek to the last offset for each of the given partitions.
func (cs *Consumers) SeekToEnd(consumerOffsetsPartitions *ConsumerOffsetsPartitions, consumerName string, consumerGroup ...string) error {
	cg, err := getConsumerGroup(cs, consumerGroup)
	if err != nil {
		return err
	}

	client := cs.Kafka.HTTPClient()
	url, err := URLJoin(cs.Kafka.URL, "consumers", cg, "instances", consumerName, "positions", "end")
	if err != nil {
		return err
	}

	b := &bytes.Buffer{}
	json.NewEncoder(b).Encode(consumerOffsetsPartitions)

	req, err := http.NewRequest("POST", url, b)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", cs.Kafka.Accept)
	req.Header.Set("Content-Type", cs.Kafka.ContentType)

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer closeBody(res)

	err = validateStatusCode(res, http.StatusNoContent)
	if err != nil {
		return err
	}

	return nil
}

// Records fetch message for the topics or partitions specified via API v2.
// Records arguments include Timeout (optional) MaxBytes (optional) ConsumerName  ConsumerGroup.
// Timeout is the number of milliseconds for the underlying request to fetch the records. Default to 5000ms.
// MaxBytes is the maximum number of bytes of unencoded keys and values that should be included in the response. Default is unlimited.
func (cs *Consumers) Records(recordsArg Argument) ([]Message, error) {
	timeout := recordsArg.Timeout
	maxBytes := recordsArg.MaxBytes
	consumerName := recordsArg.ConsumerName
	consumerGroup := []string{recordsArg.ConsumerGroup}

	if consumerName == "" {
		return nil, errors.New("Error: empty ConsumerName")
	}

	cg, err := getConsumerGroup(cs, consumerGroup)
	if err != nil {
		return nil, err
	}

	client := cs.Kafka.HTTPClient()
	url, err := URLJoin(cs.Kafka.URL, "consumers", cg, "instances", consumerName, "records")
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", cs.Kafka.Accept)

	if !(timeout == 0 && maxBytes == 0) {
		q := req.URL.Query()
		if timeout != 0 {
			q.Add("timeout", strconv.Itoa(timeout))
		}
		if maxBytes != 0 {
			q.Add("max_bytes", strconv.Itoa(maxBytes))
		}
		req.URL.RawQuery = q.Encode()
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer closeBody(res)

	err = validateStatusCode(res)
	if err != nil {
		return nil, err
	}

	m := []Message{}

	err = json.NewDecoder(res.Body).Decode(m)
	if err != nil && err != io.EOF {
		return nil, err
	}

	return m, nil
}

// Messages consume messages from a topic via API v1.
// Messages arguments include MaxBytes (optional) TopicName ConsumerName  ConsumerGroup.
// MaxBytes is the maximum number of bytes of unencoded keys and values that should be included in the response. Default is unlimited.
func (cs *Consumers) Messages(messagesArg Argument) ([]Message, error) {
	topicName := messagesArg.TopicName
	maxBytes := messagesArg.MaxBytes
	consumerName := messagesArg.ConsumerName
	consumerGroup := []string{messagesArg.ConsumerGroup}

	if consumerName == "" {
		return nil, errors.New("Error: empty ConsumerName")
	}

	cg, err := getConsumerGroup(cs, consumerGroup)
	if err != nil {
		return nil, err
	}

	client := cs.Kafka.HTTPClient()
	url, err := URLJoin(cs.Kafka.URL, "consumers", cg, "instances", consumerName, "topics", topicName)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", cs.Kafka.Accept)

	if maxBytes != 0 {
		q := req.URL.Query()
		q.Add("max_bytes", strconv.Itoa(maxBytes))
		req.URL.RawQuery = q.Encode()
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer closeBody(res)

	err = validateStatusCode(res)
	if err != nil {
		return nil, err
	}

	m := []Message{}

	err = json.NewDecoder(res.Body).Decode(m)
	if err != nil && err != io.EOF {
		return nil, err
	}

	return m, nil
}

// Poll keep polling messages from a topic.
// the interval to poll messages is every interval ms, onMessage to handle polled messages.
// returned func is for cancellation.
func (cs *Consumers) Poll(interval time.Duration, messagesArg Argument, onMessage func(error, []Message)) func() {

	timedout := make(chan struct{})
	cancelFunc := func() {
		close(timedout)
	}

	go func() {
		t := time.NewTicker((time.Duration(interval) * time.Millisecond))
		var messages []Message
		var err error
		for {
			select {
			case <-timedout:
				t.Stop()
				return
			case <-t.C:
				break
			}

			if cs.Kafka.Version == V1 {
				messages, err = cs.Messages(messagesArg)
			} else {
				messages, err = cs.Records(messagesArg)
			}

			if err != nil {
				onMessage(err, nil)
				return
			}
			onMessage(nil, messages)
		}
	}()

	return cancelFunc
}
