package cache

import (
    "Gym_StrongCode/internal/models"
    "fmt"
    "time"
)

type DataCache struct {
    redisClient *RedisClient
}

func NewDataCache(redisClient *RedisClient) *DataCache {
    return &DataCache{
        redisClient: redisClient,
    }
}

// Cache методы для классов

func (dc *DataCache) SetClasses(classes []models.Class) error {
    return dc.redisClient.Set("classes:all", classes, 5*time.Minute)
}

func (dc *DataCache) GetClasses() ([]models.Class, error) {
    var classes []models.Class
    err := dc.redisClient.Get("classes:all", &classes)
    return classes, err
}

func (dc *DataCache) SetClass(id int, class *models.Class) error {
    return dc.redisClient.Set(
        fmt.Sprintf("class:%d", id),
        class,
        10*time.Minute,
    )
}

func (dc *DataCache) GetClass(id int) (*models.Class, error) {
    var class models.Class
    err := dc.redisClient.Get(
        fmt.Sprintf("class:%d", id),
        &class,
    )
    if err != nil {
        return nil, err
    }
    return &class, nil
}

func (dc *DataCache) DeleteClass(id int) error {
    return dc.redisClient.Delete(fmt.Sprintf("class:%d", id))
}

func (dc *DataCache) InvalidateClasses() error {
    return dc.redisClient.Delete("classes:all")
}

// Cache методы для тренеров

func (dc *DataCache) SetTrainers(trainers []models.Trainer) error {
    return dc.redisClient.Set("trainers:all", trainers, 10*time.Minute)
}

func (dc *DataCache) GetTrainers() ([]models.Trainer, error) {
    var trainers []models.Trainer
    err := dc.redisClient.Get("trainers:all", &trainers)
    return trainers, err
}

func (dc *DataCache) SetTrainer(id int, trainer *models.Trainer) error {
    return dc.redisClient.Set(
        fmt.Sprintf("trainer:%d", id),
        trainer,
        30*time.Minute,
    )
}

func (dc *DataCache) GetTrainer(id int) (*models.Trainer, error) {
    var trainer models.Trainer
    err := dc.redisClient.Get(
        fmt.Sprintf("trainer:%d", id),
        &trainer,
    )
    if err != nil {
        return nil, err
    }
    return &trainer, nil
}

func (dc *DataCache) DeleteTrainer(id int) error {
    return dc.redisClient.Delete(fmt.Sprintf("trainer:%d", id))
}

func (dc *DataCache) InvalidateTrainers() error {
    return dc.redisClient.Delete("trainers:all")
}

// Cache методы для подписок

func (dc *DataCache) SetMemberships(memberships []models.Membership) error {
    return dc.redisClient.Set("memberships:all", memberships, 30*time.Minute)
}

func (dc *DataCache) GetMemberships() ([]models.Membership, error) {
    var memberships []models.Membership
    err := dc.redisClient.Get("memberships:all", &memberships)
    return memberships, err
}

func (dc *DataCache) SetMembership(id int, membership *models.Membership) error {
    return dc.redisClient.Set(
        fmt.Sprintf("membership:%d", id),
        membership,
        1*time.Hour,
    )
}

func (dc *DataCache) GetMembership(id int) (*models.Membership, error) {
    var membership models.Membership
    err := dc.redisClient.Get(
        fmt.Sprintf("membership:%d", id),
        &membership,
    )
    if err != nil {
        return nil, err
    }
    return &membership, nil
}

func (dc *DataCache) DeleteMembership(id int) error {
    return dc.redisClient.Delete(fmt.Sprintf("membership:%d", id))
}

func (dc *DataCache) InvalidateMemberships() error {
    return dc.redisClient.Delete("memberships:all")
}

// Cache методы для пользователей

func (dc *DataCache) SetUser(id int, user *models.User) error {
    return dc.redisClient.Set(
        fmt.Sprintf("user:%d", id),
        user,
        15*time.Minute,
    )
}

func (dc *DataCache) GetUser(id int) (*models.User, error) {
    var user models.User
    err := dc.redisClient.Get(
        fmt.Sprintf("user:%d", id),
        &user,
    )
    if err != nil {
        return nil, err
    }
    return &user, nil
}

func (dc *DataCache) DeleteUser(id int) error {
    return dc.redisClient.Delete(fmt.Sprintf("user:%d", id))
}

func (dc *DataCache) SetUserMembership(userID int, membership *models.UserMembership) error {
    return dc.redisClient.Set(
        fmt.Sprintf("user_membership:%d", userID),
        membership,
        5*time.Minute,
    )
}

func (dc *DataCache) GetUserMembership(userID int) (*models.UserMembership, error) {
    var membership models.UserMembership
    err := dc.redisClient.Get(
        fmt.Sprintf("user_membership:%d", userID),
        &membership,
    )
    if err != nil {
        return nil, err
    }
    return &membership, nil
}

func (dc *DataCache) DeleteUserMembership(userID int) error {
    return dc.redisClient.Delete(fmt.Sprintf("user_membership:%d", userID))
}

// Cache методы для бронирований

func (dc *DataCache) SetUserBookings(userID int, bookings []models.Booking) error {
    return dc.redisClient.Set(
        fmt.Sprintf("user_bookings:%d", userID),
        bookings,
        2*time.Minute,
    )
}

func (dc *DataCache) GetUserBookings(userID int) ([]models.Booking, error) {
    var bookings []models.Booking
    err := dc.redisClient.Get(
        fmt.Sprintf("user_bookings:%d", userID),
        &bookings,
    )
    return bookings, err
}

func (dc *DataCache) DeleteUserBookings(userID int) error {
    return dc.redisClient.Delete(fmt.Sprintf("user_bookings:%d", userID))
}

// Rate limiting

func (dc *DataCache) CheckRateLimit(key string, limit int, window time.Duration) (bool, error) {
    current, err := dc.redisClient.Get(key, &struct{}{})
    if err != nil && err.Error() != "redis: nil" {
        return false, err
    }

    count := 0
    if err == nil {
        // Ключ существует, парсим значение
        // Здесь упрощенная логика - в реальности нужно хранить счетчик
        count = 1
    }

    if count >= limit {
        return false, nil
    }

    // Увеличиваем счетчик
    err = dc.redisClient.Set(key, count+1, window)
    return err == nil, err
}

// Invalidate all cache
func (dc *DataCache) InvalidateAll() error {
    // Удаляем все ключи с префиксами кэша
    patterns := []string{
        "classes:*",
        "trainers:*",
        "memberships:*",
        "user:*",
        "user_membership:*",
        "user_bookings:*",
        "class:*",
        "trainer:*",
        "membership:*",
    }

    for _, pattern := range patterns {
        keys, err := dc.redisClient.Keys(pattern)
        if err != nil {
            continue
        }
        for _, key := range keys {
            dc.redisClient.Delete(key)
        }
    }

    return nil
}

// Статистика кэша
func (dc *DataCache) GetStats() (map[string]interface{}, error) {
    stats := make(map[string]interface{})
    
    // Считаем ключи по типам
    cacheTypes := []string{
        "classes:*",
        "trainers:*",
        "memberships:*",
        "user:*",
        "session:*",
    }

    for _, pattern := range cacheTypes {
        keys, err := dc.redisClient.Keys(pattern)
        if err != nil {
            continue
        }
        stats[pattern] = len(keys)
    }

    // Получаем общее количество ключей
    allKeys, _ := dc.redisClient.Keys("*")
    stats["total_keys"] = len(allKeys)

    return stats, nil
}