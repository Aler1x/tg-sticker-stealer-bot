package db

const Schema = `
CREATE TABLE IF NOT EXISTS packs (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	user_id INTEGER NOT NULL,
	pack_name TEXT NOT NULL,
	pack_title TEXT NOT NULL,
	pack_type TEXT NOT NULL,
	pack_link TEXT NOT NULL,
	sticker_count INTEGER NOT NULL,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	UNIQUE(user_id, pack_name)
);

CREATE INDEX IF NOT EXISTS idx_user_id ON packs(user_id);
CREATE INDEX IF NOT EXISTS idx_created_at ON packs(created_at);

CREATE TABLE IF NOT EXISTS users (
	user_id INTEGER PRIMARY KEY,
	username TEXT,
	first_name TEXT,
	last_name TEXT,
	language_code TEXT,
	is_active INTEGER DEFAULT 1,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	last_seen_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_users_last_seen ON users(last_seen_at);
CREATE INDEX IF NOT EXISTS idx_users_is_active ON users(is_active);

CREATE TABLE IF NOT EXISTS subscription_prices (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	subscription_type TEXT NOT NULL UNIQUE,
	price_stars INTEGER NOT NULL,
	description TEXT NOT NULL,
	value INTEGER NOT NULL,
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS user_subscriptions (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	user_id INTEGER NOT NULL,
	subscription_type TEXT NOT NULL,
	remaining_count INTEGER,
	expires_at DATETIME,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY(user_id) REFERENCES users(user_id)
);

CREATE INDEX IF NOT EXISTS idx_user_subscriptions_user_id ON user_subscriptions(user_id);
CREATE INDEX IF NOT EXISTS idx_user_subscriptions_expires_at ON user_subscriptions(expires_at);

CREATE TABLE IF NOT EXISTS payment_history (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	user_id INTEGER NOT NULL,
	subscription_type TEXT NOT NULL,
	price_stars INTEGER NOT NULL,
	payment_charge_id TEXT,
	payment_provider TEXT,
	status TEXT NOT NULL,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY(user_id) REFERENCES users(user_id)
);

CREATE INDEX IF NOT EXISTS idx_payment_history_user_id ON payment_history(user_id);
CREATE INDEX IF NOT EXISTS idx_payment_history_status ON payment_history(status);

INSERT OR IGNORE INTO subscription_prices (subscription_type, price_stars, description, value) VALUES
	('one_steal', 10, 'Single emoji steal', 1),
	('ten_steals', 80, '10 emoji steals', 10),
	('week', 50, '1 week unlimited', 7),
	('month', 150, '1 month unlimited', 30),
	('year', 1200, '1 year unlimited', 365);
`

