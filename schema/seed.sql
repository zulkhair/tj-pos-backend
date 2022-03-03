-------------------------------------
---------- DATABASE SEED ------------
-------------------------------------

--------------- MENU ----------------

INSERT INTO menu(id, name, seq_order, path, icon)
VALUES ('web:user', 'User', 0, 'user', 'fas fa-users'),
       ('web:role', 'Role', 1, 'role', 'fas fa-sitemap'),
       ('web:masterdata', 'Master Data', 2, 'master', 'fas fa-database')
;

INSERT INTO sub_menu(id, menu_id, name, seq_order, outcome, icon)
VALUES ('web:user:createUser', 'web:user', 'Tambah Data', 0, '/user/register-user.html', 'fas fa-plus'),
       ('web:user:editUser', 'web:user', 'Ubah Data', 1, '/user/edit-user.html', 'fas fa-pen'),
       ('web:role:createRole', 'web:role', 'Tambah Data', 0, '/role/create-role.html', 'fas fa-plus'),
       ('web:role:editRole', 'web:role', 'Ubah Data', 1, '/role/edit-role.html', 'fas fa-pen'),
       ('web:masterdata:product', 'web:masterdata', 'Produk', 0, '/master/product.html', 'fas fa-seedling')
;

------------ PERMISSION -------------
INSERT INTO permission(id, sub_menu_id, name, seq_order, apis)
VALUES ('web:user:createUser', 'web:user:createUser', 'Registrasi', 0, '/api/role/active-list;/api/user/register-user'),
       ('web:user:editUser', 'web:user:editUser', 'Ubah Data', 1,
        '/api/role/active-list;/api/user/find-all;/api/user/force-change-password;/api/user/change-status'),
       ('web:role:createRole', 'web:role:createRole', 'Tambah Data', 0, '/api/role/permissions;/api/role/create'),
       ('web:role:editRole', 'web:role:editRole', 'Ubah Data', 1, '/api/role/find-all;/api/role/edit'),
       ('web:masterdata:product:add', 'web:masterdata:product', 'Tambah Data Produk', 0, '/api/product/find;/api/product/create'),
       ('web:masterdata:product:view', 'web:masterdata:product', 'Lihat Data Produk', 1, '/api/product/find'),
       ('web:masterdata:product:edit', 'web:masterdata:product', 'Perbarui Data Produk', 2, '/api/product/find;/api/product/edit')
;

----------------- ROLE ---------------

INSERT INTO role(id, active, name)
VALUES ('735c7b8b96a8463c8493037d4c8ff085', true, 'Super Admin')
;

----------- ROLE_PERMISSION ----------

INSERT INTO role_permission(role_id, permission_id)
VALUES ('735c7b8b96a8463c8493037d4c8ff085', 'web:user:createUser'),
       ('735c7b8b96a8463c8493037d4c8ff085', 'web:user:editUser'),
       ('735c7b8b96a8463c8493037d4c8ff085', 'web:role:createRole'),
       ('735c7b8b96a8463c8493037d4c8ff085', 'web:role:editRole')
;

----------------- USER ---------------
INSERT INTO web_user(id, username, name, password_hash, password_salt, email, role_id, active, registration_timestamp,
                     created_by)
VALUES ('b11dd364dd714d5a8279123426bb92e5', 'super', 'Super Admin',
        'ff6b81d4a3803f8e5863c0d3dd9cdcab7d2bffebef17079ea832a77169912ac5', '0238daa119ab995155f346bf52cdb727',
        'email@gmail.com', '735c7b8b96a8463c8493037d4c8ff085', true, now(), null);

INSERT INTO public.config(id, value)
VALUES ('LOGIN_URL', 'http://localhost/login.html'),
       ('FORBIDDEN_URL', 'http://localhost/forbidden.html'),
       ('UNAUTHORIZED_URL', 'http://localhost/unauthorized.html'),
       ('SESSION_TIMEOUT_MINUTE', '30')
;