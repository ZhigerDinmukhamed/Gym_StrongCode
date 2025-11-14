package router

import (
	"encoding/json"
	"net/http"

	"gym-api/middleware"
	"gym-api/store"
	"gym-api/utils"

	"github.com/gorilla/mux"
)

func NewRouter(st *store.Store) *mux.Router {
	r := mux.NewRouter()

	// health
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// AUTH
	r.HandleFunc("/api/register", func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Name     string `json:"name"`
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		json.NewDecoder(r.Body).Decode(&body)

		hash, _ := utils.HashPassword(body.Password)
		_, err := st.DB.Exec("INSERT INTO users(name,email,password_hash) VALUES(?,?,?)",
			body.Name, body.Email, hash)

		if err != nil {
			http.Error(w, "user exists?", 400)
			return
		}

		w.Write([]byte("registered"))
	}).Methods("POST")

	r.HandleFunc("/api/login", func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		json.NewDecoder(r.Body).Decode(&body)

		var id int
		var hash string
		var admin int

		err := st.DB.QueryRow("SELECT id,password_hash,is_admin FROM users WHERE email=?", body.Email).Scan(&id, &hash, &admin)
		if err != nil {
			http.Error(w, "invalid login", 400)
			return
		}

		if !utils.CheckPassword(body.Password, hash) {
			http.Error(w, "wrong password", 400)
			return
		}

		token, _ := utils.CreateToken(id, body.Email, admin == 1)
		json.NewEncoder(w).Encode(map[string]string{"token": token})
	}).Methods("POST")

	// Protected routes
	auth := r.PathPrefix("/api").Subrouter()
	auth.Use(middleware.AuthMiddleware)

	auth.HandleFunc("/classes", func(w http.ResponseWriter, r *http.Request) {
		rows, _ := st.DB.Query("SELECT id,title,description,start_time FROM classes")
		var out []map[string]interface{}
		for rows.Next() {
			var id int
			var t, d, s string
			rows.Scan(&id, &t, &d, &s)
			out = append(out, map[string]interface{}{
				"id":          id,
				"title":       t,
				"description": d,
				"start_time":  s,
			})
		}
		json.NewEncoder(w).Encode(out)
	}).Methods("GET")

	auth.HandleFunc("/book", func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			ClassID int `json:"class_id"`
		}
		json.NewDecoder(r.Body).Decode(&body)

		user := utils.GetUserFromContext(r.Context())

		_, err := st.DB.Exec("INSERT INTO bookings(user_id,class_id) VALUES(?,?)", user.UserID, body.ClassID)
		if err != nil {
			http.Error(w, "book error", 400)
			return
		}

		w.Write([]byte("booked"))
	}).Methods("POST")

	// Admin
	admin := r.PathPrefix("/api/admin").Subrouter()
	admin.Use(middleware.AuthMiddleware)

	admin.HandleFunc("/class", func(w http.ResponseWriter, r *http.Request) {
		user := utils.GetUserFromContext(r.Context())
		if !user.IsAdmin {
			http.Error(w, "forbidden", 403)
			return
		}

		var body struct {
			Title       string `json:"title"`
			Description string `json:"description"`
			StartTime   string `json:"start_time"`
			TrainerID   int    `json:"trainer_id"`
			Duration    int    `json:"duration_min"`
			Capacity    int    `json:"capacity"`
		}
		json.NewDecoder(r.Body).Decode(&body)

		_, err := st.DB.Exec("INSERT INTO classes(title,description,start_time,trainer_id,duration_min,capacity) VALUES(?,?,?,?,?,?)",
			body.Title, body.Description, body.StartTime, body.TrainerID, body.Duration, body.Capacity)

		if err != nil {
			http.Error(w, "error", 400)
			return
		}

		w.Write([]byte("class added"))
	}).Methods("POST")

	return r
}
