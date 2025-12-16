-- +goose Down
DROP TABLE IF EXISTS payments;
DROP TABLE IF EXISTS bookings;
DROP TABLE IF EXISTS classes;
DROP TABLE IF EXISTS trainers;
DROP TABLE IF EXISTS user_memberships;
DROP TABLE IF EXISTS memberships;
DROP TABLE IF EXISTS users;