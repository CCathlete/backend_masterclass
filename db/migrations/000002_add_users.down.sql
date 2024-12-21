alter table if exists "accounts" drop constraint "accounts_owner_currency_unique";

alter table if exists "accounts" drop constraint "accounts_owner_fkey";

drop table if exists "users";