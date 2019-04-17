/*
Copyright 2019 The Knative Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package rabbitmq

import (
	"context"
	"encoding/json"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/client"
	"github.com/knative/eventing-contrib/pkg/kncloudevents"
	"github.com/sbcd90/wabbit/amqp"
	"github.com/sbcd90/wabbit/amqptest"
	"github.com/sbcd90/wabbit/amqptest/server"
	origamqp "github.com/streadway/amqp"
	"io/ioutil"
	"knative.dev/pkg/logging"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPostMessage_ServeHttp(t *testing.T) {
	testCases := map[string]struct{
		sink              func(http.ResponseWriter, *http.Request)
		reqBody           string
		attributes        map[string]string
		expectedEventType string
		error             bool
	}{
		"accepted": {
			sink:    sinkAccepted,
			reqBody: `{"key":"value"}`,
		},
		"rejected": {
			sink:    sinkRejected,
			reqBody: `{"key":"value"}`,
			error:   true,
		},
	}

	for n, tc := range testCases {
		t.Run(n, func(t *testing.T) {
			h := &fakeHandler{
				handler: tc.sink,
			}
			sinkServer := httptest.NewServer(h)
			defer sinkServer.Close()

			a := &Adapter{
				Topic:          "topic",
				Brokers:        "amqp://guest:guest@localhost:5672/",
				ExchangeConfig: ExchangeConfig{
					TypeOf:      "topic",
					Durable:     true,
					AutoDeleted: false,
					Internal:    false,
					NoWait:      false,
				},
				QueueConfig:   QueueConfig{
					Name:             "",
					Durable:          false,
					DeleteWhenUnused: false,
					Exclusive:        true,
					NoWait:           false,
				},
				SinkURI:       sinkServer.URL,
				client: func() client.Client {
					c, _ := kncloudevents.NewDefaultClient(sinkServer.URL)
					return c
				}(),
			}

			data, err := json.Marshal(map[string]string{"key": "value"})
			if err != nil {
				t.Errorf("unexpected error, %v", err)
			}

			m := &amqp.Delivery{}
			m.Delivery = &origamqp.Delivery{
				MessageId: "id",
				Body:      []byte(data),
			}
			err = a.postMessage(context.TODO(), m)

			if tc.error && err == nil {
				t.Errorf("expected error, but got %v", err)
			}

			et := h.header.Get("Ce-Type")

			expectedEventType := eventType
			if tc.expectedEventType != "" {
				expectedEventType = tc.expectedEventType
			}

			if et != expectedEventType {
				t.Errorf("Expected eventtype '%q', but got '%q'", tc.expectedEventType, et)
			}
			if tc.reqBody != string(h.body) {
				t.Errorf("Expected request body '%q', but got '%q'", tc.reqBody, h.body)
			}
		})
	}
}

func TestAdapter_CreateConn(t *testing.T) {
	fakeServer := server.NewServer("amqp://localhost:5672/%2f")
	err := fakeServer.Start()
	if err != nil {
		t.Errorf("%s: %s", "Failed to connect to RabbitMQ", err)
	}

	a := &Adapter{
		Topic:          "",
		Brokers:        "amqp://localhost:5672/%2f",
		ExchangeConfig: ExchangeConfig{
			TypeOf:      "direct",
			Durable:     true,
			AutoDeleted: false,
			Internal:    false,
			NoWait:      false,
		},
		QueueConfig:   QueueConfig{
			Name:             "",
			Durable:          false,
			DeleteWhenUnused: false,
			Exclusive:        true,
			NoWait:           false,
		},
	}

	conn, _ := a.CreateConn("", "", logging.FromContext(context.TODO()))
	if conn != nil {
		t.Errorf("Failed to connect to RabbitMQ")
	}

	conn, _ = a.CreateConn("guest", "guest", logging.FromContext(context.TODO()))
	if conn != nil {
		t.Errorf("Failed to connect to RabbitMQ")
	}
	fakeServer.Stop()
}

func TestAdapter_CreateChannel(t *testing.T) {
	fakeServer := server.NewServer("amqp://localhost:5672/%2f")
	err := fakeServer.Start()
	if err != nil {
		t.Errorf("%s: %s", "Failed to connect to RabbitMQ", err)
	}

	conn, err := amqptest.Dial("amqp://localhost:5672/%2f")
	if err != nil {
		t.Errorf("%s: %s", "Failed to connect to RabbitMQ", err)
	}

	a := &Adapter{
		Topic:          "",
		Brokers:        "amqp://localhost:5672/%2f",
		ExchangeConfig: ExchangeConfig{
			TypeOf:      "direct",
			Durable:     true,
			AutoDeleted: false,
			Internal:    false,
			NoWait:      false,
		},
		QueueConfig:   QueueConfig{
			Name:             "",
			Durable:          false,
			DeleteWhenUnused: false,
			Exclusive:        true,
			NoWait:           false,
		},
	}

	for i := 1; i <= 10000; i++ {
		channel, _ := a.CreateChannel(nil, conn, logging.FromContext(context.TODO()))
		if channel == nil {
			t.Logf("Failed to open a channel")
			break
		}
	}
	fakeServer.Stop()
}

func TestAdapter_StartAmqpClient(t *testing.T) {

	/**
	Test for exchange type "direct"
	 */
	fakeServer := server.NewServer("amqp://localhost:5672/%2f")
	err := fakeServer.Start()
	if err != nil {
		t.Errorf("%s: %s", "Failed to connect to RabbitMQ", err)
	}

	conn, err := amqptest.Dial("amqp://localhost:5672/%2f")
	if err != nil {
		t.Errorf("%s: %s", "Failed to connect to RabbitMQ", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		t.Errorf("Failed to open a channel")
	}

	a := &Adapter{
		Topic:          "",
		Brokers:        "amqp://localhost:5672/%2f",
		ExchangeConfig: ExchangeConfig{
			TypeOf:      "direct",
			Durable:     true,
			AutoDeleted: false,
			Internal:    false,
			NoWait:      false,
		},
		QueueConfig:   QueueConfig{
			Name:             "",
			Durable:          false,
			DeleteWhenUnused: false,
			Exclusive:        true,
			NoWait:           false,
		},
	}

	_, err = a.StartAmqpClient(context.TODO(), &channel)
	if err != nil {
		t.Errorf("Failed to start RabbitMQ")
	}

	/**
	Test for exchange type "fanout"
	 */
	a = &Adapter{
		Topic:          "",
		Brokers:        "amqp://localhost:5672/%2f",
		ExchangeConfig: ExchangeConfig{
			TypeOf:      "fanout",
			Durable:     true,
			AutoDeleted: false,
			Internal:    false,
			NoWait:      false,
		},
		QueueConfig:   QueueConfig{
			Name:             "",
			Durable:          false,
			DeleteWhenUnused: false,
			Exclusive:        true,
			NoWait:           false,
		},
	}
	_, err = a.StartAmqpClient(context.TODO(), &channel)
	if err != nil {
		t.Errorf("Failed to start RabbitMQ")
	}

	/**
	Test for exchange type "topic"
	 */
	a = &Adapter{
		Topic:          "",
		Brokers:        "amqp://localhost:5672/%2f",
		ExchangeConfig: ExchangeConfig{
			TypeOf:      "topic",
			Durable:     true,
			AutoDeleted: false,
			Internal:    false,
			NoWait:      false,
		},
		QueueConfig:   QueueConfig{
			Name:             "",
			Durable:          false,
			DeleteWhenUnused: false,
			Exclusive:        true,
			NoWait:           false,
		},
	}
	_, err = a.StartAmqpClient(context.TODO(), &channel)
	if err != nil {
		t.Errorf("Failed to start RabbitMQ")
	}

	fakeServer.Stop()
}

func TestAdapter_StartAmqpClient_InvalidExchangeType(t *testing.T) {
	/**
	Test for invalid exchange type
	 */
	fakeServer := server.NewServer("amqp://localhost:5674/%2f")
	err := fakeServer.Start()
	if err != nil {
		t.Errorf("%s: %s", "Failed to connect to RabbitMQ", err)
	}

	conn, err := amqptest.Dial("amqp://localhost:5674/%2f")
	if err != nil {
		t.Errorf("%s: %s", "Failed to connect to RabbitMQ", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		t.Errorf("Failed to open a channel")
	}

	a := &Adapter{
		Topic:          "",
		Brokers:        "amqp://localhost:5674/%2f",
		ExchangeConfig: ExchangeConfig{
			TypeOf:      "",
			Durable:     true,
			AutoDeleted: false,
			Internal:    false,
			NoWait:      false,
		},
		QueueConfig:   QueueConfig{
			Name:             "",
			Durable:          false,
			DeleteWhenUnused: false,
			Exclusive:        true,
			NoWait:           false,
		},
	}
	_, err = a.StartAmqpClient(context.TODO(), &channel)
	if err != nil {
		t.Logf("Failed to start RabbitMQ")
	}
	fakeServer.Stop()
}

func TestAdapter_ConsumeMessages(t *testing.T) {
	fakeServer := server.NewServer("amqp://localhost:5672/%2f")
	err := fakeServer.Start()
	if err != nil {
		t.Errorf("%s: %s", "Failed to connect to RabbitMQ", err)
	}

	conn, err := amqptest.Dial("amqp://localhost:5672/%2f")
	if err != nil {
		t.Errorf("%s: %s", "Failed to connect to RabbitMQ", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		t.Errorf("Failed to open a channel")
	}

	a := &Adapter{
		Topic:          "",
		Brokers:        "amqp://localhost:5672/%2f",
		ExchangeConfig: ExchangeConfig{
			TypeOf:      "direct",
			Durable:     true,
			AutoDeleted: false,
			Internal:    false,
			NoWait:      false,
		},
		QueueConfig:   QueueConfig{
			Name:             "",
			Durable:          false,
			DeleteWhenUnused: false,
			Exclusive:        true,
			NoWait:           false,
		},
	}
	queue, err := a.StartAmqpClient(context.TODO(), &channel)
	if err != nil {
		t.Errorf("Failed to start RabbitMQ")
	}
	_, err = a.ConsumeMessages(&channel, queue, logging.FromContext(context.TODO()))
	if err != nil {
		t.Errorf("Failed to consume from RabbitMQ")
	}
}

func TestAdapter_JsonEncode(t *testing.T) {
	a := &Adapter{
		Topic:          "",
		Brokers:        "amqp://localhost:5672/%2f",
		ExchangeConfig: ExchangeConfig{
			TypeOf:      "direct",
			Durable:     true,
			AutoDeleted: false,
			Internal:    false,
			NoWait:      false,
		},
		QueueConfig:   QueueConfig{
			Name:             "",
			Durable:          false,
			DeleteWhenUnused: false,
			Exclusive:        true,
			NoWait:           false,
		},
	}
	data := a.JsonEncode(context.TODO(), []byte("test json"))
	if data == nil {
		t.Errorf("Json decoded incorrectly")
	}
}

type fakeHandler struct {
	body   []byte
	header http.Header

	handler func(http.ResponseWriter, *http.Request)
}

func (h *fakeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request)  {
	h.header = r.Header
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "can not read body", http.StatusBadRequest)
		return
	}
	h.body = body

	defer r.Body.Close()
	h.handler(w, r)
}

func sinkAccepted(writer http.ResponseWriter, req *http.Request) {
	writer.WriteHeader(http.StatusOK)
}

func sinkRejected(writer http.ResponseWriter, _ *http.Request) {
	writer.WriteHeader(http.StatusRequestTimeout)
}
