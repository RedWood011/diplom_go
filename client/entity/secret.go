package entity

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/gob"
	"strings"
)

const (
	TEXT    = "text|"
	BIN     = "bin|"
	LOGPASS = "logPas|"
	CARD    = "card|"
)

type Secret struct {
	ID         string
	UserID     string
	Data       []byte
	Name       []byte
	TypeSecret string
	Card       Card
	TextFile   []byte
	BinFile    []byte
	PassLogin  PasswordLogin
}

type Card struct {
	Number string
	Date   string
	Owner  string
	CVV    string
}

type PasswordLogin struct {
	Password string
	Login    string
}

func (s *Secret) EncryptSecret(key, nonce []byte) error {
	var buf bytes.Buffer
	switch s.TypeSecret {
	case "text":
		s.Data = s.TextFile
		name := string(s.Name)
		name = TEXT + name
		s.Name = []byte(name)

	case "bin":
		s.Data = s.BinFile
		name := string(s.Name)
		name = BIN + name
		s.Name = []byte(name)

	case "logPas":
		enc := gob.NewEncoder(&buf)
		err := enc.Encode(s.PassLogin)
		if err != nil {
			return err
		}
		s.Data = buf.Bytes()
		name := string(s.Name)
		name = LOGPASS + name
		s.Name = []byte(name)

	case "card":
		enc := gob.NewEncoder(&buf)
		err := enc.Encode(s.Card)
		if err != nil {
			return err
		}
		s.Data = buf.Bytes()
		name := string(s.Name)
		name = CARD + name
		s.Name = []byte(name)
	}

	aesblock, err := aes.NewCipher(key)
	if err != nil {
		return err
	}
	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		return err
	}

	encryptData := aesgcm.Seal(nil, nonce, s.Data, nil)
	encryptName := aesgcm.Seal(nil, nonce, s.Name, nil)

	s.Data = encryptData
	s.Name = encryptName

	return nil
}

func (s *Secret) DecryptSecret(key, nonce []byte) error {
	aesblock, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		return err
	}

	decryptData, err := aesgcm.Open(nil, nonce, s.Data, nil)
	if err != nil {
		return err
	}

	s.Data = decryptData

	decryptName, err := aesgcm.Open(nil, nonce, s.Name, nil)
	if err != nil {
		return err
	}
	var typeSecret string
	index := strings.Index(string(decryptName), "|")
	if index != -1 {
		typeSecret = string(decryptName)[:index]
		name := string(decryptName)[index+1:]
		s.Name = []byte(name)
	}

	switch typeSecret {
	case "text":
		s.TextFile = s.Data

	case "bin":
		s.BinFile = s.Data

	case "logPas":
		buf := bytes.NewBuffer(s.Data)
		dec := gob.NewDecoder(buf)
		err = dec.Decode(&s.PassLogin)
		if err != nil {
			return err
		}

	case "card":
		buf := bytes.NewBuffer(s.Data)
		dec := gob.NewDecoder(buf)
		err = dec.Decode(&s.Card)
		if err != nil {
			return err
		}
	}

	return nil
}
