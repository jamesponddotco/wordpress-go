package wordpress

import (
	"testing"
)

func TestConfig_Init(t *testing.T) {
	t.Parallel()

	var (
		endpoint = "https://example.com"
		username = "testuser"
		password = "testpassword"
	)

	tests := []struct {
		name     string
		input    *Config
		expected *Config
	}{
		{
			name: "Config with missing fields",
			input: &Config{
				Application: nil,
				Endpoint:    endpoint,
				Username:    username,
				Password:    password,
			},
			expected: &Config{
				Application: DefaultApplication(),
				Endpoint:    endpoint,
				Username:    username,
				Password:    password,
				MaxRetries:  DefaultMaxRetries,
				Timeout:     DefaultTimeout,
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tt.input.init()

			if tt.input.Application == nil || tt.input.Application.Name != DefaultApplication().Name ||
				tt.input.Application.Version != DefaultApplication().Version ||
				tt.input.Application.Contact != DefaultApplication().Contact {
				t.Errorf("Expected application %+v, got %+v", DefaultApplication(), tt.input.Application)
			}

			if tt.input.Endpoint != endpoint {
				t.Errorf("Expected endpoint %s, got %s", endpoint, tt.input.Endpoint)
			}

			if tt.input.Username != username {
				t.Errorf("Expected username %s, got %s", username, tt.input.Username)
			}

			if tt.input.Password != password {
				t.Errorf("Expected password %s, got %s", password, tt.input.Password)
			}

			if tt.input.MaxRetries != DefaultMaxRetries {
				t.Errorf("Expected max retries %d, got %d", DefaultMaxRetries, tt.input.MaxRetries)
			}

			if tt.input.Timeout != DefaultTimeout {
				t.Errorf("Expected timeout %v, got %v", DefaultTimeout, tt.input.Timeout)
			}
		})
	}
}

func TestConfig_Validate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name: "Config with invalid Application",
			config: &Config{
				Application: &Application{
					Name:    "",
					Version: "1.0.0",
					Contact: "https://example.com",
				},
				Endpoint:   "https://example.com",
				Username:   "testuser",
				Password:   "testpassword",
				MaxRetries: 3,
				Timeout:    DefaultTimeout,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.config.validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
