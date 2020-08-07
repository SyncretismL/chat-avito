package main

import (
	"chat/internal/chat"
	"chat/internal/message"
	"chat/internal/user"
	"chat/pkg/logger"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type Handler struct {
	logger   logger.Logger
	repoUser user.Users
	repoMsg  message.Messages
	repoChat chat.Chats
	//	templates map[string]*template.Template
}

func newHandler(newLogger logger.Logger, repoUser user.Users, repoMsg message.Messages, repoChat chat.Chats,
) *Handler {
	return &Handler{
		logger:   newLogger,
		repoUser: repoUser,
		repoMsg:  repoMsg,
		repoChat: repoChat,
		//		templates: templates,
	}
}

func (h *Handler) Routers(r *chi.Mux) *chi.Mux {
	r.Use(middleware.Recoverer)

	r.Route("/", func(r chi.Router) {
		r.Post("/users/add", h.UserCreate)
		r.Route("/chats", func(r chi.Router) {
			r.Post("/get", h.GetChats)
			r.Post("/add", h.AddChat)
		})
		r.Route("/messages", func(r chi.Router) {
			r.Post("/get", h.GetMsg)
			r.Post("/add", h.AddMsg)
		})
	})

	return r
}

func (h *Handler) UserCreate(w http.ResponseWriter, r *http.Request) {
	var u *user.User

	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = u.ValidCreate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := h.repoUser.Create(u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	js, err := json.Marshal(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(js)
}

func (h *Handler) AddChat(w http.ResponseWriter, r *http.Request) {
	var c *chat.Chat

	err := json.NewDecoder(r.Body).Decode(&c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = c.ValidCreate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := h.repoChat.Create(c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	js, err := json.Marshal(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(js)
}

func (h *Handler) AddMsg(w http.ResponseWriter, r *http.Request) {
	var m *message.Message

	err := json.NewDecoder(r.Body).Decode(&m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = m.ValidCreate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := h.repoMsg.Create(m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	js, err := json.Marshal(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(js)
}

func (h *Handler) GetChats(w http.ResponseWriter, r *http.Request) {
	var u user.User

	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = u.ValidGet(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	chs, err := h.repoChat.Find(u.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	js, err := json.Marshal(chs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(js)
}

func (h *Handler) GetMsg(w http.ResponseWriter, r *http.Request) {
	var ms message.Message

	err := json.NewDecoder(r.Body).Decode(&ms)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = ms.ValidGet(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	msgs, err := h.repoMsg.Find(ms.Chat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	js, err := json.Marshal(msgs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(js)
}
