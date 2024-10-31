CREATE TABLE auth (
                      "login" text not null primary key,
                      "hash_pass" text not null);
CREATE TABLE orders (
                        "orders_id" varchar primary key,
                        "login" text not null REFERENCES auth (login),
                        "accrual"  double precision,
                        "status" text not null default 'NEW', CHECK (status in ('NEW', 'PROCESSING', 'INVALID', 'PROCESSED')),
                        "uploaded_at" timestamp not null default now());
CREATE TABLE withdrawn_balance (
                                   "login" varchar REFERENCES auth (login),
                                   "order" varchar,
                                   "processed_at" timestamp not null default now(),
                                   "withdrawn" double precision not null, CHECK (withdrawn > 0));
CREATE TABLE current_balance (
                                 "login" varchar primary key REFERENCES auth (login),
                                 "current" double precision not null, CHECK (current >= 0));
CREATE INDEX withdrawn_login_idx ON withdrawn_balance (login);