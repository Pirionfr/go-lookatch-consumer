package core

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

var api *Api

// equals fails the test if exp is not equal to act.
func equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}

// ok fails the test if an err is not nil.
func ok(tb testing.TB, err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: unexpected error: %s\033[39m\n\n", filepath.Base(file), line, err.Error())
		tb.FailNow()
	}
}

func TestBasicAuth(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Test request parameters
		equals(t, "/auth/token", req.URL.String())
		equals(t, "Basic dXNlcjp0ZXN0", req.Header.Get("Authorization"))
		// Send response to be tested
		rw.WriteHeader(200)
		rw.Write([]byte(`"token"`))
	}))

	defer server.Close()

	api = NewClient(server.URL, "")

	err := api.BasicAuth("user", "test")

	ok(t, err)

	equals(t, "token", api.Token)

}

func TestGetSinks(t *testing.T) {
	agentId := "48a7e6f5-fe4e-4579-a12b-c7d39729d546"

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Test request parameters
		equals(t, "/collectors/"+agentId+"/sinks", req.URL.String())
		equals(t, "Bearer dXNlcjp0ZXN0", req.Header.Get("Authorization"))
		// Send response to be tested
		rw.WriteHeader(200)
		rw.Write([]byte(`["querysink"]`))
	}))

	defer server.Close()

	api = NewClient(server.URL, "dXNlcjp0ZXN0")

	res, err := api.GetSinks(agentId)
	ok(t, err)

	equals(t, "querysink", res[0])

}

func TestGetSink(t *testing.T) {
	agentId := "48a7e6f5-fe4e-4579-a12b-c7d39729d546"
	sinkName := "querysink"
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Test request parameters
		equals(t, "/collectors/"+agentId+"/sinks/"+sinkName, req.URL.String())
		equals(t, "Bearer dXNlcjp0ZXN0", req.Header.Get("Authorization"))
		// Send response to be tested
		rw.WriteHeader(200)
		rw.Write([]byte(`{"brokers": [ "broker1:9093", "broker2:9093"], "provider": "ovh", "max_message_bytes": 10485760, "topic": "lookatch.test", "tls": true, "type": "kafka", "nb_producer": 3, "enabled": true, "consumer": { "user": "lookatch.test", "password": "test"}}`))
	}))

	defer server.Close()

	api = NewClient(server.URL, "dXNlcjp0ZXN0")

	res, err := api.GetSink(agentId, sinkName)
	ok(t, err)

	equals(t, "broker2:9093", res["brokers"].([]interface{})[1])
	equals(t, "ovh", res["provider"])
	equals(t, "lookatch.test", res["topic"])
	equals(t, 3, int(res["nb_producer"].(float64)))
	equals(t, true, res["tls"])
	equals(t, "kafka", res["type"])
	equals(t, "lookatch.test", res["consumer"].(map[string]interface{})["user"])
	equals(t, "test", res["consumer"].(map[string]interface{})["password"])

}
