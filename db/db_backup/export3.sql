--
-- PostgreSQL database dump
--

-- Dumped from database version 14.15
-- Dumped by pg_dump version 14.15

-- Started on 2024-12-30 10:14:14 UTC

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- TOC entry 6 (class 2615 OID 41293)
-- Name: bank; Type: SCHEMA; Schema: -; Owner: ccat
--

CREATE SCHEMA bank;


ALTER SCHEMA bank OWNER TO ccat;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- TOC entry 212 (class 1259 OID 41090)
-- Name: accounts; Type: TABLE; Schema: bank; Owner: ccat
--

CREATE TABLE bank.accounts (
    id bigint NOT NULL,
    owner character varying NOT NULL,
    balance bigint NOT NULL,
    currency character varying NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE bank.accounts OWNER TO ccat;

--
-- TOC entry 211 (class 1259 OID 41089)
-- Name: accounts_id_seq; Type: SEQUENCE; Schema: bank; Owner: ccat
--

CREATE SEQUENCE bank.accounts_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE bank.accounts_id_seq OWNER TO ccat;

--
-- TOC entry 3445 (class 0 OID 0)
-- Dependencies: 211
-- Name: accounts_id_seq; Type: SEQUENCE OWNED BY; Schema: bank; Owner: ccat
--

ALTER SEQUENCE bank.accounts_id_seq OWNED BY bank.accounts.id;


--
-- TOC entry 214 (class 1259 OID 41100)
-- Name: entries; Type: TABLE; Schema: bank; Owner: ccat
--

CREATE TABLE bank.entries (
    id bigint NOT NULL,
    account_id bigint NOT NULL,
    amount bigint NOT NULL,
    currency character varying NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE bank.entries OWNER TO ccat;

--
-- TOC entry 3446 (class 0 OID 0)
-- Dependencies: 214
-- Name: COLUMN entries.amount; Type: COMMENT; Schema: bank; Owner: ccat
--

COMMENT ON COLUMN bank.entries.amount IS 'can be positive or negative';


--
-- TOC entry 213 (class 1259 OID 41099)
-- Name: entries_id_seq; Type: SEQUENCE; Schema: bank; Owner: ccat
--

CREATE SEQUENCE bank.entries_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE bank.entries_id_seq OWNER TO ccat;

--
-- TOC entry 3447 (class 0 OID 0)
-- Dependencies: 213
-- Name: entries_id_seq; Type: SEQUENCE OWNED BY; Schema: bank; Owner: ccat
--

ALTER SEQUENCE bank.entries_id_seq OWNED BY bank.entries.id;


--
-- TOC entry 210 (class 1259 OID 24640)
-- Name: schema_migrations; Type: TABLE; Schema: bank; Owner: ccat
--

CREATE TABLE bank.schema_migrations (
    version bigint NOT NULL,
    dirty boolean NOT NULL
);


ALTER TABLE bank.schema_migrations OWNER TO ccat;

--
-- TOC entry 218 (class 1259 OID 41256)
-- Name: sessions; Type: TABLE; Schema: bank; Owner: ccat
--

CREATE TABLE bank.sessions (
    id uuid NOT NULL,
    username character varying NOT NULL,
    refresh_token character varying NOT NULL,
    user_agent character varying NOT NULL,
    client_ip character varying NOT NULL,
    is_blocked boolean DEFAULT false NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    expires_at timestamp with time zone NOT NULL
);


ALTER TABLE bank.sessions OWNER TO ccat;

--
-- TOC entry 216 (class 1259 OID 41115)
-- Name: transfers; Type: TABLE; Schema: bank; Owner: ccat
--

CREATE TABLE bank.transfers (
    id bigint NOT NULL,
    from_account_id bigint NOT NULL,
    to_account_id bigint NOT NULL,
    amount bigint NOT NULL,
    currency character varying NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    CONSTRAINT transfers_amount_check CHECK ((amount > 0))
);


ALTER TABLE bank.transfers OWNER TO ccat;

--
-- TOC entry 3448 (class 0 OID 0)
-- Dependencies: 216
-- Name: COLUMN transfers.amount; Type: COMMENT; Schema: bank; Owner: ccat
--

COMMENT ON COLUMN bank.transfers.amount IS 'absolute value';


--
-- TOC entry 215 (class 1259 OID 41114)
-- Name: transfers_id_seq; Type: SEQUENCE; Schema: bank; Owner: ccat
--

CREATE SEQUENCE bank.transfers_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE bank.transfers_id_seq OWNER TO ccat;

--
-- TOC entry 3449 (class 0 OID 0)
-- Dependencies: 215
-- Name: transfers_id_seq; Type: SEQUENCE OWNED BY; Schema: bank; Owner: ccat
--

ALTER SEQUENCE bank.transfers_id_seq OWNED BY bank.transfers.id;


--
-- TOC entry 217 (class 1259 OID 41179)
-- Name: users; Type: TABLE; Schema: bank; Owner: ccat
--

CREATE TABLE bank.users (
    username character varying NOT NULL,
    hashed_password character varying NOT NULL,
    full_name character varying NOT NULL,
    email character varying NOT NULL,
    password_changed_at timestamp with time zone DEFAULT '0001-01-01 00:00:00+00'::timestamp with time zone NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE bank.users OWNER TO ccat;

--
-- TOC entry 3258 (class 2604 OID 41093)
-- Name: accounts id; Type: DEFAULT; Schema: bank; Owner: ccat
--

ALTER TABLE ONLY bank.accounts ALTER COLUMN id SET DEFAULT nextval('bank.accounts_id_seq'::regclass);


--
-- TOC entry 3260 (class 2604 OID 41103)
-- Name: entries id; Type: DEFAULT; Schema: bank; Owner: ccat
--

ALTER TABLE ONLY bank.entries ALTER COLUMN id SET DEFAULT nextval('bank.entries_id_seq'::regclass);


--
-- TOC entry 3262 (class 2604 OID 41118)
-- Name: transfers id; Type: DEFAULT; Schema: bank; Owner: ccat
--

ALTER TABLE ONLY bank.transfers ALTER COLUMN id SET DEFAULT nextval('bank.transfers_id_seq'::regclass);


--
-- TOC entry 3272 (class 2606 OID 41196)
-- Name: accounts accounts_owner_currency_unique; Type: CONSTRAINT; Schema: bank; Owner: ccat
--

ALTER TABLE ONLY bank.accounts
    ADD CONSTRAINT accounts_owner_currency_unique UNIQUE (owner, currency);


--
-- TOC entry 3275 (class 2606 OID 41098)
-- Name: accounts accounts_pkey; Type: CONSTRAINT; Schema: bank; Owner: ccat
--

ALTER TABLE ONLY bank.accounts
    ADD CONSTRAINT accounts_pkey PRIMARY KEY (id);


--
-- TOC entry 3278 (class 2606 OID 41108)
-- Name: entries entries_pkey; Type: CONSTRAINT; Schema: bank; Owner: ccat
--

ALTER TABLE ONLY bank.entries
    ADD CONSTRAINT entries_pkey PRIMARY KEY (id);


--
-- TOC entry 3270 (class 2606 OID 24644)
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: bank; Owner: ccat
--

ALTER TABLE ONLY bank.schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);


--
-- TOC entry 3291 (class 2606 OID 41264)
-- Name: sessions sessions_pkey; Type: CONSTRAINT; Schema: bank; Owner: ccat
--

ALTER TABLE ONLY bank.sessions
    ADD CONSTRAINT sessions_pkey PRIMARY KEY (id);


--
-- TOC entry 3282 (class 2606 OID 41124)
-- Name: transfers transfers_pkey; Type: CONSTRAINT; Schema: bank; Owner: ccat
--

ALTER TABLE ONLY bank.transfers
    ADD CONSTRAINT transfers_pkey PRIMARY KEY (id);


--
-- TOC entry 3285 (class 2606 OID 41189)
-- Name: users users_email_key; Type: CONSTRAINT; Schema: bank; Owner: ccat
--

ALTER TABLE ONLY bank.users
    ADD CONSTRAINT users_email_key UNIQUE (email);


--
-- TOC entry 3287 (class 2606 OID 41187)
-- Name: users users_pkey; Type: CONSTRAINT; Schema: bank; Owner: ccat
--

ALTER TABLE ONLY bank.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (username);


--
-- TOC entry 3273 (class 1259 OID 41135)
-- Name: accounts_owner_idx; Type: INDEX; Schema: bank; Owner: ccat
--

CREATE INDEX accounts_owner_idx ON bank.accounts USING btree (owner);


--
-- TOC entry 3276 (class 1259 OID 41136)
-- Name: entries_account_id_idx; Type: INDEX; Schema: bank; Owner: ccat
--

CREATE INDEX entries_account_id_idx ON bank.entries USING btree (account_id);


--
-- TOC entry 3288 (class 1259 OID 41268)
-- Name: sessions_created_at_idx; Type: INDEX; Schema: bank; Owner: ccat
--

CREATE INDEX sessions_created_at_idx ON bank.sessions USING btree (created_at);


--
-- TOC entry 3289 (class 1259 OID 41267)
-- Name: sessions_expires_at_idx; Type: INDEX; Schema: bank; Owner: ccat
--

CREATE INDEX sessions_expires_at_idx ON bank.sessions USING btree (expires_at);


--
-- TOC entry 3292 (class 1259 OID 41265)
-- Name: sessions_token_idx; Type: INDEX; Schema: bank; Owner: ccat
--

CREATE UNIQUE INDEX sessions_token_idx ON bank.sessions USING btree (refresh_token);


--
-- TOC entry 3293 (class 1259 OID 41269)
-- Name: sessions_username_expires_at_idx; Type: INDEX; Schema: bank; Owner: ccat
--

CREATE INDEX sessions_username_expires_at_idx ON bank.sessions USING btree (username, expires_at);


--
-- TOC entry 3294 (class 1259 OID 41270)
-- Name: sessions_username_id_idx; Type: INDEX; Schema: bank; Owner: ccat
--

CREATE INDEX sessions_username_id_idx ON bank.sessions USING btree (username, id);


--
-- TOC entry 3295 (class 1259 OID 41266)
-- Name: sessions_username_idx; Type: INDEX; Schema: bank; Owner: ccat
--

CREATE UNIQUE INDEX sessions_username_idx ON bank.sessions USING btree (username);


--
-- TOC entry 3279 (class 1259 OID 41137)
-- Name: transfers_from_account_id_idx; Type: INDEX; Schema: bank; Owner: ccat
--

CREATE INDEX transfers_from_account_id_idx ON bank.transfers USING btree (from_account_id);


--
-- TOC entry 3280 (class 1259 OID 41139)
-- Name: transfers_from_account_id_to_account_id_idx; Type: INDEX; Schema: bank; Owner: ccat
--

CREATE INDEX transfers_from_account_id_to_account_id_idx ON bank.transfers USING btree (from_account_id, to_account_id);


--
-- TOC entry 3283 (class 1259 OID 41138)
-- Name: transfers_to_account_id_idx; Type: INDEX; Schema: bank; Owner: ccat
--

CREATE INDEX transfers_to_account_id_idx ON bank.transfers USING btree (to_account_id);


--
-- TOC entry 3296 (class 2606 OID 41190)
-- Name: accounts accounts_owner_fkey; Type: FK CONSTRAINT; Schema: bank; Owner: ccat
--

ALTER TABLE ONLY bank.accounts
    ADD CONSTRAINT accounts_owner_fkey FOREIGN KEY (owner) REFERENCES bank.users(username);


--
-- TOC entry 3297 (class 2606 OID 41109)
-- Name: entries entries_account_id_fkey; Type: FK CONSTRAINT; Schema: bank; Owner: ccat
--

ALTER TABLE ONLY bank.entries
    ADD CONSTRAINT entries_account_id_fkey FOREIGN KEY (account_id) REFERENCES bank.accounts(id) ON DELETE CASCADE;


--
-- TOC entry 3300 (class 2606 OID 41271)
-- Name: sessions sessions_username_fkey; Type: FK CONSTRAINT; Schema: bank; Owner: ccat
--

ALTER TABLE ONLY bank.sessions
    ADD CONSTRAINT sessions_username_fkey FOREIGN KEY (username) REFERENCES bank.users(username) ON DELETE CASCADE;


--
-- TOC entry 3298 (class 2606 OID 41125)
-- Name: transfers transfers_from_account_id_fkey; Type: FK CONSTRAINT; Schema: bank; Owner: ccat
--

ALTER TABLE ONLY bank.transfers
    ADD CONSTRAINT transfers_from_account_id_fkey FOREIGN KEY (from_account_id) REFERENCES bank.accounts(id) ON DELETE CASCADE;


--
-- TOC entry 3299 (class 2606 OID 41130)
-- Name: transfers transfers_to_account_id_fkey; Type: FK CONSTRAINT; Schema: bank; Owner: ccat
--

ALTER TABLE ONLY bank.transfers
    ADD CONSTRAINT transfers_to_account_id_fkey FOREIGN KEY (to_account_id) REFERENCES bank.accounts(id) ON DELETE CASCADE;


--
-- TOC entry 2048 (class 826 OID 41294)
-- Name: DEFAULT PRIVILEGES FOR TABLES; Type: DEFAULT ACL; Schema: bank; Owner: ccat
--

ALTER DEFAULT PRIVILEGES FOR ROLE ccat IN SCHEMA bank GRANT ALL ON TABLES  TO ccat;


--
-- TOC entry 2049 (class 826 OID 41311)
-- Name: DEFAULT PRIVILEGES FOR TABLES; Type: DEFAULT ACL; Schema: public; Owner: ccat
--

ALTER DEFAULT PRIVILEGES FOR ROLE ccat IN SCHEMA public GRANT ALL ON TABLES  TO PUBLIC;


-- Completed on 2024-12-30 10:14:14 UTC

--
-- PostgreSQL database dump complete
--

