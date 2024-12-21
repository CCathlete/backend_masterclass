--
-- PostgreSQL database dump
--

-- Dumped from database version 14.15
-- Dumped by pg_dump version 14.15

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

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: accounts; Type: TABLE; Schema: public; Owner: ccat
--

CREATE TABLE public.accounts (
    id bigint NOT NULL,
    owner character varying NOT NULL,
    balance bigint NOT NULL,
    currency character varying NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.accounts OWNER TO ccat;

--
-- Name: accounts_id_seq; Type: SEQUENCE; Schema: public; Owner: ccat
--

CREATE SEQUENCE public.accounts_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.accounts_id_seq OWNER TO ccat;

--
-- Name: accounts_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: ccat
--

ALTER SEQUENCE public.accounts_id_seq OWNED BY public.accounts.id;


--
-- Name: entries; Type: TABLE; Schema: public; Owner: ccat
--

CREATE TABLE public.entries (
    id bigint NOT NULL,
    account_id bigint NOT NULL,
    amount bigint NOT NULL,
    currency character varying NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.entries OWNER TO ccat;

--
-- Name: COLUMN entries.amount; Type: COMMENT; Schema: public; Owner: ccat
--

COMMENT ON COLUMN public.entries.amount IS 'can be positive or negative';


--
-- Name: entries_id_seq; Type: SEQUENCE; Schema: public; Owner: ccat
--

CREATE SEQUENCE public.entries_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.entries_id_seq OWNER TO ccat;

--
-- Name: entries_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: ccat
--

ALTER SEQUENCE public.entries_id_seq OWNED BY public.entries.id;


--
-- Name: schema_migrations; Type: TABLE; Schema: public; Owner: ccat
--

CREATE TABLE public.schema_migrations (
    version bigint NOT NULL,
    dirty boolean NOT NULL
);


ALTER TABLE public.schema_migrations OWNER TO ccat;

--
-- Name: transfers; Type: TABLE; Schema: public; Owner: ccat
--

CREATE TABLE public.transfers (
    id bigint NOT NULL,
    from_account_id bigint NOT NULL,
    to_account_id bigint NOT NULL,
    amount bigint NOT NULL,
    currency character varying NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    CONSTRAINT transfers_amount_check CHECK ((amount > 0))
);


ALTER TABLE public.transfers OWNER TO ccat;

--
-- Name: COLUMN transfers.amount; Type: COMMENT; Schema: public; Owner: ccat
--

COMMENT ON COLUMN public.transfers.amount IS 'absolute value';


--
-- Name: transfers_id_seq; Type: SEQUENCE; Schema: public; Owner: ccat
--

CREATE SEQUENCE public.transfers_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.transfers_id_seq OWNER TO ccat;

--
-- Name: transfers_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: ccat
--

ALTER SEQUENCE public.transfers_id_seq OWNED BY public.transfers.id;


--
-- Name: accounts id; Type: DEFAULT; Schema: public; Owner: ccat
--

ALTER TABLE ONLY public.accounts ALTER COLUMN id SET DEFAULT nextval('public.accounts_id_seq'::regclass);


--
-- Name: entries id; Type: DEFAULT; Schema: public; Owner: ccat
--

ALTER TABLE ONLY public.entries ALTER COLUMN id SET DEFAULT nextval('public.entries_id_seq'::regclass);


--
-- Name: transfers id; Type: DEFAULT; Schema: public; Owner: ccat
--

ALTER TABLE ONLY public.transfers ALTER COLUMN id SET DEFAULT nextval('public.transfers_id_seq'::regclass);


--
-- Name: accounts accounts_pkey; Type: CONSTRAINT; Schema: public; Owner: ccat
--

ALTER TABLE ONLY public.accounts
    ADD CONSTRAINT accounts_pkey PRIMARY KEY (id);


--
-- Name: entries entries_pkey; Type: CONSTRAINT; Schema: public; Owner: ccat
--

ALTER TABLE ONLY public.entries
    ADD CONSTRAINT entries_pkey PRIMARY KEY (id);


--
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: ccat
--

ALTER TABLE ONLY public.schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);


--
-- Name: transfers transfers_pkey; Type: CONSTRAINT; Schema: public; Owner: ccat
--

ALTER TABLE ONLY public.transfers
    ADD CONSTRAINT transfers_pkey PRIMARY KEY (id);


--
-- Name: accounts_owner_idx; Type: INDEX; Schema: public; Owner: ccat
--

CREATE INDEX accounts_owner_idx ON public.accounts USING btree (owner);


--
-- Name: entries_account_id_idx; Type: INDEX; Schema: public; Owner: ccat
--

CREATE INDEX entries_account_id_idx ON public.entries USING btree (account_id);


--
-- Name: transfers_from_account_id_idx; Type: INDEX; Schema: public; Owner: ccat
--

CREATE INDEX transfers_from_account_id_idx ON public.transfers USING btree (from_account_id);


--
-- Name: transfers_from_account_id_to_account_id_idx; Type: INDEX; Schema: public; Owner: ccat
--

CREATE INDEX transfers_from_account_id_to_account_id_idx ON public.transfers USING btree (from_account_id, to_account_id);


--
-- Name: transfers_to_account_id_idx; Type: INDEX; Schema: public; Owner: ccat
--

CREATE INDEX transfers_to_account_id_idx ON public.transfers USING btree (to_account_id);


--
-- Name: entries entries_account_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: ccat
--

ALTER TABLE ONLY public.entries
    ADD CONSTRAINT entries_account_id_fkey FOREIGN KEY (account_id) REFERENCES public.accounts(id) ON DELETE CASCADE;


--
-- Name: transfers transfers_from_account_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: ccat
--

ALTER TABLE ONLY public.transfers
    ADD CONSTRAINT transfers_from_account_id_fkey FOREIGN KEY (from_account_id) REFERENCES public.accounts(id) ON DELETE CASCADE;


--
-- Name: transfers transfers_to_account_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: ccat
--

ALTER TABLE ONLY public.transfers
    ADD CONSTRAINT transfers_to_account_id_fkey FOREIGN KEY (to_account_id) REFERENCES public.accounts(id) ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

