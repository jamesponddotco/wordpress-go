package wordpress_test

import (
	"fmt"
	"reflect"
	"testing"

	"git.sr.ht/~jamesponddotco/httpx-go"
	"git.sr.ht/~jamesponddotco/wordpress-go"
)

func TestApplication_NewApplication(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		giveName    string
		giveVersion string
		giveContact string
		wantErr     bool
	}{
		{
			name:        "Valid application",
			giveName:    "WordPress",
			giveVersion: "5.4.1",
			giveContact: "https://wordpress.org/",
			wantErr:     false,
		},
		{
			name:        "Missing application name",
			giveName:    "",
			giveVersion: "5.4.1",
			giveContact: "https://wordpress.org/",
			wantErr:     true,
		},
		{
			name:        "Missing application version",
			giveName:    "WordPress",
			giveVersion: "",
			giveContact: "https://wordpress.org/",
			wantErr:     true,
		},
		{
			name:        "Missing application contact",
			giveName:    "WordPress",
			giveVersion: "5.4.1",
			giveContact: "",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := wordpress.NewApplication(tt.giveName, tt.giveVersion, tt.giveContact)
			if (err != nil) != tt.wantErr {
				t.Error("Application.UserAgent(): want error, got nothing")
				return
			}
		})
	}
}

func TestApplication_UserAgent(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		giveApplication *wordpress.Application
		wantUserAgent   *httpx.UserAgent
	}{
		{
			name: "Valid application",
			giveApplication: &wordpress.Application{
				Name:    "WordPress",
				Version: "5.4.1",
				Contact: "https://wordpress.org/",
			},
			wantUserAgent: &httpx.UserAgent{
				Token:   "WordPress",
				Version: "5.4.1",
				Comment: []string{"https://wordpress.org/"},
			},
		},
		{
			name: "Invalid application",
			giveApplication: &wordpress.Application{
				Name:    "WordPress",
				Version: "5.4.1",
			},
			wantUserAgent: nil,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gotUserAgent := tt.giveApplication.UserAgent()
			if !reflect.DeepEqual(gotUserAgent, tt.wantUserAgent) {
				t.Errorf("Application.UserAgent(): want %v, got %v", tt.wantUserAgent, gotUserAgent)
			}
		})
	}
}

func TestConfig_NewConfig(t *testing.T) {
	t.Parallel()

	app, err := wordpress.NewApplication("Test App", "1.0.0", "https://testapp.com")
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name     string
		app      *wordpress.Application
		endpoint string
		username string
		password string
		err      error
	}{
		{
			name:     "ValidConfig",
			app:      app,
			endpoint: "https://example.com",
			username: "username",
			password: "password",
			err:      nil,
		},
		{
			name:     "NilApplication",
			app:      nil,
			endpoint: "https://example.com",
			username: "username",
			password: "password",
			err:      wordpress.ErrApplicationRequired,
		},
		{
			name:     "EmptyEndpoint",
			app:      app,
			endpoint: "",
			username: "username",
			password: "password",
			err:      wordpress.ErrEndpointRequired,
		},
		{
			name:     "InvalidEndpoint",
			app:      app,
			endpoint: "://example.com",
			username: "username",
			password: "password",
			err:      fmt.Errorf("invalid endpoint: parse \"://example.com\": missing protocol scheme"),
		},
		{
			name:     "EmptyUsername",
			app:      app,
			endpoint: "https://example.com",
			username: "",
			password: "password",
			err:      wordpress.ErrUsernameRequired,
		},
		{
			name:     "EmptyPassword",
			app:      app,
			endpoint: "https://example.com",
			username: "username",
			password: "",
			err:      wordpress.ErrPasswordRequired,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			config, err := wordpress.NewConfig(tt.app, tt.endpoint, tt.username, tt.password)
			switch {
			case tt.err != nil:
				if err == nil || err.Error() != tt.err.Error() {
					t.Errorf("expected error %v, got %v", tt.err, err)
				}
			case err != nil:
				t.Errorf("unexpected error: %v", err)
			default:
				if config.Application != tt.app || config.Endpoint != tt.endpoint ||
					config.Username != tt.username || config.Password != tt.password {
					t.Errorf("expected config with app: %v, endpoint: %s, username: %s, "+
						"password: %s, got app: %v, endpoint: %s, username: %s, password: %s",
						tt.app, tt.endpoint, tt.username, tt.password,
						config.Application, config.Endpoint, config.Username, config.Password)
				}
			}
		})
	}
}
