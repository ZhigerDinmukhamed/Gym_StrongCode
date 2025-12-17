package cache

import (
    "context"
    "encoding/json"
    "fmt"
    "time"

    "github.com/go-redis/redis/v8"
)

type RedisClient struct {
    client *redis.Client
    ctx    context.Context
}

func NewRedisClient(addr, password string, db int) (*RedisClient, error) {
    client := redis.NewClient(&redis.Options{
        Addr:     addr,
        Password: password,
        DB:       db,
        PoolSize: 100,
    })

    ctx := context.Background()
    
    // Проверяем соединение
    _, err := client.Ping(ctx).Result()
    if err != nil {
        return nil, fmt.Errorf("failed to connect to Redis: %w", err)
    }

    return &RedisClient{
        client: client,
        ctx:    ctx,
    }, nil
}

// Set сохраняет значение с TTL
func (r *RedisClient) Set(key string, value interface{}, ttl time.Duration) error {
    jsonData, err := json.Marshal(value)
    if err != nil {
        return err
    }

    return r.client.Set(r.ctx, key, jsonData, ttl).Err()
}

// Get получает значение по ключу
func (r *RedisClient) Get(key string, dest interface{}) error {
    val, err := r.client.Get(r.ctx, key).Result()
    if err != nil {
        return err
    }

    return json.Unmarshal([]byte(val), dest)
}

// Delete удаляет ключ
func (r *RedisClient) Delete(key string) error {
    return r.client.Del(r.ctx, key).Err()
}

// Exists проверяет существование ключа
func (r *RedisClient) Exists(key string) (bool, error) {
    result, err := r.client.Exists(r.ctx, key).Result()
    return result > 0, err
}

// Increment увеличивает значение счетчика
func (r *RedisClient) Increment(key string) (int64, error) {
    return r.client.Incr(r.ctx, key).Result()
}

// Decrement уменьшает значение счетчика
func (r *RedisClient) Decrement(key string) (int64, error) {
    return r.client.Decr(r.ctx, key).Result()
}

// SetNX устанавливает значение если ключ не существует
func (r *RedisClient) SetNX(key string, value interface{}, ttl time.Duration) (bool, error) {
    jsonData, err := json.Marshal(value)
    if err != nil {
        return false, err
    }

    return r.client.SetNX(r.ctx, key, jsonData, ttl).Result()
}

// HSet устанавливает значение в хэше
func (r *RedisClient) HSet(key, field string, value interface{}) error {
    jsonData, err := json.Marshal(value)
    if err != nil {
        return err
    }

    return r.client.HSet(r.ctx, key, field, jsonData).Err()
}

// HGet получает значение из хэша
func (r *RedisClient) HGet(key, field string, dest interface{}) error {
    val, err := r.client.HGet(r.ctx, key, field).Result()
    if err != nil {
        return err
    }

    return json.Unmarshal([]byte(val), dest)
}

// HDel удаляет поле из хэша
func (r *RedisClient) HDel(key, field string) error {
    return r.client.HDel(r.ctx, key, field).Err()
}

// HMSet устанавливает несколько полей в хэше
func (r *RedisClient) HMSet(key string, fields map[string]interface{}) error {
    redisFields := make(map[string]interface{})
    for field, value := range fields {
        jsonData, err := json.Marshal(value)
        if err != nil {
            return err
        }
        redisFields[field] = jsonData
    }

    return r.client.HMSet(r.ctx, key, redisFields).Err()
}

// HMGet получает несколько полей из хэша
func (r *RedisClient) HMGet(key string, fields ...string) (map[string]string, error) {
    result := make(map[string]string)
    
    values, err := r.client.HMGet(r.ctx, key, fields...).Result()
    if err != nil {
        return nil, err
    }

    for i, field := range fields {
        if values[i] != nil {
            result[field] = values[i].(string)
        }
    }

    return result, nil
}

// LPush добавляет элемент в начало списка
func (r *RedisClient) LPush(key string, values ...interface{}) error {
    redisValues := make([]interface{}, len(values))
    for i, value := range values {
        jsonData, err := json.Marshal(value)
        if err != nil {
            return err
        }
        redisValues[i] = jsonData
    }

    return r.client.LPush(r.ctx, key, redisValues...).Err()
}

// RPop удаляет и возвращает последний элемент списка
func (r *RedisClient) RPop(key string, dest interface{}) error {
    val, err := r.client.RPop(r.ctx, key).Result()
    if err != nil {
        return err
    }

    return json.Unmarshal([]byte(val), dest)
}

// SAdd добавляет элемент в множество
func (r *RedisClient) SAdd(key string, members ...interface{}) error {
    redisMembers := make([]interface{}, len(members))
    for i, member := range members {
        jsonData, err := json.Marshal(member)
        if err != nil {
            return err
        }
        redisMembers[i] = jsonData
    }

    return r.client.SAdd(r.ctx, key, redisMembers...).Err()
}

// SIsMember проверяет наличие элемента в множестве
func (r *RedisClient) SIsMember(key string, member interface{}) (bool, error) {
    jsonData, err := json.Marshal(member)
    if err != nil {
        return false, err
    }

    return r.client.SIsMember(r.ctx, key, jsonData).Result()
}

// ZAdd добавляет элемент в отсортированное множество
func (r *RedisClient) ZAdd(key string, score float64, member interface{}) error {
    jsonData, err := json.Marshal(member)
    if err != nil {
        return err
    }

    return r.client.ZAdd(r.ctx, key, &redis.Z{
        Score:  score,
        Member: jsonData,
    }).Err()
}

// ZRange получает диапазон элементов из отсортированного множества
func (r *RedisClient) ZRange(key string, start, stop int64, dest interface{}) error {
    vals, err := r.client.ZRange(r.ctx, key, start, stop).Result()
    if err != nil {
        return err
    }

    // Для простоты возвращаем первую строку
    if len(vals) > 0 {
        return json.Unmarshal([]byte(vals[0]), dest)
    }

    return fmt.Errorf("no elements in sorted set")
}

// Keys возвращает все ключи по шаблону
func (r *RedisClient) Keys(pattern string) ([]string, error) {
    return r.client.Keys(r.ctx, pattern).Result()
}

// FlushDB очищает всю базу
func (r *RedisClient) FlushDB() error {
    return r.client.FlushDB(r.ctx).Err()
}

// Close закрывает соединение
func (r *RedisClient) Close() error {
    return r.client.Close()
}

// GetTTL получает оставшееся время жизни ключа
func (r *RedisClient) GetTTL(key string) (time.Duration, error) {
    return r.client.TTL(r.ctx, key).Result()
}

// Expire устанавливает TTL для ключа
func (r *RedisClient) Expire(key string, ttl time.Duration) error {
    return r.client.Expire(r.ctx, key, ttl).Err()
}

// Publish публикует сообщение в канал
func (r *RedisClient) Publish(channel string, message interface{}) error {
    jsonData, err := json.Marshal(message)
    if err != nil {
        return err
    }

    return r.client.Publish(r.ctx, channel, jsonData).Err()
}

// Subscribe подписывается на канал
func (r *RedisClient) Subscribe(channel string) *redis.PubSub {
    return r.client.Subscribe(r.ctx, channel)
}

// Pipeline выполняет пайплайн операций
func (r *RedisClient) Pipeline() redis.Pipeliner {
    return r.client.Pipeline()
}