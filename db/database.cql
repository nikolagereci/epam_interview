-- Create a keyspace
CREATE KEYSPACE IF NOT EXISTS companies WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : '1' };

-- Create a table
CREATE TABLE IF NOT EXISTS companies.company (
   id uuid PRIMARY KEY,
   name text,
   description text,
   employees int,
   registered boolean,
   type text
);
CREATE INDEX IF NOT EXISTS index_name ON companies.company (name);

-- Create a test keyspace
CREATE KEYSPACE IF NOT EXISTS companies_test WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : '1' };

-- Create a test table
CREATE TABLE IF NOT EXISTS companies_test.company (
   id uuid PRIMARY KEY,
   name text,
   description text,
   employees int,
   registered boolean,
   type text
);
CREATE INDEX IF NOT EXISTS index_name ON companies_test.company (name);

--empty test data
TRUNCATE companies_test.company;

