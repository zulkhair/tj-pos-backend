-------------------------------------
---------- DATABASE SEED ------------
-------------------------------------

--------------- MENU ----------------

INSERT INTO menu(id, name, seq_order, path, icon)
VALUES ('web:user', 'User', 0, 'user', 'fas fa-users'),
       ('web:role', 'Role', 1, 'role', 'fas fa-sitemap'),
       ('web:masterdata', 'Master Data', 2, 'master;price', 'fas fa-database'),
--        ('web:price', 'Harga', 3, 'price', 'fas fa-rupiah-sign'),
       ('web:transaction', 'Transaksi', 4, 'transaction;price', 'fas fa-money-bill')
--        ('web:report', 'Laporan', 5, 'report', 'fas fa-book')
;

INSERT INTO sub_menu(id, menu_id, name, seq_order, outcome, icon)
VALUES ('web:user:createUser', 'web:user', 'Tambah Data', 0, '/user/register-user.html', 'fas fa-plus'),
       ('web:user:editUser', 'web:user', 'Ubah Data', 1, '/user/edit-user.html', 'fas fa-pen'),
       ('web:role:createRole', 'web:role', 'Tambah Data', 0, '/role/create-role.html', 'fas fa-plus'),
       ('web:role:editRole', 'web:role', 'Ubah Data', 1, '/role/edit-role.html', 'fas fa-pen'),
       ('web:masterdata:product', 'web:masterdata', 'Produk', 0, '/master/product.html', 'fas fa-seedling'),
       ('web:masterdata:supplier', 'web:masterdata', 'Supplier', 1, '/master/supplier.html', 'fas fa-people-carry-box'),
       ('web:masterdata:customer', 'web:masterdata', 'Customer', 2, '/master/customer.html', 'fas fa-user-tag'),
       ('web:price:buy', 'web:masterdata', 'Harga Beli', 3, '/price/buy.html', 'fas fa-cart-shopping'),
       ('web:price:sell', 'web:masterdata', 'Harga Jual', 4, '/price/sell.html', 'fas fa-cash-register'),
       ('web:price:template', 'web:masterdata', 'Template Harga', 4, '/price/template.html', 'fas fa-cash-register'),
       ('web:transaction:sell', 'web:transaction', 'Penjualan', 2, '/transaction/sell.html', 'fas fa-money-check'),
       ('web:transaction:status', 'web:transaction', 'Status', 3, '/transaction/status.html', 'fas fa-clipboard-list'),
       ('web:transaction:kontrabon', 'web:transaction', 'Kontrabon', 3, '/transaction/kontrabon.html', 'fas fa-clipboard-list'),
       ('web:transaction:buy', 'web:transaction', 'Pembelian', 1, '/transaction/buy.html', 'fas fa-bag-shopping'),
       ('web:transaction:report', 'web:transaction', 'Laporan Transaksi', 4, '/transaction/report.html', 'fas fa-coins')
;

------------ PERMISSION -------------
INSERT INTO permission(id, sub_menu_id, name, seq_order, apis)
VALUES ('web:user:createUser', 'web:user:createUser', 'Registrasi', 0, '/api/role/active-list;/api/user/register-user'),
       ('web:user:editUser', 'web:user:editUser', 'Ubah Data', 1,
        '/api/role/active-list;/api/user/find-all;/api/user/force-change-password;/api/user/change-status;/api/user/edit-user'),
       ('web:role:createRole', 'web:role:createRole', 'Tambah Data', 0, '/api/role/permissions;/api/role/create'),
       ('web:role:editRole', 'web:role:editRole', 'Ubah Data', 1, '/api/role/find-all;/api/role/edit'),
       ('web:masterdata:product:add', 'web:masterdata:product', 'Tambah Data Produk', 0, '/api/product/find;/api/product/create;/api/auth/check;/api/unit/findActive'),
       ('web:masterdata:product:view', 'web:masterdata:product', 'Lihat Data Produk', 1, '/api/product/find;/api/auth/check'),
       ('web:masterdata:product:edit', 'web:masterdata:product', 'Perbarui Data Produk', 2, '/api/product/find;/api/product/edit;/api/auth/check;/api/unit/findActive'),
       ('web:masterdata:supplier:add', 'web:masterdata:supplier', 'Tambah Data Supplier', 0, '/api/supplier/find;/api/supplier/create;/api/auth/check'),
       ('web:masterdata:supplier:view', 'web:masterdata:supplier', 'Lihat Data Supplier', 1, '/api/supplier/find;/api/auth/check'),
       ('web:masterdata:supplier:edit', 'web:masterdata:supplier', 'Perbarui Data Supplier', 2, '/api/supplier/find;/api/supplier/edit;/api/auth/check'),
       ('web:masterdata:customer:add', 'web:masterdata:customer', 'Tambah Data Customer', 0, '/api/customer/find;/api/customer/create;/api/auth/check'),
       ('web:masterdata:customer:view', 'web:masterdata:customer', 'Lihat Data Customer', 1, '/api/customer/find;/api/auth/check'),
       ('web:masterdata:customer:edit', 'web:masterdata:customer', 'Perbarui Data Customer', 2, '/api/customer/find;/api/customer/edit;/api/auth/check'),
       ('web:price:buy:manage', 'web:price:buy', 'Kelola Harga Beli', 0, '/api/unit/find;/api/unit/edit;/api/unit/create;/api/supplier/find;/api/product/findActive;/api/supplier/add-price;/api/supplier/find-latest-price;/api/supplier/find-price'),
       ('web:price:sell:manage', 'web:price:sell', 'Kelola Harga Jual', 1, '/api/unit/find;/api/unit/edit;/api/unit/create;/api/customer/find;/api/product/findActive;/api/customer/add-price;/api/customer/find-latest-price;/api/customer/find-price'),
       ('web:price:template:manage', 'web:price:template', 'Template Harga', 3, '/api/customer/findActive;/api/product/findActive;/api/price/template/find;/api/price/template/create;/api/price/template/edit-price;/api/price/template/apply;/api/price/template/findDetail;/api/price/template/delete;/api/price/template/copy;/api/price/template/download'),
       ('web:transaction:sell:add', 'web:transaction:sell', 'Penjualan', 2, '/api/transaction/create;/api/unit/find;/api/customer/find;/api/product/findActive;/api/customer/sell-price'),
       ('web:transaction:kontrabon:manage', 'web:transaction:kontrabon', 'Kelola Kontrabon', 3, '/api/kontrabon/find;/api/kontrabon/findTransaction;/api/kontrabon/create;/api/kontrabon/add;/api/kontrabon/remove;/api/kontrabon/update-lunas'),
       ('web:transaction:kontrabon:view', 'web:transaction:kontrabon', 'Lihat Kontrabon', 4, '/api/kontrabon/find'),
       ('web:transaction:status:viewstatus', 'web:transaction:status', 'Lihat Status Penjualan', 0, '/api/transaction/find'),
       ('web:transaction:status:managestatus', 'web:transaction:status', 'Perbarui Status Penjualan', 1, '/api/transaction/find;/api/transaction/updateStatus;/api/transaction/updateBuyPrice;/api/transaction/cancelTrx;/api/transaction/update'),
       ('web:transaction:report:view', 'web:transaction:report', 'Laporan', 0, '/api/transaction/find;/api/transaction/report'),
       ('web:transaction:report:updatebuyprice', 'web:transaction:report', 'Ubah Harga Beli', 0, '/api/transaction/find;/api/transaction/report;/api/transaction/updateHargaBeli'),
       ('web:transaction:buy:add', 'web:transaction:buy', 'Pembelian', 0, '/api/transaction/insertTransactionBuy;/api/unit/find;/api/supplier/find;/api/product/find;')
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
       ('735c7b8b96a8463c8493037d4c8ff085', 'web:masterdata:customer:add'),
       ('735c7b8b96a8463c8493037d4c8ff085', 'web:masterdata:customer:view'),
       ('735c7b8b96a8463c8493037d4c8ff085', 'web:masterdata:customer:edit'),
       ('735c7b8b96a8463c8493037d4c8ff085', 'web:price:buy:manage'),
       ('735c7b8b96a8463c8493037d4c8ff085', 'web:price:sell:manage'),
       ('735c7b8b96a8463c8493037d4c8ff085', 'web:transaction:sell:add'),
       ('735c7b8b96a8463c8493037d4c8ff085', 'web:transaction:status:viewstatus'),
       ('735c7b8b96a8463c8493037d4c8ff085', 'web:transaction:status:managestatus')
;

----------------- USER ---------------
INSERT INTO web_user(id, username, name, password_hash, password_salt, email, role_id, active, registration_timestamp,
                     created_by)
VALUES ('b11dd364dd714d5a8279123426bb92e5', 'super', 'Super Admin',
        'ZuuV7rk4-Nr9PJ_TXqlUpxMoLfKA4WIcm5enyhdJlSE=', '0238daa119ab995155f346bf52cdb727',
        'email@gmail.com', '735c7b8b96a8463c8493037d4c8ff085', true, now(), null);

INSERT INTO public.config(id, value)
VALUES ('LOGIN_URL', 'http://localhost/login.html'),
       ('FORBIDDEN_URL', 'http://localhost/forbidden.html'),
       ('UNAUTHORIZED_URL', 'http://localhost/unauthorized.html'),
       ('SESSION_TIMEOUT_MINUTE', '30')
;

INSERT INTO public.unit(id, code, description, active) VALUES ('KG', 'Kg', 'Kilogram', true);
INSERT INTO public.unit(id, code, description, active) VALUES ('PACK', 'Pack', 'Pack', true);
INSERT INTO public.unit(id, code, description, active) VALUES ('PCS', 'Pcs', 'Pcs', true);
INSERT INTO public.unit(id, code, description, active) VALUES ('BOX', 'Box', 'Box', true);
INSERT INTO public.unit(id, code, description, active) VALUES ('PAPAN', 'Papan', 'Papan', true);
INSERT INTO public.unit(id, code, description, active) VALUES ('LTR', 'Ltr', 'Ltr', true);
INSERT INTO public.unit(id, code, description, active) VALUES ('SACHET', 'Sachet', 'Sachet', true);
INSERT INTO public.unit(id, code, description, active) VALUES ('CAN', 'Can', 'Can', true);

INSERT INTO public.product(id, unit_id, code, name, description, active)VALUES ('1','KG','Anggur Import Hijau','Anggur Import Hijau','Anggur Import Hijau',TRUE);
INSERT INTO public.product(id, unit_id, code, name, description, active)VALUES ('2','KG','Anggur Import Merah','Anggur Import Merah','Anggur Import Merah',TRUE);
INSERT INTO public.product(id, unit_id, code, name, description, active)VALUES ('3','KG','Apel Green Smith','Apel Green Smith','Apel Green Smith',TRUE);
INSERT INTO public.product(id, unit_id, code, name, description, active)VALUES ('4','KG','Apel Malang','Apel Malang','Apel Malang',TRUE);
INSERT INTO public.product(id, unit_id, code, name, description, active)VALUES ('5','KG','Apel Merah USA','Apel Merah USA','Apel Merah USA',TRUE);
INSERT INTO public.product(id, unit_id, code, name, description, active)VALUES ('6','KG','Apel Fuji Import','Apel Fuji Import','Apel Fuji Import',TRUE);
INSERT INTO public.product(id, unit_id, code, name, description, active)VALUES ('7','KG','Avocado','Avocado','Avocado',TRUE);
INSERT INTO public.product(id, unit_id, code, name, description, active)VALUES ('8','KG','Belimbing','Belimbing','Belimbing',TRUE);
INSERT INTO public.product(id, unit_id, code, name, description, active)VALUES ('9','KG','Bengkuang','Bengkuang','Bengkuang',TRUE);
INSERT INTO public.product(id, unit_id, code, name, description, active)VALUES ('10','PACK','Blueberry','Blueberry','Blueberry',TRUE);
INSERT INTO public.product(id, unit_id, code, name, description, active)VALUES ('11','PACK','Cecenet','Cecenet','Cecenet',TRUE);
INSERT INTO public.product(id, unit_id, code, name, description, active)VALUES ('12','PACK','Cendol Elizabeth','Cendol Elizabeth','Cendol Elizabeth',TRUE);
INSERT INTO public.product(id, unit_id, code, name, description, active)VALUES ('13','PACK','Cincau hitam','Cincau hitam','Cincau hitam',TRUE);
INSERT INTO public.product(id, unit_id, code, name, description, active)VALUES ('14','PACK','Cincau hijau','Cincau hijau','Cincau hijau',TRUE);
INSERT INTO public.product(id, unit_id, code, name, description, active)VALUES ('15','KG','Duku','Duku','Duku',TRUE);
INSERT INTO public.product(id, unit_id, code, name, description, active)VALUES ('16','KG','Durian montong utuh','Durian montong utuh','Durian montong utuh',TRUE);
INSERT INTO public.product(id, unit_id, code, name, description, active)VALUES ('17','PACK','Durian lokal utuh beku','Durian lokal utuh beku','Durian lokal utuh beku',TRUE);
INSERT INTO public.product(id, unit_id, code, name, description, active)VALUES ('18','KG','Durian frozen montong @ 400gr','Durian frozen montong @ 400gr','Durian frozen montong @ 400gr',TRUE);
INSERT INTO public.product(id, unit_id, code, name, description, active)VALUES ('19','KG','Durian musang king @ 400 gr','Durian musang king @ 400 gr','Durian musang king @ 400 gr',TRUE);
INSERT INTO public.product(id, unit_id, code, name, description, active)VALUES ('20','PACK','Fruit Blackberry','Fruit Blackberry','Fruit Blackberry',TRUE);
INSERT INTO public.product(id, unit_id, code, name, description, active)VALUES ('21','KG','Jambu Air Cincalo','Jambu Air Cincalo','Jambu Air Cincalo',TRUE);
INSERT INTO public.product(id, unit_id, code, name, description, active)VALUES ('22','KG','Jambu Air Biasa','Jambu Air Biasa','Jambu Air Biasa',TRUE);
INSERT INTO public.product(id, unit_id, code, name, description, active)VALUES ('23','KG','Jambu Merah','Jambu Merah','Jambu Merah',TRUE);
INSERT INTO public.product(id, unit_id, code, name, description, active)VALUES ('24','KG','Jeruk bali / pomelo ','Jeruk bali / pomelo ','Jeruk bali / pomelo ',TRUE);
INSERT INTO public.product(id, unit_id, code, name, description, active)VALUES ('25','KG','Jeruk Mandarin','Jeruk Mandarin','Jeruk Mandarin',TRUE);
INSERT INTO public.product(id, unit_id, code, name, description, active)VALUES ('26','KG','Jeruk Medan','Jeruk Medan','Jeruk Medan',TRUE);
INSERT INTO public.product(id, unit_id, code, name, description, active)VALUES ('27','KG','Jeruk Medan Super ','Jeruk Medan Super ','Jeruk Medan Super ',TRUE);
INSERT INTO public.product(id, unit_id, code, name, description, active)VALUES ('28','KG','Jeruk Pontianak / Peres','Jeruk Pontianak / Peres','Jeruk Pontianak / Peres',TRUE);
INSERT INTO public.product(id, unit_id, code, name, description, active)VALUES ('29','KG','Jeruk Sunkist','Jeruk Sunkist','Jeruk Sunkist',TRUE);
INSERT INTO public.product(id, unit_id, code, name, description, active)VALUES ('30','KG','Kedondong','Kedondong','Kedondong',TRUE);
INSERT INTO public.product(id, unit_id, code, name, description, active)VALUES ('31','KG','Kelapa Muda Kerok','Kelapa Muda Kerok','Kelapa Muda Kerok',TRUE);
INSERT INTO public.product(id, unit_id, code, name, description, active)VALUES ('32','KG','Kelapa Muda Utuh','Kelapa Muda Utuh','Kelapa Muda Utuh',TRUE);
INSERT INTO public.product(id, unit_id, code, name, description, active)VALUES ('33','KG','Kiwi','Kiwi','Kiwi',TRUE);
INSERT INTO public.product(id, unit_id, code, name, description, active)VALUES ('34','KG','Kiwi zespri','Kiwi zespri','Kiwi zespri',TRUE);
INSERT INTO public.product(id, unit_id, code, name, description, active)VALUES ('35','KG','Kolang Kaling','Kolang Kaling','Kolang Kaling',TRUE);
INSERT INTO public.product(id, unit_id, code, name, description, active)VALUES ('36','KG','Konyal','Konyal','Konyal',TRUE);
INSERT INTO public.product(id, unit_id, code, name, description, active)VALUES ('37','KG','Kurma Palm bertangkai','Kurma Palm bertangkai','Kurma Palm bertangkai',TRUE);
INSERT INTO public.product(id, unit_id, code, name, description, active)VALUES ('38','KG','Lemon Import','Lemon Import','Lemon Import',TRUE);
INSERT INTO public.product(id, unit_id, code, name, description, active)VALUES ('39','KG','Lemon Lokal','Lemon Lokal','Lemon Lokal',TRUE);
INSERT INTO public.product(id, unit_id, code, name, description, active)VALUES ('40','KG','Lemon lembang','Lemon lembang','Lemon lembang',TRUE);
INSERT INTO public.product(id, unit_id, code, name, description, active)VALUES ('41','KG','Lengkeng','Lengkeng','Lengkeng',TRUE);
INSERT INTO public.product(id, unit_id, code, name, description, active)VALUES ('42','KG','Mangga Gedong','Mangga Gedong','Mangga Gedong',TRUE);
INSERT INTO public.product(id, unit_id, code, name, description, active)VALUES ('43','KG','Mangga Gedong Gincu','Mangga Gedong Gincu','Mangga Gedong Gincu',TRUE);
INSERT INTO public.product(id, unit_id, code, name, description, active)VALUES ('44','KG','Mangga Harum Manis','Mangga Harum Manis','Mangga Harum Manis',TRUE);
INSERT INTO public.product(id, unit_id, code, name, description, active)VALUES ('45','KG','Mangga Kweni','Mangga Kweni','Mangga Kweni',TRUE);
INSERT INTO public.product(id, unit_id, code, name, description, active)VALUES ('46','KG','Mangga Muda cengkir','Mangga Muda cengkir','Mangga Muda cengkir',TRUE);
INSERT INTO public.product(id, unit_id, code, name, description, active)VALUES ('47','KG','Mangga muda biasa','Mangga muda biasa','Mangga muda biasa',TRUE);
INSERT INTO public.product(id, unit_id, code, name, description, active)VALUES ('48','KG','Manggis','Manggis','Manggis',TRUE);
INSERT INTO public.product(id, unit_id, code, name, description, active)VALUES ('49','KG','Melon Honey Dew','Melon Honey Dew','Melon Honey Dew',TRUE);
INSERT INTO public.product(id, unit_id, code, name, description, active)VALUES ('50','KG','Melon Merah / Melon Rock','Melon Merah / Melon Rock','Melon Merah / Melon Rock',TRUE);
INSERT INTO public.product(id, unit_id, code, name, description, active)VALUES ('51','KG','Naga Merah','Naga Merah','Naga Merah',TRUE);
INSERT INTO public.product(id, unit_id, code, name, description, active)VALUES ('52','PCS','Nanas Jawa','Nanas Jawa','Nanas Jawa',TRUE);
INSERT INTO public.product(id, unit_id, code, name, description, active)VALUES ('53','PCS','Nanas Madu','Nanas Madu','Nanas Madu',TRUE);
INSERT INTO public.product(id, unit_id, code, name, description, active)VALUES ('54','KG','Nangka Matang Kupas','Nangka Matang Kupas','Nangka Matang Kupas',TRUE);
INSERT INTO public.product(id, unit_id, code, name, description, active)VALUES ('55','KG','Pear century','Pear century','Pear century',TRUE);
INSERT INTO public.product(id, unit_id, code, name, description, active)VALUES ('56','KG','Pear Hijau Aust/ Pear Jambu','Pear Hijau Aust/ Pear Jambu','Pear Hijau Aust/ Pear Jambu',TRUE);
INSERT INTO public.product(id, unit_id, code, name, description, active)VALUES ('57','KG','Pear Xiang Lie','Pear Xiang Lie','Pear Xiang Lie',TRUE);
INSERT INTO public.product(id, unit_id, code, name, description, active)VALUES ('58','KG','Pepaya','Pepaya','Pepaya',TRUE);
INSERT INTO public.product(id, unit_id, code, name, description, active)VALUES ('59','KG','Pepaya muda','Pepaya muda','Pepaya muda',TRUE);
INSERT INTO public.product(id, unit_id, code, name, description, active)VALUES ('60','KG','Pisang Ambon ','Pisang Ambon ','Pisang Ambon ',TRUE);
INSERT INTO public.product(id, unit_id, code, name, description, active)VALUES ('61','KG','Pisang raja cere','Pisang raja cere','Pisang raja cere',TRUE);
INSERT INTO public.product(id, unit_id, code, name, description, active)VALUES ('62','KG','Pisang raja bulu','Pisang raja bulu','Pisang raja bulu',TRUE);
INSERT INTO public.product(id, unit_id, code, name, description, active)VALUES ('63','KG','Pisang Muli','Pisang Muli','Pisang Muli',TRUE);
INSERT INTO public.product(id, unit_id, code, name, description, active)VALUES ('64','KG','Pisang Tanduk','Pisang Tanduk','Pisang Tanduk',TRUE);
INSERT INTO public.product(id, unit_id, code, name, description, active)VALUES ('65','KG','Pisang Nangka','Pisang Nangka','Pisang Nangka',TRUE);
INSERT INTO public.product(id, unit_id, code, name, description, active)VALUES ('66','KG','Rambutan','Rambutan','Rambutan',TRUE);
INSERT INTO public.product(id, unit_id, code, name, description, active)VALUES ('67','PACK','Raspberry','Raspberry','Raspberry',TRUE);
INSERT INTO public.product(id, unit_id, code, name, description, active)VALUES ('68','KG','Salak Pondoh','Salak Pondoh','Salak Pondoh',TRUE);
INSERT INTO public.product(id, unit_id, code, name, description, active)VALUES ('69','KG','Sawo','Sawo','Sawo',TRUE);
INSERT INTO public.product(id, unit_id, code, name, description, active)VALUES ('70','KG','Semangka Kuning','Semangka Kuning','Semangka Kuning',TRUE);
INSERT INTO public.product(id, unit_id, code, name, description, active)VALUES ('71','KG','Semangka Merah','Semangka Merah','Semangka Merah',TRUE);
INSERT INTO public.product(id, unit_id, code, name, description, active)VALUES ('72','KG','Sirsak','Sirsak','Sirsak',TRUE);
INSERT INTO public.product(id, unit_id, code, name, description, active)VALUES ('73','PACK','Strawberry','Strawberry','Strawberry',TRUE);
INSERT INTO public.product(id, unit_id, code, name, description, active)VALUES ('74','KG','Strawberry','Strawberry','Strawberry',TRUE);
INSERT INTO public.product(id, unit_id, code, name, description, active)VALUES ('75','KG','Sukun mentah','Sukun mentah','Sukun mentah',TRUE);
INSERT INTO public.product(id, unit_id, code, name, description, active)VALUES ('76','KG','Timun Suri','Timun Suri','Timun Suri',TRUE);
INSERT INTO public.product(id, unit_id, code, name, description, active)VALUES ('77','KG','Ubi Cilembu','Ubi Cilembu','Ubi Cilembu',TRUE);
INSERT INTO public.product(id, unit_id, code, name, description, active)VALUES ('78','KG','Ubi cilembu Mentah','Ubi cilembu Mentah','Ubi cilembu Mentah',TRUE);
INSERT INTO public.product(id, unit_id, code, name, description, active)VALUES ('79','KG','Ubi Manis Ungu  ','Ubi Manis Ungu  ','Ubi Manis Ungu  ',TRUE);
INSERT INTO public.product(id, unit_id, code, name, description, active)VALUES ('80','KG','Ubi Manis Biasa ','Ubi Manis Biasa ','Ubi Manis Biasa ',TRUE);
INSERT INTO public.product(id, unit_id, code, name, description, active)VALUES ('81','KG','Tape singkong','Tape singkong','Tape singkong',TRUE);
INSERT INTO public.product(id, unit_id, code, name, description, active)VALUES ('82','KG','Tape ketan hitam','Tape ketan hitam','Tape ketan hitam',TRUE);

INSERT INTO public.customer(id, code, name, description, active)VALUES ('1', 'HOTEL HARRIS', 'HOTEL HARRIS', 'HOTEL HARRIS', TRUE);
INSERT INTO public.customer(id, code, name, description, active)VALUES ('2', 'HOTEL HILTON', 'HOTEL HILTON', 'HOTEL HILTON', TRUE);
INSERT INTO public.customer(id, code, name, description, active)VALUES ('3', 'HOTEL SAVOY HOMMAN', 'HOTEL SAVOY HOMMAN', 'HOTEL SAVOY HOMMAN', TRUE);
INSERT INTO public.customer(id, code, name, description, active)VALUES ('4', 'HOTEL ART DECO', 'HOTEL ART DECO', 'HOTEL ART DECO', TRUE);
INSERT INTO public.customer(id, code, name, description, active)VALUES ('5', 'HOTEL CALIFORNIA', 'HOTEL CALIFORNIA', 'HOTEL CALIFORNIA', TRUE);
INSERT INTO public.customer(id, code, name, description, active)VALUES ('6', 'HOTEL ARYA DUTA', 'HOTEL ARYA DUTA', 'HOTEL ARYA DUTA', TRUE);
INSERT INTO public.customer(id, code, name, description, active)VALUES ('7', 'RUMAH SAKIT BORROMEUS', 'RUMAH SAKIT BORROMEUS', 'RUMAH SAKIT BORROMEUS', TRUE);
INSERT INTO public.customer(id, code, name, description, active)VALUES ('8', 'RUMAH SAKIT IMMANUEL', 'RUMAH SAKIT IMMANUEL', 'RUMAH SAKIT IMMANUEL', TRUE);

