package service

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"io"
	"log"
	"net/http"
	"serverdb/models"
	"strings"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type repository interface {
	FindUser(userID string) (models.User, error)
	Insert(user models.User) error
	Update(user, order, object string, objectChange interface{}) error
	RemoveAll(userID string) error
}

type Service struct {
	repo repository
}

func NewService(repo repository) *Service {
	return &Service{repo: repo}
}

func march(s models.Answer)[]byte {
	answer, err := json.Marshal(s)
	if err != nil {
		log.Println(err)
	}
	return answer
}

func (s *Service) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		content, err := io.ReadAll(r.Body)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			_, err := w.Write(march(models.Answer{Error:err.Error()}))
			if err != nil {
				fmt.Println(err)
			}
			return
		}
		defer r.Body.Close()

		var u models.User
		if err := json.Unmarshal(content, &u); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			_, err := w.Write(march(models.Answer{Error:err.Error()}))
			if err != nil {
				fmt.Println(err)
			}
			return
		}

		if u.Name == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_, err = w.Write(march(models.Answer{Error:"User name not specified"}))
			if err != nil {
				fmt.Println(err)
			}
			return
		}

		err = s.repo.Insert(u)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, err = w.Write(march(models.Answer{Message: "User was created " + u.Name, Id: u.Name}))
		if err != nil {
			fmt.Println(err)
		}
		return
	}
	w.WriteHeader(http.StatusBadRequest)
}

func (s *Service) MakeFriends(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		contentFriends, err := io.ReadAll(r.Body)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			_, err := w.Write(march(models.Answer{Error:err.Error()}))
			if err != nil {
				fmt.Println(err)
			}
			return
		}
		defer r.Body.Close()

		var f models.MakeFriends
		if err := json.Unmarshal(contentFriends, &f); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			_, err := w.Write(march(models.Answer{Error:err.Error()}))
			if err != nil {
				fmt.Println(err)
			}
			return
		}

		_, errTarget := s.repo.FindUser(f.TargetId)
		if errTarget != nil {
			fmt.Println(errTarget)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_, err = w.Write(march(models.Answer{Error:"Пользователь с именем " + f.TargetId + " не найден"}))
			if err != nil {
				fmt.Println(err)
			}
			return
		}

		_, errSource := s.repo.FindUser(f.SourceId)
		if errSource != nil {
			fmt.Println(errSource)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_, err = w.Write(march(models.Answer{Error:"Пользователь с именем " + f.SourceId + " не найден"}))
			if err != nil {
				fmt.Println(err)
			}
			return
		}

		if errTarget == nil && errSource == nil {
			err = s.repo.Update(f.TargetId, "$addToSet", "friends", f.SourceId)
			if err != nil {
				fmt.Println(err)
			}

			err = s.repo.Update(f.SourceId, "$addToSet", "friends", f.TargetId)
			if err != nil {
				fmt.Println(err)
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, err = w.Write(march(models.Answer{Message:f.SourceId + " и " + f.TargetId + " теперь друзья"}))
			if err != nil {
				fmt.Println(err)
			}
			return
		}
	}
	w.WriteHeader(http.StatusBadRequest)
}

func (s *Service) UserDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method == "DELETE" {
		userID := chi.URLParam(r, "user_id")

		userDelete, err := s.repo.FindUser(userID)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write(march(models.Answer{Error:"Пользователь с именем " + userID + " не найден"}))
			return
		}

		for _, friendName := range userDelete.Friends {
			err := s.repo.Update(friendName, "$pull", "friends", userID)
			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				w.Write(march(models.Answer{Error:"Пользователь с именем " + friendName + " не найден"}))
				return
			}
		}

		err = s.repo.RemoveAll(userID)
		if err != nil {
			fmt.Println(err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(march(models.Answer{Error:"Ошибка на стороне приложения. Не удалось удалить пользователя - " + userID}))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(march(models.Answer{Message:"Аккаунт " + userID + " был удалён"}))
		if err != nil {
			fmt.Println(err)
		}
		return
	}
	w.WriteHeader(http.StatusBadRequest)
}

func (s *Service) GetFriends(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		userID := chi.URLParam(r, "user_id")
		response := ""
		friendsUser, err := s.repo.FindUser(userID)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write(march(models.Answer{Error:"такого пользователя не существует"}))
			return
		}

		for _, friend := range friendsUser.Friends {
			response += friend + ", "
		}
		response = strings.Trim(response, ", ")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(march(models.Answer{Message:response, Id: userID}))
		if err != nil {
			fmt.Println(err)
		}
		return
	}
	w.WriteHeader(http.StatusBadRequest)
}

func (s *Service) UserUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method == "PUT" {
		userID := chi.URLParam(r, "user_id")
		contentUpdate, err := io.ReadAll(r.Body)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			_, err := w.Write(march(models.Answer{Error:err.Error()}))
			if err != nil {
				fmt.Println(err)
			}
			return
		}
		defer r.Body.Close()

		var u models.UpdateUser
		if err := json.Unmarshal(contentUpdate, &u); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			_, err := w.Write(march(models.Answer{Error:err.Error()}))
			if err != nil {
				fmt.Println(err)
			}
			return
		}
		err = s.repo.Update(userID, "$set", "age", u.NewAge)
		if err != nil {
			fmt.Println(err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(march(models.Answer{Error:"Ошибка на стороне приложения. Не можем найти пользователя - " + userID}))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(march(models.Answer{Message:"Возраст пользователя успешно обновлен", Id: userID}))
		if err != nil {
			fmt.Println(err)
		}
		return
	}
	w.WriteHeader(http.StatusBadRequest)
}
