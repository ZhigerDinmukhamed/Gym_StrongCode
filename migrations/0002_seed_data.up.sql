-- +goose Up
INSERT OR IGNORE INTO gyms(id, name, address) VALUES
(1, 'Central Gym', 'Алматы, Абай 123'),
(2, 'North Branch', 'Астана, Мәңгілік Ел 55');

INSERT OR IGNORE INTO memberships(id, name, duration_days, price_cents) VALUES
(1, 'Айлық', 30, 1500000),
(2, 'Кварталдық', 90, 4000000),
(3, 'Жылдық', 365, 15000000),
(4, 'VIP Жылдық', 365, 30000000);

INSERT OR IGNORE INTO trainers(id, name, bio) VALUES
(1, 'Жігер Дінмұхамед', 'Мастер спорта по кроссфиту'),
(2, 'Ақжол Кадырбаев', 'Сертифицированный йога-инструктор');