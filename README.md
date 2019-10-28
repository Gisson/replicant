# Replicant

Replicant is a synthetic transaction execution framework named after the bioengineered androids from Blade Runner. (all synthetics came from Blade Runner :)

It defines a common interface for transactions and results, provides a transaction manager, execution scheduler, api and facilities for emitting result data to external systems.

## Requirements

* Go 1.13
* External URL for API tests that require webhook based callbaks
* Chrome with remote debugging (CDP) either in headless mode or in foreground (useful for testing)

## Examples

Running the server
```bash
CALLBACK_URL="https://external.name.net" LISTEN_ADDRESS=localhost:8080 EMITTER=stdout,prometheus go run cmd/replicant/main.go
```

### Web application testing (local development)

* Start the Chrome browser with Chrome DevTools Protocol enabled:
`/Applications/Google\ Chrome.app/Contents/MacOS/Google\ Chrome --remote-debugging-port=9222 &`

```yaml
POST http://127.0.0.1:8080/v1/run
content-type: application/yaml

name: duckduckgo-search
type: web
schedule: '@every 1m'
timeout: 200s
retry_count: 2
inputs:
  url: "https://duckduckgo.com"
  cdp_address: "http://127.0.0.1:9222"
  user_agent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3809.87 Safari/537.36"
  timeout: 5000000
  text: "blade runner"
metadata:
  application: web-search
  environment: production
  component: web
script: |
  LET doc = DOCUMENT('{{ index . "url" }}', { driver: "cdp", userAgent: "{{ index . "user_agent" }}"})
  INPUT(doc, '#search_form_input_homepage', "{{ index . "text" }}")
  CLICK(doc, '#search_button_homepage')
  WAIT_NAVIGATION(doc)
  LET result = ELEMENT(doc, '#r1-0 > div > div.result__snippet.js-result-snippet').innerText
  RETURN {
    failed: result == "",
    message: "search result",
    data: result,
  }

```

### API testing (local development)

```yaml
POST http://127.0.0.1:8080/v1/run
content-type: application/yaml

name: duckduckgo-search
type: go
schedule: '@every 1m'
timeout: 200s
retry_count: 2
inputs:
  url: "https://api.duckduckgo.com/"
  text: "blade runner"
metadata:
  application: api-search
  environment: production
  component: api
script: |
  package transaction
  import "bytes"
  import "context"
  import "fmt"
  import "net/http"
  import "io/ioutil"
  import "net/http"
  import "regexp"
  func Run(ctx context.Context) (m string, d string, err error) {
    req, err := http.NewRequest(http.MethodGet, "{{ index . "url" }}", nil)
      if err != nil {
        return "request build failed", "", err
    }
    req.Header.Add("Accept-Charset","utf-8")
    q := req.URL.Query()
    q.Add("q", "{{ index . "text" }}")
    q.Add("format", "json")
    q.Add("pretty", "1")
    q.Add("no_redirect", "1")
    req.URL.RawQuery = q.Encode()
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
      return "failed to send request", "", err
    }
    buf, err := ioutil.ReadAll(resp.Body)
    if err != nil {
      return "failed to read response", "", err
    }
    rx, err := regexp.Compile(`"Text"\s*:\s*"(.*?)"`)
    if err != nil {
      return "failed to compile regexp", "", err
    }
    s := rx.FindSubmatch(buf)
    if len(s) < 2 {
      return "failed to find data", "", fmt.Errorf("failed to find data")
    }
    return "search result", fmt.Sprintf("%s", s[1]), nil
  }
```

## TODO

* Tests
* Developer Documentation
* Architecture and API documentation
* Javascript transaction support

## Related Projects

* [Yaegi is Another Elegant Go Interpreter](https://github.com/containous/yaegi)
* [Ferret Declarative web scraping](https://github.com/MontFerret/ferret)

## Contact
Bruno Moura [brunotm@gmail.com](mailto:brunotm@gmail.com)

## License
Replicant source code is available under the Apache Version 2.0 [License](/LICENSE)