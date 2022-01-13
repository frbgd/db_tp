DROP SCHEMA IF EXISTS public CASCADE;
CREATE SCHEMA public;

CREATE EXTENSION IF NOT EXISTS citext WITH SCHEMA public;
COMMENT ON EXTENSION citext IS 'data type for case-insensitive character strings';

-- functions

CREATE FUNCTION public.insert_forum_user() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    INSERT INTO forums_users_nicknames (forum_slug, user_nickname) VALUES (NEW.forum_slug, NEW.user_nickname) ON CONFLICT DO NOTHING;
    RETURN NULL;
END;
$$;
ALTER FUNCTION public.insert_forum_user() OWNER TO postgres;

CREATE FUNCTION public.insert_post_path() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
DECLARE
    parent_path BIGINT[];
BEGIN
    IF NEW.parent IS NULL THEN
        NEW.path := NEW.path || NEW.id;
    ELSE
        SELECT path FROM posts WHERE id = NEW.parent AND thread_id = NEW.thread_id INTO parent_path;

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
ALTER FUNCTION public.insert_post_path() OWNER TO postgres;

CREATE FUNCTION public.inc_forum_posts() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    UPDATE forums SET posts = posts + 1 WHERE slug = NEW.forum_slug;
    RETURN NEW;
END;
$$;
ALTER FUNCTION public.inc_forum_posts() OWNER TO postgres;

CREATE FUNCTION public.inc_forum_threads() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    UPDATE forums SET threads = threads + 1 WHERE slug = NEW.forum_slug;
    RETURN NULL;
END;
$$;
ALTER FUNCTION public.inc_forum_threads() OWNER TO postgres;

CREATE FUNCTION public.insert_vote() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    UPDATE threads SET votes = votes + NEW.voice WHERE id = NEW.thread_id;
    RETURN NULL;
END;
$$;
ALTER FUNCTION public.insert_vote() OWNER TO postgres;

CREATE FUNCTION public.update_vote() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    UPDATE threads SET votes = votes + NEW.voice - OLD.voice WHERE id = NEW.thread_id;
    RETURN NULL;
END;
$$;
ALTER FUNCTION public.update_vote() OWNER TO postgres;

SET default_tablespace = '';
SET default_table_access_method = heap;

-- tables

CREATE UNLOGGED TABLE public.forums (
    slug public.citext NOT NULL,
    title text NOT NULL,
    threads integer NOT NULL,
    posts bigint NOT NULL,
    owner_nickname public.citext NOT NULL
);
ALTER TABLE public.forums OWNER TO postgres;

CREATE UNLOGGED TABLE public.forums_users_nicknames (
    forum_slug public.citext NOT NULL,
    user_nickname public.citext NOT NULL COLLATE pg_catalog.ucs_basic
);
ALTER TABLE public.forums_users_nicknames OWNER TO postgres;

CREATE UNLOGGED TABLE public.users (
    nickname public.citext NOT NULL COLLATE pg_catalog.ucs_basic,
    email public.citext NOT NULL,
    fullname text NOT NULL,
    about text NOT NULL
);
ALTER TABLE public.users OWNER TO postgres;

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

CREATE UNLOGGED TABLE public.posts (
    id bigint NOT NULL,
    thread_id integer NOT NULL,
    user_nickname public.citext NOT NULL,
    forum_slug public.citext NOT NULL,
    is_edited boolean DEFAULT false NOT NULL,
    message text NOT NULL,
    parent bigint,
    created timestamp(3) with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    path bigint[] NOT NULL
);
ALTER TABLE public.posts OWNER TO postgres;
CREATE SEQUENCE public.posts_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
ALTER TABLE public.posts_id_seq OWNER TO postgres;
ALTER SEQUENCE public.posts_id_seq OWNED BY public.posts.id;

CREATE UNLOGGED TABLE public.threads (
    id integer NOT NULL,
    slug public.citext,
    forum_slug public.citext NOT NULL,
    user_nickname public.citext NOT NULL,
    title text NOT NULL,
    message text NOT NULL,
    votes integer DEFAULT 0 NOT NULL,
    created timestamp(3) with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);
ALTER TABLE public.threads OWNER TO postgres;
CREATE SEQUENCE public.threads_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
ALTER TABLE public.threads_id_seq OWNER TO postgres;
ALTER SEQUENCE public.threads_id_seq OWNED BY public.threads.id;

CREATE UNLOGGED TABLE public.votes (
    thread_id integer NOT NULL,
    user_nickname public.citext NOT NULL,
    voice smallint NOT NULL
);
ALTER TABLE public.votes OWNER TO postgres;

ALTER TABLE ONLY public.posts ALTER COLUMN id SET DEFAULT nextval('public.posts_id_seq'::regclass);
ALTER TABLE ONLY public.threads ALTER COLUMN id SET DEFAULT nextval('public.threads_id_seq'::regclass);
ALTER TABLE ONLY public.forums
    ADD CONSTRAINT forums_pkey PRIMARY KEY (slug);
ALTER TABLE ONLY public.forums_users_nicknames
    ADD CONSTRAINT forums_users_nicknames_pk PRIMARY KEY (forum_slug, user_nickname);
ALTER TABLE ONLY public.posts
    ADD CONSTRAINT posts_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.threads
    ADD CONSTRAINT threads_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.threads
    ADD CONSTRAINT threads_slug_key UNIQUE (slug);
ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_email_key UNIQUE (email);
ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (nickname);
ALTER TABLE ONLY public.votes
    ADD CONSTRAINT votes_pk PRIMARY KEY (thread_id, user_nickname);

-- indices

CREATE INDEX idx_posts_id_created_thread_id ON public.posts USING btree (id, created, thread_id);
CREATE INDEX idx_posts_id_path ON public.posts USING btree (id, path);
CREATE INDEX idx_posts_parent_id ON public.posts USING btree (parent, id);
CREATE INDEX idx_posts_thread_id_id ON public.posts USING btree (thread_id, id);
CREATE INDEX idx_posts_thread_id_parent_path ON public.posts USING btree (thread_id, parent, path);
CREATE INDEX idx_posts_thread_id_path1_id ON public.posts USING btree (thread_id, (path[1]), id);
CREATE INDEX idx_posts_thread_id_path ON public.posts USING btree (thread_id, path);
CREATE INDEX idx_threads_forum_slug_created ON public.threads USING btree (forum_slug, created);
CREATE INDEX idx_users_nickname_email_include_about_fullname ON public.users USING btree (nickname, email) INCLUDE (about, fullname);

-- triggers

CREATE TRIGGER insert_forum_user_after_insert_in_posts AFTER INSERT ON public.posts FOR EACH ROW EXECUTE FUNCTION public.insert_forum_user();
CREATE TRIGGER insert_forum_user_after_insert_in_threads AFTER INSERT ON public.threads FOR EACH ROW EXECUTE FUNCTION public.insert_forum_user();
CREATE TRIGGER insert_post_path BEFORE INSERT ON public.posts FOR EACH ROW EXECUTE FUNCTION public.insert_post_path();
CREATE TRIGGER inc_forum_posts_after_insert_on_threads BEFORE INSERT ON public.posts FOR EACH ROW EXECUTE FUNCTION public.inc_forum_posts();
CREATE TRIGGER inc_forum_threads_after_insert_on_threads AFTER INSERT ON public.threads FOR EACH ROW EXECUTE FUNCTION public.inc_forum_threads();
CREATE TRIGGER insert_vote_after_insert_on_threads AFTER INSERT ON public.votes FOR EACH ROW EXECUTE FUNCTION public.insert_vote();
CREATE TRIGGER update_vote_after_insert_on_threads AFTER UPDATE ON public.votes FOR EACH ROW EXECUTE FUNCTION public.update_vote();

-- foreign keys

ALTER TABLE ONLY public.forums
    ADD CONSTRAINT forums_owner_nickname_fkey FOREIGN KEY (owner_nickname) REFERENCES public.users(nickname) ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE ONLY public.forums_users_nicknames
    ADD CONSTRAINT forums_users_nicknames_forum_slug_fkey FOREIGN KEY (forum_slug) REFERENCES public.forums(slug) ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE ONLY public.forums_users_nicknames
    ADD CONSTRAINT forums_users_nicknames_user_nickname_fkey FOREIGN KEY (user_nickname) REFERENCES public.users(nickname) ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE ONLY public.posts
    ADD CONSTRAINT posts_user_nickname_fkey FOREIGN KEY (user_nickname) REFERENCES public.users(nickname) ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE ONLY public.posts
    ADD CONSTRAINT posts_forum_slug_fkey FOREIGN KEY (forum_slug) REFERENCES public.forums(slug) ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE ONLY public.posts
    ADD CONSTRAINT posts_thread_id_fkey FOREIGN KEY (thread_id) REFERENCES public.threads(id) ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE ONLY public.threads
    ADD CONSTRAINT threads_user_nickname_fkey FOREIGN KEY (user_nickname) REFERENCES public.users(nickname) ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE ONLY public.threads
    ADD CONSTRAINT threads_forum_slug_fkey FOREIGN KEY (forum_slug) REFERENCES public.forums(slug) ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE ONLY public.votes
    ADD CONSTRAINT votes_user_nickname_fkey FOREIGN KEY (user_nickname) REFERENCES public.users(nickname) ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE ONLY public.votes
    ADD CONSTRAINT votes_thread_id_fkey FOREIGN KEY (thread_id) REFERENCES public.threads(id) ON UPDATE CASCADE ON DELETE CASCADE;
