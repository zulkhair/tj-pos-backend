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
    id        VARCHAR(32) PRIMARY KEY,
    name      VARCHAR(32)  NOT NULL,
    seq_order SMALLINT     NOT NULL,
    path      VARCHAR(128) NOT NULL,
    icon      VARCHAR(32)
);
CREATE TABLE public.sub_menu
(
    id        VARCHAR(32) PRIMARY KEY,
    menu_id   VARCHAR(32)  NOT NULL,
    name      VARCHAR(32)  NOT NULL,
    seq_order SMALLINT     NOT NULL,
    outcome   VARCHAR(128) NOT NULL,
    icon      VARCHAR(32),
    FOREIGN KEY (menu_id) REFERENCES public.menu (id)
);
CREATE TABLE public.permission
(
    id          VARCHAR(32) PRIMARY KEY,
    sub_menu_id VARCHAR(32)  NOT NULL,
    name        VARCHAR(32)  NOT NULL,
    seq_order   SMALLINT     NOT NULL,
    apis        VARCHAR(256) NOT NULL,
    FOREIGN KEY (sub_menu_id) REFERENCES public.sub_menu (id)
);
CREATE TABLE public.role_permission
(
    role_id       VARCHAR(32) PRIMARY KEY,
    permission_id VARCHAR(32) PRIMARY KEY,
    FOREIGN KEY (role_id) REFERENCES public.role (id),
    FOREIGN KEY (permission_id) REFERENCES public.permission (id),
    PRIMARY KEY (role_id, permission_id)
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
CREATE TABLE public.unit
(
    id          VARCHAR(32),
    code        VARCHAR(32) NOT NULL,
    description VARCHAR(32),
    active      boolean     NOT NULL DEFAULT true,
    PRIMARY KEY (id)
);
CREATE TABLE public.buy_price
(
    id               VARCHAR(32) PRIMARY KEY,
    date             TIMESTAMP WITHOUT TIME ZONE,
    supplier_id      VARCHAR(32),
    unit_id          VARCHAR(32) NOT NULL,
    product_id       VARCHAR(32) NOT NULL,
    price            NUMERIC     NOT NULL,
    web_user_id      VARCHAR(32) NOT NULL,
    latest           BOOLEAN DEFAULT FALSE,
    transaction_code VARCHAR(128),
    FOREIGN KEY (web_user_id) REFERENCES public.web_user (id),
    FOREIGN KEY (supplier_id) REFERENCES public.supplier (id),
    FOREIGN KEY (unit_id) REFERENCES public.unit (id),
    FOREIGN KEY (product_id) REFERENCES public.product (id)
);
CREATE TABLE public.sell_price
(
    id          VARCHAR(32) PRIMARY KEY,
    date        TIMESTAMP WITHOUT TIME ZONE,
    customer_id VARCHAR(32) NOT NULL,
    unit_id     VARCHAR(32) NOT NULL,
    product_id  VARCHAR(32) NOT NULL,
    price       NUMERIC     NOT NULL,
    web_user_id VARCHAR(32) NOT NULL,
    latest      BOOLEAN DEFAULT FALSE,
    transaction_code VARCHAR(128),
    FOREIGN KEY (web_user_id) REFERENCES public.web_user (id),
    FOREIGN KEY (customer_id) REFERENCES public.customer (id),
    FOREIGN KEY (unit_id) REFERENCES public.unit (id),
    FOREIGN KEY (product_id) REFERENCES public.product (id)
);
CREATE TABLE public.transaction
(
    id               VARCHAR(32)                 NOT NULL PRIMARY KEY,
    code             VARCHAR(128)                NOT NULL UNIQUE,
    date             TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    stakeholder_id   VARCHAR(32)                 NOT NULL,
    transaction_type VARCHAR(16)                 NOT NULL,
    status           VARCHAR(16)                 NOT NULL,
    reference_code   VARCHAR(128)
);
CREATE TABLE public.transaction_detail
(
    transaction_id VARCHAR(128),
    unit_id        VARCHAR(32),
    product_id     VARCHAR(32),
    price          NUMERIC  NOT NULL,
    quantity       SMALLINT NOT NULL,
    date           TIMESTAMP WITHOUT TIME ZONE,
    web_user_id    VARCHAR(32),
    FOREIGN KEY (web_user_id) REFERENCES public.web_user (id),
    FOREIGN KEY (transaction_id) REFERENCES public.transaction (id),
    FOREIGN KEY (unit_id) REFERENCES public.unit (id),
    FOREIGN KEY (product_id) REFERENCES public.product (id),
    PRIMARY KEY (transaction_id, unit_id, product_id)
);
CREATE TABLE public.SEQUENCE
(
    id         VARCHAR(128) PRIMARY KEY,
    next_value INTEGER NOT NULL DEFAULT 0
);

