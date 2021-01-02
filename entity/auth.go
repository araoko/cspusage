package entity

type Authenticator interface {
	Authenticate() (bool, error)
}
