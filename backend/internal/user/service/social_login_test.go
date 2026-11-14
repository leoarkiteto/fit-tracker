package service

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"fittracker-api/internal/user"
	"fittracker-api/pkg/auth"
)

func seedSocialAccount(t *testing.T, f *socialFixture, id, email, name, password string) {
	t.Helper()

	hash := user.NoPasswordHash
	if password != "" {
		var err error
		hash, err = auth.HashPassword(password)
		if err != nil {
			t.Fatalf("hashing seed password: %v", err)
		}
	}
	profileID := "profile-" + id
	f.users.Seed(&user.User{ID: id, Email: email, Name: name, ProfileID: &profileID}, hash)
}

func googleIdentity(sub, email string, verified bool) *user.ProviderIdentity {
	return &user.ProviderIdentity{
		Provider: user.ProviderGoogle, Subject: sub, Email: email, EmailVerified: verified, Name: "Ana Google",
	}
}

func requireError(t *testing.T, err error, want error) {
	t.Helper()
	appErr := requireAppError(t, err)
	wantErr := requireAppError(t, want)
	if appErr.Status != wantErr.Status || appErr.Message != wantErr.Message {
		t.Errorf("error = %d %q, want %d %q", appErr.Status, appErr.Message, wantErr.Status, wantErr.Message)
	}
}

// AC-1, AC-8 (FR-8): a new verified email creates a password-less account.
func TestLoginWithGoogle_CreatesAccount(t *testing.T) {
	f := newSocialTestService(t)
	f.google.Identity = googleIdentity("sub-1", "  Ana@X.com ", true)

	result, err := f.svc.LoginWithGoogle(context.Background(), "id-token")
	if err != nil {
		t.Fatalf("LoginWithGoogle() error = %v", err)
	}

	if f.google.LastToken != "id-token" {
		t.Errorf("verifier token = %q", f.google.LastToken)
	}
	if result.User.Email != "ana@x.com" || result.User.Name != "Ana Google" {
		t.Errorf("user = %+v", result.User)
	}
	if result.User.ProfileID == nil || *result.User.ProfileID == "" {
		t.Errorf("ProfileID = %v, want set", result.User.ProfileID)
	}
	if f.identities.Creates != 1 || f.identities.Links != 0 {
		t.Errorf("creates/links = %d/%d, want 1/0", f.identities.Creates, f.identities.Links)
	}
	if !user.IsPasswordless(f.users.Hash(result.User.ID)) {
		t.Error("new account must be password-less")
	}
	stored, ok := f.identities.Identity(user.ProviderGoogle, "sub-1")
	if !ok || stored.UserID != result.User.ID || stored.CreatedAt.IsZero() {
		t.Errorf("stored identity = %+v, %v", stored, ok)
	}
	if f.users.LastLoginUpdates != 1 {
		t.Errorf("LastLoginUpdates = %d, want 1", f.users.LastLoginUpdates)
	}

	claims, err := f.signer.ValidateToken(result.Token)
	if err != nil {
		t.Fatalf("token does not validate: %v", err)
	}
	if got := (*claims)["profileId"]; got != *result.User.ProfileID {
		t.Errorf("token profileId = %v, want %q", got, *result.User.ProfileID)
	}
	if got := (*claims)["sub"]; got != result.User.ID {
		t.Errorf("token sub = %v, want %q", got, result.User.ID)
	}
}

// AC-3: GitHub name falls back to login, then to the email local part.
func TestLoginWithGitHub_NameFallbacks(t *testing.T) {
	tests := []struct {
		name     string
		identity user.ProviderIdentity
		want     string
	}{
		{"name", user.ProviderIdentity{Subject: "1", Email: "a@x.com", EmailVerified: true, Name: " Ana ", Login: "ana-gh"}, "Ana"},
		{"login", user.ProviderIdentity{Subject: "2", Email: "b@x.com", EmailVerified: true, Name: " ", Login: " bob-gh "}, "bob-gh"},
		{"local part", user.ProviderIdentity{Subject: "3", Email: "Carl@x.com", EmailVerified: true}, "carl"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := newSocialTestService(t)
			identity := tt.identity
			identity.Provider = user.ProviderGitHub
			f.github.Identity = &identity

			result, err := f.svc.LoginWithGitHub(context.Background(), "code", "verifier", "app://cb")
			if err != nil {
				t.Fatalf("LoginWithGitHub() error = %v", err)
			}
			if result.User.Name != tt.want {
				t.Errorf("Name = %q, want %q", result.User.Name, tt.want)
			}
			if f.github.LastCode != "code" || f.github.LastVerifier != "verifier" || f.github.LastRedirectURI != "app://cb" {
				t.Errorf("exchanger args = %q %q %q", f.github.LastCode, f.github.LastVerifier, f.github.LastRedirectURI)
			}
		})
	}
}

// AC-4, AC-5, FR-6: no verified email and no linked identity creates nothing.
func TestSocialLogin_NoVerifiedEmail(t *testing.T) {
	f := newSocialTestService(t)
	f.google.Identity = googleIdentity("sub-1", "ana@x.com", false)
	f.github.Identity = &user.ProviderIdentity{Provider: user.ProviderGitHub, Subject: "9", Email: "", EmailVerified: false, Login: "x"}

	_, err := f.svc.LoginWithGoogle(context.Background(), "t")
	requireError(t, err, user.ErrSocialNoVerifiedEmail(user.ProviderGoogle))

	_, err = f.svc.LoginWithGitHub(context.Background(), "c", "v", "r")
	requireError(t, err, user.ErrSocialNoVerifiedEmail(user.ProviderGitHub))

	if f.identities.Creates != 0 || f.identities.Links != 0 || len(f.identities.Identities()) != 0 {
		t.Error("nothing may be created or linked")
	}
	if f.users.LastLoginUpdates != 0 {
		t.Errorf("LastLoginUpdates = %d, want 0", f.users.LastLoginUpdates)
	}
}

// AC-6: automatic link by case-insensitive email keeps the name, lowercases the email.
func TestLoginWithGoogle_LinksExistingAccount(t *testing.T) {
	f := newSocialTestService(t)
	seedSocialAccount(t, f, "u1", "Ana@X.com", "Ana Original", "secret123")
	f.google.Identity = googleIdentity("sub-1", "ana@x.com", true)

	result, err := f.svc.LoginWithGoogle(context.Background(), "t")
	if err != nil {
		t.Fatalf("LoginWithGoogle() error = %v", err)
	}

	if result.User.ID != "u1" || result.User.Name != "Ana Original" || result.User.Email != "ana@x.com" {
		t.Errorf("user = %+v", result.User)
	}
	if f.identities.Links != 1 || f.identities.Creates != 0 {
		t.Errorf("links/creates = %d/%d, want 1/0", f.identities.Links, f.identities.Creates)
	}
	stored, ok := f.identities.Identity(user.ProviderGoogle, "sub-1")
	if !ok || stored.UserID != "u1" {
		t.Errorf("stored identity = %+v, %v", stored, ok)
	}
	if f.users.LastLoginUpdates != 1 {
		t.Errorf("LastLoginUpdates = %d, want 1", f.users.LastLoginUpdates)
	}
}

// AC-7 (FR-3): a linked identity wins over the email, and changes nothing.
func TestLoginWithGoogle_LinkedIdentityIgnoresEmail(t *testing.T) {
	f := newSocialTestService(t)
	seedSocialAccount(t, f, "a", "a@x.com", "A", "")
	f.identities.SeedIdentity(user.LinkedIdentity{UserID: "a", Provider: user.ProviderGoogle, Subject: "sub-1"})
	f.google.Identity = googleIdentity("sub-1", "other@y.com", false)

	result, err := f.svc.LoginWithGoogle(context.Background(), "t")
	if err != nil {
		t.Fatalf("LoginWithGoogle() error = %v", err)
	}

	if result.User.ID != "a" || result.User.Email != "a@x.com" {
		t.Errorf("user = %+v", result.User)
	}
	if f.identities.Creates != 0 || f.identities.Links != 0 {
		t.Error("nothing may be created or linked")
	}
	if f.users.LastLoginUpdates != 1 {
		t.Errorf("LastLoginUpdates = %d, want 1", f.users.LastLoginUpdates)
	}
}

// AC-8: more than one account with the verified email is refused as FR-7.
func TestSocialLogin_AmbiguousEmail(t *testing.T) {
	f := newSocialTestService(t)
	seedSocialAccount(t, f, "u1", "Ana@x.com", "A1", "secret123")
	seedSocialAccount(t, f, "u2", "ana@x.com", "A2", "secret123")
	f.google.Identity = googleIdentity("sub-1", "ana@x.com", true)

	_, err := f.svc.LoginWithGoogle(context.Background(), "t")

	requireError(t, err, user.ErrSocialLoginFailed(user.ProviderGoogle))
	if f.identities.Links != 0 || f.identities.Creates != 0 || len(f.identities.Identities()) != 0 {
		t.Error("nothing may be linked or created")
	}
	if user.IsPasswordless(f.users.Hash("u1")) || user.IsPasswordless(f.users.Hash("u2")) {
		t.Error("passwords must be untouched")
	}
}

// FR-7: provider failures and repository failures yield the generic error.
func TestSocialLogin_FailuresAreGeneric(t *testing.T) {
	boom := errors.New("boom secret-payload")

	t.Run("google error", func(t *testing.T) {
		f := newSocialTestService(t)
		f.google.Err = boom
		_, err := f.svc.LoginWithGoogle(context.Background(), "t")
		requireError(t, err, user.ErrSocialLoginFailed(user.ProviderGoogle))
	})
	t.Run("google nil identity", func(t *testing.T) {
		f := newSocialTestService(t)
		_, err := f.svc.LoginWithGoogle(context.Background(), "t")
		requireError(t, err, user.ErrSocialLoginFailed(user.ProviderGoogle))
	})
	t.Run("github error", func(t *testing.T) {
		f := newSocialTestService(t)
		f.github.Err = boom
		_, err := f.svc.LoginWithGitHub(context.Background(), "c", "v", "r")
		requireError(t, err, user.ErrSocialLoginFailed(user.ProviderGitHub))
	})
	t.Run("find identity error", func(t *testing.T) {
		f := newSocialTestService(t)
		f.google.Identity = googleIdentity("s", "a@x.com", true)
		f.identities.FindUserByIdentityErr = boom
		_, err := f.svc.LoginWithGoogle(context.Background(), "t")
		requireError(t, err, user.ErrSocialLoginFailed(user.ProviderGoogle))
	})
	t.Run("find by email error", func(t *testing.T) {
		f := newSocialTestService(t)
		f.google.Identity = googleIdentity("s", "a@x.com", true)
		f.identities.FindUsersByEmailFoldErr = boom
		_, err := f.svc.LoginWithGoogle(context.Background(), "t")
		requireError(t, err, user.ErrSocialLoginFailed(user.ProviderGoogle))
	})
	t.Run("create error", func(t *testing.T) {
		f := newSocialTestService(t)
		f.google.Identity = googleIdentity("s", "a@x.com", true)
		f.identities.CreateErr = boom
		_, err := f.svc.LoginWithGoogle(context.Background(), "t")
		requireError(t, err, user.ErrSocialLoginFailed(user.ProviderGoogle))
	})
}

// AC-14 (service level): unconfigured ports answer 503.
func TestSocialLogin_Unavailable(t *testing.T) {
	svc, _ := newTestService(t, nil)

	_, err := svc.LoginWithGoogle(context.Background(), "t")
	requireError(t, err, user.ErrSocialUnavailable(user.ProviderGoogle))
	if appErr := requireAppError(t, err); appErr.Status != http.StatusServiceUnavailable {
		t.Errorf("status = %d, want 503", appErr.Status)
	}

	_, err = svc.LoginWithGitHub(context.Background(), "c", "v", "r")
	requireError(t, err, user.ErrSocialUnavailable(user.ProviderGitHub))

	// Identities missing makes both unavailable even with providers set.
	f := newSocialTestService(t)
	f.svc.identities = nil
	_, err = f.svc.LoginWithGoogle(context.Background(), "t")
	requireError(t, err, user.ErrSocialUnavailable(user.ProviderGoogle))
	if f.google.Calls != 0 {
		t.Error("verifier must not be called when unavailable")
	}

	// Only one provider configured.
	f = newSocialTestService(t)
	f.svc.github = nil
	_, err = f.svc.LoginWithGitHub(context.Background(), "c", "v", "r")
	requireError(t, err, user.ErrSocialUnavailable(user.ProviderGitHub))
}

// The same verified email from Google and GitHub lands on one account.
func TestSocialLogin_TwoProvidersOneAccount(t *testing.T) {
	f := newSocialTestService(t)
	f.google.Identity = googleIdentity("g-1", "ana@x.com", true)
	f.github.Identity = &user.ProviderIdentity{Provider: user.ProviderGitHub, Subject: "42", Email: "ana@x.com", EmailVerified: true, Login: "ana"}

	first, err := f.svc.LoginWithGoogle(context.Background(), "t")
	if err != nil {
		t.Fatalf("google: %v", err)
	}
	second, err := f.svc.LoginWithGitHub(context.Background(), "c", "v", "r")
	if err != nil {
		t.Fatalf("github: %v", err)
	}

	if first.User.ID != second.User.ID {
		t.Errorf("accounts differ: %q vs %q", first.User.ID, second.User.ID)
	}
	if f.identities.Creates != 1 || f.identities.Links != 1 || len(f.identities.Identities()) != 2 {
		t.Errorf("creates/links/identities = %d/%d/%d", f.identities.Creates, f.identities.Links, len(f.identities.Identities()))
	}
}

// AC-11: a password-less account answers like a wrong password, even with an empty one.
func TestLogin_PasswordlessAccount(t *testing.T) {
	f := newSocialTestService(t)
	seedSocialAccount(t, f, "u1", "ana@x.com", "Ana", "")

	for _, password := range []string{"", "anything"} {
		_, err := f.svc.Login(context.Background(), "ana@x.com", password)
		requireError(t, err, user.ErrInvalidCredentials)
	}
	if f.users.LastLoginUpdates != 0 {
		t.Errorf("LastLoginUpdates = %d, want 0", f.users.LastLoginUpdates)
	}
}

// AC-12: change password on a password-less account is refused and changes nothing.
func TestChangePassword_PasswordlessAccount(t *testing.T) {
	f := newSocialTestService(t)
	seedSocialAccount(t, f, "u1", "ana@x.com", "Ana", "")

	err := f.svc.ChangePassword(context.Background(), "u1", "", "newsecret")

	requireError(t, err, user.ErrNoPassword)
	if appErr := requireAppError(t, err); appErr.Status != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", appErr.Status)
	}
	if f.users.PasswordUpdates != 0 || !user.IsPasswordless(f.users.Hash("u1")) {
		t.Error("nothing may change")
	}
}

// AC-13: a password-less account sets a password through recovery and still logs in socially.
func TestResetPassword_PasswordlessAccount(t *testing.T) {
	f := newSocialTestService(t)
	ctx := context.Background()
	seedSocialAccount(t, f, "u1", "ana@x.com", "Ana", "")
	f.identities.SeedIdentity(user.LinkedIdentity{UserID: "u1", Provider: user.ProviderGoogle, Subject: "sub-1"})
	f.google.Identity = googleIdentity("sub-1", "ana@x.com", true)

	if err := f.svc.RequestPasswordReset(ctx, "ana@x.com"); err != nil {
		t.Fatalf("RequestPasswordReset() error = %v", err)
	}
	code := f.mailer.Sent[len(f.mailer.Sent)-1].Code
	if err := f.svc.ResetPassword(ctx, "ana@x.com", code, "brandnew1"); err != nil {
		t.Fatalf("ResetPassword() error = %v", err)
	}

	if _, err := f.svc.Login(ctx, "ana@x.com", "brandnew1"); err != nil {
		t.Fatalf("Login() with new password error = %v", err)
	}
	if _, err := f.svc.LoginWithGoogle(ctx, "t"); err != nil {
		t.Fatalf("LoginWithGoogle() error = %v", err)
	}
	if _, err := f.svc.Login(ctx, "ana@x.com", "brandnew1"); err != nil {
		t.Fatalf("password must survive social login: %v", err)
	}
}

// AC-23: the account pre-hijacking scenario.
func TestSocialLogin_AttackerScenario(t *testing.T) {
	f := newSocialTestService(t)
	ctx := context.Background()
	now := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	f.svc.now = func() time.Time { return now }

	if _, err := f.svc.Register(ctx, "vitima@gmail.com", "attacker-pass", "Atacante"); err != nil {
		t.Fatalf("Register() error = %v", err)
	}
	if err := f.svc.RequestPasswordReset(ctx, "vitima@gmail.com"); err != nil {
		t.Fatalf("RequestPasswordReset() error = %v", err)
	}
	oldCode := f.mailer.Sent[len(f.mailer.Sent)-1].Code

	f.google.Identity = googleIdentity("victim-sub", "vitima@gmail.com", true)
	if _, err := f.svc.LoginWithGoogle(ctx, "t"); err != nil {
		t.Fatalf("LoginWithGoogle() error = %v", err)
	}
	if f.identities.Links != 1 {
		t.Fatalf("Links = %d, want 1", f.identities.Links)
	}

	_, err := f.svc.Login(ctx, "vitima@gmail.com", "attacker-pass")
	requireError(t, err, user.ErrInvalidCredentials)

	now = now.Add(user.ResetCodeCooldown + time.Second)
	err = f.svc.ResetPassword(ctx, "vitima@gmail.com", oldCode, "attacker-new")
	requireError(t, err, user.ErrInvalidResetCode)

	if err := f.svc.RequestPasswordReset(ctx, "vitima@gmail.com"); err != nil {
		t.Fatalf("second RequestPasswordReset() error = %v", err)
	}
	newCode := f.mailer.Sent[len(f.mailer.Sent)-1].Code
	if err := f.svc.ResetPassword(ctx, "vitima@gmail.com", newCode, "victim-pass"); err != nil {
		t.Fatalf("ResetPassword() error = %v", err)
	}
	if _, err := f.svc.Login(ctx, "vitima@gmail.com", "victim-pass"); err != nil {
		t.Fatalf("Login() with new password error = %v", err)
	}

	if _, err := f.svc.LoginWithGoogle(ctx, "t"); err != nil {
		t.Fatalf("second LoginWithGoogle() error = %v", err)
	}
	if _, err := f.svc.Login(ctx, "vitima@gmail.com", "victim-pass"); err != nil {
		t.Fatalf("password must still be valid after social login: %v", err)
	}
	if f.identities.Links != 1 {
		t.Errorf("Links = %d, want 1 (second login uses the identity)", f.identities.Links)
	}
}

// AC-24: a failed link leaves identity, password and reset code untouched.
func TestSocialLogin_LinkFailureChangesNothing(t *testing.T) {
	f := newSocialTestService(t)
	ctx := context.Background()
	if _, err := f.svc.Register(ctx, "vitima@gmail.com", "attacker-pass", "Atacante"); err != nil {
		t.Fatalf("Register() error = %v", err)
	}
	if err := f.svc.RequestPasswordReset(ctx, "vitima@gmail.com"); err != nil {
		t.Fatalf("RequestPasswordReset() error = %v", err)
	}
	account, _, _ := f.users.GetByEmail(ctx, "vitima@gmail.com")
	hashBefore := f.users.Hash(account.ID)

	f.identities.LinkErr = errors.New("write failed")
	f.google.Identity = googleIdentity("victim-sub", "vitima@gmail.com", true)

	_, err := f.svc.LoginWithGoogle(ctx, "t")

	requireError(t, err, user.ErrSocialLoginFailed(user.ProviderGoogle))
	if len(f.identities.Identities()) != 0 {
		t.Error("no identity may be linked")
	}
	if f.users.Hash(account.ID) != hashBefore {
		t.Error("password hash changed")
	}
	if _, ok := f.resets.Code(account.ID); !ok {
		t.Error("reset code was removed")
	}
	if _, err := f.svc.Login(ctx, "vitima@gmail.com", "attacker-pass"); err != nil {
		t.Errorf("password login must still work: %v", err)
	}
}

// AC-25: a mixed-case account is linked, its email lowercased, and recovery works in lowercase.
func TestSocialLogin_MixedCaseAccountThenRecovery(t *testing.T) {
	f := newSocialTestService(t)
	ctx := context.Background()
	if _, err := f.svc.Register(ctx, "Vitima@Gmail.com", "attacker-pass", "Vitima"); err != nil {
		t.Fatalf("Register() error = %v", err)
	}
	f.google.Identity = googleIdentity("victim-sub", "vitima@gmail.com", true)

	result, err := f.svc.LoginWithGoogle(ctx, "t")
	if err != nil {
		t.Fatalf("LoginWithGoogle() error = %v", err)
	}
	if result.User.Email != "vitima@gmail.com" {
		t.Errorf("response email = %q, want lowercase", result.User.Email)
	}
	stored, found, _ := f.users.GetByEmail(ctx, "vitima@gmail.com")
	if !found || stored.Email != "vitima@gmail.com" {
		t.Fatalf("stored account = %+v, found = %v", stored, found)
	}

	if err := f.svc.RequestPasswordReset(ctx, "vitima@gmail.com"); err != nil {
		t.Fatalf("RequestPasswordReset() error = %v", err)
	}
	if len(f.mailer.Sent) != 1 {
		t.Fatalf("codes sent = %d, want 1", len(f.mailer.Sent))
	}
	if err := f.svc.ResetPassword(ctx, "vitima@gmail.com", f.mailer.Sent[0].Code, "victim-pass"); err != nil {
		t.Fatalf("ResetPassword() error = %v", err)
	}
	if _, err := f.svc.Login(ctx, "vitima@gmail.com", "victim-pass"); err != nil {
		t.Fatalf("Login() error = %v", err)
	}
}
