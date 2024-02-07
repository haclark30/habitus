CREATE TABLE users (
  id integer primary key autoincrement,
  username text unique NOT NULL,
  passwordHash text NOT NULL
);

CREATE TABLE sessions (
  token text primary key,
  userId integer NOT NULL,
  FOREIGN KEY(userId) REFERENCES users(id)
);

CREATE TABLE habits (
  id integer primary key autoincrement,
  userId integer NOT NULL,
  name text NOT NULL,
  hasUp boolean NOT NULL,
  hasDown boolean NOT NULL,
  FOREIGN KEY(userId) REFERENCES users(id)
);

CREATE TABLE habitLog (
  id integer PRIMARY KEY autoincrement,
  habitId integer NOT NULL,
  upCount integer NOT NULL,
  downCount integer NOT NULL,
  dateTime integer NOT NULL,
  FOREIGN KEY(habitId) REFERENCES habits(id)
);

CREATE TABLE dailys (
  id integer primary key autoincrement,
  userId integer NOT NULL,
  name text NOT NULL,
  FOREIGN KEY(userId) REFERENCES users(id)
);

CREATE TABLE dailyLog(
  id integer primary key autoincrement,
  dailyId integer NOT NULL,
  done boolean NOT NULL,
  dateTime integer NOT NULL,
  FOREIGN KEY(dailyId) REFERENCES dailys(id)
);
