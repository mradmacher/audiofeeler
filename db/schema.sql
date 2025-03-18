CREATE TABLE accounts (
  id INTEGER PRIMARY KEY,
  name TEXT,
  dir TEXT
);

CREATE TABLE deploys (
  id INTEGER PRIMARY KEY,
  account_id INTEGER,
  server TEXT,
  username TEXT,
  username_iv TEXT,
  password TEXT,
  password_iv TEXT,
  dir TEXT,
  FOREIGN KEY(account_id) REFERENCES accounts(id)
);

CREATE TABLE events (
  id INTEGER PRIMARY KEY,
  account_id INTEGER,
  date TEXT,
  hour TEXT,
  venue TEXT,
  place TEXT,
  city TEXT,
  address TEXT,
  FOREIGN KEY(account_id) REFERENCES accounts(id)
);
