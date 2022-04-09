-------------------------------------
---------- DATABASE SEED ------------
-------------------------------------

--------------- MENU ----------------

INSERT INTO menu(id, name, seq_order, path, icon)
VALUES ('web:user', 'User', 0, 'user', 'fas fa-users'),
       ('web:role', 'Role', 1, 'role', 'fas fa-sitemap'),
       ('web:masterdata', 'Master Data', 2, 'master', 'fas fa-database'),
       ('web:price', 'Harga', 3, 'price', 'fas fa-rupiah-sign'),
       ('web:transaction', 'Transaksi', 4, 'transaction', 'fas fa-money-bill'),
       ('web:report', 'Laporan', 5, 'report', 'fas fa-book')
;

INSERT INTO sub_menu(id, menu_id, name, seq_order, outcome, icon)
VALUES ('web:user:createUser', 'web:user', 'Tambah Data', 0, '/user/register-user.html', 'fas fa-plus'),
       ('web:user:editUser', 'web:user', 'Ubah Data', 1, '/user/edit-user.html', 'fas fa-pen'),
       ('web:role:createRole', 'web:role', 'Tambah Data', 0, '/role/create-role.html', 'fas fa-plus'),
       ('web:role:editRole', 'web:role', 'Ubah Data', 1, '/role/edit-role.html', 'fas fa-pen'),
       ('web:masterdata:product', 'web:masterdata', 'Produk', 0, '/master/product.html', 'fas fa-seedling'),
       ('web:masterdata:supplier', 'web:masterdata', 'Supplier', 1, '/master/supplier.html', 'fas fa-people-carry-box'),
       ('web:masterdata:customer', 'web:masterdata', 'Customer', 2, '/master/customer.html', 'fas fa-user-tag'),
       ('web:price:buy', 'web:price', 'Harga Beli', 0, '/price/buy.html', 'fas fa-cart-shopping'),
       ('web:price:sell', 'web:price', 'Harga Jual', 1, '/price/sell.html', 'fas fa-cash-register'),
       ('web:transaction:sell', 'web:transaction', 'Penjualan', 0, '/transaction/sell.html', 'fas fa-cash-register'),
       ('web:transaction:buy', 'web:transaction', 'Pembelian', 1, '/transaction/buy.html', 'fas fa-bag-shopping'),
       ('web:report:view', 'web:report', 'Laporan Transaksi', 1, '/report/transaction.html', 'fas fa-coins')
;

------------ PERMISSION -------------
INSERT INTO permission(id, sub_menu_id, name, seq_order, apis)
VALUES ('web:user:createUser', 'web:user:createUser', 'Registrasi', 0, '/api/role/active-list;/api/user/register-user'),
       ('web:user:editUser', 'web:user:editUser', 'Ubah Data', 1,
        '/api/role/active-list;/api/user/find-all;/api/user/force-change-password;/api/user/change-status'),
       ('web:role:createRole', 'web:role:createRole', 'Tambah Data', 0, '/api/role/permissions;/api/role/create'),
       ('web:role:editRole', 'web:role:editRole', 'Ubah Data', 1, '/api/role/find-all;/api/role/edit'),
       ('web:masterdata:product:add', 'web:masterdata:product', 'Tambah Data Produk', 0, '/api/product/find;/api/product/create;/api/auth/check'),
       ('web:masterdata:product:view', 'web:masterdata:product', 'Lihat Data Produk', 1, '/api/product/find;/api/auth/check'),
       ('web:masterdata:product:edit', 'web:masterdata:product', 'Perbarui Data Produk', 2, '/api/product/find;/api/product/edit;/api/auth/check'),
       ('web:masterdata:supplier:add', 'web:masterdata:supplier', 'Tambah Data Supplier', 0, '/api/supplier/find;/api/supplier/create;/api/auth/check'),
       ('web:masterdata:supplier:view', 'web:masterdata:supplier', 'Lihat Data Supplier', 1, '/api/supplier/find;/api/auth/check'),
       ('web:masterdata:supplier:edit', 'web:masterdata:supplier', 'Perbarui Data Supplier', 2, '/api/supplier/find;/api/supplier/edit;/api/auth/check'),
       ('web:masterdata:customer:add', 'web:masterdata:customer', 'Tambah Data Customer', 0, '/api/customer/find;/api/customer/create;/api/auth/check'),
       ('web:masterdata:customer:view', 'web:masterdata:customer', 'Lihat Data Customer', 1, '/api/customer/find;/api/auth/check'),
       ('web:masterdata:customer:edit', 'web:masterdata:customer', 'Perbarui Data Customer', 2, '/api/customer/find;/api/customer/edit;/api/auth/check'),
       ('web:price:buy:manage', 'web:price:buy', 'Kelola Harga Beli', 0, '/api/unit/find;/api/unit/edit;/api/unit/create;/api/supplier/find;/api/product/find;/api/supplier/buy-price;/api/supplier/update-buy-price'),
       ('web:price:sell:manage', 'web:price:sell', 'Kelola Harga Jual', 0, '/api/unit/find;/api/unit/edit;/api/unit/create;/api/supplier/find;/api/product/find;/api/customer/sell-price;/api/supplier/update-sell-price'),
       ('web:transaction:sell:add', 'web:transaction:sell', 'Penjualan', 0, '/api/transaction/create;/api/unit/find;/api/customer/find;/api/product/find;/api/customer/sell-price')
       ('web:transaction:buy:add', 'web:transaction:buy', 'Pembelian', 0, '/api/transaction/create;/api/unit/find;/api/supplier/find;/api/product/find;/api/customer/buy-price')
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
       ('735c7b8b96a8463c8493037d4c8ff085', 'web:role:editRole'),
       ('735c7b8b96a8463c8493037d4c8ff085', 'web:masterdata:product:add'),
       ('735c7b8b96a8463c8493037d4c8ff085', 'web:masterdata:product:view'),
       ('735c7b8b96a8463c8493037d4c8ff085', 'web:masterdata:product:edit'),
       ('735c7b8b96a8463c8493037d4c8ff085', 'web:masterdata:supplier:add'),
       ('735c7b8b96a8463c8493037d4c8ff085', 'web:masterdata:supplier:view'),
       ('735c7b8b96a8463c8493037d4c8ff085', 'web:masterdata:supplier:edit'),
       ('735c7b8b96a8463c8493037d4c8ff085', 'web:masterdata:customer:add'),
       ('735c7b8b96a8463c8493037d4c8ff085', 'web:masterdata:customer:view'),
       ('735c7b8b96a8463c8493037d4c8ff085', 'web:masterdata:customer:edit'),
       ('735c7b8b96a8463c8493037d4c8ff085', 'web:price:buy:manage'),
       ('735c7b8b96a8463c8493037d4c8ff085', 'web:price:sell:manage')
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