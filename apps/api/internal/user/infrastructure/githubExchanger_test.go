package infrastructure

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type githubFake struct {
	tokenBody   string
	tokenStatus int
	userBody    string
	userStatus  int
	emailsBody  string

	form       map[string][]string
	tokenAccep string
	userAuth   string
	userAccept string
	apiVersion string
}

func githubNewServer(t *testing.T, f *githubFake) *httptest.Server {
	t.Helper()
	if f.tokenStatus == 0 {
		f.tokenStatus = 200
	}
	if f.userStatus == 0 {
		f.userStatus = 200
	}
	mux := http.NewServeMux()
	mux.HandleFunc("POST /login/oauth/access_token", func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm()
		f.form = r.PostForm
		f.tokenAccep = r.Header.Get("Accept")
		w.WriteHeader(f.tokenStatus)
		_, _ = w.Write([]byte(f.tokenBody))
	})
	mux.HandleFunc("GET /user", func(w http.ResponseWriter, r *http.Request) {
		f.userAuth = r.Header.Get("Authorization")
		f.userAccept = r.Header.Get("Accept")
		f.apiVersion = r.Header.Get("X-GitHub-Api-Version")
		w.WriteHeader(f.userStatus)
		_, _ = w.Write([]byte(f.userBody))
	})
	mux.HandleFunc("GET /user/emails", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(f.emailsBody))
	})
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	return srv
}

func githubNewExchanger(srv *httptest.Server) *GitHubExchanger {
	return NewGitHubExchanger("cid", "csecret", WithGitHubBaseURLs(srv.URL, srv.URL))
}

func githubOKFake() *githubFake {
	return &githubFake{
		tokenBody:  `{"access_token":"gho_tok","token_type":"bearer"}`,
		userBody:   `{"id":4242,"login":"octo","name":"Octo Cat"}`,
		emailsBody: `[{"email":"o@example.com","primary":true,"verified":true}]`,
	}
}

func TestGitHubExchanger_Success(t *testing.T) {
	f := githubOKFake()
	srv := githubNewServer(t, f)
	id, err := githubNewExchanger(srv).Exchange(context.Background(), "the-code", "the-verifier", "app://cb")
	if err != nil {
		t.Fatal(err)
	}
	if id.Provider != "github" || id.Subject != "4242" || id.Login != "octo" || id.Name != "Octo Cat" ||
		id.Email != "o@example.com" || !id.EmailVerified {
		t.Fatalf("unexpected identity %+v", id)
	}
	for k, want := range map[string]string{
		"client_id": "cid", "client_secret": "csecret", "code": "the-code",
		"redirect_uri": "app://cb", "code_verifier": "the-verifier",
	} {
		if got := f.form[k]; len(got) != 1 || got[0] != want {
			t.Errorf("form[%s] = %v, want %s", k, got, want)
		}
	}
	if f.tokenAccep != "application/json" {
		t.Errorf("token Accept = %q", f.tokenAccep)
	}
	if f.userAuth != "Bearer gho_tok" || f.userAccept != "application/vnd.github+json" || f.apiVersion != "2022-11-28" {
		t.Errorf("api headers: %q %q %q", f.userAuth, f.userAccept, f.apiVersion)
	}
}

func TestGitHubExchanger_NullName(t *testing.T) {
	f := githubOKFake()
	f.userBody = `{"id":7,"login":"octo","name":null}`
	srv := githubNewServer(t, f)
	id, err := githubNewExchanger(srv).Exchange(context.Background(), "c", "v", "r")
	if err != nil {
		t.Fatal(err)
	}
	if id.Name != "" || id.Login != "octo" {
		t.Fatalf("got %+v", id)
	}
}

func TestGitHubExchanger_PrimaryUnverified(t *testing.T) {
	f := githubOKFake()
	f.emailsBody = `[{"email":"p@example.com","primary":true,"verified":false},{"email":"s@example.com","primary":false,"verified":true}]`
	srv := githubNewServer(t, f)
	id, err := githubNewExchanger(srv).Exchange(context.Background(), "c", "v", "r")
	if err != nil {
		t.Fatal(err)
	}
	if id.Email != "" || id.EmailVerified {
		t.Fatalf("got %+v", id)
	}
}

func TestGitHubExchanger_ExchangeRejected(t *testing.T) {
	for name, body := range map[string]string{
		"error json":      `{"error":"bad_verification_code","error_description":"The code passed is incorrect or expired."}`,
		"no access_token": `{}`,
	} {
		t.Run(name, func(t *testing.T) {
			f := githubOKFake()
			f.tokenBody = body
			srv := githubNewServer(t, f)
			if _, err := githubNewExchanger(srv).Exchange(context.Background(), "c", "v", "r"); err == nil {
				t.Fatal("expected error")
			}
		})
	}
}

func TestGitHubExchanger_UserUnauthorized(t *testing.T) {
	f := githubOKFake()
	f.userStatus = 401
	f.userBody = `{"message":"Bad credentials gho_tok"}`
	srv := githubNewServer(t, f)
	_, err := githubNewExchanger(srv).Exchange(context.Background(), "c", "v", "r")
	if err == nil {
		t.Fatal("expected error")
	}
	if strings.Contains(err.Error(), "gho_tok") || strings.Contains(err.Error(), "csecret") {
		t.Fatalf("error leaks secret: %v", err)
	}
}
