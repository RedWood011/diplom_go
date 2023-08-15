package cli

import (
	"context"
	"errors"
	"fmt"
	"os"

	"RedWood011/client/apperrors"
	"RedWood011/client/entity"
	"RedWood011/client/service/secret"
	"RedWood011/client/service/user"

	"github.com/manifoldco/promptui"
)

const (
	welcome      = "Добро пожаловать! Что вы хотите сделать?"
	registration = "Создать учетную запись"
	auth         = "Войти в учетную запись"

	enterLogin = "Введите логин" // #nosec G101

	enterPassword = "Введите пароль" // #nosec G101

	mySecrets       = "Мои секреты"                   // #nosec G101
	oneSecretCreate = "Создать секрет"                // #nosec G101
	secretNameList  = "Получить список имен секретов" // #nosec G101
	oneSecretGet    = "Получить секрет"               // #nosec G101
	oneSecretDelete = "Удалить секрет"                // #nosec G101

	secretName = "Введите тип секрета" // #nosec G101

	secretCardSave        = "Сохранить кредитную карту" // #nosec G101
	secretTextSave        = "Сохранить текстовый файл"  // #nosec G101
	secretBinarySave      = "Сохранить бинарный файл"   // #nosec G101
	secretLogPasswordSave = "Сохранить логин и пароль"  // #nosec G101

	back = "Назад"
)

func displayCli(label string, items ...string) string {
	startPrompt := promptui.Select{
		Label: label,
		Items: items,
	}

	_, action, err := startPrompt.Run()

	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
		os.Exit(1)
	}

	return action
}

//nolint:funlen,gocognit,gocyclo,cyclop // test Function
func RunCli(ctx context.Context,
	userService *user.Service,
	secretService *secret.Service) {
	action := displayCli(welcome, registration, auth)

	for {
		switch action {
		case registration:
			user := labelLoginPassword()
			accessToken, _, err := userService.RegisterUser(ctx, user)
			if err != nil {
				fmt.Println("Произошла ошибка при регистрации.Пожалуйста попробуйте снова.")
				action = displayCli(welcome, registration, auth)
				continue
			}
			fmt.Println("Вы успешно зарегистрировались.")
			secretService.AccessToken = accessToken
			action = displayCli(mySecrets, oneSecretCreate, secretNameList, oneSecretGet, oneSecretDelete)

		case auth:
			user := labelLoginPassword()
			accessToken, _, err := userService.AuthUser(ctx, user)
			if err != nil {
				fmt.Println("Произошла ошибка при входе в учетную запись. Пожалуйста попробуйте снова.")
				action = displayCli(welcome, registration, auth)
				continue
			}
			fmt.Println("Вы успешно вошли в аккаунт.")
			secretService.AccessToken = accessToken
			action = displayCli(mySecrets, oneSecretCreate, secretNameList, oneSecretGet, oneSecretDelete)

		case oneSecretCreate:
			for action == oneSecretCreate && (action != back && action != registration && action != auth) {
				action = displayCli(secretName, secretCardSave, secretTextSave, secretBinarySave, secretLogPasswordSave, back)
				if action == back {
					action = displayCli(mySecrets, oneSecretCreate, secretNameList, oneSecretGet, oneSecretDelete)
				}
				var secretID string
				var err error
				switch action {
				case secretCardSave:
					secretID, err = createSecretCard(ctx, secretService)
					if err != nil {
						if errors.Is(err, apperrors.ErrAuth) {
							fmt.Println("Ошибка авторизации. Войдите в аккаунт и попробуйте снова.")
							action = displayCli(welcome, registration, auth)
							continue
						} else {
							fmt.Println("Произошла ошибка при создании секрета. Пожалуйста попробуйте снова.")
							action = displayCli(mySecrets, oneSecretCreate, secretNameList, oneSecretGet, oneSecretDelete)
							continue
						}
					}

				case secretTextSave:
					secretID, err = createSecretData(ctx, secretService, "text")
					if err != nil {
						if errors.Is(err, apperrors.ErrAuth) {
							fmt.Println("Ошибка авторизации. Войдите в аккаунт и попробуйте снова.")
							action = displayCli(welcome, registration, auth)
							continue
						} else {
							fmt.Println("Произошла ошибка при создании секрета. Пожалуйста попробуйте снова.")
							action = displayCli(mySecrets, oneSecretCreate, secretNameList, oneSecretGet, oneSecretDelete)
							continue
						}
					}

				case secretBinarySave:
					secretID, err = createSecretData(ctx, secretService, "bin")
					if err != nil {
						if errors.Is(err, apperrors.ErrAuth) {
							fmt.Println("Ошибка авторизации. Войдите в аккаунт и попробуйте снова.")
							action = displayCli(welcome, registration, auth)
							continue
						} else {
							fmt.Println("Произошла ошибка при создании секрета. Пожалуйста попробуйте снова.")
							action = displayCli(secretName, secretCardSave, secretTextSave, secretBinarySave, secretLogPasswordSave, back)
							continue
						}
					}

				case secretLogPasswordSave:
					secretID, err = createSecretLogPass(ctx, secretService)
					if err != nil {
						if errors.Is(err, apperrors.ErrAuth) {
							fmt.Println("Ошибка авторизации. Войдите в аккаунт и попробуйте снова.")
							action = displayCli(welcome, registration, auth)
							continue
						} else {
							fmt.Println("Произошла ошибка при создании секрета. Пожалуйста попробуйте снова.")
							action = displayCli(mySecrets, oneSecretCreate, secretNameList, oneSecretGet, oneSecretDelete)
							continue
						}
					}
				}
				fmt.Println("Ваш секрет создан. Секрет ID: ", secretID)
				action = displayCli(mySecrets, oneSecretCreate, secretNameList, oneSecretGet, oneSecretDelete)
			}

		case secretNameList:
			for action == secretNameList && (action != back) {
				secrets, err := secretService.ListSecrets(ctx)
				if err != nil {
					if errors.Is(err, apperrors.ErrAuth) {
						fmt.Println("Ошибка авторизации. Войдите в аккаунт и попробуйте снова.")
						action = displayCli(welcome, registration, auth)
						continue
					} else {
						fmt.Println("Произошла ошибка при запросе секретов. Пожалуйста попробуйте снова.")
						action = displayCli(mySecrets, oneSecretCreate, secretNameList, oneSecretGet, oneSecretDelete)
						continue
					}
				}

				fmt.Println("Список секретов:")
				for _, secret := range secrets {
					fmt.Printf("ID: %s  Имя секрета: %s \n", secret.ID, string(secret.Name))
				}
				action = displayCli("Меню секретов", back)
				if action == back {
					action = displayCli(mySecrets, oneSecretCreate, secretNameList, oneSecretGet, oneSecretDelete)
				}
			}

		case oneSecretGet:
			for action == oneSecretGet && action != back && action != registration && action != auth {
				secretID, err := getInput("Введите ID секрета", 0, "", 0)
				if err != nil {
					fmt.Println("Произошла ошибка при вводе ID секрета.Пожалуйста попробуйте снова.")
					action = displayCli(mySecrets, oneSecretCreate, secretNameList, oneSecretGet, oneSecretDelete)
					continue
				}
				secret, err := secretService.GetSecret(ctx, secretID)
				if err != nil {
					if errors.Is(err, apperrors.ErrAuth) {
						fmt.Println("Ошибка авторизации. Войдите в аккаунт и попробуйте снова.")
						action = displayCli(welcome, registration, auth)
						continue
					} else {
						fmt.Println("Произошла ошибка при получении секрета. Пожалуйста попробуйте снова.")
						action = displayCli(mySecrets, oneSecretCreate, secretNameList, oneSecretGet, oneSecretDelete)
						continue
					}
				}
				displaySecret(secret)
				action = displayCli(mySecrets, oneSecretCreate, secretNameList, oneSecretGet, oneSecretDelete)
			}

		case oneSecretDelete:
			fmt.Println("Удалить секрет:")
			secretID, err := getInput("Введите ID секрета", 0, "", 0)
			if err != nil {
				fmt.Println("Произошла ошибка при вводе ID секрета.Пожалуйста попробуйте снова.")
				action = displayCli(mySecrets, oneSecretCreate, secretNameList, oneSecretGet, oneSecretDelete)
				continue
			}
			err = secretService.DeleteSecret(ctx, secretID)
			if err != nil {
				if errors.Is(err, apperrors.ErrAuth) {
					fmt.Println("Ошибка авторизации. Войдите в аккаунт и попробуйте снова.")
					action = displayCli(welcome, registration, auth)
					continue
				} else {
					fmt.Println("Произошла ошибка при удалении секрета. Пожалуйста попробуйте снова.")
					action = displayCli(mySecrets, oneSecretCreate, secretNameList, oneSecretGet, oneSecretDelete)
					continue
				}
			}
			fmt.Println("Секрет удален")
			action = displayCli(mySecrets, oneSecretCreate, secretNameList, oneSecretGet, oneSecretDelete)
		}
	}
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

func labelLoginPassword() entity.User {
	loginPrompt := promptui.Prompt{
		Label: enterLogin,
	}
	login, err := loginPrompt.Run()
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
	}

	passwordPrompt := promptui.Prompt{
		Label: enterPassword,
		Mask:  '*',
	}
	password, err := passwordPrompt.Run()
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
	}
	return entity.User{Password: password, Login: login}
}

func createSecretCard(ctx context.Context, secretService *secret.Service) (string, error) {
	var secret entity.Secret
	const maskCVV = 42
	name, err := getInput("Введите название секрета", 0, "мойСекрет", 0)
	if err != nil {
		return "", err
	}
	secret.Name = []byte(name)

	secret.Card.Number, err = getInput("Введите номер карты", 0, "", 0)
	if err != nil {
		return "", err
	}
	secret.Card.CVV, err = getInput("Введите СVV", 0, "", maskCVV)
	if err != nil {
		return "", err
	}
	secret.Card.Date, err = getInput("Введите дату истечения срока действия карты", 0, "", 0)
	if err != nil {
		return "", err
	}
	secret.Card.Owner, err = getInput("Введите владельца карты", 0, "", 0)
	if err != nil {
		return "", err
	}
	secret.TypeSecret = "card"
	var secretID string
	secretID, err = secretService.CreateSecret(ctx, &secret)
	if err != nil {
		return "", err
	}
	return secretID, nil
}

func createSecretLogPass(ctx context.Context, secretService *secret.Service) (string, error) {
	var secret entity.Secret
	const lenMinName = 5
	name, err := getInput("Введите название секрета", lenMinName, "", 0)
	if err != nil {
		return "", err
	}
	secret.Name = []byte(name)

	secret.PassLogin.Login, err = getInput("Введите логин", 1, "", 0)
	if err != nil {
		return "", err
	}

	secret.PassLogin.Password, err = getInput("Введите пароль", 1, "", 0)
	if err != nil {
		return "", err
	}

	secret.TypeSecret = "logPas"
	var secretID string
	secretID, err = secretService.CreateSecret(ctx, &secret)
	if err != nil {
		return "", err
	}
	return secretID, nil
}

func createSecretData(ctx context.Context, secretService *secret.Service, typeSecret string) (string, error) {
	var secret entity.Secret
	name, err := getInput("Введите название секрета", 0, "", 0)
	if err != nil {
		return "", err
	}
	patch, err := getInput("Укажите абсолютный путь до файла", 0, "", 0)
	if err != nil {
		return "", err
	}
	secret.Name = []byte(name)
	var fileBytes []byte
	fileBytes, err = os.ReadFile(patch)
	if err != nil {
		fmt.Println("Ошибка при чтении файла")
	}
	secret.Data = fileBytes
	secret.TypeSecret = typeSecret

	var secretID string
	secretID, err = secretService.CreateSecret(ctx, &secret)
	if err != nil {
		return "", err
	}
	return secretID, nil
}

func displaySecret(secret *entity.Secret) {
	switch secret.TypeSecret {
	case "card":
		fmt.Println("Название секрета: ", string(secret.Name))
		fmt.Println("Номер карты: ", secret.Card.Number)
		fmt.Println("Срок действия карты: ", secret.Card.Date)
		fmt.Println("Владелец карты: ", secret.Card.Owner)
		fmt.Println("CVV: ", secret.Card.CVV)
	case "logPas":
		fmt.Println("Название секрета: ", string(secret.Name))
		fmt.Println("Логин: ", secret.PassLogin.Login)
		fmt.Println("Пароль: ", secret.PassLogin.Password)
	case "text", "bin":
		var patch string
		fmt.Println("Введите путь, куда хотите сохранить файл")
		_, err := fmt.Scan(&patch)
		if err != nil {
			fmt.Println("Ошибка ввода директории")
			return
		}
		err = os.WriteFile(patch, secret.Data, 0600)
		if err != nil {
			fmt.Println("Ошибка записи файла")
			return
		}
		fmt.Println("Файл сохранен")
	}
}
