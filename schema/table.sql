CREATE TABLE public.role
(
    id     VARCHAR(32) PRIMARY KEY,
    active boolean     NOT NULL DEFAULT true,
    name   VARCHAR(32) NOT NULL
);
CREATE TABLE public.web_user
(
    id                     VARCHAR(32) PRIMARY KEY,
    name                   VARCHAR(128)        NOT NULL,
    username               VARCHAR(64) UNIQUE  NOT NULL,
    password_hash          VARCHAR(128)        NOT NULL,
    password_salt          VARCHAR(32)         NOT NULL,
    email                  VARCHAR(255) UNIQUE NOT NULL,
    role_id                VARCHAR(32)         NOT NULL,
    active                 boolean             NOT NULL DEFAULT true,
    registration_timestamp TIMESTAMP           NOT NULL,
    created_by             VARCHAR(32),
    FOREIGN KEY (role_id) REFERENCES public.role (id),
    FOREIGN KEY (created_by) REFERENCES public.web_user (id)
);
CREATE TABLE public.menu
(
    id    VARCHAR(32) PRIMARY KEY,
    name  VARCHAR(32)  NOT NULL,
    seq_order SMALLINT     NOT NULL,
    path  VARCHAR(128) NOT NULL,
    icon  VARCHAR(32)
);
CREATE TABLE public.sub_menu
(
    id      VARCHAR(32) PRIMARY KEY,
    menu_id VARCHAR(32)  NOT NULL,
    name    VARCHAR(32)  NOT NULL,
    seq_order   SMALLINT     NOT NULL,
    outcome VARCHAR(128) NOT NULL,
    icon    VARCHAR(32),
    FOREIGN KEY (menu_id) REFERENCES public.menu (id)
);
CREATE TABLE public.permission
(
    id          VARCHAR(32) PRIMARY KEY,
    sub_menu_id VARCHAR(32)  NOT NULL,
    name        VARCHAR(32)  NOT NULL,
    seq_order       SMALLINT     NOT NULL,
    apis        VARCHAR(256) NOT NULL,
    FOREIGN KEY (sub_menu_id) REFERENCES public.sub_menu (id)
);
CREATE TABLE public.role_permission
(
    role_id       VARCHAR(32) NOT NULL,
    permission_id VARCHAR(32) NOT NULL,
    FOREIGN KEY (role_id) REFERENCES public.role (id),
    FOREIGN KEY (permission_id) REFERENCES public.permission (id)
);
CREATE TABLE public.config
(
    id    VARCHAR(32) PRIMARY KEY,
    value VARCHAR(512) NOT NULL
);
CREATE TABLE public.product
(
    id          VARCHAR(32) PRIMARY KEY,
    code        VARCHAR(32)  NOT NULL,
    name        VARCHAR(128) NOT NULL,
    description VARCHAR(256),
    active      boolean      NOT NULL DEFAULT true
);
CREATE TABLE public.supplier
(
    id          VARCHAR(32) PRIMARY KEY,
    code        VARCHAR(32)  NOT NULL,
    name        VARCHAR(128) NOT NULL,
    description VARCHAR(256),
    active      boolean      NOT NULL DEFAULT true
);
CREATE TABLE public.customer
(
    id          VARCHAR(32) PRIMARY KEY,
    code        VARCHAR(32)  NOT NULL,
    name        VARCHAR(128) NOT NULL,
    description VARCHAR(256),
    active      boolean      NOT NULL DEFAULT true
);

