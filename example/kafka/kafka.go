package main

import (
	"encoding/json"
	"fmt"
	"time"

	K "github.com/andy2046/kafka-rest-go/kafka"
)

func main() {
	const topicName = "kafka-topic"
	const consumerGroup = "consumer-group"

	// create new Kafka instance with provided URL
	k, _ := K.New(K.SetURL("http://example.com:8082"))
	str, _ := K.Stringer(*k)
	fmt.Println(str)

	// list all brokers by id for Kafka instance
	brokers, err := k.Broker()
	if err != nil {
		panic(err)
	}
	fmt.Println(*brokers)

	// create new Topics instance
	topics := k.NewTopics()
	// get Topic instance by topic name
	topic, _ := topics.Topic(topicName)
	topicPointer := *topic
	fmt.Println(topicPointer.Name, string(topicPointer.Configs))
	// list Topic's Partitions
	for _, partition := range topicPointer.Partitions {
		str, _ := K.Stringer(*partition)
		fmt.Println(str)
	}

	// list Topics by name
	names, _ := topics.Names()
	for _, name := range *names {
		fmt.Println(name)
	}

	// create new Partitions instance
	ps := topics.NewPartitions()
	// get Partitions by topic name
	pss, err := ps.Partitions(topicName)
	if err != nil {
		panic(err)
	}
	// list Partition
	for _, p := range *pss {
		str, _ := K.Stringer(p)
		fmt.Println(str)
	}

	// get the Partition with provided partitionID
	p, _ := ps.Partition(0, topicName)
	s, _ := K.Stringer(*p)
	fmt.Println(s)

	// create new Record and Message to publish
	record := []*K.ProducerRecord{{Key: json.RawMessage(`"v5"`), Value: json.RawMessage(`"a2Fma2E="`)}, {Value: json.RawMessage(`"Z28="`)}}
	message := K.ProducerMessage{Records: record}

	// Produce message to partition
	pr, err := ps.Produce(0, &message, topicName)
	if err != nil {
		panic(err)
	}
	prstr, _ := K.Stringer(*pr)
	fmt.Println(prstr)

	pr, err = ps.Produce(1, &message, topicName)
	if err != nil {
		panic(err)
	}
	prstr, _ = K.Stringer(*pr)
	fmt.Println(prstr)

	// create Consumers instance
	c := k.NewConsumers()

	// create ConsumerRequest
	cr := &K.ConsumerRequest{
		Format:     K.Binary,
		Offset:     K.Smallest,
		AutoCommit: "true",
		Name:       "", // blank consumer name for auto generated one
	}

	// get new ConsumerInstance
	ci, err := c.NewConsumer(cr, consumerGroup)
	if err != nil {
		panic(err)
	}
	cistr, _ := K.Stringer(*ci)
	fmt.Println(cistr)

	cr2 := &K.ConsumerRequest{
		Format:     K.Binary,
		Offset:     K.Smallest,
		AutoCommit: "true",
		Name:       "",
	}

	ci2, err := c.NewConsumer(cr2, consumerGroup)
	if err != nil {
		panic(err)
	}
	cistr, _ = K.Stringer(*ci2)
	fmt.Println(cistr)

	// new Argument to consumer topic messages
	mArg := K.Argument{
		TopicName:     topicName,
		ConsumerGroup: consumerGroup,
		ConsumerName:  ci.ConsumerName,
	}

	// keep polling messages until cancelled or error
	cancelFunc := c.Poll(3000, mArg, func(err error, msg *[]K.Message) {
		if err != nil {
			panic(err)
		}
		fmt.Println("Records:")
		for _, m := range *msg {
			s, _ := K.Stringer(m)
			fmt.Println(s)
		}
	})

	time.Sleep(10 * time.Second)
	// cancel polling
	cancelFunc()

	// new Argument to consumer topic messages
	mArg2 := K.Argument{
		TopicName:     topicName,
		ConsumerGroup: consumerGroup,
		ConsumerName:  ci2.ConsumerName,
	}

	// consume topic messages one time
	msg2, err := c.Messages(mArg2)
	if err != nil {
		panic(err)
	}
	fmt.Println("Records:")
	for _, m2 := range *msg2 {
		s, _ := K.Stringer(m2)
		fmt.Println(s)
	}

	// delete consumer
	err = c.DeleteConsumer(ci.ConsumerName, consumerGroup)
	if err != nil {
		panic(err)
	}

	err = c.DeleteConsumer(ci2.ConsumerName, consumerGroup)
	if err != nil {
		panic(err)
	}

}
