package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"
	"syscall"
	"testing"
	"time"
)

func TestDataStore(t *testing.T) {
	s := launchServer(t)

	t.Cleanup(func() {
		_ = s.Shutdown()
	})

	testPut(t)
	testGet(t)
	testDelete(t)
}

func testPut(t *testing.T) {
	payload1 := strings.NewReader("something")
	res1 := putBlob(t, payload1)

	payload2 := strings.NewReader("other")
	res2 := putBlob(t, payload2)

	if res1.OID == res2.OID {
		t.Errorf("expected to have unique oid")
	}

	if res1.Size != payload1.Size() {
		t.Errorf("expected a size of %d, got %d", payload1.Size(), res1.Size)
	}

	if res2.Size != payload2.Size() {
		t.Errorf("expected a size of %d, got %d", payload2.Size(), res2.Size)
	}
}

func testGet(t *testing.T) {
	content1 := "something"
	payload1 := strings.NewReader(content1)
	res1 := putBlob(t, payload1)

	content2 := "other"
	payload2 := strings.NewReader(content2)
	res2 := putBlob(t, payload2)

	body, status := getBlob(t, res1.OID)
	if status != http.StatusOK {
		t.Errorf("expected HTTP status of %d, got %d", http.StatusOK, status)
	}
	if body != content1 {
		t.Errorf("expected content of %s, got %s", content1, body)
	}

	body, status = getBlob(t, res2.OID)
	if status != http.StatusOK {
		t.Errorf("expected HTTP status of %d, got %d", http.StatusOK, status)
	}
	if body != content2 {
		t.Errorf("expected content of %s, got %s", content2, body)
	}
}

func testDelete(t *testing.T) {
	content := "something"
	payload := strings.NewReader(content)
	res := putBlob(t, payload)

	status := deleteBlob(t, res.OID)
	if status != http.StatusOK {
		t.Errorf("expected HTTP status of %d, got %d", http.StatusOK, status)
	}

	_, status = getBlob(t, res.OID)
	if status != http.StatusNotFound {
		t.Errorf("expected HTTP status of %d, got %d", http.StatusNotFound, status)
	}

	status = deleteBlob(t, res.OID)
	if status != http.StatusNotFound {
		t.Errorf("expected HTTP status of %d, got %d", http.StatusNotFound, status)
	}
}

type _response struct {
	OID  string `json:"oid"`
	Size int64  `json:"size"`
}

type testServer struct {
	cmd *exec.Cmd
}

func (s *testServer) Shutdown() error {
	if err := s.cmd.Process.Signal(syscall.SIGKILL); err != nil {
		return err
	}
	if err := s.cmd.Wait(); err != nil {
		return err
	}
	return nil
}

// closes an io.Closer and ignores the returned error
func closeIgnore(closer io.Closer) {
	_ = closer.Close()
}
func launchServer(t *testing.T) *testServer {
	cmd := exec.Command("go", "run", ".")
	err := cmd.Start()
	if err != nil {
		t.Fatal("unable to start server", err)
	}

	time.Sleep(time.Millisecond * 400) // Wait until we boot up

	return &testServer{cmd: cmd}
}

func getBlob(t *testing.T, oid string) (string, int) {
	objURL := fmt.Sprintf("http://localhost:8282/data/codingtest/%s", oid)
	req, _ := http.NewRequest("GET", objURL, nil)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("error making GET request: %s", err)
	}
	defer closeIgnore(res.Body)

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("error reading GET response: %s", err)
	}

	return string(body), res.StatusCode
}

func putBlob(t *testing.T, payload io.Reader) *_response {
	req, _ := http.NewRequest("PUT", "http://localhost:8282/data/codingtest", payload)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("error making PUT request: %s", err)
	}
	defer closeIgnore(res.Body)

	if res.StatusCode != http.StatusCreated {
		t.Errorf("expected response code 201, got %d", res.StatusCode)
	}

	contentType := res.Header.Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("expected Content-Type 'application/json', got %s", contentType)
	}

	var data _response
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		t.Fatalf("error decoding response: %s", err)
	}

	return &data
}

func deleteBlob(t *testing.T, oid string) int {
	objURL := fmt.Sprintf("http://localhost:8282/data/codingtest/%s", oid)
	req, _ := http.NewRequest("DELETE", objURL, nil)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("error making DELETE request: %s", err)
	}
	res.Body.Close() // ignore error

	return res.StatusCode
}
