package controller

import (
	"context"
	"fmt"
	"testing"

	"gitlab.com/Fratbe/addglee/src/server/i18n"
	"gitlab.com/Fratbe/addglee/src/server/models"
)

func (testDB *TestDB) CreateUserAccount(userAccount *models.UserAccount) error {
	if userAccount.Pseudo == "make this failed" {
		return fmt.Errorf("error")
	}

	key := userAccount.Pseudo + "/" + userAccount.Password
	(*testDB)[key] = userAccount
	return nil
}

func (testDB *TestDB) GetUserAccount(condition, val string) *models.UserAccount {
	if _, ok := (*testDB)[val]; !ok {
		return nil
	}
	return (*testDB)[val].(*models.UserAccount)
}

func (testDB *TestDB) SetUserAccountFieldData(pseudonym, field string, val interface{}) error {
	if val.(string) == "make this failed" {
		return fmt.Errorf("error")
	}

	userAccount := (*testDB)[pseudonym].(*models.UserAccount)
	switch field {
	case "pseudo":
		(*userAccount).Pseudo = val.(string)
	case "password":
		(*userAccount).Password = val.(string)
	}
	return nil
}

func TestSignUpAccount(t *testing.T) {
	h := &Handler{mockDB}

	fmt.Println("------- TestSignUpUserAccount Success Case -------")
	testSuccess := []struct {
		pseudo      string
		lastName    string
		firstName   string
		password    string
		languageId  uint32
		activated   bool
		fakeAccount bool
		want        string
	}{
		{
			pseudo:      "test1",
			lastName:    "testtest1",
			firstName:   "testtesttest1",
			password:    "1234",
			languageId:  1,
			activated:   true,
			fakeAccount: true,
			want:        i18n.T("account_created"),
		},
		{
			pseudo:      "test2",
			lastName:    "testtest2",
			firstName:   "testtesttest2",
			password:    "5678",
			languageId:  2,
			activated:   false,
			fakeAccount: false,
			want:        i18n.T("account_created"),
		},
	}
	for _, tt := range testSuccess {
		req := &models.UserAccount{
			Pseudo:      tt.pseudo,
			LastName:    tt.lastName,
			FirstName:   tt.firstName,
			Password:    tt.password,
			LanguageId:  tt.languageId,
			Activated:   tt.activated,
			FakeAccount: tt.fakeAccount,
		}
		resp, err := h.SignUpAccount(context.Background(), req)
		if err != nil {
			t.Errorf("SignUpAccount got unexpected error: %s", err)
		}
		if resp.Message != tt.want {
			t.Errorf("SignUpAccount(%v)=\"%v\", wanted \"%v\"", tt, resp.Message, tt.want)
		}
	}

	fmt.Println("------- TestSignUpUserAccount Failure Case -------")
	testFailure := []struct {
		pseudo      string
		lastName    string
		firstName   string
		password    string
		languageId  uint32
		activated   bool
		fakeAccount bool
	}{
		{
			pseudo:      "make this failed",
			lastName:    "testtest2",
			firstName:   "testtesttest2",
			password:    "5678",
			languageId:  2,
			activated:   false,
			fakeAccount: false,
		},
	}
	for _, tt := range testFailure {
		req := &models.UserAccount{
			Pseudo:      tt.pseudo,
			LastName:    tt.lastName,
			FirstName:   tt.firstName,
			Password:    tt.password,
			LanguageId:  tt.languageId,
			Activated:   tt.activated,
			FakeAccount: tt.fakeAccount,
		}
		_, err := h.SignUpAccount(context.Background(), req)
		if err == nil {
			t.Errorf("SignUpAccount got nil error")
		}
	}
}

func TestUpdateAccount(t *testing.T) {
	h := &Handler{mockDB}

	fmt.Println("------- TestUpdateAccount Success Case -------")
	testSuccess := []struct {
		pseudo   string
		password string
		field    string
		value    string
		want     string
	}{
		{
			pseudo:   "king",
			password: "1234",
			field:    "first_name",
			value:    "Jon2",
			want:     i18n.T("account_updated"),
		},
	}
	for _, tt := range testSuccess {
		req := &SelfCareInformation{Pseudo: tt.pseudo, Password: tt.password, Field: tt.field, Value: tt.value}
		resp, err := h.UpdateAccount(context.Background(), req)
		if err != nil {
			t.Errorf("UpdateAccount got unexpected error: %s", err)
		}
		if resp.Message != tt.want {
			t.Errorf("UpdateAccount(%v)=\"%v\", wanted \"%v\"", req, resp.Message, tt.want)
		}
	}

	fmt.Println("------- TestUpdateAccount Failure Case -------")
	testFailure := []struct {
		pseudo   string
		password string
		field    string
		value    string
	}{
		{
			pseudo:   "queen",
			password: "1234",
			field:    "first_name",
			value:    "Jon2",
		},
	}
	for _, tt := range testFailure {
		req := &SelfCareInformation{Pseudo: tt.pseudo, Password: tt.password, Field: tt.field, Value: tt.value}
		_, err := h.UpdateAccount(context.Background(), req)
		if err != nil {
			t.Errorf("UpdateAccount got unexpected error: %s", err)
		}
	}
}

func TestLoginAccount(t *testing.T) {
	h := &Handler{mockDB}

	fmt.Println("------- TestLoginAccount Success Case -------")
	testSuccess := []struct {
		pseudo   string
		password string
	}{
		{
			pseudo:   "king",
			password: "1234",
		},
	}
	for _, tt := range testSuccess {
		req := &LoginInformation{Pseudo: tt.pseudo, Password: tt.password}
		_, err := h.LoginAccount(context.Background(), req)
		if err != nil {
			t.Errorf("LoginAccount got unexpected error: %s", err)
		}
	}

	fmt.Println("------- TestLoginAccount Failure Case -------")
	testFailure := []struct {
		pseudo   string
		password string
	}{
		{
			pseudo:   "king",
			password: "5678",
		},
	}
	for _, tt := range testFailure {
		req := &LoginInformation{Pseudo: tt.pseudo, Password: tt.password}
		_, err := h.LoginAccount(context.Background(), req)
		if err == nil {
			t.Errorf("LoginAccount got nil error")
		}
	}

}
