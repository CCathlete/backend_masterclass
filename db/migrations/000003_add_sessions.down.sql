-- Removing foreign key constraints.
alter table if exists sessions drop constraint sessions_username_fkey;

-- Removing indexes.
drop index if exists sessions_expires_at_idx;
drop index if exists sessions_created_at_idx;
drop index if exists sessions_username_expires_at_idx;
drop index if exists sessions_username_id_idx;
drop index if exists sessions_token_idx;
drop index if exists sessions_username_idx;

-- Dropping the table.
drop table if exists sessions;