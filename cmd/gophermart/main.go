package main

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/golang-jwt/jwt/v4"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
)

var log = Log()
var dbStorage DBStorage
var DatabaseDsn string
var RunAddress string
var AccrualSystemAddress string

func Log() *zap.SugaredLogger {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	return logger.Sugar()
}

func InitServerConfig() {

	flag.StringVar(&DatabaseDsn, "d", fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", `localhost`, 5434, `docker`, `docker`, `postgres`), "Aдрес подключения к базе данных")
	flag.StringVar(&RunAddress, "a", ":8080", "Aдрес и порт запуска сервиса")
	flag.StringVar(&AccrualSystemAddress, "r", "", "адрес системы расчёта начислений")

	flag.Parse()
	if envDatabaseDsn := os.Getenv("DATABASE_URI"); envDatabaseDsn != "" {
		DatabaseDsn = envDatabaseDsn
	}
	if envRunAddress := os.Getenv("RUN_ADDRESS"); envRunAddress != "" {
		RunAddress = envRunAddress
	}
	if envAccrualSystemAddress := os.Getenv("ACCRUAL_SYSTEM_ADDRESS"); envAccrualSystemAddress != "" {
		AccrualSystemAddress = envAccrualSystemAddress
	}
}

func main() {
	InitServerConfig()
	dbStorage = InitDB(context.Background())
	r := chi.NewRouter()

	r.Use(middleware.RealIP, middleware.Recoverer)
	r.Post("/api/user/register", UserRegister)
	//r.Post("/api/user/login", handlers.SaveMetDataForJSON(iStorage))
	//r.Post("/api/user/orders", handlers.SaveArraysMetDataForJSON(iStorage))
	//r.Get("/api/user/orders", handlers.GetDataForJSON(iStorage))
	//
	//r.Get("/api/user/balance", handlers.GetDataByPathParams(iStorage))
	//r.Post("/api/user/balance/withdraw", handlers.GetPing(iStorage))
	//r.Get("/api/user/withdrawals", handlers.GetAllData(iStorage))

	log.Fatal(http.ListenAndServe(RunAddress, r))
}

type UserModel struct {
	ID       *int64  `json:"id,omitempty"`
	Login    *string `json:"login"`
	Password *string `json:"password"`
}

func (u UserModel) isValid() error {
	if u.Login == nil || u.Password == nil || *u.Login == "" || *u.Password == "" {
		return fmt.Errorf("ошибка валидации! Обязательные поля: password и login, не могут быть пустыми или null: %+v", u)
	}
	return nil
}

var sampleSecretKey = []byte("GoDiplomKey")

func generateJWT(username string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["username"] = username
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	tokenString, err := token.SignedString(sampleSecretKey)

	if err != nil {
		log.Error("Ошибка генерации токена: %s", err.Error())
		return "", err
	}
	return tokenString, nil
}

func UserRegister(w http.ResponseWriter, req *http.Request) {
	var tModel UserModel
	err := json.NewDecoder(req.Body).Decode(&tModel)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, "Ошибка десериализации!", http.StatusBadRequest)
		return
	}
	err = tModel.isValid()
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	auth, err := dbStorage.GetAuth(*tModel.Login)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err := dbStorage.CreateAuth(tModel)
			if err != nil {
				log.Error(err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			token, err := generateJWT(*tModel.Login)
			if err != nil {
				log.Error(err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			cookie := &http.Cookie{
				Name:     "Token",
				Value:    token,
				Secure:   false,
				HttpOnly: true,
				MaxAge:   300,
			}
			http.SetCookie(w, cookie)
			w.WriteHeader(http.StatusOK)
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}
	if auth != nil {
		http.Error(w, "Пользователь уже существует", http.StatusConflict)
		return
	}

}

type DBStorage struct {
	DB  *sql.DB
	ctx context.Context
}

func InitDB(ctx context.Context) DBStorage {
	db, err := sql.Open("pgx", DatabaseDsn)
	if err != nil {
		log.Error(err.Error())
	}
	tCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	if err = db.PingContext(tCtx); err != nil {
		log.Error(err.Error())
	}
	dbStorage := DBStorage{DB: db, ctx: ctx}
	dbStorage.InitTable(tCtx)
	return dbStorage
}

func encodeHash(value string) string {
	h := hmac.New(sha256.New, []byte("KEY123!"))
	h.Write([]byte(value))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func (db DBStorage) CreateAuth(authModel UserModel) error {
	_, err := db.DB.ExecContext(db.ctx, `INSERT INTO auth(login, hash_pass) values ($1,$2)`, authModel.Login, encodeHash(*authModel.Password))
	return err
}

func (db DBStorage) GetAuth(login string) (*UserModel, error) {
	var data UserModel
	rows := db.DB.QueryRow("SELECT * from auth where login = $1", login)
	err := rows.Scan(&data.ID, &data.Login, &data.Password)
	if err != nil {
		return nil, err
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (db DBStorage) InitTable(ctx context.Context) {
	_, err := db.DB.ExecContext(ctx, `DROP TABLE IF EXISTS auth`)
	if err != nil {
		log.Error(err.Error())
	}
	_, err = db.DB.ExecContext(ctx, `CREATE TABLE auth (
        "id" bigserial primary key,
        "login" text not null unique,
        "hash_pass" text not null
      )`)
	if err != nil {
		log.Error(err.Error())
	}
}
