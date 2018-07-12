

# kafka
`import "github.com/andy2046/kafka-rest-go/kafka"`

* [Overview](#pkg-overview)
* [Index](#pkg-index)

## <a name="pkg-overview">Overview</a>
Package kafka provides a thin wrapper around the REST API,
providing a more convenient interface for accessing cluster metadata and producing and consuming data.




## <a name="pkg-index">Index</a>
* [Constants](#pkg-constants)
* [Variables](#pkg-variables)
* [func AvroFormat(k *Kafka) error](#AvroFormat)
* [func BinaryFormat(k *Kafka) error](#BinaryFormat)
* [func EarliestOffset(k *Kafka) error](#EarliestOffset)
* [func JSONFormat(k *Kafka) error](#JSONFormat)
* [func LargestOffset(k *Kafka) error](#LargestOffset)
* [func LatestOffset(k *Kafka) error](#LatestOffset)
* [func SetAccept(accept string) func(*Kafka) error](#SetAccept)
* [func SetContentType(contentType string) func(*Kafka) error](#SetContentType)
* [func SetTimeout(timeout time.Duration) func(*Kafka) error](#SetTimeout)
* [func SetURL(url string) func(*Kafka) error](#SetURL)
* [func SmallestOffset(k *Kafka) error](#SmallestOffset)
* [func Stringer(v interface{}) (string, error)](#Stringer)
* [func URLJoin(urlstr string, pathstrs ...string) (string, error)](#URLJoin)
* [func V1Version(k *Kafka) error](#V1Version)
* [func V2Version(k *Kafka) error](#V2Version)
* [type Argument](#Argument)
* [type Broker](#Broker)
* [type ConsumerInstance](#ConsumerInstance)
* [type ConsumerOffset](#ConsumerOffset)
* [type ConsumerOffsets](#ConsumerOffsets)
* [type ConsumerOffsetsPartitions](#ConsumerOffsetsPartitions)
* [type ConsumerPartitions](#ConsumerPartitions)
* [type ConsumerRequest](#ConsumerRequest)
* [type Consumers](#Consumers)
  * [func (cs *Consumers) Assign(consumerOffsetsPartitions *ConsumerOffsetsPartitions, consumerName string, consumerGroup ...string) error](#Consumers.Assign)
  * [func (cs *Consumers) Assignments(consumerName string, consumerGroup ...string) (*ConsumerOffsetsPartitions, error)](#Consumers.Assignments)
  * [func (cs *Consumers) CommitOffsets(consumerOffsets *ConsumerOffsets, consumerName string, consumerGroup ...string) error](#Consumers.CommitOffsets)
  * [func (cs *Consumers) DeleteConsumer(consumerName string, consumerGroup ...string) error](#Consumers.DeleteConsumer)
  * [func (cs *Consumers) Messages(messagesArg Argument) (*[]Message, error)](#Consumers.Messages)
  * [func (cs *Consumers) NewConsumer(consumerRequest *ConsumerRequest, consumerGroup ...string) (*ConsumerInstance, error)](#Consumers.NewConsumer)
  * [func (cs *Consumers) Offsets(consumerOffsetsPartitions *ConsumerOffsetsPartitions, consumerName string, consumerGroup ...string) (*ConsumerOffsets, error)](#Consumers.Offsets)
  * [func (cs *Consumers) Poll(interval time.Duration, messagesArg Argument, onMessage func(error, *[]Message)) func()](#Consumers.Poll)
  * [func (cs *Consumers) Records(recordsArg Argument) (*[]Message, error)](#Consumers.Records)
  * [func (cs *Consumers) Seek(consumerOffsets *ConsumerOffsets, consumerName string, consumerGroup ...string) error](#Consumers.Seek)
  * [func (cs *Consumers) SeekToBeginning(consumerOffsetsPartitions *ConsumerOffsetsPartitions, consumerName string, consumerGroup ...string) error](#Consumers.SeekToBeginning)
  * [func (cs *Consumers) SeekToEnd(consumerOffsetsPartitions *ConsumerOffsetsPartitions, consumerName string, consumerGroup ...string) error](#Consumers.SeekToEnd)
  * [func (cs *Consumers) Subscribe(topicSubscription *TopicSubscription, useTopicPattern bool, consumerName string, consumerGroup ...string) error](#Consumers.Subscribe)
  * [func (cs *Consumers) Subscriptions(consumerName string, consumerGroup ...string) (*TopicsSubscription, error)](#Consumers.Subscriptions)
  * [func (cs *Consumers) Unsubscribe(consumerName string, consumerGroup ...string) error](#Consumers.Unsubscribe)
* [type ErrorMessage](#ErrorMessage)
* [type Format](#Format)
* [type Kafka](#Kafka)
  * [func New(options ...func(*Kafka) error) (*Kafka, error)](#New)
  * [func (k *Kafka) Broker() (*Broker, error)](#Kafka.Broker)
  * [func (k *Kafka) HTTPClient() *http.Client](#Kafka.HTTPClient)
  * [func (k *Kafka) NewConsumers(consumerGroup ...string) *Consumers](#Kafka.NewConsumers)
  * [func (k *Kafka) NewTopics() *Topics](#Kafka.NewTopics)
  * [func (k *Kafka) SetOption(options ...func(*Kafka) error) error](#Kafka.SetOption)
* [type Message](#Message)
* [type Offset](#Offset)
* [type Partition](#Partition)
* [type Partitions](#Partitions)
  * [func (ps *Partitions) Partition(partitionID int, topicName ...string) (*Partition, error)](#Partitions.Partition)
  * [func (ps *Partitions) Partitions(topicName ...string) (*[]Partition, error)](#Partitions.Partitions)
  * [func (ps *Partitions) Produce(id int, message *ProducerMessage, topicName ...string) (*ProducerResponse, error)](#Partitions.Produce)
* [type ProducerMessage](#ProducerMessage)
* [type ProducerOffsets](#ProducerOffsets)
* [type ProducerRecord](#ProducerRecord)
* [type ProducerResponse](#ProducerResponse)
* [type Replica](#Replica)
* [type Topic](#Topic)
* [type TopicNames](#TopicNames)
* [type TopicPatternSubscription](#TopicPatternSubscription)
* [type TopicSubscription](#TopicSubscription)
* [type Topics](#Topics)
  * [func (ts *Topics) Names() (*TopicNames, error)](#Topics.Names)
  * [func (ts *Topics) NewPartitions(t ...*Topic) *Partitions](#Topics.NewPartitions)
  * [func (ts *Topics) Produce(topicName string, message *ProducerMessage) (*ProducerResponse, error)](#Topics.Produce)
  * [func (ts *Topics) Topic(topicName string) (*Topic, error)](#Topics.Topic)
  * [func (ts *Topics) Topics() ([]*Topic, error)](#Topics.Topics)
* [type TopicsSubscription](#TopicsSubscription)
* [type Version](#Version)


#### <a name="pkg-files">Package files</a>
[consumer.go](./kafka/consumer.go) [kafka.go](./kafka/kafka.go) [partition.go](./kafka/partition.go) [topic.go](./kafka/topic.go) 


## <a name="pkg-constants">Constants</a>
``` go
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
```

## <a name="pkg-variables">Variables</a>
``` go
var Defaults = Kafka{
    URL:         "http://localhost:8082",
    Timeout:     60 * time.Second,
    Accept:      "application/vnd.kafka+json, application/json",
    ContentType: "application/vnd.kafka+json",
    Format:      Binary,
    Offset:      Largest,
    Version:     V1,
}
```
Defaults for Kafka



## <a name="AvroFormat">func</a> [AvroFormat](./kafka/kafka.go?s=3096:3127#L159)
``` go
func AvroFormat(k *Kafka) error
```
AvroFormat set Format to Avro



## <a name="BinaryFormat">func</a> [BinaryFormat](./kafka/kafka.go?s=2704:2737#L135)
``` go
func BinaryFormat(k *Kafka) error
```
BinaryFormat set Format to Binary



## <a name="EarliestOffset">func</a> [EarliestOffset](./kafka/kafka.go?s=2261:2296#L111)
``` go
func EarliestOffset(k *Kafka) error
```
EarliestOffset set Offset to Earliest



## <a name="JSONFormat">func</a> [JSONFormat](./kafka/kafka.go?s=2997:3028#L153)
``` go
func JSONFormat(k *Kafka) error
```
JSONFormat set Format to JSON



## <a name="LargestOffset">func</a> [LargestOffset](./kafka/kafka.go?s=2595:2629#L129)
``` go
func LargestOffset(k *Kafka) error
```
LargestOffset set Offset to Latest



## <a name="LatestOffset">func</a> [LatestOffset](./kafka/kafka.go?s=2372:2405#L117)
``` go
func LatestOffset(k *Kafka) error
```
LatestOffset set Offset to Latest



## <a name="SetAccept">func</a> [SetAccept](./kafka/kafka.go?s=3513:3561#L181)
``` go
func SetAccept(accept string) func(*Kafka) error
```
SetAccept applies Accept to Kafka.



## <a name="SetContentType">func</a> [SetContentType](./kafka/kafka.go?s=3682:3740#L189)
``` go
func SetContentType(contentType string) func(*Kafka) error
```
SetContentType applies ContentType to Kafka.



## <a name="SetTimeout">func</a> [SetTimeout](./kafka/kafka.go?s=3202:3259#L165)
``` go
func SetTimeout(timeout time.Duration) func(*Kafka) error
```
SetTimeout applies Timeout to Kafka.



## <a name="SetURL">func</a> [SetURL](./kafka/kafka.go?s=3366:3408#L173)
``` go
func SetURL(url string) func(*Kafka) error
```
SetURL applies URL to Kafka.



## <a name="SmallestOffset">func</a> [SmallestOffset](./kafka/kafka.go?s=2483:2518#L123)
``` go
func SmallestOffset(k *Kafka) error
```
SmallestOffset set Offset to Earliest



## <a name="Stringer">func</a> [Stringer](./kafka/topic.go?s=4679:4723#L220)
``` go
func Stringer(v interface{}) (string, error)
```
Stringer returns formated string.



## <a name="URLJoin">func</a> [URLJoin](./kafka/kafka.go?s=5637:5700#L270)
``` go
func URLJoin(urlstr string, pathstrs ...string) (string, error)
```
URLJoin joins url with path and return the whole url string.



## <a name="V1Version">func</a> [V1Version](./kafka/kafka.go?s=2805:2835#L141)
``` go
func V1Version(k *Kafka) error
```
V1Version set Version to V1



## <a name="V2Version">func</a> [V2Version](./kafka/kafka.go?s=2900:2930#L147)
``` go
func V2Version(k *Kafka) error
```
V2Version set Version to V2




## <a name="Argument">type</a> [Argument](./kafka/consumer.go?s=2110:2239#L85)
``` go
type Argument struct {
    Timeout       int
    TopicName     string
    MaxBytes      int
    ConsumerName  string
    ConsumerGroup string
}
```
Argument is the argument for both method Records and Messages










## <a name="Broker">type</a> [Broker](./kafka/kafka.go?s=597:648#L36)
``` go
type Broker struct {
    Brokers []int `json:"brokers"`
}
```
Broker data










## <a name="ConsumerInstance">type</a> [ConsumerInstance](./kafka/consumer.go?s=439:550#L23)
``` go
type ConsumerInstance struct {
    ConsumerName string `json:"instance_id"`
    BaseURI      string `json:"base_uri"`
}
```
ConsumerInstance data










## <a name="ConsumerOffset">type</a> [ConsumerOffset](./kafka/consumer.go?s=723:903#L36)
``` go
type ConsumerOffset struct {
    Partition int    `json:"partition"`
    Offset    int64  `json:"offset"`
    Topic     string `json:"topic"`
    Metadata  string `json:"metadata,omitempty"`
}
```
ConsumerOffset are the offsets to commit










## <a name="ConsumerOffsets">type</a> [ConsumerOffsets](./kafka/consumer.go?s=931:1003#L44)
``` go
type ConsumerOffsets struct {
    Offsets []*ConsumerOffset `json:"offsets"`
}
```
ConsumerOffsets data










## <a name="ConsumerOffsetsPartitions">type</a> [ConsumerOffsetsPartitions](./kafka/consumer.go?s=1086:1178#L49)
``` go
type ConsumerOffsetsPartitions struct {
    Partitions []*ConsumerPartitions `json:"partitions"`
}
```
ConsumerOffsetsPartitions are the partitions for consumer committed offsets










## <a name="ConsumerPartitions">type</a> [ConsumerPartitions](./kafka/consumer.go?s=1236:1338#L54)
``` go
type ConsumerPartitions struct {
    Partition int    `json:"partition"`
    Topic     string `json:"topic"`
}
```
ConsumerPartitions are the partitions for consumer










## <a name="ConsumerRequest">type</a> [ConsumerRequest](./kafka/consumer.go?s=191:410#L15)
``` go
type ConsumerRequest struct {
    Format     Format `json:"format"`
    Offset     Offset `json:"auto.offset.reset"`
    AutoCommit string `json:"auto.commit.enable"` // true or false
    Name       string `json:"name,omitempty"`
}
```
ConsumerRequest is the metadata needed to create a consumer instance










## <a name="Consumers">type</a> [Consumers](./kafka/consumer.go?s=572:675#L29)
``` go
type Consumers struct {
    Kafka         *Kafka
    ConsumerGroup string
    List          []*ConsumerInstance
}
```
Consumers data










### <a name="Consumers.Assign">func</a> (\*Consumers) [Assign](./kafka/consumer.go?s=8397:8530#L365)
``` go
func (cs *Consumers) Assign(consumerOffsetsPartitions *ConsumerOffsetsPartitions, consumerName string, consumerGroup ...string) error
```
Assign manually assign a list of partitions to this consumer.




### <a name="Consumers.Assignments">func</a> (\*Consumers) [Assignments](./kafka/consumer.go?s=9305:9419#L402)
``` go
func (cs *Consumers) Assignments(consumerName string, consumerGroup ...string) (*ConsumerOffsetsPartitions, error)
```
Assignments get the list of partitions currently manually assigned to this consumer.




### <a name="Consumers.CommitOffsets">func</a> (\*Consumers) [CommitOffsets](./kafka/consumer.go?s=4014:4134#L172)
``` go
func (cs *Consumers) CommitOffsets(consumerOffsets *ConsumerOffsets, consumerName string, consumerGroup ...string) error
```
CommitOffsets commits a list of offsets for the consumer.




### <a name="Consumers.DeleteConsumer">func</a> (\*Consumers) [DeleteConsumer](./kafka/consumer.go?s=3267:3354#L138)
``` go
func (cs *Consumers) DeleteConsumer(consumerName string, consumerGroup ...string) error
```
DeleteConsumer destroy the consumer instance.




### <a name="Consumers.Messages">func</a> (\*Consumers) [Messages](./kafka/consumer.go?s=14717:14788#L619)
``` go
func (cs *Consumers) Messages(messagesArg Argument) (*[]Message, error)
```
Messages consume messages from a topic via API v1.
Messages arguments include MaxBytes (optional) TopicName ConsumerName  ConsumerGroup.
MaxBytes is the maximum number of bytes of unencoded keys and values that should be included in the response. Default is unlimited.




### <a name="Consumers.NewConsumer">func</a> (\*Consumers) [NewConsumer](./kafka/consumer.go?s=2313:2431#L95)
``` go
func (cs *Consumers) NewConsumer(consumerRequest *ConsumerRequest, consumerGroup ...string) (*ConsumerInstance, error)
```
NewConsumer creates a new consumer instance in the consumer group.




### <a name="Consumers.Offsets">func</a> (\*Consumers) [Offsets](./kafka/consumer.go?s=4852:5006#L208)
``` go
func (cs *Consumers) Offsets(consumerOffsetsPartitions *ConsumerOffsetsPartitions, consumerName string, consumerGroup ...string) (*ConsumerOffsets, error)
```
Offsets get the last committed offsets for the given partitions.




### <a name="Consumers.Poll">func</a> (\*Consumers) [Poll](./kafka/consumer.go?s=16008:16121#L677)
``` go
func (cs *Consumers) Poll(interval time.Duration, messagesArg Argument, onMessage func(error, *[]Message)) func()
```
Poll keep polling messages from a topic.
the interval to poll messages is every interval ms, onMessage to handle polled messages.
returned func is for cancellation.




### <a name="Consumers.Records">func</a> (\*Consumers) [Records](./kafka/consumer.go?s=13230:13299#L556)
``` go
func (cs *Consumers) Records(recordsArg Argument) (*[]Message, error)
```
Records fetch message for the topics or partitions specified via API v2.
Records arguments include Timeout (optional) MaxBytes (optional) ConsumerName  ConsumerGroup.
Timeout is the number of milliseconds for the underlying request to fetch the records. Default to 5000ms.
MaxBytes is the maximum number of bytes of unencoded keys and values that should be included in the response. Default is unlimited.




### <a name="Consumers.Seek">func</a> (\*Consumers) [Seek](./kafka/consumer.go?s=10210:10321#L442)
``` go
func (cs *Consumers) Seek(consumerOffsets *ConsumerOffsets, consumerName string, consumerGroup ...string) error
```
Seek overrides the fetch offsets that the consumer will use for the next set of records to fetch.




### <a name="Consumers.SeekToBeginning">func</a> (\*Consumers) [SeekToBeginning](./kafka/consumer.go?s=11074:11216#L479)
``` go
func (cs *Consumers) SeekToBeginning(consumerOffsetsPartitions *ConsumerOffsetsPartitions, consumerName string, consumerGroup ...string) error
```
SeekToBeginning seek to the first offset for each of the given partitions.




### <a name="Consumers.SeekToEnd">func</a> (\*Consumers) [SeekToEnd](./kafka/consumer.go?s=11985:12121#L516)
``` go
func (cs *Consumers) SeekToEnd(consumerOffsetsPartitions *ConsumerOffsetsPartitions, consumerName string, consumerGroup ...string) error
```
SeekToEnd seek to the last offset for each of the given partitions.




### <a name="Consumers.Subscribe">func</a> (\*Consumers) [Subscribe](./kafka/consumer.go?s=5820:5962#L250)
``` go
func (cs *Consumers) Subscribe(topicSubscription *TopicSubscription, useTopicPattern bool, consumerName string, consumerGroup ...string) error
```
Subscribe to the given list of topics or a topic pattern.




### <a name="Consumers.Subscriptions">func</a> (\*Consumers) [Subscriptions](./kafka/consumer.go?s=6830:6939#L292)
``` go
func (cs *Consumers) Subscriptions(consumerName string, consumerGroup ...string) (*TopicsSubscription, error)
```
Subscriptions get the current subscribed list of topics.




### <a name="Consumers.Unsubscribe">func</a> (\*Consumers) [Unsubscribe](./kafka/consumer.go?s=7687:7771#L332)
``` go
func (cs *Consumers) Unsubscribe(consumerName string, consumerGroup ...string) error
```
Unsubscribe from topics currently subscribed.




## <a name="ErrorMessage">type</a> [ErrorMessage](./kafka/kafka.go?s=847:966#L50)
``` go
type ErrorMessage struct {
    ErrorCode int    `json:"error_code,omitempty"`
    Message   string `json:"message,omitempty"`
}
```
ErrorMessage for API response










## <a name="Format">type</a> [Format](./kafka/kafka.go?s=693:706#L41)
``` go
type Format string
```
Format is one of json, binary or avro










## <a name="Kafka">type</a> [Kafka](./kafka/kafka.go?s=364:536#L21)
``` go
type Kafka struct {
    URL         string
    Timeout     time.Duration
    Accept      string
    ContentType string
    Format      Format
    Offset      Offset
    Version     Version
}
```
Kafka represents a Kafka REST API.







### <a name="New">func</a> [New](./kafka/kafka.go?s=5919:5974#L281)
``` go
func New(options ...func(*Kafka) error) (*Kafka, error)
```
New returns a Kafka instance with default setting.





### <a name="Kafka.Broker">func</a> (\*Kafka) [Broker](./kafka/kafka.go?s=6489:6530#L313)
``` go
func (k *Kafka) Broker() (*Broker, error)
```
Broker returns the brokers.




### <a name="Kafka.HTTPClient">func</a> (\*Kafka) [HTTPClient](./kafka/kafka.go?s=1851:1892#L93)
``` go
func (k *Kafka) HTTPClient() *http.Client
```
HTTPClient creates a new http.Client with timeout.




### <a name="Kafka.NewConsumers">func</a> (\*Kafka) [NewConsumers](./kafka/kafka.go?s=6271:6335#L300)
``` go
func (k *Kafka) NewConsumers(consumerGroup ...string) *Consumers
```
NewConsumers returns a Consumers instance.




### <a name="Kafka.NewTopics">func</a> (\*Kafka) [NewTopics](./kafka/kafka.go?s=6138:6173#L292)
``` go
func (k *Kafka) NewTopics() *Topics
```
NewTopics returns a Topics instance.




### <a name="Kafka.SetOption">func</a> (\*Kafka) [SetOption](./kafka/kafka.go?s=2055:2117#L101)
``` go
func (k *Kafka) SetOption(options ...func(*Kafka) error) error
```
SetOption takes one or more option function and applies them in order to Kafka.




## <a name="Message">type</a> [Message](./kafka/consumer.go?s=1804:2041#L76)
``` go
type Message struct {
    Topic     string          `json:"topic"`
    Key       json.RawMessage `json:"key"`
    Value     json.RawMessage `json:"value"`
    Partition int             `json:"partition"`
    Offset    int64           `json:"offset"`
}
```
Message is a single Kafka message










## <a name="Offset">type</a> [Offset](./kafka/kafka.go?s=749:762#L44)
``` go
type Offset string
```
Offset is either earliest or latest










## <a name="Partition">type</a> [Partition](./kafka/partition.go?s=266:406#L22)
``` go
type Partition struct {
    Partition int       `json:"partition"`
    Leader    int       `json:"leader"`
    Replicas  []Replica `json:"replicas"`
}
```
Partition data










## <a name="Partitions">type</a> [Partitions](./kafka/partition.go?s=429:502#L29)
``` go
type Partitions struct {
    Kafka *Kafka
    Topic *Topic
    List  *[]Partition
}
```
Partitions data










### <a name="Partitions.Partition">func</a> (\*Partitions) [Partition](./kafka/partition.go?s=1329:1418#L77)
``` go
func (ps *Partitions) Partition(partitionID int, topicName ...string) (*Partition, error)
```
Partition returns the Partition with provided partitionID.




### <a name="Partitions.Partitions">func</a> (\*Partitions) [Partitions](./kafka/partition.go?s=552:627#L37)
``` go
func (ps *Partitions) Partitions(topicName ...string) (*[]Partition, error)
```
Partitions lists partitions for the topic.




### <a name="Partitions.Produce">func</a> (\*Partitions) [Produce](./kafka/partition.go?s=2136:2247#L117)
``` go
func (ps *Partitions) Produce(id int, message *ProducerMessage, topicName ...string) (*ProducerResponse, error)
```
Produce post message to the Partition with provided id.




## <a name="ProducerMessage">type</a> [ProducerMessage](./kafka/topic.go?s=456:852#L30)
``` go
type ProducerMessage struct {
    KeySchema   string `json:"key_schema,omitempty"`
    KeySchemaID int    `json:"key_schema_id,omitempty"`
    //either value schema or value schema id must be provided for avro messages
    ValueSchema   string            `json:"value_schema,omitempty"`
    ValueSchemaID int               `json:"value_schema_id,omitempty"`
    Records       []*ProducerRecord `json:"records"`
}
```
ProducerMessage is the wrapper for the Topic / Partition data










## <a name="ProducerOffsets">type</a> [ProducerOffsets](./kafka/topic.go?s=1425:1598#L54)
``` go
type ProducerOffsets struct {
    Partition int    `json:"partition"`
    Offset    int64  `json:"offset"`
    ErrorCode int64  `json:"error_code"`
    Error     string `json:"error"`
}
```
ProducerOffsets are the resulting offsets for Topic / Partition










## <a name="ProducerRecord">type</a> [ProducerRecord](./kafka/topic.go?s=921:1098#L40)
``` go
type ProducerRecord struct {
    Key       json.RawMessage `json:"key,omitempty"`
    Value     json.RawMessage `json:"value"`
    Partition int             `json:"partition,omitempty"`
}
```
ProducerRecord is an individual message for Topic / Partition










## <a name="ProducerResponse">type</a> [ProducerResponse](./kafka/topic.go?s=1156:1354#L47)
``` go
type ProducerResponse struct {
    KeySchemaID   int                `json:"key_schema_id"`
    ValueSchemaID int                `json:"value_schema_id"`
    Offsets       []*ProducerOffsets `json:"offsets"`
}
```
ProducerResponse is the Topic / Partition response










## <a name="Replica">type</a> [Replica](./kafka/partition.go?s=134:244#L15)
``` go
type Replica struct {
    Broker int  `json:"broker"`
    Leader bool `json:"leader"`
    InSync bool `json:"in_sync"`
}
```
Replica data










## <a name="Topic">type</a> [Topic](./kafka/topic.go?s=121:276#L14)
``` go
type Topic struct {
    Name       string          `json:"name"`
    Configs    json.RawMessage `json:"configs"`
    Partitions []*Partition    `json:"partitions"`
}
```
Topic data










## <a name="TopicNames">type</a> [TopicNames](./kafka/topic.go?s=368:387#L27)
``` go
type TopicNames []string
```
TopicNames data










## <a name="TopicPatternSubscription">type</a> [TopicPatternSubscription](./kafka/consumer.go?s=1682:1763#L71)
``` go
type TopicPatternSubscription struct {
    TopicPattern string `json:"topic_pattern"`
}
```
TopicPatternSubscription with topic pattern










## <a name="TopicSubscription">type</a> [TopicSubscription](./kafka/consumer.go?s=1424:1529#L60)
``` go
type TopicSubscription struct {
    Topics       *TopicsSubscription
    TopicPattern *TopicPatternSubscription
}
```
TopicSubscription data, topic_pattern and topics fields are mutually exclusive










## <a name="Topics">type</a> [Topics](./kafka/topic.go?s=295:345#L21)
``` go
type Topics struct {
    Kafka *Kafka
    List  []*Topic
}
```
Topics data










### <a name="Topics.Names">func</a> (\*Topics) [Names](./kafka/topic.go?s=1924:1970#L81)
``` go
func (ts *Topics) Names() (*TopicNames, error)
```
Names lists all topic names.




### <a name="Topics.NewPartitions">func</a> (\*Topics) [NewPartitions](./kafka/topic.go?s=4459:4515#L205)
``` go
func (ts *Topics) NewPartitions(t ...*Topic) *Partitions
```
NewPartitions returns a Partitions instance.




### <a name="Topics.Produce">func</a> (\*Topics) [Produce](./kafka/topic.go?s=3210:3306#L151)
``` go
func (ts *Topics) Produce(topicName string, message *ProducerMessage) (*ProducerResponse, error)
```
Produce post message to the Topic with provided topicName.




### <a name="Topics.Topic">func</a> (\*Topics) [Topic](./kafka/topic.go?s=2555:2612#L116)
``` go
func (ts *Topics) Topic(topicName string) (*Topic, error)
```
Topic returns the Topic with provided topicName.




### <a name="Topics.Topics">func</a> (\*Topics) [Topics](./kafka/topic.go?s=1630:1674#L63)
``` go
func (ts *Topics) Topics() ([]*Topic, error)
```
Topics lists all topics.




## <a name="TopicsSubscription">type</a> [TopicsSubscription](./kafka/consumer.go?s=1567:1631#L66)
``` go
type TopicsSubscription struct {
    Topics []string `json:"topics"`
}
```
TopicsSubscription with topics










## <a name="Version">type</a> [Version](./kafka/kafka.go?s=796:810#L47)
``` go
type Version string
```
Version is the API version












