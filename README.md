# Client Authorization

## Middleware usage

#### func  WithClientIDAndPassKeyAuthorization

```go
func WithClientIDAndPassKeyAuthorization(authenticator ClientAuthenticator) middleware
```

#### type ClientAuthenticator

```go
type ClientAuthenticator interface {
	Authenticate(clientID, passKey string) error
}
```

#### type ClientAuthentication

```go
type ClientAuthentication struct {
}
```

#### func  NewClientAuthentication

```go
func NewClientAuthentication(authConfig *Config) *ClientAuthentication
```

#### func (*ClientAuthentication) Authenticate

```go
func (ca *ClientAuthentication) Authenticate(clientID, passKey string) error
```

#### type Config

```go
type Config struct {
}
```

#### func  NewConfig

```go
func NewConfig(dbDriver, dbConnURL string) *Config
```
