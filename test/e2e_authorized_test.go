package test

import (
	"net/http"
	"testing"
	"time"
)

func TestAuthorized(t *testing.T) {

	users := []User{
		{
			Username: "foo",
			Password: "bar",
		},
		{
			Username: "lobby",
			Password: "niu",
		},
	}
	startEnvoy(users)
	time.Sleep(5 * time.Second)

	//Unauthorized
	req, err := http.NewRequest(http.MethodGet, "http://localhost:10000/", nil)
	resp1, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp1.Body.Close()
	if resp1.StatusCode != http.StatusUnauthorized {
		t.Fatalf("unexpected status code: %v", resp1.StatusCode)
	}

	//Unauthorized
	req.SetBasicAuth(users[0].Username, "wrong password")
	resp2, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp2.Body.Close()
	if resp2.StatusCode != http.StatusUnauthorized {
		t.Fatalf("unexpected status code: %v", resp2.StatusCode)
	}

	//Authorized
	req.SetBasicAuth(users[0].Username, users[0].Password)
	resp3, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp3.Body.Close()
	if resp3.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status code: %v", resp2.StatusCode)
	}

	//Authorized
	req.SetBasicAuth(users[1].Username, users[1].Password)
	resp4, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp4.Body.Close()
	if resp4.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status code: %v", resp2.StatusCode)
	}

	t.Log("TestAuthorized pass")
}
