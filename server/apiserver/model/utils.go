package model

import (
	"encoding/hex"

	"github.com/asif10388/synctab/apiserver/controller"
	"golang.org/x/crypto/bcrypt"
)

func (model *Model) isValidField(field string, minLength, maxLength int) bool {
	if len(field) < minLength || len(field) > maxLength {
		return false
	}

	return true
}

func (model *Model) GetAdaptiveHashString(str string) (string, error) {
	if len(str) == 0 {
		return "", controller.ErrInvalidInput
	}

	hashString, err := bcrypt.GenerateFromPassword([]byte(str), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(hashString), nil
}

func (model *Model) GetPasswordHash(password string) (string, error) {
	if !model.isValidField(password, model.MinPasswordLen, model.MaxPasswordLen) {
		return "", controller.ErrInvalidPassword
	}

	passwordHash, err := model.GetAdaptiveHashString(password)
	if err != nil {
		return "", err
	}

	return passwordHash, nil
}

func (model *Model) CompareAdaptiveHashString(hashedPassword, plaintextPassword string) bool {
	if len(hashedPassword) == 0 || len(plaintextPassword) == 0 {
		return false
	}

	decodedHex, err := hex.DecodeString(hashedPassword)
	if err != nil {
		return false
	}

	err = bcrypt.CompareHashAndPassword(decodedHex, []byte(plaintextPassword))
	return err == nil
}
