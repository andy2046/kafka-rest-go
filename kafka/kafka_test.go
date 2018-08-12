package kafka_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	K "github.com/andy2046/kafka-rest-go/kafka"
)

func TestKafkaBroker(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		b := K.Broker{
			Brokers: []int{1, 2, 3},
		}
		m, _ := json.Marshal(b)
		io.WriteString(w, string(m))
	}))
	defer ts.Close()

	// create new Kafka instance with provided URL
	k, _ := K.New(K.SetURL(ts.URL))
	str, _ := K.Stringer(*k)
	fmt.Println(str)

	// list all brokers by id for Kafka instance
	brokers, err := k.Broker()
	if err != nil {
		t.Fatalf("Expected no error got %v", err)
	}
	fmt.Println(*brokers)
}
