package javascript

import (
	"context"
	"net/http"
	"testing"

	"github.com/Unbabel/replicant/internal/tmpl"
	"github.com/Unbabel/replicant/transaction"
)

func TestDriverTransaction(t *testing.T) {

	mux := http.NewServeMux()
	mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {

		if r.URL.Query().Get("q") != "blade runner" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"reason": "query parameter q not found"}`))
			return
		}

		if r.Header.Get("X-Auth") != "Joi" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"reason": "request header X-Auth not found"}`))
			return
		}

		w.WriteHeader(200)
		w.Write([]byte(`{"reason": "test successful"}`))
		return
	})

	server := &http.Server{}
	server.Addr = "localhost:8080"
	server.Handler = mux
	go server.ListenAndServe()
	defer server.Close()

	d, err := New()
	if err != nil {
		t.Fatalf("error creating driver: %s", err)
	}

	cfg, err := tmpl.Parse(config)
	if err != nil {
		t.Fatalf("error parsing template: %s", err)
	}

	txn, err := d.New(cfg)
	if err != nil {
		t.Fatalf("error creating transaction: %s\n, %#v", err, cfg)
	}

	ctx := context.WithValue(context.Background(), "transaction_uuid", "test-test-test")

	result := txn.Run(ctx)
	if result.Error != nil {
		t.Fatalf("error running transaction: %s", result.Error)
	}

	if result.Failed {
		t.Fatalf("transaction failed:\n%#v", result)
	}

	t.Logf("%#v", result)
}

// test transaction
var config transaction.Config = transaction.Config{
	Name:       "test-transaction",
	Driver:     "javascript",
	Schedule:   "@every 60s",
	Timeout:    "5s",
	RetryCount: 1,
	Inputs: map[string]interface{}{
		"url":   "http://localhost:8080/test",
		"text":  "blade runner",
		"xauth": "Joi",
	},
	Metadata: map[string]string{
		"transaction": "api-test",
		"application": "test",
		"environment": "test",
		"component":   "api",
	},
	Script: `function Run(ctx) {
	//replicant.Log("test started")
	req = replicant.http.NewRequest()
	req.Method = "GET"
	req.URL = "{{ index . "url" }}"
	req.Params.q = "{{ index . "text" }}"
	req.Header["X-Auth"] = "{{ index . "xauth" }}"
	//replicant.Log("going to perform request")
	resp = replicant.http.Do(req)
	data = JSON.parse(resp.Body)
	//replicant.Log(data)
	rr = replicant.NewResult()
	rr.Message = resp.Status
	switch(resp.StatusCode > 200) {
		case true:
		rr.Error = data.reason
		rr.Failed = true
		break
	case false:
		rr.Data = data.reason
		rr.Failed = false
		break
	}
	return rr.JSON()
}`,
}
