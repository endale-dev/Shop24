package services

import (
	"context"
	"strings"
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"

	"shop24/internal/domain/models"
	"shop24/internal/repositories/interfaces"
	"shop24/internal/token"
)

type memUsers struct {
	users map[string]models.User
	seq   int
}

func newMemUsers() *memUsers {
	return &memUsers{users: map[string]models.User{}}
}

func (m *memUsers) Create(_ context.Context, user models.User) (models.User, error) {
	for _, existing := range m.users {
		if existing.Email == user.Email {
			return models.User{}, &interfaces.DuplicateKey{Field: "email"}
		}
	}
	m.seq++
	user.ID = strings.Repeat("a", m.seq)
	m.users[user.ID] = user
	return user, nil
}

func (m *memUsers) FindByEmail(_ context.Context, email string) (models.User, error) {
	for _, user := range m.users {
		if user.Email == email {
			return user, nil
		}
	}
	return models.User{}, interfaces.ErrNotFound
}

func (m *memUsers) FindByID(_ context.Context, id string) (models.User, error) {
	user, ok := m.users[id]
	if !ok {
		return models.User{}, interfaces.ErrNotFound
	}
	return user, nil
}

func (m *memUsers) FindByGoogleSub(_ context.Context, sub string) (models.User, error) {
	if sub == "" {
		return models.User{}, interfaces.ErrNotFound
	}
	for _, user := range m.users {
		if user.GoogleSub == sub {
			return user, nil
		}
	}
	return models.User{}, interfaces.ErrNotFound
}

func (m *memUsers) FindByResetHash(_ context.Context, hash string) (models.User, error) {
	if hash == "" {
		return models.User{}, interfaces.ErrNotFound
	}
	for _, user := range m.users {
		if user.ResetHash == hash {
			return user, nil
		}
	}
	return models.User{}, interfaces.ErrNotFound
}

func (m *memUsers) Update(_ context.Context, user models.User) error {
	if _, ok := m.users[user.ID]; !ok {
		return interfaces.ErrNotFound
	}
	m.users[user.ID] = user
	return nil
}

type captureMail struct {
	to   string
	body string
	sent int
}

func (c *captureMail) Send(_ context.Context, to, _, body string) error {
	c.sent++
	c.to = to
	c.body = body
	return nil
}

func TestGoogleLinksTheSameEmailAndRejectsADifferentSubject(t *testing.T) {
	users := newMemUsers()
	mailer := &captureMail{}
	auth := NewAuthService(users, token.NewIssuer("test-secret-value-32", time.Hour), mailer, "http://localhost:5173")
	ctx := context.Background()

	created, _, _, err := auth.Register(ctx, "Ada North", "ada.north@example.com", "wool-linen-shop")
	if err != nil {
		t.Fatal(err)
	}
	linked, _, _, err := auth.LoginGoogle(ctx, "sub-1", "ada.north@example.com", "Ada North")
	if err != nil {
		t.Fatal(err)
	}
	if linked.ID != created.ID || linked.GoogleSub != "sub-1" || linked.PasswordHash == "" {
		t.Fatalf("google login did not keep the password account: %+v", linked)
	}
	again, _, _, err := auth.LoginGoogle(ctx, "sub-1", "ada.north@example.com", "Ada North")
	if err != nil || again.ID != created.ID {
		t.Fatal(err)
	}
	if _, _, _, err := auth.LoginGoogle(ctx, "sub-2", "ada.north@example.com", "Ada North"); err == nil {
		t.Fatal("a second Google subject was accepted for the same email")
	}
	fresh, _, _, err := auth.LoginGoogle(ctx, "sub-3", "new.seller@example.com", "")
	if err != nil {
		t.Fatal(err)
	}
	if fresh.PasswordHash != "" || fresh.Name != "new.seller" {
		t.Fatalf("google-only account: %+v", fresh)
	}
	if _, _, _, err := auth.Login(ctx, fresh.Email, "wool-linen-shop"); err == nil {
		t.Fatal("an empty password hash accepted a password")
	}
}

func TestResetAndChangePassword(t *testing.T) {
	users := newMemUsers()
	mailer := &captureMail{}
	auth := NewAuthService(users, token.NewIssuer("test-secret-value-32", time.Hour), mailer, "http://localhost:5173")
	ctx := context.Background()
	user, _, _, err := auth.Register(ctx, "Ada North", "ada.north@example.com", "wool-linen-shop")
	if err != nil {
		t.Fatal(err)
	}
	if err := auth.RequestReset(ctx, "missing@example.com"); err != nil {
		t.Fatal(err)
	}
	if mailer.sent != 0 {
		t.Fatal("unknown email sent mail")
	}
	if err := auth.RequestReset(ctx, "ada.north@example.com"); err != nil {
		t.Fatal(err)
	}
	if mailer.sent != 1 || !strings.Contains(mailer.body, "/reset?token=") {
		t.Fatalf("reset mail: %s", mailer.body)
	}
	tokenValue := mailer.body[strings.Index(mailer.body, "token=")+6:]
	tokenValue = strings.TrimSpace(strings.Split(tokenValue, "\n")[0])
	if _, _, _, err := auth.ResetPassword(ctx, "not-a-token", "new-password"); err == nil {
		t.Fatal("bad token was accepted")
	}
	reset, _, _, err := auth.ResetPassword(ctx, tokenValue, "new-password")
	if err != nil {
		t.Fatal(err)
	}
	if bcrypt.CompareHashAndPassword([]byte(reset.PasswordHash), []byte("new-password")) != nil || reset.ResetHash != "" {
		t.Fatal("password was not replaced")
	}
	if _, _, _, err := auth.ResetPassword(ctx, tokenValue, "another-password"); err == nil {
		t.Fatal("reset token worked twice")
	}
	if _, err := auth.ChangePassword(ctx, user.ID, "wool-linen-shop", "third-password"); err == nil {
		t.Fatal("old password still worked")
	}
	changed, err := auth.ChangePassword(ctx, user.ID, "new-password", "third-password")
	if err != nil {
		t.Fatal(err)
	}
	if bcrypt.CompareHashAndPassword([]byte(changed.PasswordHash), []byte("third-password")) != nil {
		t.Fatal("change password did not stick")
	}
	named, err := auth.UpdateName(ctx, user.ID, "Ada N")
	if err != nil || named.Name != "Ada N" || named.Email != user.Email {
		t.Fatal(err)
	}
}
