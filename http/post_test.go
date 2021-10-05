package main

import (
	"bytes"
	"time"
	//"context"
	"encoding/json"
	//"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"

	//"os"
	//"path/filepath"
	"testing"
)

type User struct {
    First string
    Last string
}

func handlePostUser(t *testing.T) func(http.ResponseWriter, *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        defer func(r io.ReadCloser) {
            _, _ = io.Copy(ioutil.Discard, r)
            _ = r.Close()
        }(r.Body)
        
        if r.Method != http.MethodPost {
            http.Error(w, "", http.StatusMethodNotAllowed)
            return
        }

        var u User
        err := json.NewDecoder(r.Body).Decode(&u)
        if err != nil {
            t.Error(err)
            http.Error(w, "decode dailed", http.StatusBadRequest)
            return
        }

        w.WriteHeader(http.StatusAccepted)
    }
}

func TestPostUser(t *testing.T) {
    ts := httptest.NewServer(http.HandlerFunc(handlePostUser(t)))
    defer ts.Close()

    resp, err := http.Get(ts.URL)
    if err != nil {
        t.Fatal(err)
    }
    if resp.StatusCode != http.StatusMethodNotAllowed {
        t.Fatalf("expected status %d; actual status %d", 
        http.StatusMethodNotAllowed, resp.StatusCode)
    }

    buf := new(bytes.Buffer)
    u := User{First: "Neil", Last: "Ulises"}
    err = json.NewEncoder(buf).Encode(&u)
    if err != nil {
        t.Fatal(err)
    }

    resp, err = http.Post(ts.URL, "application/json", buf)
    if err != nil {
        t.Fatal(err)
    }
    if resp.StatusCode != http.StatusAccepted {
        t.Fatalf("expected status %d; actual status %d",
        http.StatusAccepted, resp.StatusCode)
    } 
    _ = resp.Body.Close()
}

func TestMultiPartPost(t *testing.T) {
    reqBody := new(bytes.Buffer)
    w := multipart.NewWriter(reqBody)

    for k, v := range map[string]string {
        "date": time.Now().Format(time.RFC3339),
        "description": "Form values with attached files",
    } {
        err := w.WriteField(k, v)
        if err != nil {
            t.Fatal(err)
        }
    }
}
