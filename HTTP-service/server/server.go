package server

import (
	"serverdb/server/service"

	"github.com/go-chi/chi/v5"
)

func NewMux(service *service.Service) *chi.Mux {
	mux := chi.NewMux()
	mux.HandleFunc("/create", service.Create)
	mux.HandleFunc("/make_friends", service.MakeFriends)
	mux.HandleFunc("/delete/{user_id}", service.UserDelete)
	mux.HandleFunc("/get/{user_id}", service.GetFriends)
	mux.HandleFunc("/update/{user_id}", service.UserUpdate)

	return mux
}
