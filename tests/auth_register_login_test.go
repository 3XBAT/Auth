package tests

import (
	"auth/tests/suite"
	authv1 "github.com/3XBAT/protos/gen/go"
	"github.com/brianvoe/gofakeit"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

const (
	passDefaultLen = 10
)

func TestRegisterLogin_Login_HappyPath(t *testing.T) {
	ctx, st := suite.New(t)

	name := gofakeit.Name()
	username := gofakeit.Username()
	password := RandomFakePassword()

	respReg, err := st.AuthClient.Register(ctx, &authv1.RegisterRequest{
		Name:     name,
		Username: username,
		Password: password,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, respReg.GetUserId())

	respLogin, err := st.AuthClient.Login(ctx, &authv1.LoginRequest{
		Username: username,
		Password: password,
	})
	require.NoError(t, err)

	loginTime := time.Now()

	token := respLogin.GetToken()

	tokenParsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})

	require.NoError(t, err)

	claims, ok := tokenParsed.Claims.(jwt.MapClaims)
	assert.True(t, ok)

	assert.Equal(t, respReg.GetUserId(), int64(claims["uid"].(float64)))
	assert.Equal(t, username, claims["username"].(string))

	const deltaSeconds = 1

	assert.InDelta(t, loginTime.Add(st.Cfg.TokenTTL).Unix(), claims["exp"].(float64), deltaSeconds)
}

func TestRegisterLogin_DuplicateRegistration(t *testing.T) {
	ctx, st := suite.New(t)

	name := gofakeit.Name()
	username := gofakeit.Username()
	password := RandomFakePassword()

	respReg, err := st.AuthClient.Register(ctx, &authv1.RegisterRequest{
		Name:     name,
		Username: username,
		Password: password,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, respReg.GetUserId())

	respReg, err = st.AuthClient.Register(ctx, &authv1.RegisterRequest{
		Name:     name,
		Username: username,
		Password: password,
	})

	require.Error(t, err)
	assert.Empty(t, respReg.GetUserId())
	assert.ErrorContains(t, err, "user already exists")
}

func TestRegister_FailCase(t *testing.T) {
	ctx, st := suite.New(t)

	tests := []struct {
		nameTest    string
		name        string
		username    string
		password    string
		expectedErr string
	}{
		{
			nameTest:    "Register with Empty Password",
			name:        gofakeit.Name(),
			username:    gofakeit.Username(),
			password:    "",
			expectedErr: "password is empty",
		},
		{
			nameTest:    "Register with Empty Username",
			name:        gofakeit.Name(),
			username:    "",
			password:    RandomFakePassword(),
			expectedErr: "username is empty",
		},
		{
			nameTest:    "Register with Both Empty",
			name:        gofakeit.Name(),
			username:    "",
			password:    "",
			expectedErr: "username is empty",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := st.AuthClient.Register(ctx, &authv1.RegisterRequest{
				Name:     test.name,
				Username: test.username,
				Password: test.password,
			})
			require.Error(t, err)
			assert.Contains(t, err.Error(), test.expectedErr)
		})
	}
}

func TestLogin_FailCase(t *testing.T) {
	ctx, st := suite.New(t)

	tests := []struct {
		nameTest    string
		username    string
		password    string
		expectedErr string
	}{
		{
			nameTest:    "Login with Empty Username",
			password:    RandomFakePassword(),
			username:    "",
			expectedErr: "username is empty",
		},
		{
			nameTest:    "Login with Empty Password",
			password:    "",
			username:    gofakeit.Username(),
			expectedErr: "password is empty",
		},
		{
			nameTest:    "Login with Both Empty",
			password:    "",
			username:    "",
			expectedErr: "password is empty",
		},
		{
			nameTest:    "Login with Non-Matching Password",
			password:    RandomFakePassword(),
			username:    gofakeit.Username(),
			expectedErr: "user not found",
		},
	}

	for _, test := range tests {
		_, err := st.AuthClient.Login(ctx, &authv1.LoginRequest{
			Username: test.username,
			Password: test.password,
		})
		require.Error(t, err)
		assert.Contains(t, err.Error(), test.expectedErr)
	}
}

func RandomFakePassword() string {
	return gofakeit.Password(true, true, true, true, false, passDefaultLen)
}
