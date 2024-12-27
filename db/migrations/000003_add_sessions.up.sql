create table sessions (
  id uuid primary key,
  username varchar not null,
  refresh_token varchar not null,
  user_agent varchar not null,
  client_ip varchar not null,
  is_blocked boolean not null default false,
  created_at timestamptz not null default now(),
  expires_at timestamptz not null
)
;

create unique index sessions_token_idx on sessions (refresh_token);

create unique index sessions_username_idx on sessions (username);

create index sessions_expires_at_idx on sessions (expires_at);

create index sessions_created_at_idx on sessions (created_at);

create index sessions_username_expires_at_idx on sessions (username, expires_at);

create index sessions_username_id_idx on sessions (username, id);

alter table sessions
add foreign key (username) references users (username) on delete cascade
;