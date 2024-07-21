create table "users" (
  "username" varchar PRIMARY KEY,
  "hashed_password" varchar not null,
  "full_name" varchar not null,
  "email" varchar unique not null,
  "password_changed_at" timestamptz not null default '0001-01-01 00:00:00Z',
  "created_at" timestamptz not null default (now())
);

ALTER TABLE "accounts" ADD FOREIGN KEY ("owner") REFERENCES "users" ("username");
ALTER TABLE "accounts" ADD CONSTRAINT "owner_currency_key" UNIQUE ("owner", "currency");