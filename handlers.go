// handlers.go
package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Gym API
// @version 1.0
// @description API для управления фитнес-клубом: бронирование занятий, подписки, тренеры, админка.
// @contact.name API Support
// @contact.email support@gym.example.com
// @license.name MIT
// @host localhost:8080
// @BasePath /api
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Введите JWT токен в формате: Bearer <ваш_токен>

// NewRouter builds router with all endpoints
func NewRouter(store *Store) *mux.Router {
	r := mux.NewRouter()
	h := &Handler{Store: store}

	// Swagger UI
	r.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)

	// public
	r.HandleFunc("/api/health", h.Health).Methods("GET")
	r.HandleFunc("/api/users/register", h.Register).Methods("POST")
	r.HandleFunc("/api/users/login", h.Login).Methods("POST")

	// open endpoints
	r.HandleFunc("/api/classes", h.GetClasses).Methods("GET")
	r.HandleFunc("/api/memberships", h.GetMemberships).Methods("GET")

	// auth required
	api := r.PathPrefix("/api").Subrouter()
	api.Use(AuthMiddleware)
	api.HandleFunc("/me", h.Me).Methods("GET")
	api.HandleFunc("/bookings", h.CreateBooking).Methods("POST")
	api.HandleFunc("/bookings", h.ListBookings).Methods("GET")
	api.HandleFunc("/memberships/buy", h.BuyMembership).Methods("POST")
	api.HandleFunc("/payments", h.Pay).Methods("POST")

	// admin
	admin := r.PathPrefix("/api/admin").Subrouter()
	admin.Use(AuthMiddleware)
	admin.Use(AdminOnly)
	admin.HandleFunc("/trainers", h.CreateTrainer).Methods("POST")
	admin.HandleFunc("/classes", h.CreateClass).Methods("POST")

	return r
}

type Handler struct {
	Store *Store
}

// Helper: write JSON
func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// Health godoc
// @Summary Проверка здоровья сервера
// @Description Возвращает статус OK
// @Tags health
// @Success 200 {object} map[string]string
// @Router /api//health [get]
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

type registerReq struct {
	Name     string `json:"name" example:"Иван Иванов"`
	Email    string `json:"email" example:"ivan@example.com"`
	Password string `json:"password" example:"password123"`
}

// Register godoc
// @Summary Регистрация пользователя
// @Description Создаёт нового пользователя
// @Tags auth
// @Accept json
// @Produce json
// @Param body body registerReq true "Данные пользователя"
// @Success 201 {object} User
// @Failure 400 {object} map[string]string
// @Router /api/users/register [post]
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid"})
		return
	}
	if req.Email == "" || req.Password == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "email and password required"})
		return
	}
	pwHash, err := HashPassword(req.Password)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "server"})
		return
	}
	res, err := h.Store.DB.Exec("INSERT INTO users(name,email,password_hash) VALUES(?,?,?)", req.Name, req.Email, pwHash)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "could not create user"})
		return
	}
	id, _ := res.LastInsertId()
	u := User{ID: int(id), Name: req.Name, Email: req.Email}
	writeJSON(w, http.StatusCreated, u)
}

type loginReq struct {
	Email    string `json:"email" example:"ivan@example.com"`
	Password string `json:"password" example:"password123"`
}

// Login godoc
// @Summary Авторизация
// @Description Возвращает JWT токен
// @Tags auth
// @Accept json
// @Produce json
// @Param body body loginReq true "Логин и пароль"
// @Success 200 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/users/login [post]
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid"})
		return
	}
	var u User
	var id int
	var pwHash string
	var isAdmin int
	row := h.Store.DB.QueryRow("SELECT id,name,email,password_hash,is_admin,created_at FROM users WHERE email = ?", req.Email)
	var createdAt string
	if err := row.Scan(&id, &u.Name, &u.Email, &pwHash, &isAdmin, &createdAt); err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid credentials"})
		return
	}
	u.ID = id
	u.IsAdmin = isAdmin == 1
	if !CheckPasswordHash(req.Password, pwHash) {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid credentials"})
		return
	}
	token, err := CreateToken(u.ID, u.Email, u.IsAdmin)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "token error"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"token": token})
}

// Me godoc
// @Summary Получить текущего пользователя
// @Description Возвращает данные авторизованного пользователя
// @Tags user
// @Security Bearer
// @Success 200 {object} User
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/me [get]
func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	uid := r.Context().Value(ContextUserID).(int)
	var u User
	var createdAt string
	var isAdmin int
	row := h.Store.DB.QueryRow("SELECT id,name,email,is_admin,created_at FROM users WHERE id = ?", uid)
	if err := row.Scan(&u.ID, &u.Name, &u.Email, &isAdmin, &createdAt); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "user not found"})
		return
	}
	u.IsAdmin = isAdmin == 1
	writeJSON(w, http.StatusOK, u)
}

// GetClasses godoc
// @Summary Список занятий
// @Description Возвращает все доступные занятия
// @Tags classes
// @Success 200 {array} Class
// @Failure 500 {object} map[string]string
// @Router /api/classes [get]
func (h *Handler) GetClasses(w http.ResponseWriter, r *http.Request) {
	rows, err := h.Store.DB.Query(`SELECT id,title,description,trainer_id,start_time,duration_min,capacity,created_at FROM classes ORDER BY start_time ASC`)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "db error"})
		return
	}
	defer rows.Close()
	var out []Class
	for rows.Next() {
		var c Class
		var start string
		if err := rows.Scan(&c.ID, &c.Title, &c.Description, &c.TrainerID, &start, &c.DurationMin, &c.Capacity, &c.CreatedAt); err != nil {
			continue
		}
		c.StartTime, _ = time.Parse(time.RFC3339, start)
		out = append(out, c)
	}
	writeJSON(w, http.StatusOK, out)
}

// GetMemberships godoc
// @Summary Список подписок
// @Description Возвращает все доступные подписки
// @Tags memberships
// @Success 200 {array} Membership
// @Failure 500 {object} map[string]string
// @Router /api/memberships [get]
func (h *Handler) GetMemberships(w http.ResponseWriter, r *http.Request) {
	rows, err := h.Store.DB.Query(`SELECT id,name,duration_days,price_cents,created_at FROM memberships ORDER BY id`)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "db error"})
		return
	}
	defer rows.Close()
	var out []Membership
	for rows.Next() {
		var m Membership
		if err := rows.Scan(&m.ID, &m.Name, &m.DurationDays, &m.PriceCents, &m.CreatedAt); err != nil {
			continue
		}
		out = append(out, m)
	}
	writeJSON(w, http.StatusOK, out)
}

type buyReq struct {
	MembershipID int    `json:"membership_id" example:"1"`
	Method       string `json:"method" example:"card"`
}

// BuyMembership godoc
// @Summary Купить подписку
// @Description Создаёт оплату и активирует подписку
// @Tags memberships
// @Security Bearer
// @Accept json
// @Produce json
// @Param body body buyReq true "ID подписки и метод оплаты"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/memberships/buy [post]
func (h *Handler) BuyMembership(w http.ResponseWriter, r *http.Request) {
	uid := r.Context().Value(ContextUserID).(int)
	var req buyReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid"})
		return
	}
	var price int
	var duration int
	err := h.Store.DB.QueryRow("SELECT price_cents,duration_days FROM memberships WHERE id = ?", req.MembershipID).Scan(&price, &duration)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "membership not found"})
		return
	}
	res, err := h.Store.DB.Exec("INSERT INTO payments(user_id,amount_cents,currency,method,status) VALUES(?,?,?,?,?)", uid, price, "KZT", req.Method, "done")
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "payment failed"})
		return
	}
	pid, _ := res.LastInsertId()
	start := time.Now().UTC()
	end := start.Add(time.Duration(duration*24) * time.Hour)
	_, err = h.Store.DB.Exec("INSERT INTO user_memberships(user_id, membership_id, start_date, end_date, active) VALUES(?,?,?,?,1)", uid, req.MembershipID, start.Format("2006-01-02"), end.Format("2006-01-02"))
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not assign membership"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"payment_id": pid, "start": start.Format("2006-01-02"), "end": end.Format("2006-01-02")})
}

type payReq struct {
	AmountCents int    `json:"amount_cents" example:"10000"`
	Method      string `json:"method" example:"card"`
}

// Pay godoc
// @Summary Произвольная оплата
// @Description Создаёт запись об оплате
// @Tags payments
// @Security Bearer
// @Accept json
// @Produce json
// @Param body body payReq true "Сумма и метод"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/payments [post]
func (h *Handler) Pay(w http.ResponseWriter, r *http.Request) {
	uid := r.Context().Value(ContextUserID).(int)
	var req payReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid"})
		return
	}
	res, err := h.Store.DB.Exec("INSERT INTO payments(user_id,amount_cents,currency,method,status) VALUES(?,?,?,?,?)", uid, req.AmountCents, "KZT", req.Method, "done")
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "payment failed"})
		return
	}
	id, _ := res.LastInsertId()
	writeJSON(w, http.StatusOK, map[string]interface{}{"payment_id": id})
}

type trainerReq struct {
	Name string `json:"name" example:"Алексей Петров"`
	Bio  string `json:"bio" example:"Мастер спорта по фитнесу"`
}

// CreateTrainer godoc
// @Summary Создать тренера
// @Description Только для админа
// @Tags admin
// @Security Bearer
// @Accept json
// @Produce json
// @Param body body trainerReq true "Данные тренера"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/admin/trainers [post]
func (h *Handler) CreateTrainer(w http.ResponseWriter, r *http.Request) {
	var req trainerReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid"})
		return
	}
	res, err := h.Store.DB.Exec("INSERT INTO trainers(name,bio) VALUES(?,?)", req.Name, req.Bio)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "db"})
		return
	}
	id, _ := res.LastInsertId()
	writeJSON(w, http.StatusCreated, map[string]interface{}{"trainer_id": id})
}

type classReq struct {
	Title       string `json:"title" example:"Йога для начинающих"`
	Description string `json:"description" example:"Спокойная практика для новичков"`
	TrainerID   int    `json:"trainer_id" example:"1"`
	StartTime   string `json:"start_time" example:"2025-11-20T10:00:00Z"`
	DurationMin int    `json:"duration_min" example:"60"`
	Capacity    int    `json:"capacity" example:"15"`
}

// CreateClass godoc
// @Summary Создать занятие
// @Description Только для админа
// @Tags admin
// @Security Bearer
// @Accept json
// @Produce json
// @Param body body classReq true "Данные занятия"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/admin/classes [post]
func (h *Handler) CreateClass(w http.ResponseWriter, r *http.Request) {
	var req classReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid"})
		return
	}
	if req.Title == "" || req.StartTime == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "title and start_time required"})
		return
	}
	_, err := h.Store.DB.Exec("INSERT INTO classes(title,description,trainer_id,start_time,duration_min,capacity) VALUES(?,?,?,?,?,?)", req.Title, req.Description, req.TrainerID, req.StartTime, req.DurationMin, req.Capacity)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "db"})
		return
	}
	writeJSON(w, http.StatusCreated, map[string]string{"ok": "created"})
}

type createBookingReq struct {
	ClassID int `json:"class_id" example:"1"`
}

// CreateBooking godoc
// @Summary Забронировать занятие
// @Description Требует активную подписку
// @Tags bookings
// @Security Bearer
// @Accept json
// @Produce json
// @Param body body createBookingReq true "ID занятия"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/bookings [post]
func (h *Handler) CreateBooking(w http.ResponseWriter, r *http.Request) {
	uid := r.Context().Value(ContextUserID).(int)
	var req createBookingReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid"})
		return
	}
	var capacity int
	var count int
	var start string
	err := h.Store.DB.QueryRow("SELECT capacity, start_time FROM classes WHERE id = ?", req.ClassID).Scan(&capacity, &start)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "class not found"})
		return
	}
	err = h.Store.DB.QueryRow("SELECT COUNT(1) FROM bookings WHERE class_id = ? AND status = 'booked'", req.ClassID).Scan(&count)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "db"})
		return
	}
	if count >= capacity {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "class full"})
		return
	}
	if ok, _ := h.userHasActiveMembership(uid); !ok {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "no active membership"})
		return
	}
	var exists int
	_ = h.Store.DB.QueryRow("SELECT COUNT(1) FROM bookings WHERE user_id = ? AND class_id = ? AND status = 'booked'", uid, req.ClassID).Scan(&exists)
	if exists > 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "already booked"})
		return
	}
	res, err := h.Store.DB.Exec("INSERT INTO bookings(user_id,class_id,status) VALUES(?,?, 'booked')", uid, req.ClassID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "db"})
		return
	}
	id, _ := res.LastInsertId()
	writeJSON(w, http.StatusCreated, map[string]interface{}{"booking_id": id, "start_time": start})
}

// ListBookings godoc
// @Summary Мои бронирования
// @Description Список всех бронирований пользователя
// @Tags bookings
// @Security Bearer
// @Success 200 {array} Booking
// @Failure 500 {object} map[string]string
// @Router /api/bookings [get]
func (h *Handler) ListBookings(w http.ResponseWriter, r *http.Request) {
	uid := r.Context().Value(ContextUserID).(int)
	rows, err := h.Store.DB.Query("SELECT id,user_id,class_id,status,created_at FROM bookings WHERE user_id = ? ORDER BY created_at DESC", uid)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "db"})
		return
	}
	defer rows.Close()
	var out []Booking
	for rows.Next() {
		var b Booking
		if err := rows.Scan(&b.ID, &b.UserID, &b.ClassID, &b.Status, &b.CreatedAt); err != nil {
			continue
		}
		out = append(out, b)
	}
	writeJSON(w, http.StatusOK, out)
}

func (h *Handler) userHasActiveMembership(userID int) (bool, error) {
	var cnt int
	now := time.Now().Format("2006-01-02")
	err := h.Store.DB.QueryRow("SELECT COUNT(1) FROM user_memberships WHERE user_id = ? AND active = 1 AND start_date <= ? AND end_date >= ?", userID, now, now).Scan(&cnt)
	if err != nil {
		return false, err
	}
	return cnt > 0, nil
}
