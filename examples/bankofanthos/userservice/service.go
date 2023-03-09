package userservice

import (
	"context"
	"crypto/rsa"
	"errors"
	"fmt"
	"io/ioutil"
	"reflect"
	"regexp"
	"time"

	"github.com/ServiceWeaver/weaver"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

// CreateUserRequest contains data used for creating a new user.
type CreateUserRequest struct {
	weaver.AutoMarshal
	Username       string
	Password       string
	PasswordRepeat string
	FirstName      string
	LastName       string
	Birthday       string
	Timezone       string
	Address        string
	State          string
	Zip            string
	Ssn            string
}

// LoginRequest contains data used for logging in an existing user.
type LoginRequest struct {
	weaver.AutoMarshal
	Username string
	Password string
}

type T interface {
	// CreateUser is used to create a new user.
	CreateUser(ctx context.Context, r CreateUserRequest) error
	// Login logs in an existing user and returns a signed JWT on success.
	Login(ctx context.Context, r LoginRequest) (string, error)
}

type config struct {
	AccountDBURI       string `toml:"account_db_uri"`
	TokenExpirySeconds int    `toml:"token_expiry_seconds"`
	PrivateKeyPath     string `toml:"private_key_path"`
}

type impl struct {
	weaver.Implements[T]
	weaver.WithConfig[config]
	db         *userDB
	privateKey *rsa.PrivateKey
}

func (i *impl) Init(context.Context) error {
	privateKeyBytes, err := ioutil.ReadFile(i.Config().PrivateKeyPath)
	if err != nil {
		return err
	}
	i.privateKey, err = jwt.ParseRSAPrivateKeyFromPEM(privateKeyBytes)
	if err != nil {
		return err
	}
	i.db, err = newUserDB(i.Config().AccountDBURI)
	return err
}

func (i *impl) validateNewUser(r CreateUserRequest) error {
	v := reflect.ValueOf(r)
	// Start validation from field 1, since field 0 is weaver.AutoMarshal.
	for i := 1; i < v.NumField(); i++ {
		if v.Field(i).IsZero() {
			return fmt.Errorf("missing value for input field: %v", v.Field(i).Type().Name())
		}
	}

	// Verify username contains only 2-15 alphanumeric or underscore characters.
	if !regexp.MustCompile("[a-zA-Z0-9_]{2,15}").Match([]byte(r.Username)) {
		return errors.New("username must contain 2-15 alphanumeric characters or underscores")
	}
	if r.Password != r.PasswordRepeat {
		return errors.New("passwords do not match")
	}
	return nil
}

func (i *impl) CreateUser(ctx context.Context, r CreateUserRequest) error {
	if err := i.validateNewUser(r); err != nil {
		return err
	}
	user, err := i.db.getUser(r.Username)
	if err != nil {
		return err
	}
	if user != nil {
		err := fmt.Errorf("user %s already exists", r.Username)
		return err
	}
	i.Logger().Info("Creating password hash.")
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(r.Password), bcrypt.MinCost)
	if err != nil {
		return err
	}
	accountID := i.db.generateAccountID()

	userData := User{
		AccountID: accountID,
		Username:  r.Username,
		Passhash:  passwordHash,
		Firstname: r.FirstName,
		Lastname:  r.LastName,
		Birthday:  r.Birthday,
		Timezone:  r.Timezone,
		Address:   r.Address,
		State:     r.State,
		Zip:       r.Zip,
		SSN:       r.Ssn,
	}

	if err := i.db.addUser(userData); err != nil {
		return err
	}
	return nil
}

func (i *impl) Login(ctx context.Context, r LoginRequest) (string, error) {
	i.Logger().Debug("Getting user data.")
	user, err := i.db.getUser(r.Username)
	if err != nil {
		err = fmt.Errorf("error logging in: %v", err)
		return "", err
	}
	if user == nil {
		err = fmt.Errorf("user %s doesn't exist", r.Username)
		return "", err
	}
	i.Logger().Debug("Validating the password.")
	if err := bcrypt.CompareHashAndPassword(user.Passhash, []byte(r.Password)); err != nil {
		return "", err
	}

	i.Logger().Debug("Creating jwt token.")
	payload := jwt.MapClaims{
		"user": r.Username,
		"acct": user.AccountID,
		"name": user.Firstname + " " + user.Lastname,
		"iat":  time.Now().Unix(),
		"exp":  time.Now().Add(time.Duration(i.Config().TokenExpirySeconds) * time.Second).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, payload)
	tokenStr, err := token.SignedString(i.privateKey)
	if err != nil {
		return "", fmt.Errorf("couldn't sign jwt: %v", err)
	}
	i.Logger().Info("Login successful.")
	return tokenStr, nil
}
