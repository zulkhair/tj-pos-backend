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
    id          VARCHAR(128) PRIMARY KEY,
    sub_menu_id VARCHAR(32)  NOT NULL,
    name        VARCHAR(32)  NOT NULL,
    seq_order   SMALLINT     NOT NULL,
    apis        VARCHAR(256) NOT NULL,
    FOREIGN KEY (sub_menu_id) REFERENCES public.sub_menu (id)
);
CREATE TABLE public.role_permission
(
    role_id       VARCHAR(32),
    permission_id VARCHAR(128),
    FOREIGN KEY (role_id) REFERENCES public.role (id),
    FOREIGN KEY (permission_id) REFERENCES public.permission (id),
    PRIMARY KEY (role_id, permission_id)
);
CREATE TABLE public.config
(
    id    VARCHAR(32) PRIMARY KEY,
    value VARCHAR(512) NOT NULL
);
CREATE TABLE public.unit
(
    id          VARCHAR(32),
    code        VARCHAR(32) NOT NULL,
    description VARCHAR(32),
    active      boolean     NOT NULL DEFAULT true,
    PRIMARY KEY (id)
);
CREATE TABLE public.product
(
    id          VARCHAR(32) PRIMARY KEY,
    unit_id     VARCHAR(32)  NOT NULL,
    code        VARCHAR(32)  NOT NULL,
    name        VARCHAR(128) NOT NULL,
    description VARCHAR(256),
    active      boolean      NOT NULL DEFAULT true,
    FOREIGN KEY (unit_id) REFERENCES public.unit (id)
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
CREATE TABLE public.transaction
(
    id               VARCHAR(32)  NOT NULL PRIMARY KEY,
    code             VARCHAR(128) NOT NULL UNIQUE,
    date             TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    stakeholder_id   VARCHAR(32)  NOT NULL,
    transaction_type VARCHAR(16)  NOT NULL,
    status           VARCHAR(16)  NOT NULL,
    reference_code   VARCHAR(128),
    web_user_id      VARCHAR(32),
    created_time     TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    FOREIGN KEY (web_user_id) REFERENCES public.web_user (id)
);
CREATE TABLE public.transaction_detail
(
    id             VARCHAR(32) PRIMARY KEY,
    transaction_id VARCHAR(32),
    product_id     VARCHAR(32),
    buy_price      NUMERIC  NOT NULL,
    sell_price     NUMERIC  NOT NULL,
    quantity       SMALLINT NOT NULL,
    buy_quantity   SMALLINT NOT NULL,
    created_time   TIMESTAMP WITHOUT TIME ZONE,
    web_user_id    VARCHAR(32),
    latest         BOOLEAN DEFAULT TRUE,
    FOREIGN KEY (web_user_id) REFERENCES public.web_user (id),
    FOREIGN KEY (transaction_id) REFERENCES public.transaction (id),
    FOREIGN KEY (product_id) REFERENCES public.product (id)
);
CREATE TABLE public.buy_price
(
    id             VARCHAR(32) PRIMARY KEY,
    date           TIMESTAMP WITHOUT TIME ZONE,
    supplier_id    VARCHAR(32),
    product_id     VARCHAR(32) NOT NULL,
    price          NUMERIC     NOT NULL,
    web_user_id    VARCHAR(32) NOT NULL,
    latest         BOOLEAN DEFAULT FALSE,
    transaction_id VARCHAR(32),
    FOREIGN KEY (web_user_id) REFERENCES public.web_user (id),
    FOREIGN KEY (supplier_id) REFERENCES public.supplier (id),
    FOREIGN KEY (product_id) REFERENCES public.product (id),
    FOREIGN KEY (transaction_id) REFERENCES public.transaction (id)
);
CREATE TABLE public.sell_price
(
    id             VARCHAR(32) PRIMARY KEY,
    date           TIMESTAMP WITHOUT TIME ZONE,
    customer_id    VARCHAR(32) NOT NULL,
    product_id     VARCHAR(32) NOT NULL,
    price          NUMERIC     NOT NULL,
    web_user_id    VARCHAR(32) NOT NULL,
    latest         BOOLEAN DEFAULT FALSE,
    transaction_id VARCHAR(32),
    FOREIGN KEY (web_user_id) REFERENCES public.web_user (id),
    FOREIGN KEY (customer_id) REFERENCES public.customer (id),
    FOREIGN KEY (product_id) REFERENCES public.product (id),
    FOREIGN KEY (transaction_id) REFERENCES public.transaction (id)
);
CREATE TABLE public.SEQUENCE
(
    id         VARCHAR(128) PRIMARY KEY,
    next_value INTEGER NOT NULL DEFAULT 0
);
CREATE TABLE public.kontrabon
(
    id           VARCHAR(32) PRIMARY KEY,
    code         VARCHAR(128) NOT NULL,
    customer_id  VARCHAR(32)  NOT NULL,
    created_time TIMESTAMP WITHOUT TIME ZONE,
    status       VARCHAR(128) NOT NULL,
    FOREIGN KEY (customer_id) REFERENCES public.customer (id)
);
CREATE TABLE public.kontrabon_transaction
(
    kontrabon_id   VARCHAR(32),
    transaction_id VARCHAR(32),
    FOREIGN KEY (kontrabon_id) REFERENCES public.kontrabon (id),
    FOREIGN KEY (transaction_id) REFERENCES public.transaction (id),
    PRIMARY KEY (kontrabon_id, transaction_id)
);
CREATE TABLE public.price_template
(
    id         VARCHAR(32) PRIMARY KEY,
    name       VARCHAR(256) NOT NULL
);
CREATE TABLE public.price_template_detail
(
    id                     VARCHAR(32) PRIMARY KEY,
    price_template_id VARCHAR(32) NOT NULL,
    product_id             VARCHAR(32) NOT NULL,
    price                  NUMERIC     NOT NULL,
    FOREIGN KEY (product_id) REFERENCES public.product (id),
    FOREIGN KEY (price_template_id) REFERENCES public.price_template (id)
);

ALTER TABLE price_template
    ADD COLUMN applied_to VARCHAR (512);

ALTER TABLE price_template_detail
    ADD COLUMN checked boolean;

ALTER TABLE transaction_detail
    ADD COLUMN sorting_val SMALLINT;