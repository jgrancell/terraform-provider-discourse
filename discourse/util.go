package discourse

import (
  "encoding/json"
  "io"
  "io/ioutil"
  "net/http"
  "net/url"
  "os"
  "strings"
  "strconv"
  "time"

  "github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)


func (req *ApiRequest) Call() (string, diag.Diagnostic, bool) {
  log("API "+req.Method+" call beginning for "+req.Endpoint)

  // Generating the full request url
  reqUrl, _ := url.Parse(req.Config.Host + req.Endpoint)
  log("Request URL: "+reqUrl.String())

  var reqBody io.ReadCloser
  if req.Method != "GET" {
    reqBody = ioutil.NopCloser(strings.NewReader(req.JsonBody))
  }

  // Building our request
  request := &http.Request {
    Method: req.Method,
    URL: reqUrl,
    Header: map[string][]string {
      "Api-Key": {req.Config.Token},
      "Api-Username": {req.Config.Username},
      "Content-Type": {"application/json"},
    },
    Body: reqBody,
  }

  log("API "+req.Method+" raw http request built")

  // Creating our http.Client with a sane timeout
  var httpClient *http.Client
  httpClient = &http.Client{Timeout: 10 * time.Second}

  log("API "+req.Method+" http client built")

  resp, err := httpClient.Do(request)
  if err != nil {
    diag := diag.Diagnostic{
      Severity: diag.Error,
      Summary: "HTTP client request returned an error.",
      Detail: err.Error(),
    }
    return "", diag, false
  }

  log("Parsing API response for HTTP "+req.Method+" request to "+req.Endpoint)

  bodyBytes, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    diag := diag.Diagnostic{
      Severity: diag.Error,
      Summary: "Could not read HTTP response body",
      Detail: err.Error(),
    }
    return "", diag, false
  }

  log("Validating response HTTP Status code for HTTP "+req.Method+" request to "+req.Endpoint)
  log("Dumping body:" + string(bodyBytes))

  switch resp.StatusCode {
    case 200:
      return string(bodyBytes), diag.Diagnostic{}, true
  }

  log("Received an unexpected status code for HTTP "+req.Method+" request to "+req.Endpoint)
  log("Status code received: "+strconv.Itoa(resp.StatusCode))

  var de DiscourseError
  je := json.Unmarshal(bodyBytes, &de)
  if je != nil {
    diag := diag.Diagnostic{
      Severity: diag.Error,
      Summary: "Could not unmarshal Discourse error message.",
      Detail: je.Error(),
    }
    return "", diag, false
  } else {
    diag := diag.Diagnostic{
      Severity: diag.Error,
      Summary: "Received a Discourse API error response.",
      Detail: "Request endpoint: "+req.Endpoint+"; Error type: "+de.ErrorType+"; Error message: "+de.Message[0],
    }
    return "", diag, false
  }
}

func log(message string) {
  // If we're not in DEBUG or TRACE we're not going to do anything
  var log_level string
  log_level = os.Getenv("TF_LOG")
  if log_level == "debug" || log_level == "trace" {

    // Setting our default time string
    cur := time.Now()
    timestamp := cur.Format("2006-01-02 15:04:05")
    f, _ := os.OpenFile("./terraform-provider-discourse.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    defer f.Close()
    _, err := f.WriteString(timestamp+" - [DEBUG] - "+message+"\n")
    if err != nil {
      panic(err)
    }
    f.Sync()
  }
}
