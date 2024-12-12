-- Create the accounts table
create table
  accounts (
    id bigserial primary KEY,
    owner varchar not null,
    balance bigint not null, -- Maybe change to decimal later
    currency varchar not null,
    created_at timestamptz not null default (now ())
  )
;
-- Create the entries table
create table
  entries (
    id bigserial primary KEY,
    account_id bigint not null references accounts (id) on delete cascade,
    amount bigint, -- Can be positive or negative
    created_at timestamptz not null default (now ())
  )
;
-- Create the transfers table
create table
  transfers (
    id bigserial primary KEY,
    from_account_id bigint not null references accounts (id) on delete cascade,
    to_account_id bigint not null references accounts (id) on delete cascade,
    amount bigint check (amount > 0), -- Can be only positive (absolute value)
    created_at timestamptz not null default (now ())
  )
;
create index on "accounts" ("owner")
;
create index on "entries" ("account_id")
;
create index on "transfers" ("from_account_id")
;
create index on "transfers" ("to_account_id")
;
create index on "transfers" ("from_account_id", "to_account_id")
;
comment on column "entries"."amount" is 'can be positive or negative'
;
comment on column "transfers"."amount" is 'absolute value'
;