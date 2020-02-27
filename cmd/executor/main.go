package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/Unbabel/replicant/internal/executor"
	"github.com/Unbabel/replicant/log"
	"github.com/Unbabel/replicant/transaction"
	"github.com/julienschmidt/httprouter"
)

var (
	address      = flag.String("address", "0.0.0.0:8080", "address to listen on")
	serverURL    = flag.String("server-url", "http://0.0.0.0:8080", "replicant server base url")
	advertiseURL = flag.String("advertise-url", "https://some.external.url:8080", "advertise base url for async responses")
	interval     = flag.Duration("interval", time.Second*300, "interval to recycle chrome process")
	args         = flag.String("args", defaultArgs, "chrome command line")
	level        = flag.String("level", "INFO", "log level")

	defaultArgs = "/headless-shell/headless-shell --headless --no-zygote --no-sandbox --disable-gpu --disable-software-rasterizer --disable-dev-shm-usage --remote-debugging-address=127.0.0.1 --remote-debugging-port=9222 --incognito --disable-shared-workers --disable-remote-fonts --disable-background-networking --disable-crash-reporter --disable-default-apps --disable-domain-reliability --disable-extensions --disable-shared-workers --disable-setuid-sandbox"
)

func main() {
	flag.Parse()
	log.Init(*level)

	config := executor.Config{}
	config.ServerURL = *serverURL
	config.AdvertiseURL = *advertiseURL

	arguments := strings.Split(*args, " ")
	config.Web.BinaryPath = arguments[:1][0]
	config.Web.BinaryArgs = arguments[1:]
	config.Web.ServerURL = "http://127.0.0.1:9222"
	config.Web.RecycleInterval = *interval

	server := &http.Server{}
	server.Addr = *address
	router := httprouter.New()
	server.Handler = router

	e, err := executor.New(config)
	if err != nil {
		log.Error("error creating replicant-executor").Error("error", err).Log()
		os.Exit(1)
	}

	router.Handle(http.MethodPost, "/api/v1/run/:uuid", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		defer r.Body.Close()

		uuid := p.ByName("uuid")
		var err error
		var buf []byte
		var config transaction.Config

		if buf, err = ioutil.ReadAll(r.Body); err != nil {
			httpError(w, uuid, config, fmt.Errorf("error reading request body: %w", err), http.StatusBadRequest)
			return
		}

		if err = json.Unmarshal(buf, &config); err != nil {
			httpError(w, uuid, config, fmt.Errorf("error deserializing json request body: %w", err), http.StatusBadRequest)
			return
		}

		if err = json.Unmarshal(buf, &config); err != nil {
			httpError(w, uuid, config, err, http.StatusBadRequest)
			return
		}

		result, err := e.Run(uuid, config)
		if err != nil {
			httpError(w, uuid, config, err, http.StatusBadRequest)
			return
		}

		buf, err = json.Marshal(&result)
		if err != nil {
			httpError(w, uuid, config, fmt.Errorf("error serializing results: %w", err), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(buf)
	})

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)

	// listen for stop signals
	go func() {
		<-signalCh
		if err := server.Shutdown(context.Background()); err != nil {
			log.Error("error stopping replicant-executor").Error("error", err).Log()
			os.Exit(1)
		}
	}()

	log.Info("starting replicant-executor").Log()
	if err := server.ListenAndServe(); err != nil {
		log.Error("error running replicant-cdp").Error("error", err).Log()
		os.Exit(1)
	}

	log.Info("replicant-cdp stopped").Log()

}

// httpError wraps http status codes and error messages as json responses
func httpError(w http.ResponseWriter, uuid string, config transaction.Config, err error, code int) {
	var result transaction.Result

	result.Name = config.Name
	result.Driver = config.Driver
	result.Metadata = config.Metadata
	result.Time = time.Now()
	result.DurationSeconds = 0
	result.Failed = true
	result.Error = err

	res, _ := json.Marshal(&result)

	w.WriteHeader(code)
	w.Write(res)

	log.Error("handling web transaction request").Error("error", err).Log()
}
