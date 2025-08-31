--
-- PostgreSQL database dump
--

\restrict dB7f2gCPrJ7iuEFGtqE2c6gieMveeCCyiU6fgjOpiObNeoficTqcDSKCqFEgOeC

-- Dumped from database version 15.14 (Debian 15.14-1.pgdg13+1)
-- Dumped by pg_dump version 15.14 (Debian 15.14-1.pgdg13+1)

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
-- Name: posts_search_vector_trigger(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.posts_search_vector_trigger() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
begin
  new.search_vector :=
    to_tsvector('english', coalesce(new.content,'') || ' ' || coalesce(new.content_html,''));
  return new;
end
$$;


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: connections; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.connections (
    user_a text NOT NULL,
    user_b text NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    status smallint DEFAULT 1 NOT NULL
);


--
-- Name: follows; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.follows (
    follower_id text NOT NULL,
    followee_id text NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    status smallint DEFAULT 1 NOT NULL
);


--
-- Name: post_media; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.post_media (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    post_id uuid NOT NULL,
    url text NOT NULL,
    media_type character varying(32) NOT NULL,
    width integer,
    height integer,
    size_bytes bigint,
    meta jsonb DEFAULT '{}'::jsonb,
    sort_index integer DEFAULT 0,
    created_at timestamp with time zone DEFAULT now() NOT NULL
);


--
-- Name: posts; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.posts (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    author_id text NOT NULL,
    content text NOT NULL,
    content_html text,
    visibility smallint DEFAULT 0 NOT NULL,
    reply_to_post_id uuid,
    original_post_id uuid,
    language character varying(8),
    meta jsonb DEFAULT '{}'::jsonb,
    likes_count bigint DEFAULT 0 NOT NULL,
    comments_count bigint DEFAULT 0 NOT NULL,
    shares_count bigint DEFAULT 0 NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    deleted_at timestamp with time zone,
    search_vector tsvector
);


--
-- Name: schema_migrations; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.schema_migrations (
    version bigint NOT NULL,
    dirty boolean NOT NULL
);


--
-- Name: connections connections_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.connections
    ADD CONSTRAINT connections_pkey PRIMARY KEY (user_a, user_b);


--
-- Name: follows follows_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.follows
    ADD CONSTRAINT follows_pkey PRIMARY KEY (follower_id, followee_id);


--
-- Name: post_media post_media_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.post_media
    ADD CONSTRAINT post_media_pkey PRIMARY KEY (id);


--
-- Name: posts posts_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.posts
    ADD CONSTRAINT posts_pkey PRIMARY KEY (id);


--
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);


--
-- Name: idx_follows_followee; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_follows_followee ON public.follows USING btree (followee_id);


--
-- Name: idx_follows_follower; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_follows_follower ON public.follows USING btree (follower_id);


--
-- Name: idx_post_media_post; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_post_media_post ON public.post_media USING btree (post_id);


--
-- Name: idx_posts_search_vector; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_posts_search_vector ON public.posts USING gin (search_vector);


--
-- Name: posts tsvectorupdate; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER tsvectorupdate BEFORE INSERT OR UPDATE ON public.posts FOR EACH ROW EXECUTE FUNCTION public.posts_search_vector_trigger();


--
-- Name: post_media post_media_post_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.post_media
    ADD CONSTRAINT post_media_post_id_fkey FOREIGN KEY (post_id) REFERENCES public.posts(id) ON DELETE CASCADE;


--
-- Name: posts posts_original_post_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.posts
    ADD CONSTRAINT posts_original_post_id_fkey FOREIGN KEY (original_post_id) REFERENCES public.posts(id) ON DELETE SET NULL;


--
-- Name: posts posts_reply_to_post_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.posts
    ADD CONSTRAINT posts_reply_to_post_id_fkey FOREIGN KEY (reply_to_post_id) REFERENCES public.posts(id) ON DELETE SET NULL;


--
-- PostgreSQL database dump complete
--

\unrestrict dB7f2gCPrJ7iuEFGtqE2c6gieMveeCCyiU6fgjOpiObNeoficTqcDSKCqFEgOeC

