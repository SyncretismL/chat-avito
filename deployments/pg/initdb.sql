CREATE TABLE public.users (
    id bigserial PRIMARY KEY,
    username text UNIQUE NOT NULL,
    created_at timestamp NOT NULL
);

CREATE TABLE public.chats (
    id bigserial PRIMARY KEY,
    name text UNIQUE NOT NULL,
    created_at timestamp NOT NULL
);

CREATE TABLE public.users_chats (
    id bigserial PRIMARY KEY,
    user_id integer NOT NULL,
    chat_id integer NOT NULL,
    FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (chat_id) REFERENCES public.chats(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE public.messages (
    id bigserial PRIMARY KEY,
    chat bigserial NOT NULL,
    author bigserial NOT NULL,
    text text NOT NULL,
    created_at timestamp NOT NULL,
    FOREIGN KEY (author) REFERENCES public.users(id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (chat) REFERENCES public.chats(id) ON DELETE CASCADE ON UPDATE CASCADE
);