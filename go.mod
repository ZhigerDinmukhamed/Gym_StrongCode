module Gym-StrongCode

go 1.23

require (
    github.com/gin-gonic/gin v1.10.0
    github.com/go-sqlite3 v1.14.22 // или github.com/mattn/go-sqlite3
    github.com/golang-jwt/jwt/v5 v5.2.1
    github.com/golang-migrate/migrate/v4 v4.17.1
    github.com/swaggo/gin-swagger v1.6.0
    github.com/swaggo/files v1.0.1
    github.com/spf13/viper v1.19.0
    go.uber.org/zap v1.27.0
    golang.org/x/crypto v0.27.0
    net/smtp v0.0.0 // встроен, но для ясности
    github.com/jordan-wright/email v1.3.0 // для удобной отправки email
    github.com/gin-contrib/rate-limit v0.0.0-20240618041941-aca2e7d8870f // bonus rate limit
)

require (
    // indirect dependencies
    github.com/bytedance/sonic v1.11.6 // indirect
    github.com/chenzhuoyu/base64x v0.0.0-20221115062448-fe3a3abad311 // indirect
    // ... остальные indirect
)