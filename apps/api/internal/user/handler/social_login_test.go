package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"fittracker-api/internal/user"
	"fittracker-api/internal/user/mock"
	"fittracker-api/internal/user/service"
	"fittracker-api/pkg/auth"
)

const (
	googleFailed  = "Não foi possível entrar com Google. Tente novamente."
	githubFailed  = "Não foi possível entrar com GitHub. Tente novamente."
	googlePath    = "/api/auth/oauth/google"
	githubPath    = "/api/auth/oauth/github"
	githubBodyOK  = `{"code":"c","codeVerifier":"v","redirectUri":"app://cb"}`
	googleBodyOK  = `{"idToken":"tok"}`
	changePwdPath = "/api/auth/change-password"
)

type socialFixture struct {
	mux    *http.ServeMux
	signer *auth.Signer
	users  *mock.UserRepository
	google *mock.GoogleVerifier
	github *mock.GitHubExchanger
}

// newSocialFixture wires the router to a service with fake providers.
func newSocialFixture(t *testing.T, configured bool) *socialFixture {
	t.Helper()

	users := mock.NewUserRepository()
	resets := mock.NewPasswordResetRepository()
	signer, err := auth.NewSigner("test-secret")
	if err != nil {
		t.Fatalf("NewSigner: %v", err)
	}
	svc := service.New(users, resets, mock.NewMailer(), signer, auth.DeriveKey("test-secret", "password-reset"))

	f := &socialFixture{
		signer: signer,
		users:  users,
		google: &mock.GoogleVerifier{Identity: &user.ProviderIdentity{
			Provider: user.ProviderGoogle, Subject: "g1", Email: "Mixed@Case.com", EmailVerified: true, Name: "Gina",
		}},
		github: &mock.GitHubExchanger{Identity: &user.ProviderIdentity{
			Provider: user.ProviderGitHub, Subject: "42", Email: "Octo@Cat.com", EmailVerified: true, Login: "octo",
		}},
	}
	if configured {
		svc = svc.WithSocialLogin(mock.NewIdentityRepository(users, resets), f.google, f.github)
	}

	f.mux = http.NewServeMux()
	Router(f.mux, signer, svc)
	return f
}

func (f *socialFixture) do(method, path, body, token string) response {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	rec := httptest.NewRecorder()
	f.mux.ServeHTTP(rec, req)
	return response{status: rec.Code, body: strings.TrimSpace(rec.Body.String())}
}

func TestSocialLogin_UnconfiguredReturns503(t *testing.T) {
	f := newSocialFixture(t, false)

	if res := f.do(http.MethodPost, googlePath, googleBodyOK, ""); res.status != http.StatusServiceUnavailable ||
		res.body != "Login com Google indisponível." {
		t.Errorf("google = %+v, want 503 unavailable", res)
	}
	if res := f.do(http.MethodPost, githubPath, githubBodyOK, ""); res.status != http.StatusServiceUnavailable ||
		res.body != "Login com GitHub indisponível." {
		t.Errorf("github = %+v, want 503 unavailable", res)
	}

	// Existing endpoints keep working.
	if res := f.do(http.MethodPost, "/api/auth/login", `{"email":"x@y.com","password":"nope"}`, ""); res.status != http.StatusUnauthorized {
		t.Errorf("login = %+v, want 401", res)
	}
}

func TestSocialLogin_NilPortsReturn503(t *testing.T) {
	users := mock.NewUserRepository()
	resets := mock.NewPasswordResetRepository()
	signer, _ := auth.NewSigner("test-secret")
	svc := service.New(users, resets, mock.NewMailer(), signer, auth.DeriveKey("test-secret", "password-reset")).
		WithSocialLogin(mock.NewIdentityRepository(users, resets), nil, nil)
	mux := http.NewServeMux()
	Router(mux, signer, svc)

	for path, want := range map[string]string{
		googlePath: "Login com Google indisponível.",
		githubPath: "Login com GitHub indisponível.",
	} {
		body := googleBodyOK
		if path == githubPath {
			body = githubBodyOK
		}
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, path, strings.NewReader(body)))
		if rec.Code != http.StatusServiceUnavailable || strings.TrimSpace(rec.Body.String()) != want {
			t.Errorf("%s = %d %q, want 503 %q", path, rec.Code, rec.Body.String(), want)
		}
	}
}

func TestSocialLogin_BadBodyReturns400(t *testing.T) {
	f := newSocialFixture(t, true)

	cases := []struct {
		path, body, want string
	}{
		{googlePath, `{"idToken":`, googleFailed},
		{googlePath, `{}`, googleFailed},
		{googlePath, `{"idToken":"   "}`, googleFailed},
		{githubPath, `{"code":`, githubFailed},
		{githubPath, `{"codeVerifier":"v","redirectUri":"r"}`, githubFailed},
		{githubPath, `{"code":"c","redirectUri":"r"}`, githubFailed},
		{githubPath, `{"code":"c","codeVerifier":"v","redirectUri":" "}`, githubFailed},
	}
	for _, tc := range cases {
		res := f.do(http.MethodPost, tc.path, tc.body, "")
		if res.status != http.StatusBadRequest || res.body != tc.want {
			t.Errorf("%s %s = %+v, want 400 %q", tc.path, tc.body, res, tc.want)
		}
	}
	if f.google.Calls != 0 || f.github.Calls != 0 {
		t.Errorf("providers called for invalid bodies: google=%d github=%d", f.google.Calls, f.github.Calls)
	}
}

func decodeAuth(t *testing.T, res response) user.AuthResponse {
	t.Helper()
	if res.status != http.StatusOK {
		t.Fatalf("status = %d (%s), want 200", res.status, res.body)
	}
	var out user.AuthResponse
	if err := json.Unmarshal([]byte(res.body), &out); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if out.Token == "" || out.ExpiresAt.IsZero() {
		t.Errorf("missing token/expiresAt: %+v", out)
	}
	return out
}

func TestGoogleLogin_Success(t *testing.T) {
	f := newSocialFixture(t, true)

	out := decodeAuth(t, f.do(http.MethodPost, googlePath, googleBodyOK, ""))
	if out.User.Email != "mixed@case.com" {
		t.Errorf("email = %q, want lowercase", out.User.Email)
	}
	if f.google.LastToken != "tok" {
		t.Errorf("verifier token = %q, want tok", f.google.LastToken)
	}
}

func TestGitHubLogin_Success(t *testing.T) {
	f := newSocialFixture(t, true)

	out := decodeAuth(t, f.do(http.MethodPost, githubPath, githubBodyOK, ""))
	if out.User.Email != "octo@cat.com" {
		t.Errorf("email = %q, want lowercase", out.User.Email)
	}
	if f.github.LastCode != "c" || f.github.LastVerifier != "v" || f.github.LastRedirectURI != "app://cb" {
		t.Errorf("exchange args = %q %q %q", f.github.LastCode, f.github.LastVerifier, f.github.LastRedirectURI)
	}
}

func TestSocialLogin_ProviderFailureReturns401(t *testing.T) {
	f := newSocialFixture(t, true)
	f.google.Err = http.ErrAbortHandler
	f.github.Err = http.ErrAbortHandler

	if res := f.do(http.MethodPost, googlePath, googleBodyOK, ""); res.status != http.StatusUnauthorized || res.body != googleFailed {
		t.Errorf("google = %+v, want 401 %q", res, googleFailed)
	}
	if res := f.do(http.MethodPost, githubPath, githubBodyOK, ""); res.status != http.StatusUnauthorized || res.body != githubFailed {
		t.Errorf("github = %+v, want 401 %q", res, githubFailed)
	}
}

func TestSocialLogin_UnverifiedEmailReturns403(t *testing.T) {
	f := newSocialFixture(t, true)
	f.google.Identity.EmailVerified = false
	f.github.Identity.EmailVerified = false

	wantG := "Não foi possível entrar: sua conta Google não tem um e-mail verificado."
	wantH := "Não foi possível entrar: sua conta GitHub não tem um e-mail verificado."
	if res := f.do(http.MethodPost, googlePath, googleBodyOK, ""); res.status != http.StatusForbidden || res.body != wantG {
		t.Errorf("google = %+v, want 403 %q", res, wantG)
	}
	if res := f.do(http.MethodPost, githubPath, githubBodyOK, ""); res.status != http.StatusForbidden || res.body != wantH {
		t.Errorf("github = %+v, want 403 %q", res, wantH)
	}
}

func TestChangePassword_PasswordlessAccount(t *testing.T) {
	f := newSocialFixture(t, true)
	f.users.Seed(&user.User{ID: "u1", Email: "p@q.com", Name: "P"}, user.NoPasswordHash)

	token, _, err := f.signer.GenerateToken(auth.Claims{UserID: "u1", Email: "p@q.com", Name: "P"})
	if err != nil {
		t.Fatalf("GenerateToken: %v", err)
	}

	body := `{"currentPassword":"old","newPassword":"new-password"}`
	res := f.do(http.MethodPatch, changePwdPath, body, token)
	if res.status != http.StatusBadRequest || res.body != "Sua conta não tem senha; use Esqueci minha senha" {
		t.Errorf("PATCH = %+v, want 400 no-password message", res)
	}
	if got := f.users.Hash("u1"); !user.IsPasswordless(got) {
		t.Errorf("account hash changed to %q", got)
	}

	if res := f.do(http.MethodPost, changePwdPath, body, token); res.status != http.StatusMethodNotAllowed {
		t.Errorf("POST = %d, want 405", res.status)
	}
}
