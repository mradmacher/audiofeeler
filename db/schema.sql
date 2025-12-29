DROP TABLE IF EXISTS events;
DROP TABLE IF EXISTS deployments;
DROP TABLE IF EXISTS accounts;

CREATE TABLE accounts (
  id INTEGER PRIMARY KEY,
  name TEXT,
  source_dir TEXT
);

CREATE TABLE deployments (
  id INTEGER PRIMARY KEY,
  account_id INTEGER,
  server TEXT,
  username TEXT,
  username_iv TEXT,
  password TEXT,
  password_iv TEXT,
  remote_dir TEXT,
  FOREIGN KEY(account_id) REFERENCES accounts(id)
);

CREATE TABLE events (
  id INTEGER PRIMARY KEY,
  account_id INTEGER,
  name TEXT,
  date TEXT,
  hour TEXT,
  venue TEXT,
  town TEXT,
  location TEXT,
  description TEXT,
  status INTEGER,
  FOREIGN KEY(account_id) REFERENCES accounts(id)
);
