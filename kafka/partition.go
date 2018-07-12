package kafka

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/pkg/errors"
)

type (
	// Replica data
	Replica struct {
		Broker int  `json:"broker"`
		Leader bool `json:"leader"`
		InSync bool `json:"in_sync"`
	}

	// Partition data
	Partition struct {
		Partition int       `json:"partition"`
		Leader    int       `json:"leader"`
		Replicas  []Replica `json:"replicas"`
	}

	// Partitions data
	Partitions struct {
		Kafka *Kafka
		Topic *Topic
		List  *[]Partition
	}
)

// Partitions lists partitions for the topic.
func (ps *Partitions) Partitions(topicName ...string) (*[]Partition, error) {
	tn, err := getTopicName(ps.Topic, topicName)
	if err != nil {
		return nil, err
	}

	client := ps.Kafka.HTTPClient()
	url, err := URLJoin(ps.Kafka.URL, "topics", tn, "partitions")
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", ps.Kafka.Accept)

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	err = validateStatusCode(res)
	if err != nil {
		return nil, err
	}

	pss := &[]Partition{}

	err = json.NewDecoder(res.Body).Decode(pss)
	if err != nil {
		return nil, err
	}

	return pss, nil
}

// Partition returns the Partition with provided partitionID.
func (ps *Partitions) Partition(partitionID int, topicName ...string) (*Partition, error) {
	tn, err := getTopicName(ps.Topic, topicName)
	if err != nil {
		return nil, err
	}

	client := ps.Kafka.HTTPClient()
	url, err := URLJoin(ps.Kafka.URL, "topics", tn, "partitions", strconv.Itoa(partitionID))
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", ps.Kafka.Accept)

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	err = validateStatusCode(res)
	if err != nil {
		return nil, err
	}

	p := &Partition{}

	err = json.NewDecoder(res.Body).Decode(p)
	if err != nil {
		return nil, err
	}

	return p, nil
}

// Produce post message to the Partition with provided id.
func (ps *Partitions) Produce(id int, message *ProducerMessage, topicName ...string) (*ProducerResponse, error) {
	if ps.Kafka.Format == Avro && message.ValueSchema == "" && message.ValueSchemaID == 0 {
		return nil, fmt.Errorf("Must provide a value schema or value schema id for Avro format")
	}

	tn, err := getTopicName(ps.Topic, topicName)
	if err != nil {
		return nil, err
	}

	client := ps.Kafka.HTTPClient()
	url, err := URLJoin(ps.Kafka.URL, "topics", tn, "partitions", strconv.Itoa(id))
	if err != nil {
		return nil, err
	}

	b := &bytes.Buffer{}
	json.NewEncoder(b).Encode(message)
	req, err := http.NewRequest("POST", url, b)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", ps.Kafka.Accept)
	req.Header.Set("Content-Type", ps.Kafka.ContentType)

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	err = validateStatusCode(res)
	if err != nil {
		return nil, err
	}

	pr := &ProducerResponse{}

	err = json.NewDecoder(res.Body).Decode(pr)
	if err != nil {
		return nil, err
	}

	cause, hasError := errors.New("Error: produce messages to partition "+strconv.Itoa(id)+"of topic"+tn), false
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
