package user

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"server/internal/handlers"
)

var _ handlers.Handler = &handler{}


type handler struct {
	Service
}

func NewHandler(service *Service) handlers.Handler {
	return &handler{
		*service,
	}
}

func (h *handler) Register(router *chi.Mux) {

	router.MethodFunc(http.MethodGet, "/users", handlers.Middleware(h.GetList))
	router.MethodFunc(http.MethodGet, "/friends/{articleID}", handlers.Middleware(h.GetFriendsByID))
	router.MethodFunc(http.MethodPost, "/create", handlers.Middleware(h.CreateUser))
	router.MethodFunc(http.MethodPut, "/{articleID}", handlers.Middleware(h.UpdateUser))
	router.MethodFunc(http.MethodPost, "/make_friends", handlers.Middleware(h.MakeFriends))
	router.MethodFunc(http.MethodDelete, "/delete", handlers.Middleware(h.DeleteUser))
}

func (h *handler) GetList(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		response, err := h.storage.GetAll()
		if err != nil {
			return err
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
		return nil
	}
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte("Пользователи не найдены"))
	return nil
}

func (h *handler) GetFriendsByID(w http.ResponseWriter, r *http.Request) error {
	// жду реализацию этой ручки
	id := strings.Replace(r.URL.Path, "/friends/", "", 1)
	response, err := h.storage.Get(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(id + " не является ID"))
		return err
	}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
		return nil
}

func (h *handler) CreateUser(w http.ResponseWriter, r *http.Request) error {
	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return err
	}
	defer r.Body.Close()

	var u User

	if err := json.Unmarshal(content, &u); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return err
	}
	response, err := h.storage.Create(&u)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return err
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(response))
	return nil
}

func (h *handler) UpdateUser(w http.ResponseWriter, r *http.Request) error {
	id := r.URL.Path
	id = id[1:]

	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return err
	}
	defer r.Body.Close()

	var u UpdateUser

	if err := json.Unmarshal(content, &u); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return err
	}

	response, err := h.storage.Update(id, u.New_age)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	//отправка ответа
	w.WriteHeader(http.StatusCreated)
		w.Write([]byte(response))
	return nil
}

func (h *handler) MakeFriends(w http.ResponseWriter, r *http.Request) error {
	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return err
	}
	defer r.Body.Close()

	var d ID

	if err := json.Unmarshal(content, &d); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return err
	}
	response, err := h.storage.MakeFriends(&d)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
	return nil
}

func (h *handler) DeleteUser(w http.ResponseWriter, r *http.Request) error {
	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return err
	}
	defer r.Body.Close()

	var u User

	if err := json.Unmarshal(content, &u); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return err
	}
	response, err := h.storage.Delete(&u)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
	return nil
}