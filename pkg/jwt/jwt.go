package jwt

import (
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/smiletrl/micro_ecommerce/pkg/constants"
	"time"
)

type Service interface {
	ParseCustomerToken(c echo.Context) (customerID int64, err error)
	NewCustomerToken(customerID int64) (token string, err error)
}

type service struct {
	JwtSecret string
}

func NewService(secret string) Service {
	return &service{secret}
}

// Authorization header
var authScheme string = "Bearer"

// ParseToken get jwt token from request header, and then get user id signed inside the token
func (s *service) ParseCustomerToken(c echo.Context) (int64, error) {
	var (
		token *jwt.Token
		err   error
	)
	auth := c.Request().Header.Get("Authorization")
	l := len(authScheme)
	if len(auth) <= l+1 || auth[:l] != authScheme {
		return 0, errors.New("Missing auth bearer token in header")
	}
	tokenString := auth[l+1:]
	token, err = jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		// Check the signing method
		if t.Method.Alg() != constants.AlgorithmHS256 {
			return nil, errors.Errorf("header sign incorrect: %v", t.Header["alg"])
		}
		return []byte(s.JwtSecret), nil
	})
	if err != nil || !token.Valid {
		return 0, errors.New("Incorrect or outdated auth parameter")
	}

	claims := token.Claims.(jwt.MapClaims)
	customerID, ok := claims[constants.AuthCustomerID].(int64)
	if !ok {
		return 0, errors.Errorf("Unrecognized customer id: %v", claims[constants.AuthCustomerID])
	}

	return customerID, nil
}

// NewToken from user id. `user` is `staff` in this case.
// We may want to add other info into jwt token for other purposes.
func (s *service) NewCustomerToken(customerID int64) (string, error) {
	// Create token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims[constants.AuthCustomerID] = customerID

	// Set the token to be valid for 3 days.
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	// Generate encoded token and send it as response.
	tokenString, err := token.SignedString([]byte(s.JwtSecret))
	return tokenString, err
}

type mockService struct{}

func NewMockService() Service {
	return &mockService{}
}

func (m *mockService) ParseCustomerToken(c echo.Context) (int64, error) {
	return int64(0), nil
}

func (m *mockService) NewCustomerToken(userID int64) (string, error) {
	return "secret_token", nil
}
