package tests_test

import (
	"context"
	"math/rand"
	"testing"

	"RedWood011/client/adapter/secret"
	"RedWood011/client/adapter/user"
	"RedWood011/client/entity"
	secretservice "RedWood011/client/service/secret"
	userservice "RedWood011/client/service/user"

	"github.com/stretchr/testify/require"
)

func startClient() (context.Context, *userservice.Service, *secretservice.Service) {
	ctx := context.Background()
	const defaultAddress = "localhost:5050"
	userAdapter := user.NewUserAdapter(defaultAddress)
	userService := userservice.NewUserService(userAdapter)
	secretAdapter := secret.NewSecretAdapter(defaultAddress)
	secretService := secretservice.NewSecretService(secretAdapter)

	return ctx, userService, secretService
}

func TestCreateUserTest(t *testing.T) {
	ctx, userService, _ := startClient()

	user := entity.User{
		Login:    randomString(),
		Password: randomString(),
	}
	accessToken, _, err := userService.RegisterUser(ctx, user)
	require.NoError(t, err)
	require.NotEqual(t, "", accessToken)
}

func TestAuthUserTest(t *testing.T) {
	ctx, userService, _ := startClient()
	user := entity.User{
		Login:    randomString(),
		Password: randomString(),
	}
	accessToken, _, err := userService.RegisterUser(ctx, user)
	require.NoError(t, err)
	require.NotEqual(t, "", accessToken)

	accessToken, _, err = userService.AuthUser(ctx, user)
	require.NoError(t, err)
	require.NotEqual(t, "", accessToken)
}

func TestCreateSecret(t *testing.T) {
	ctx, userService, secretService := startClient()
	user := entity.User{
		Login:    randomString(),
		Password: randomString(),
	}
	accessToken, _, err := userService.RegisterUser(ctx, user)
	require.NoError(t, err)
	require.NotEqual(t, "", accessToken)
	secretService.AccessToken = accessToken

	testTable := []struct {
		name     string
		secret   entity.Secret
		wantName []byte
		err      error
	}{
		{
			name: "secretCardSave",
			secret: entity.Secret{
				Name: []byte("Card"),
				Card: entity.Card{
					Number: "1234567890123456",
					Owner:  "Test",
					CVV:    "123",
					Date:   "2020-01-01",
				},
				TypeSecret: "card",
			},
			wantName: []byte("Card"),
			err:      nil,
		}, {
			name: "secretLoginPassSave",
			secret: entity.Secret{
				Name:       []byte("LoginPass"),
				PassLogin:  entity.PasswordLogin{Login: "testLogin", Password: "testPassword"},
				TypeSecret: "logPas",
			},
			wantName: []byte("LoginPass"),
			err:      nil,
		}, {
			name: "secretTextSave",
			secret: entity.Secret{
				Name:       []byte("Text"),
				Data:       []byte("sadhjadkjasjkdkjhashkjdkjhasdjhkasjhkdkjhsadqwewqeiqwosadds[fiadosfoidasoipf[adasdasvcvmxzweorp"),
				TypeSecret: "text",
			},
			wantName: []byte("Text"),
			err:      nil,
		}, {
			name: "secretBinSave",
			secret: entity.Secret{
				Name:       []byte("Bin"),
				Data:       []byte("sadhjadkjasjkdkjhashkjdkjhasdjhkasjhkdkjhsadqwewqeiqwosadds[fiadosfoidasoipf[adasdasvcvmxzweorp"),
				TypeSecret: "bin",
			},
			wantName: []byte("Bin"),
			err:      nil,
		},
	}

	for _, testCases := range testTable {
		t.Run(testCases.name, func(t *testing.T) {
			var id string
			id, err = secretService.CreateSecret(ctx, &testCases.secret)
			require.Equal(t, testCases.err, err)
			var newSecret *entity.Secret
			newSecret, err = secretService.GetSecret(ctx, id)
			require.Equal(t, testCases.err, err)
			require.Equal(t, id, newSecret.ID)
			require.Equal(t, testCases.wantName, newSecret.Name)
		})
	}
}

func TestDeleteSecret(t *testing.T) {
	ctx, userService, secretService := startClient()
	user := entity.User{
		Login:    randomString(),
		Password: randomString(),
	}
	accessToken, _, err := userService.RegisterUser(ctx, user)
	require.NoError(t, err)
	require.NotEqual(t, "", accessToken)
	secretService.AccessToken = accessToken

	name := []byte("Card")
	id, err := secretService.CreateSecret(ctx, &entity.Secret{
		Name: name,
		Card: entity.Card{
			Number: "1234567890123456",
			Owner:  "Test",
			CVV:    "123",
			Date:   "2020-01-01",
		},
		TypeSecret: "card",
	})
	require.NoError(t, err)

	err = secretService.DeleteSecret(ctx, id)
	require.NoError(t, err)
}

func TestListSecrets(t *testing.T) {
	ctx, userService, secretService := startClient()
	user := entity.User{
		Login:    randomString(),
		Password: randomString(),
	}
	accessToken, _, err := userService.RegisterUser(ctx, user)
	require.NoError(t, err)
	require.NotEqual(t, "", accessToken)
	secretService.AccessToken = accessToken

	name := []byte("Card")
	id, err := secretService.CreateSecret(ctx, &entity.Secret{
		Name: name,
		Card: entity.Card{
			Number: "1234567890123456",
			Owner:  "Test",
			CVV:    "123",
			Date:   "2020-01-01",
		},
		TypeSecret: "card",
	})
	require.NoError(t, err)

	var secrets []entity.Secret
	secrets, err = secretService.ListSecrets(ctx)
	require.NoError(t, err)
	require.Equal(t, id, secrets[0].ID)
	require.Equal(t, name, secrets[0].Name)
}

func randomString() string {
	letters := []rune("abcdefghijklmnopqrstuvwxyz")
	length := rand.Intn(3) + 6
	result := make([]rune, length)
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}
