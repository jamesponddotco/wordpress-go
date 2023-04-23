package wordpress

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"sync"
	"time"

	"git.sr.ht/~jamesponddotco/httpx-go"
	"git.sr.ht/~jamesponddotco/wordpress-go/internal/build"
	"git.sr.ht/~jamesponddotco/xstd-go/xerrors"
)

const (
	// ErrInvalidApplication is returned when an application is invalid.
	ErrInvalidApplication xerrors.Error = "invalid application"

	// ErrInvalidEndpoint is returned when an endpoint is invalid.
	ErrInvalidEndpoint xerrors.Error = "invalid endpoint"

	// ErrEndpointRequired is returned when a Config is created without an
	// endpoint.
	ErrEndpointRequired xerrors.Error = "endpoint required"

	// ErrUsernameRequired is returned when a Config is created without an username.
	ErrUsernameRequired xerrors.Error = "username required"

	// ErrPasswordRequired is returned when a Config is created without a password.
	ErrPasswordRequired xerrors.Error = "password required"

	// ErrApplicationRequired is returned when a Config is created without an
	// application.
	ErrApplicationRequired xerrors.Error = "application required"

	// ErrApplicationNameRequired is returned when an application is created
	// without a name.
	ErrApplicationNameRequired xerrors.Error = "application name required"

	// ErrApplicationVersionRequired is returned when an application is created
	// without a version.
	ErrApplicationVersionRequired xerrors.Error = "application version required"

	// ErrApplicationContactRequired is returned when an application is created
	// without contact information.
	ErrApplicationContactRequired xerrors.Error = "application contact required"
)

// Default values for the Config struct.
const (
	DefaultMaxRetries int           = 3
	DefaultTimeout    time.Duration = 60 * time.Second
)

// Logger defines the interface for logging. It is basically a thin wrapper
// around the standard logger which implements only a subset of the logger API.
type Logger interface {
	Printf(format string, v ...any)
}

// Application represents the application that is making requests to the API.
type Application struct {
	// Name is the name of the application.
	Name string

	// Version is the version of the application.
	Version string

	// Contact is the contact information for the application. Either an email
	// or an URL.
	Contact string
}

// NewApplication returns a new Application with the given values.
func NewApplication(name, version, contact string) (*Application, error) {
	app := &Application{
		Name:    name,
		Version: version,
		Contact: contact,
	}

	if err := app.validate(); err != nil {
		return nil, err
	}

	return app, nil
}

// DefaultApplication returns a new Application with default values.
func DefaultApplication() *Application {
	return &Application{
		Name:    build.Name,
		Version: build.Version,
		Contact: build.URL,
	}
}

// UserAgent returns the user agent string for the application.
func (a *Application) UserAgent() *httpx.UserAgent {
	if err := a.validate(); err != nil {
		return nil
	}

	return &httpx.UserAgent{
		Token:   a.Name,
		Version: a.Version,
		Comment: []string{
			a.Contact,
		},
	}
}

// Validate returns an error if the application is invalid.
func (a *Application) validate() error {
	if a.Name == "" {
		return fmt.Errorf("%w: %w", ErrInvalidApplication, ErrApplicationNameRequired)
	}

	if a.Version == "" {
		return fmt.Errorf("%w: %w", ErrInvalidApplication, ErrApplicationVersionRequired)
	}

	if a.Contact == "" {
		return fmt.Errorf("%w: %w", ErrInvalidApplication, ErrApplicationContactRequired)
	}

	return nil
}

// Config holds the basic configuration for the WordPress API.
type Config struct {
	// Application is the application that is making requests to the API.
	Application *Application

	// Logger is the logger to use for logging requests when debugging.
	Logger Logger

	// Endpoint is the URL of the WordPress website to make requests to.
	Endpoint string

	// Username is the username to use when making requests to the API.
	Username string

	// Password is the password used to authenticate with the API.
	Password string

	// MaxRetries specifies the maximum number of times to retry a request if it
	// fails due to rate limiting.
	//
	// This field is optional.
	MaxRetries int

	// Timeout is the time limit for requests made by the client to the  API.
	//
	// This field is optional.
	Timeout time.Duration

	// Debug specifies whether or not to enable debug logging.
	//
	// This field is optional.
	Debug bool

	// mu protects Config initialization.
	mu sync.Mutex
}

// NewConfig returns a new Config with the given values.
func NewConfig(app *Application, endpoint, username, password string) (*Config, error) {
	c := &Config{
		Application: app,
		Endpoint:    endpoint,
		Username:    username,
		Password:    password,
		MaxRetries:  DefaultMaxRetries,
		Timeout:     DefaultTimeout,
		Debug:       false,
	}

	if err := c.validate(); err != nil {
		return nil, err
	}

	return c, nil
}

// init initializes missing Config fields with their default values.
func (c *Config) init() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.Application == nil {
		c.Application = DefaultApplication()
	}

	if c.MaxRetries < 1 {
		c.MaxRetries = DefaultMaxRetries
	}

	if c.Timeout < 1 {
		c.Timeout = DefaultTimeout
	}

	if c.Logger == nil && c.Debug {
		c.Logger = log.New(os.Stderr, "", log.LstdFlags)
	}
}

// validate returns an error if the config is invalid.
func (c *Config) validate() error {
	if c.Application == nil {
		return ErrApplicationRequired
	}

	if err := c.Application.validate(); err != nil {
		return err
	}

	if c.Endpoint == "" {
		return ErrEndpointRequired
	}

	if _, err := url.Parse(c.Endpoint); err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidEndpoint, err)
	}

	if c.Username == "" {
		return ErrUsernameRequired
	}

	if c.Password == "" {
		return ErrPasswordRequired
	}

	return nil
}
