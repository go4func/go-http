package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

type Req struct {
	URL    string `json:"-"`
	Header string `json:"header"`
	Body   string `json:"body"`
}

type Resp struct {
	Header string `json:"header"`
	Body   string `json:"body"`
}

type helloHandler struct {
}

const (
	contentType     = "Content-Type"
	contentTypeJson = "application/json"
)

func (h *helloHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req Req
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		panic(err)
	}

	resp := Resp{Header: "this is respose header", Body: "this is response body"}
	response, err := json.Marshal(resp)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set(contentType, contentTypeJson)
	w.Write(response)
}

func contextHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	u, ok := ctx.Value("user").(string)

	if !ok {
		panic("invalid context type")
	}

	resp := `{hello:"` + u + `"}`

	w.WriteHeader(http.StatusOK)
	w.Header().Set(contentType, contentTypeJson)
	w.Write([]byte(resp))
}

func main() {

}

func Send(req Req) {
	// req, err := http.NewRequest("GET", "/api/projects",
	//     // Note: url.Values is a map[string][]string
	//     url.Values{"page": {"1"}, "per_page": {"100"}})

	// req.Header.Set("Authorization", "Bearer abc123")
	// b, err := json.Marshal(req)
	// if err != nil {
	// 	panic(err)
	// }
	// client := &http.Client{Timeout: 3 * time.Second}
	// request, err := http.NewRequest(http.MethodPost, req.URL, bytes.NewBuffer(b))
	// if err != nil {
	// 	panic(err)
	// }

	// var b bytes.Buffer
	// client := &http.Client{Timeout: 3 * time.Second}
	// err := json.NewEncoder(&b).Encode(&req)
	// fmt.Println(b)
	// request, err := http.NewRequest(http.MethodPost, req.URL, &b)
	// if err != nil {
	// 	panic(err)
	// }

	var err error
	pr, pw := io.Pipe()
	client := &http.Client{Timeout: 3 * time.Second}
	go func() {
		err = json.NewEncoder(pw).Encode(&req)
		if err != nil {
			panic(err)
		}
		defer pw.Close()
	}()

	request, err := http.NewRequest(http.MethodPost, req.URL, pr)
	if err != nil {
		panic(err)
	}

	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	r, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	var resp Resp
	err = json.Unmarshal(r, &resp)
}
