package kitserver

import (
	"testing"

	"golang.org/x/net/context"

	"github.com/kolide/kolide-ose/datastore"
	"github.com/kolide/kolide-ose/kolide"
)

func TestCreateUser(t *testing.T) {
	ds, _ := datastore.New("mock", "")
	svc, _ := NewService(ds)

	var createUserTests = []struct {
		Username           *string
		Password           *string
		Email              *string
		NeedsPasswordReset *bool
		Admin              *bool
		Err                error
	}{
		{
			Username: stringPtr("admin1"),
			Password: stringPtr("foobar"),
			Err:      errInvalidArgument,
		},
		{
			Username:           stringPtr("admin1"),
			Password:           stringPtr("foobar"),
			Email:              stringPtr("admin1@example.com"),
			NeedsPasswordReset: boolPtr(true),
			Admin:              boolPtr(false),
		},
	}

	ctx := context.Background()
	for _, tt := range createUserTests {
		payload := kolide.UserPayload{
			Username:           tt.Username,
			Password:           tt.Password,
			Email:              tt.Email,
			Admin:              tt.Admin,
			NeedsPasswordReset: tt.NeedsPasswordReset,
		}
		user, err := svc.NewUser(ctx, payload)
		if err != nil {
			if err != tt.Err {
				t.Fatalf("got %q, want %q", err, tt.Err)
			}
			continue
		}

		if user.ID == 0 {
			t.Errorf("expected a user ID, got 0")
		}

		if err := user.ValidatePassword(*tt.Password); err != nil {
			t.Errorf("expected nil, got %q", err)
		}

		if err := user.ValidatePassword("different_password!"); err == nil {
			t.Errorf("expected err, got nil")
		}

		if have, want := user.NeedsPasswordReset, *tt.NeedsPasswordReset; have != want {
			t.Errorf("have %v want %v", have, want)
		}

		if have, want := user.NeedsPasswordReset, *tt.NeedsPasswordReset; have != want {
			t.Errorf("have %v want %v", have, want)
		}

		if have, want := user.Admin, *tt.Admin; have != want {
			t.Errorf("have %v want %v", have, want)
		}

		// check duplicate creation
		_, err = svc.NewUser(ctx, payload)
		if err != datastore.ErrExists {
			t.Errorf("have %q, want %q", err, datastore.ErrExists)
		}
	}
}

func stringPtr(s string) *string {
	return &s
}

func boolPtr(b bool) *bool {
	return &b
}
