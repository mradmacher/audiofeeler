CREATE TABLE event_status (
  value TEXT NOT NULL PRIMARY KEY
);

INSERT INTO event_status (value)
  VALUES ('current'), ('archived');

CREATE TABLE account (
  id TEXT PRIMARY KEY,
  name TEXT
);

CREATE TABLE deployment (
  id TEXT PRIMARY KEY,
  account_id INTEGER NOT NULL,
  server TEXT,
  username TEXT,
  username_iv TEXT,
  password TEXT,
  password_iv TEXT,
  remote_dir TEXT,
  FOREIGN KEY(account_id) REFERENCES account(id) ON UPDATE RESTRICT ON DELETE RESTRICT
);

CREATE TABLE event (
  id TEXT PRIMARY KEY,
  account_id INTEGER NOT NULL,
  name TEXT,
  date TEXT,
  hour TEXT,
  venue TEXT,
  town TEXT,
  location TEXT,
  description TEXT,
  status TEXT NOT NULL,
  FOREIGN KEY(account_id) REFERENCES account(id) ON UPDATE RESTRICT ON DELETE RESTRICT,
  FOREIGN KEY(status) REFERENCES event_status(value) ON UPDATE RESTRICT ON DELETE RESTRICT
);

