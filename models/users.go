package models

import (
	"errors"

	"github.com/alexander-emelyanenko/go-web-server/hash"
	"github.com/alexander-emelyanenko/go-web-server/rand"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"golang.org/x/crypto/bcrypt"
)

const hmacSecretKey = "secret-hmac-key"

var (
	// Verification that we implemented UserDB interface correctly
	_ UserDB = &userGorm{}
	// Verification that we implemented UserService interface correctly
	_                  UserService = &userService{}
	userPwPepper                   = "secret-random-string"
	ErrInvalidPassword             = errors.New("models: incorrect password provided")
	ErrNotFound                    = errors.New("models: resource not found")
	ErrInvalidID                   = errors.New("models: ID provided was invalid")
)

// User struct represents our user model
type User struct {
	gorm.Model
	Name         string
	Email        string `gorm:"not null;unique_index"`
	Password     string `gorm:"-"`
	PasswordHash string `gorm:"not null"`
	Remember     string `gorm:"-"`
	RememberHash string `gorm:"not null;unique_index"`
}

// userValFn is a type for validation functions
type userValFn func(*User) error

// runUserValFns runs all userValFns put in it
// and returns an error if something went wrong
func runUserValFns(user *User, fns ...userValFn) error {
	for _, fn := range fns {
		if err := fn(user); err != nil {
			return err
		}
	}
	return nil
}

// UserDB is used to interact with the users database.
//
// For pretty much all single user queries:
// If the user is found, we will return a nil error
// If the user is not found, we will return ErrNotFound
// If there is another error, we will return an error with
// more information about what went wrong. This may not be
// an error generated by the models package.
//
// For single user queries, any error but ErrNotFound should
// probably result
type UserDB interface {
	// Methods for querying for single users
	ByID(id uint) (*User, error)
	ByEmail(email string) (*User, error)
	ByRemember(token string) (*User, error)

	// Methods for altering users
	Create(user *User) error
	Update(user *User) error
	Delete(id uint) error

	// Used to close a DB connection
	Close() error

	// Migration helpers
	AutoMigrate() error
	DestructiveReset() error
}

// newUserGorm is the private constructor for our userGorm
func newUserGorm(connectionInfo string) (*userGorm, error) {
	db, err := gorm.Open("postgres", connectionInfo)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)
	return &userGorm{
		db: db,
	}, nil
}

// userGorm represents our database interface layer
// and implements the UserDB interface fully
type userGorm struct {
	db *gorm.DB
}

// ByID method fetchs single user by unique id
func (ug *userGorm) ByID(id uint) (*User, error) {
	var user User
	db := ug.db.Where("id = ?", id)
	err := first(db, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// ByEmail method fetchs single user by unique email
func (ug *userGorm) ByEmail(email string) (*User, error) {
	var user User
	db := ug.db.Where("email = ?", email)
	err := first(db, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// ByRemember method fetchs single user by remember token
func (ug *userGorm) ByRemember(rememberHash string) (*User, error) {
	var user User
	err := first(ug.db.Where("remember_hash = ?", rememberHash), &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Create method allows us to create a new user
func (ug *userGorm) Create(user *User) error {
	return ug.db.Create(user).Error
}

// Update method allows us to update user using User model struct
func (ug *userGorm) Update(user *User) error {
	return ug.db.Save(user).Error
}

// Delete method allow us to delete user by id
func (ug *userGorm) Delete(id uint) error {
	user := User{Model: gorm.Model{ID: id}}
	return ug.db.Delete(&user).Error
}

// Close method closes connection with UserDB
func (ug *userGorm) Close() error {
	return ug.db.Close()
}

// DestructiveReset method makes our life easier with DB migrations
func (ug *userGorm) DestructiveReset() error {
	if err := ug.db.DropTableIfExists(&User{}).Error; err != nil {
		return err
	}
	return ug.AutoMigrate()
}

// AutoMigrate method handles with migrations
func (ug *userGorm) AutoMigrate() error {
	if err := ug.db.AutoMigrate(&User{}).Error; err != nil {
		return err
	}
	return nil
}

// userValidator is our validation layer that validates
// and normalizes data before passing it on to the next
// UserDB in our interface chain.
type userValidator struct {
	UserDB
	hmac hash.HMAC
}

func (uv *userValidator) hmacRemember(user *User) error {
	if user.Remember == "" {
		return nil
	}
	user.RememberHash = uv.hmac.Hash(user.Remember)
	return nil
}

func (uv *userValidator) bcryptPassword(user *User) error {
	if user.Password == "" {
		return nil
	}

	pwBytes := []byte(user.Password + userPwPepper)
	hashedBytes, err := bcrypt.GenerateFromPassword(pwBytes, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hashedBytes)
	user.Password = ""
	return nil
}

// ByRemember validation method
func (uv *userValidator) ByRemember(token string) (*User, error) {
	user := User{
		Remember: token,
	}
	if err := runUserValFns(&user, uv.hmacRemember); err != nil {
		return nil, err
	}
	return uv.UserDB.ByRemember(user.RememberHash)
}

// Create validation method
func (uv *userValidator) Create(user *User) error {
	if user.Remember == "" {
		token, err := rand.RememberToken()
		if err != nil {
			return err
		}
		user.Remember = token
	}

	if err := runUserValFns(
		user,
		uv.bcryptPassword,
		uv.hmacRemember,
	); err != nil {
		return err
	}

	return uv.UserDB.Create(user)
}

// Update validation method
func (uv *userValidator) Update(user *User) error {
	if err := runUserValFns(
		user,
		uv.bcryptPassword,
		uv.hmacRemember,
	); err != nil {
		return err
	}

	return uv.UserDB.Update(user)
}

// Delete validation method
func (uv *userValidator) Delete(id uint) error {
	if id == 0 {
		return ErrInvalidID
	}
	return uv.UserDB.Delete(id)
}

// UserService is a set of methods used to manipulate and
// work with the user model
type UserService interface {
	// Authenticate will verify the provided email address and
	// password are correct. If they are correct, the user
	// corresponding to that email will be returned. Otherwise
	// You will receive either:
	// ErrNotFound, ErrInvalidPassword, or another error if
	// something goes wrong.
	Authenticate(email, password string) (*User, error)
	UserDB
}

// NewUserService is the public constructor for UserService
func NewUserService(connectionInfo string) (UserService, error) {
	ug, err := newUserGorm(connectionInfo)
	if err != nil {
		return nil, err
	}

	return &userService{
		UserDB: &userValidator{
			UserDB: ug,
			hmac:   hash.NewHMAC(hmacSecretKey),
		},
	}, nil
}

// userService defines our struct to work with a user model
type userService struct {
	UserDB
}

// Authenticate method for users
func (us *userService) Authenticate(email, password string) (*User, error) {
	foundUser, err := us.ByEmail(email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(foundUser.PasswordHash), []byte(password+userPwPepper))
	switch err {
	case nil:
		return foundUser, nil
	case bcrypt.ErrMismatchedHashAndPassword:
		return nil, ErrInvalidPassword
	default:
		return nil, err
	}
}

// first is a helper function to find first item in DB
func first(db *gorm.DB, dist interface{}) error {
	err := db.First(dist).Error
	if err == gorm.ErrRecordNotFound {
		return ErrNotFound
	}
	return err
}
