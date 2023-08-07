package tests

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoginSuccess_and_TestVerifySuccess(t *testing.T) {
	// TestLoginSuccess
	username, password := "Olenka", "dAmNmO!nAoBiZPi"
	client := &http.Client{}
	req, err := http.NewRequest("POST", "http://localhost:3000/login", nil)
	if err != nil {
		t.Fatalf("failed: %s", err)
	}
	req.SetBasicAuth(username, password)
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("failed: %s", err)
	}
	require.Equal(t, http.StatusOK, resp.StatusCode)

	// TestVerifySuccess
	req, err = http.NewRequest("POST", "http://localhost:3000/verify", nil)
	if err != nil {
		t.Fatalf("failed: %s", err)
	}
	for _, cookie := range resp.Cookies() {
		req.AddCookie(cookie)
	}
	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("failed: %s", err)
	}
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestLoginForbidden(t *testing.T) {
	var username string = "123"
	var password string = "dAmNmO!nAoBiZPi"
	client := &http.Client{}
	req, err := http.NewRequest("POST", "http://localhost:3000/login", nil)
	if err != nil {
		t.Fatalf("failed: %s", err)
	}
	req.SetBasicAuth(username, password)
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("failed: %s", err)
	}
	require.Equal(t, http.StatusForbidden, resp.StatusCode)
}

func TestVerifyForbidden(t *testing.T) {
	client := &http.Client{}
	req, err := http.NewRequest("POST", "http://localhost:3000/verify", nil)
	if err != nil {
		t.Fatalf("failed: %s", err)
	}
	cookie1 := &http.Cookie{
		Name:  "access_token",
		Value: "1",
	}
	cookie2 := &http.Cookie{
		Name:  "refresh_token",
		Value: "2",
	}
	req.AddCookie(cookie1)
	req.AddCookie(cookie2)

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("failed: %s", err)
	}
	require.Equal(t, http.StatusForbidden, resp.StatusCode)
}
