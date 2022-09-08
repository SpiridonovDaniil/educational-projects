package server

import (
	"bytes"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"serverdb/models"
	"serverdb/server/service"
	"serverdb/server/service/mocks"
	"testing"
)

func TestService_Create(t *testing.T) {

	type mockBehavior func(s *mock_service.Mockrepository, user models.User, expectedError error)

	testTable := []struct {
		name                string
		inputBody    string
		inputUser    models.User
		mockBehavior mockBehavior
		expectedError       error
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:      "OK",
			inputBody: `{"name":"test", "age": 5}`,
			inputUser: models.User{
				Name: "test",
				Age:  5,
			},
			mockBehavior: func(s *mock_service.Mockrepository, user models.User, expectedError error) {
				s.EXPECT().Insert(user).Return(expectedError)
			},
			expectedStatusCode:  201,
			expectedRequestBody: "{\"message\":\"User was created test\",\"id\":\"test\"}",
		},
		{
			name:      "nil",
			inputBody: `{"name":"", "age": 5}`,
			inputUser: models.User{
				Name: "",
				Age:  5,
			},
			mockBehavior: func(s *mock_service.Mockrepository, user models.User, expectedError error) {
			},
			expectedStatusCode:  400,
			expectedRequestBody: "{\"error\":\"User name not specified\"}",
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			create := mock_service.NewMockrepository(c)
			testCase.mockBehavior(create, testCase.inputUser, testCase.expectedError)

			handler := service.NewService(create)

			mux := NewMux(handler)
			server := httptest.NewServer(mux)
			defer server.Close()

			req, err := http.NewRequest("POST", server.URL+"/create", bytes.NewBufferString(testCase.inputBody))
			assert.NoError(t, err)
			resp, err := server.Client().Do(req)
			assert.NoError(t, err)
			assert.Equal(t, testCase.expectedStatusCode, resp.StatusCode)

			body, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)
			assert.Equal(t, testCase.expectedRequestBody, string(body))
		})
	}
}

func TestService_MakeFriends(t *testing.T) {

	type mockBehavior func(s *mock_service.Mockrepository, findUser1 models.User, findUser2 models.User, inputUsers models.MakeFriends, expectedError error)

	testTable := []struct {
		name                string
		inputBody    string
		findUser1    models.User
		findUser2    models.User
		inputUsers   models.MakeFriends
		mockBehavior mockBehavior
		expectedError       error
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:      "OK",
			inputBody: `{"SourceId":"User1", "TargetId":"User2"}`,
			findUser1:	models.User{
				Name: "User1",
				Age:  7,
			},
			findUser2:	models.User{
				Name: "User2",
				Age:  5,
			},
			inputUsers: models.MakeFriends{
				SourceId: "User1",
				TargetId: "User2",
			},
			mockBehavior: func(s *mock_service.Mockrepository, findUser1 models.User, findUser2 models.User, inputUsers models.MakeFriends, expectedError error) {
				s.EXPECT().FindUser(inputUsers.TargetId).Return(findUser2, expectedError)
				s.EXPECT().FindUser(inputUsers.SourceId).Return(findUser1, expectedError)
				s.EXPECT().Update(inputUsers.TargetId, "$addToSet", "friends", inputUsers.SourceId).Return(expectedError)
				s.EXPECT().Update(inputUsers.SourceId, "$addToSet", "friends", inputUsers.TargetId).Return(expectedError)
			},
			expectedStatusCode:  200,
			expectedRequestBody: "{\"message\":\"User1 и User2 теперь друзья\"}",
		},
		{
			name:      "User not found",
			inputBody: `{"SourceId":"User1", "TargetId":"User2"}`,
			findUser1:	models.User{
				Name: "User1",
				Age:  7,
			},
			findUser2:	models.User{
				Name: "User2",
				Age:  5,
			},
			inputUsers: models.MakeFriends{
				SourceId: "User1",
				TargetId: "User2",
			},
			mockBehavior: func(s *mock_service.Mockrepository, findUser1 models.User, findUser2 models.User, inputUsers models.MakeFriends, expectedError error) {
				s.EXPECT().FindUser(inputUsers.TargetId).Return(findUser2, errors.New(""))
			},
			expectedStatusCode:  400,
			expectedRequestBody: "{\"error\":\"Пользователь с именем User2 не найден\"}",
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			makeFriends := mock_service.NewMockrepository(c)
			testCase.mockBehavior(makeFriends, testCase.findUser1, testCase.findUser2, testCase.inputUsers, testCase.expectedError)

			handler := service.NewService(makeFriends)

			mux := NewMux(handler)
			server := httptest.NewServer(mux)
			defer server.Close()

			req, err := http.NewRequest("POST", server.URL+"/make_friends", bytes.NewBufferString(testCase.inputBody))
			assert.NoError(t, err)
			resp, err := server.Client().Do(req)
			assert.NoError(t, err)
			assert.Equal(t, testCase.expectedStatusCode, resp.StatusCode)

			body, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)
			assert.Equal(t, testCase.expectedRequestBody, string(body))
		})
	}
}

func TestService_UserDelete(t *testing.T) {
	type mockBehavior func(s *mock_service.Mockrepository, user models.User, expectedError error)

	testTable := []struct {
		name                string
		inputBody    string
		inputUser    models.User
		mockBehavior mockBehavior
		expectedError       error
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:      "OK",
			inputBody: ``,
			inputUser: models.User{
				Name: "test",
				Age:  5,
				Friends: []string{"test2"},
			},
			mockBehavior: func(s *mock_service.Mockrepository, user models.User, expectedError error) {
				s.EXPECT().FindUser("test").Return(user, expectedError)
				s.EXPECT().Update(user.Friends[0], "$pull", "friends", "test").Return(expectedError)
				s.EXPECT().RemoveAll("test").Return(expectedError)
			},
			expectedStatusCode:  200,
			expectedRequestBody: "{\"message\":\"Аккаунт test был удалён\"}",
		},
		{
			name:      "there is no such user",
			inputBody: ``,
			inputUser: models.User{
				Name: "test",
				Age:  5,
				Friends: []string{"test2"},
			},
			mockBehavior: func(s *mock_service.Mockrepository, user models.User, expectedError error) {
				s.EXPECT().FindUser("test").Return(user, errors.New(""))
			},
			expectedStatusCode:  400,
			expectedRequestBody: "{\"error\":\"Пользователь с именем test не найден\"}",
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			create := mock_service.NewMockrepository(c)
			testCase.mockBehavior(create, testCase.inputUser, testCase.expectedError)

			handler := service.NewService(create)

			mux := NewMux(handler)
			server := httptest.NewServer(mux)
			defer server.Close()

			req, err := http.NewRequest("DELETE", server.URL+"/delete/test", bytes.NewBufferString(testCase.inputBody))
			assert.NoError(t, err)
			resp, err := server.Client().Do(req)
			assert.NoError(t, err)
			assert.Equal(t, testCase.expectedStatusCode, resp.StatusCode)

			body, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)
			assert.Equal(t, testCase.expectedRequestBody, string(body))
		})
	}
}

func TestService_GetFriends(t *testing.T) {
	type mockBehavior func(s *mock_service.Mockrepository, user models.User, expectedError error)

	testTable := []struct {
		name                string
		inputBody    string
		inputUser    models.User
		mockBehavior mockBehavior
		expectedError       error
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:      "OK",
			inputBody: ``,
			inputUser: models.User{
				Name: "test",
				Age:  5,
				Friends: []string{"test2", "test3"},
			},
			mockBehavior: func(s *mock_service.Mockrepository, user models.User, expectedError error) {
				s.EXPECT().FindUser("test").Return(user, expectedError)

			},
			expectedStatusCode:  200,
			expectedRequestBody: "{\"message\":\"test2, test3\",\"id\":\"test\"}",
		},
		{
			name:      "there is no such user",
			inputBody: ``,
			inputUser: models.User{
				Name: "test",
				Age:  5,
				Friends: []string{"test2"},
			},
			mockBehavior: func(s *mock_service.Mockrepository, user models.User, expectedError error) {
				s.EXPECT().FindUser("test").Return(user, errors.New(""))
			},
			expectedStatusCode:  400,
			expectedRequestBody: "{\"error\":\"такого пользователя не существует\"}",
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			create := mock_service.NewMockrepository(c)
			testCase.mockBehavior(create, testCase.inputUser, testCase.expectedError)

			handler := service.NewService(create)

			mux := NewMux(handler)
			server := httptest.NewServer(mux)
			defer server.Close()

			req, err := http.NewRequest("GET", server.URL+"/get/test", bytes.NewBufferString(testCase.inputBody))
			assert.NoError(t, err)
			resp, err := server.Client().Do(req)
			assert.NoError(t, err)
			assert.Equal(t, testCase.expectedStatusCode, resp.StatusCode)

			body, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)
			assert.Equal(t, testCase.expectedRequestBody, string(body))
		})
	}
}

func TestService_UserUpdate(t *testing.T) {
	type mockBehavior func(s *mock_service.Mockrepository, user models.UpdateUser, expectedError error)

	testTable := []struct {
		name                string
		inputBody    string
		inputUser    models.UpdateUser
		mockBehavior mockBehavior
		expectedError       error
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:      "OK",
			inputBody: `{"age":10}`,
			inputUser: models.UpdateUser{
				NewAge: 10,
				NewName: "test",
			},
			mockBehavior: func(s *mock_service.Mockrepository, user models.UpdateUser, expectedError error) {
				s.EXPECT().Update("test", "$set", "age", user.NewAge).Return(expectedError)
			},
			expectedStatusCode:  200,
			expectedRequestBody: "{\"message\":\"Возраст пользователя успешно обновлен\",\"id\":\"test\"}",
		},
		{
			name:      "there is no such user",
			inputBody: `{"age":10}`,
			inputUser: models.UpdateUser{
				NewAge: 10,
				NewName: "test",
			},
			mockBehavior: func(s *mock_service.Mockrepository, user models.UpdateUser, expectedError error) {
				s.EXPECT().Update("test", "$set", "age", user.NewAge).Return(errors.New(""))
			},
			expectedStatusCode:  500,
			expectedRequestBody: "{\"error\":\"Ошибка на стороне приложения. Не можем найти пользователя - test\"}",
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			create := mock_service.NewMockrepository(c)
			testCase.mockBehavior(create, testCase.inputUser, testCase.expectedError)

			handler := service.NewService(create)

			mux := NewMux(handler)
			server := httptest.NewServer(mux)
			defer server.Close()

			req, err := http.NewRequest("PUT", server.URL+"/update/test", bytes.NewBufferString(testCase.inputBody))
			assert.NoError(t, err)
			resp, err := server.Client().Do(req)
			assert.NoError(t, err)
			assert.Equal(t, testCase.expectedStatusCode, resp.StatusCode)

			body, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)
			assert.Equal(t, testCase.expectedRequestBody, string(body))
		})
	}
}