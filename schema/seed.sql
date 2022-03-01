-------------------------------------
---------- DATABASE SEED ------------
-------------------------------------

--------------- MENU ----------------

INSERT INTO menu(id, name, menu_order, menu_path, icon)
VALUES ('web:user', 'User', '0', '/user/', 'fa fa-users'),
       ('web:role', 'Role', '1', '/role/', 'fa fa-sitemap')
;

------------ PERMISSION -------------

INSERT INTO permission(id, menu_id, name, permission_order, outcome, paths, icon)
VALUES ('web:user:createUser', 'web:user', 'Registrasi', 0, '/user/register-user.html',
        '/api/role/active-list;/api/user/register-user', 'fas fa-plus'),
       ('web:user:editUser', 'web:user', 'Ubah Data', 1, '/user/edit-user.html',
        '/api/role/active-list;/api/user/find-all;/api/user/force-change-password;/api/user/change-status', 'fas fa-pen'),
       ('web:role:createRole', 'web:role', 'Tambah Data', 0, '/role/create-role.html',
        '/api/role/permissions;/api/role/create', 'fas fa-plus'),
       ('web:role:editRole', 'web:role', 'Ubah Data', 1, '/role/edit-role.html',
        '/api/role/find-all;/api/role/edit', 'fas fa-pen')
;

----------------- ROLE ---------------

INSERT INTO role(id, active, name)
VALUES ('735c7b8b96a8463c8493037d4c8ff085', true, 'Super Admin')
;

----------- MENU_PERMISSION ----------

INSERT INTO menu_permission (menu_id, permission_id)
VALUES ('web:user', 'web:user:createUser'),
       ('web:user', 'web:user:editUser'),
       ('web:user', 'web:user:viewUser'),
       ('web:role', 'web:role:createRole'),
       ('web:role', 'web:role:editRole'),
       ('web:role', 'web:role:viewRole')
;

----------- ROLE_PERMISSION ----------

INSERT INTO role_permission(role_id, permission_id)
VALUES ('735c7b8b96a8463c8493037d4c8ff085', 'web:user:createUser'),
       ('735c7b8b96a8463c8493037d4c8ff085', 'web:user:editUser'),
       ('735c7b8b96a8463c8493037d4c8ff085', 'web:user:viewUser'),
       ('735c7b8b96a8463c8493037d4c8ff085', 'web:role:createRole'),
       ('735c7b8b96a8463c8493037d4c8ff085', 'web:role:editRole'),
       ('735c7b8b96a8463c8493037d4c8ff085', 'web:role:viewRole')
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