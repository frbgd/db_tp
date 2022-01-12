--
-- PostgreSQL database dump
--

BEGIN;

DROP SCHEMA IF EXISTS public CASCADE;
CREATE SCHEMA public;

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
-- Name: citext; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS citext WITH SCHEMA public;


--
-- Name: EXTENSION citext; Type: COMMENT; Schema: -; Owner:
--

COMMENT ON EXTENSION citext IS 'data type for case-insensitive character strings';


--
-- Name: add_new_forum_user(); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.add_new_forum_user() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    INSERT INTO forums_users_nicknames (forum_slug, user_nickname)
    VALUES (NEW.forum_slug, NEW.author_nickname)
    ON CONFLICT DO NOTHING;

    RETURN NULL;
END;
$$;


ALTER FUNCTION public.add_new_forum_user() OWNER TO postgres;

--
-- Name: add_path_to_post(); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.add_path_to_post() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
DECLARE
    parent_path BIGINT[];
BEGIN
    IF NEW.parent IS NULL THEN
        NEW.path := NEW.path || NEW.id;
    ELSE
        SELECT path
        FROM posts
        WHERE id = NEW.parent
          AND thread_id = NEW.thread_id
        INTO parent_path;

        IF (COALESCE(ARRAY_LENGTH(parent_path, 1), 0) = 0) THEN
            RAISE EXCEPTION
                'parent post with id=% not exists in thread with id=%',
                NEW.ID, NEW.thread_id;
        END IF;

        NEW.path := NEW.path || parent_path || NEW.id;
    END IF;
    RETURN NEW;
END;
$$;


ALTER FUNCTION public.add_path_to_post() OWNER TO postgres;

--
-- Name: increment_forum_posts(); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.increment_forum_posts() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    UPDATE forums
    SET posts = posts + 1
    WHERE slug = NEW.forum_slug;

    RETURN NEW;
END;
$$;


ALTER FUNCTION public.increment_forum_posts() OWNER TO postgres;

--
-- Name: increment_forum_threads(); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.increment_forum_threads() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    UPDATE forums
    SET threads = threads + 1
    WHERE slug = NEW.forum_slug;

    RETURN NULL;
END;
$$;


ALTER FUNCTION public.increment_forum_threads() OWNER TO postgres;

--
-- Name: insert_vote(); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.insert_vote() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    UPDATE threads
    SET votes = votes + NEW.voice
    WHERE id = NEW.thread_id;

    RETURN NULL;
END;
$$;


ALTER FUNCTION public.insert_vote() OWNER TO postgres;

--
-- Name: update_vote(); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.update_vote() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    UPDATE threads
    SET votes = votes + NEW.voice - OLD.voice
    WHERE id = NEW.thread_id;

    RETURN NULL;
END;
$$;


ALTER FUNCTION public.update_vote() OWNER TO postgres;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: forums; Type: TABLE; Schema: public; Owner: postgres
--

CREATE UNLOGGED TABLE public.forums (
    slug public.citext NOT NULL,
    title text NOT NULL,
    threads integer NOT NULL,
    posts bigint NOT NULL,
    owner_nickname public.citext NOT NULL
);


ALTER TABLE public.forums OWNER TO postgres;

--
-- Name: forums_users_nicknames; Type: TABLE; Schema: public; Owner: postgres
--

CREATE UNLOGGED TABLE public.forums_users_nicknames (
    forum_slug public.citext NOT NULL,
    user_nickname public.citext NOT NULL COLLATE pg_catalog.ucs_basic
);


ALTER TABLE public.forums_users_nicknames OWNER TO postgres;

--
-- Name: users; Type: TABLE; Schema: public; Owner: postgres
--

CREATE UNLOGGED TABLE public.users (
    nickname public.citext NOT NULL COLLATE pg_catalog.ucs_basic,
    email public.citext NOT NULL,
    fullname text NOT NULL,
    about text NOT NULL
);


ALTER TABLE public.users OWNER TO postgres;

--
-- Name: forums_users; Type: VIEW; Schema: public; Owner: postgres
--

CREATE VIEW public.forums_users AS
 SELECT fu_nicknames.forum_slug,
    fu_nicknames.user_nickname,
    u.nickname,
    u.email,
    u.fullname,
    u.about
   FROM (public.forums_users_nicknames fu_nicknames
     JOIN public.users u ON ((fu_nicknames.user_nickname OPERATOR(public.=) u.nickname)));


ALTER TABLE public.forums_users OWNER TO postgres;

--
-- Name: posts; Type: TABLE; Schema: public; Owner: postgres
--

CREATE UNLOGGED TABLE public.posts (
    id bigint NOT NULL,
    thread_id integer NOT NULL,
    author_nickname public.citext NOT NULL,
    forum_slug public.citext NOT NULL,
    is_edited boolean DEFAULT false NOT NULL,
    message text NOT NULL,
    parent bigint,
    created timestamp(3) with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    path bigint[] NOT NULL
);


ALTER TABLE public.posts OWNER TO postgres;

--
-- Name: posts_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.posts_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.posts_id_seq OWNER TO postgres;

--
-- Name: posts_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.posts_id_seq OWNED BY public.posts.id;


--
-- Name: threads; Type: TABLE; Schema: public; Owner: postgres
--

CREATE UNLOGGED TABLE public.threads (
    id integer NOT NULL,
    slug public.citext,
    forum_slug public.citext NOT NULL,
    author_nickname public.citext NOT NULL,
    title text NOT NULL,
    message text NOT NULL,
    votes integer DEFAULT 0 NOT NULL,
    created timestamp(3) with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


ALTER TABLE public.threads OWNER TO postgres;

--
-- Name: threads_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.threads_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.threads_id_seq OWNER TO postgres;

--
-- Name: threads_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.threads_id_seq OWNED BY public.threads.id;


--
-- Name: votes; Type: TABLE; Schema: public; Owner: postgres
--

CREATE UNLOGGED TABLE public.votes (
    thread_id integer NOT NULL,
    author_nickname public.citext NOT NULL,
    voice smallint NOT NULL
);


ALTER TABLE public.votes OWNER TO postgres;

--
-- Name: posts id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.posts ALTER COLUMN id SET DEFAULT nextval('public.posts_id_seq'::regclass);


--
-- Name: threads id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.threads ALTER COLUMN id SET DEFAULT nextval('public.threads_id_seq'::regclass);


--
-- Name: forums forums_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.forums
    ADD CONSTRAINT forums_pkey PRIMARY KEY (slug);


--
-- Name: forums_users_nicknames forums_users_nicknames_pk; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.forums_users_nicknames
    ADD CONSTRAINT forums_users_nicknames_pk PRIMARY KEY (forum_slug, user_nickname);


--
-- Name: posts posts_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.posts
    ADD CONSTRAINT posts_pkey PRIMARY KEY (id);


--
-- Name: threads threads_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.threads
    ADD CONSTRAINT threads_pkey PRIMARY KEY (id);


--
-- Name: threads threads_slug_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.threads
    ADD CONSTRAINT threads_slug_key UNIQUE (slug);


--
-- Name: users users_email_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_email_key UNIQUE (email);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (nickname);


--
-- Name: votes votes_pk; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.votes
    ADD CONSTRAINT votes_pk PRIMARY KEY (thread_id, author_nickname);


--
-- Name: posts_id_created_thread_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX posts_thread_created_id_idx ON public.posts USING btree (thread_id, created, id, created);

--
-- Name: threads_forum_slug_created_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX threads_forum_slug_created_idx ON public.threads USING btree (forum_slug, created);


--
-- Name: users_nickname_email_include_other_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX users_nickname_email_include_other_idx ON public.users USING btree (nickname, email) INCLUDE (about, fullname);


--
-- Name: posts add_new_forum_user_after_insert_in_posts; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER add_new_forum_user_after_insert_in_posts AFTER INSERT ON public.posts FOR EACH ROW EXECUTE FUNCTION public.add_new_forum_user();


--
-- Name: threads add_new_forum_user_after_insert_in_threads; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER add_new_forum_user_after_insert_in_threads AFTER INSERT ON public.threads FOR EACH ROW EXECUTE FUNCTION public.add_new_forum_user();


--
-- Name: posts add_path_to_post; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER add_path_to_post BEFORE INSERT ON public.posts FOR EACH ROW EXECUTE FUNCTION public.add_path_to_post();


--
-- Name: posts increment_forum_posts_after_insert_on_threads; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER increment_forum_posts_after_insert_on_threads BEFORE INSERT ON public.posts FOR EACH ROW EXECUTE FUNCTION public.increment_forum_posts();


--
-- Name: threads increment_forum_threads_after_insert_on_threads; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER increment_forum_threads_after_insert_on_threads AFTER INSERT ON public.threads FOR EACH ROW EXECUTE FUNCTION public.increment_forum_threads();


--
-- Name: votes insert_vote_after_insert_on_threads; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER insert_vote_after_insert_on_threads AFTER INSERT ON public.votes FOR EACH ROW EXECUTE FUNCTION public.insert_vote();


--
-- Name: votes update_vote_after_insert_on_threads; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER update_vote_after_insert_on_threads AFTER UPDATE ON public.votes FOR EACH ROW EXECUTE FUNCTION public.update_vote();


--
-- Name: forums forums_owner_nickname_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.forums
    ADD CONSTRAINT forums_owner_nickname_fkey FOREIGN KEY (owner_nickname) REFERENCES public.users(nickname) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: forums_users_nicknames forums_users_nicknames_forum_slug_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.forums_users_nicknames
    ADD CONSTRAINT forums_users_nicknames_forum_slug_fkey FOREIGN KEY (forum_slug) REFERENCES public.forums(slug) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: forums_users_nicknames forums_users_nicknames_user_nickname_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.forums_users_nicknames
    ADD CONSTRAINT forums_users_nicknames_user_nickname_fkey FOREIGN KEY (user_nickname) REFERENCES public.users(nickname) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: posts posts_author_nickname_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.posts
    ADD CONSTRAINT posts_author_nickname_fkey FOREIGN KEY (author_nickname) REFERENCES public.users(nickname) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: posts posts_forum_slug_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.posts
    ADD CONSTRAINT posts_forum_slug_fkey FOREIGN KEY (forum_slug) REFERENCES public.forums(slug) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: posts posts_thread_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.posts
    ADD CONSTRAINT posts_thread_id_fkey FOREIGN KEY (thread_id) REFERENCES public.threads(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: threads threads_author_nickname_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.threads
    ADD CONSTRAINT threads_author_nickname_fkey FOREIGN KEY (author_nickname) REFERENCES public.users(nickname) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: threads threads_forum_slug_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.threads
    ADD CONSTRAINT threads_forum_slug_fkey FOREIGN KEY (forum_slug) REFERENCES public.forums(slug) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: votes votes_author_nickname_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.votes
    ADD CONSTRAINT votes_author_nickname_fkey FOREIGN KEY (author_nickname) REFERENCES public.users(nickname) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: votes votes_thread_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.votes
    ADD CONSTRAINT votes_thread_id_fkey FOREIGN KEY (thread_id) REFERENCES public.threads(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

COMMIT;
