-- +goose Up
INSERT OR IGNORE INTO memberships(id, name, duration_days, price_cents) VALUES
(1, 'Месячная', 30, 15000),
(2, 'Квартальная', 90, 40000),
(3, 'Годовая', 365, 150000),
(4, 'VIP Годовая', 365, 300000);