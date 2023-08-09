package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"RedWood011/client/adapter/secret"
	"RedWood011/client/adapter/user"
	"RedWood011/client/entity"
	secretservice "RedWood011/client/service/secret"
	userservice "RedWood011/client/service/user"
	"github.com/manifoldco/promptui"
)

var (
	BuildTime  string
	AppVersion string
	address    string
)

const defaultAddress = "localhost:5050"

func init() {
	envAddress := os.Getenv("ADDRESS")
	if envAddress != "" {
		address = envAddress
	}
	address = *flag.String("a", defaultAddress, "address of gGRPC server")
}

func main() {
	ctx := context.Background()
	fmt.Println("InitApp")
	fmt.Printf("App version: %v, Date compile: %v\n", AppVersion, BuildTime)
	userAdapter := user.NewUserAdapter(defaultAddress)
	userServise := userservice.NewUserService(userAdapter)

	InitPrompt := promptui.Select{
		Label: "Добро пожаловать! Что вы хотите сделать?",
		Items: []string{"Создать учетную запись", "Авторизоваться"},
	}

	_, action, err := InitPrompt.Run()

	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
		os.Exit(1)
	}

	var tokenAcssess string
init:
	for {
		switch action {
		case "Создать учетную запись":
			user := labelLoginPassword()
			tokenAcssess, _, err = userServise.RegisterUser(ctx, user)
			if err != nil {
				fmt.Printf("Ошибка при создании учетной записи, попробуйте еще раз")
			}
			break init

		case "Авторизоваться":
			user := labelLoginPassword()
			tokenAcssess, _, err = userServise.AuthUser(ctx, user)
			if err != nil {
				fmt.Printf("Ошибка при входе в учетную запись, попробуйте еще раз")
			}
			break init
		}
	}

	secretAdapter := secret.NewSecretAdapter(defaultAddress, tokenAcssess)
	secretService := secretservice.NewSecretService(secretAdapter)
	fmt.Println(action)
	SecretPrompt := promptui.Select{
		Label: "Хранилище секретов",
		Items: []string{"Создать секрет", "Посмотреть секрет", "Получить список секретов", "Удалить секрет"},
	}
	_, action, _ = SecretPrompt.Run()
	for {
		switch action {
		case "Создать секрет":
			var secret entity.Secret
			name, err := getInput("Введите название секрета", 0, "", 0)
			if err != nil {
				fmt.Println(err)
			}
			secret.Name = []byte(name)

			secret.Card.Number, err = getInput("Введите номер карты", 0, "", 0)
			if err != nil {
				fmt.Println(err)
			}
			secret.Card.CVV, err = getInput("Введите СVV", 0, "", 0)
			if err != nil {
				fmt.Println(err)
			}
			secret.Card.Date, err = getInput("Введите дату истечения срока действия карты", 0, "", 0)
			if err != nil {
				fmt.Println(err)
			}
			secret.Card.Owner, err = getInput("Введите владельца карты", 0, "", 0)
			if err != nil {
				fmt.Println(err)
			}
			secret.TypeSecret = "card"
			var secretID string
			secretID, err = secretService.CreateSecret(ctx, &secret)
			if err != nil {
				fmt.Printf("Ошибка при создании секрета, попробуйте еще раз")
			} else {
				fmt.Println("Секрет создан cоответствующим ID:", secretID)
			}

		case "Получить список секретов":
		case "Посмотреть секрет":
		case "Удалить секрет":

		}
	}

}

func labelLoginPassword() entity.User {
	loginPrompt := promptui.Prompt{
		Label: "Введите ваш логин",
	}
	login, err := loginPrompt.Run()
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
	}

	passwordPrompt := promptui.Prompt{
		Label: "Введите ваш пароль",
		Mask:  '*',
	}
	password, err := passwordPrompt.Run()
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
	}
	return entity.User{Password: password, Login: login}

}

func getInput(label string, minLen int, defVal string, mask rune) (string, error) {
	validate := func(input string) error {
		if len(input) < minLen {
			return fmt.Errorf("%s must have more than %d characters", label, minLen)
		}
		return nil
	}
	prompt := promptui.Prompt{
		Label:   label,
		Default: defVal,
		Mask:    mask,
	}
	if minLen > 0 {
		prompt.Validate = validate
	}
	title, err := prompt.Run()
	if err != nil {
		return "", err
	}
	return title, nil
}
