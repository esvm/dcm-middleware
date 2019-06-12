CREATE KEYSPACE broker
  WITH REPLICATION = {
    'class' : 'SimpleStrategy',
    'replication_factor' : 1
};

CREATE TABLE broker.topics (
  id text PRIMARY KEY,
  name text
);

CREATE TABLE broker.indexes (
  name text PRIMARY KEY,
  value text
);
