package transactionrepo

import (
	"database/sql"
	"dromatech/pos-backend/global"
	transactiondomain "dromatech/pos-backend/internal/domain/transaction"
	stringutil "dromatech/pos-backend/internal/util/string"
	"log"
	"time"
)

func (r *Repo) FindDana(userID string, date time.Time) (*transactiondomain.DanaInquiryResponse, error) {
	var dana transactiondomain.DanaInquiryResponse

	// find to database
	var ID sql.NullString
	var SaldoAwal sql.NullFloat64
	var DanaTambahan sql.NullFloat64
	query := "SELECT id, saldo_awal, dana_tambahan FROM dana WHERE web_user_id = ? AND date = ?"
	rows, err := global.DBCON.Raw(query, userID, date).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		rows.Scan(&ID, &SaldoAwal, &DanaTambahan)
	}

	if !ID.Valid {
		return nil, nil
	}

	dana.ID = ID.String
	dana.SaldoAwal = SaldoAwal.Float64
	dana.DanaTambahan = DanaTambahan.Float64

	// find dana masuk
	var IDDanaMasuk sql.NullString
	var Sender sql.NullString
	var Amount sql.NullFloat64
	var Status sql.NullString
	query = "SELECT dt.id, wu.name, dt.amount, dt.status FROM dana_transaction dt " +
		"JOIN web_user wu ON dt.sender = wu.id " +
		"WHERE dt.receiver = ? AND dt.date = ? AND dt.status IN (?, ?) ORDER BY dt.created_time DESC"
	rows2, err := global.DBCON.Raw(query, userID, date, transactiondomain.DanaStatusApproved, transactiondomain.DanaStatusPending).Rows()
	if err != nil {
		return nil, err
	}
	defer rows2.Close()

	var danaMasuk []transactiondomain.DanaTransactionInquiryResponse
	for rows2.Next() {
		rows2.Scan(&IDDanaMasuk, &Sender, &Amount, &Status)
		danaMasuk = append(danaMasuk, transactiondomain.DanaTransactionInquiryResponse{
			ID:          IDDanaMasuk.String,
			Name:        Sender.String,
			Amount:      Amount.Float64,
			Status:      transactiondomain.DanaStatus(Status.String),
			CreatedTime: time.Now(),
		})
	}

	dana.DanaMasuk = danaMasuk

	// find dana keluar
	var IDDanaKeluar sql.NullString
	var Receiver sql.NullString
	var AmountKeluar sql.NullFloat64
	var StatusKeluar sql.NullString
	query = "SELECT dt.id, wu.name, dt.amount, dt.status FROM dana_transaction dt " +
		"JOIN web_user wu ON dt.receiver = wu.id " +
		"WHERE dt.sender = ? AND dt.date = ? AND dt.status IN (?, ?) ORDER BY dt.created_time DESC"
	rows3, err := global.DBCON.Raw(query, userID, date, transactiondomain.DanaStatusApproved, transactiondomain.DanaStatusPending).Rows()
	if err != nil {
		return nil, err
	}
	defer rows3.Close()

	var danaKeluar []transactiondomain.DanaTransactionInquiryResponse
	for rows3.Next() {
		rows3.Scan(&IDDanaKeluar, &Receiver, &AmountKeluar, &StatusKeluar)
		danaKeluar = append(danaKeluar, transactiondomain.DanaTransactionInquiryResponse{
			ID:          IDDanaKeluar.String,
			Name:        Receiver.String,
			Amount:      AmountKeluar.Float64,
			Status:      transactiondomain.DanaStatus(StatusKeluar.String),
			CreatedTime: time.Now(),
		})
	}

	dana.DanaKeluar = danaKeluar

	return &dana, nil
}

func (r *Repo) FindDanaByID(id string) (*transactiondomain.Dana, error) {
	var dana transactiondomain.Dana

	var ID sql.NullString
	var Date sql.NullTime
	var WebUserID sql.NullString
	var SaldoAwal sql.NullFloat64
	var DanaTambahan sql.NullFloat64
	var CreatedTime sql.NullTime
	query := "SELECT id, date, web_user_id, saldo_awal, dana_tambahan, created_time FROM dana WHERE id = ?"
	rows, err := global.DBCON.Raw(query, id).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		rows.Scan(&ID, &Date, &WebUserID, &SaldoAwal, &DanaTambahan, &CreatedTime)
	}

	dana.ID = ID.String
	dana.Date = Date.Time
	dana.WebUserID = WebUserID.String
	dana.SaldoAwal = SaldoAwal.Float64
	dana.DanaTambahan = DanaTambahan.Float64
	dana.CreatedTime = CreatedTime.Time

	return &dana, nil
}

func (r *Repo) FindDanaTransaction(id string) (*transactiondomain.DanaTransaction, error) {
	var danaTransaction transactiondomain.DanaTransaction

	// find dana transaction
	var ID sql.NullString
	var Date sql.NullTime
	var Sender sql.NullString
	var Receiver sql.NullString
	var Amount sql.NullFloat64
	var Status sql.NullString
	var CreatedTime sql.NullTime
	query := "SELECT id, date, sender, receiver, amount, status, created_time FROM dana_transaction WHERE id = ? AND (status = ? OR status = ?)"
	rows, err := global.DBCON.Raw(query, id, transactiondomain.DanaStatusPending, transactiondomain.DanaStatusApproved).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		rows.Scan(&ID, &Date, &Sender, &Receiver, &Amount, &Status, &CreatedTime)
	}

	danaTransaction.ID = ID.String
	danaTransaction.Date = Date.Time
	danaTransaction.Sender = Sender.String
	danaTransaction.Receiver = Receiver.String
	danaTransaction.Amount = Amount.Float64
	danaTransaction.Status = transactiondomain.DanaStatus(Status.String)
	danaTransaction.CreatedTime = CreatedTime.Time
	return &danaTransaction, nil
}

func (r *Repo) CreateDana(userID string, request transactiondomain.DanaRequest) error {
	tx := global.DBCON.Begin()

	ID := stringutil.GenerateUUID()
	date, err := time.Parse("2006-01-02", request.Date)
	if err != nil {
		return err
	}
	tx.Exec("INSERT INTO public.dana(id, date, web_user_id, saldo_awal, dana_tambahan, created_time) VALUES (?, ?, ?, ?, ?, ?);",
		ID, date, userID, request.SaldoAwal, request.DanaTambahan, time.Now())

	if tx.Error != nil {
		tx.Rollback()
		return tx.Error
	}

	tx.Commit()

	return nil
}

func (r *Repo) UpdateDana(userID string, request transactiondomain.DanaRequest) error {
	tx := global.DBCON.Begin()

	tx.Exec("UPDATE public.dana SET saldo_awal = ?, dana_tambahan = ? WHERE id = ? AND web_user_id = ?;",
		request.SaldoAwal, request.DanaTambahan, request.ID, userID)

	if tx.Error != nil {
		tx.Rollback()
		return tx.Error
	}

	tx.Commit()

	return nil
}

func (r *Repo) SendDana(userID string, request transactiondomain.DanaTransactionRequest) error {
	tx := global.DBCON.Begin()

	ID := stringutil.GenerateUUID()
	date, err := time.Parse("2006-01-02", request.Date)
	if err != nil {
		return err
	}
	tx.Exec("INSERT INTO public.dana_transaction(id, date, sender, receiver, amount, status, created_time) VALUES (?, ?, ?, ?, ?, ?, ?);",
		ID, date, userID, request.Receiver, request.Amount, transactiondomain.DanaStatusPending, time.Now())

	if tx.Error != nil {
		tx.Rollback()
		return tx.Error
	}

	tx.Commit()

	return nil
}

func (r *Repo) ApproveDana(id string) error {
	tx := global.DBCON.Begin()

	tx.Exec("UPDATE public.dana_transaction SET status = ? WHERE id = ?;", transactiondomain.DanaStatusApproved, id)

	if tx.Error != nil {
		tx.Rollback()
		return tx.Error
	}

	tx.Commit()

	return nil
}

func (r *Repo) RejectDana(id string) error {
	tx := global.DBCON.Begin()

	tx.Exec("UPDATE public.dana_transaction SET status = ? WHERE id = ?;", transactiondomain.DanaStatusRejected, id)

	if tx.Error != nil {
		tx.Rollback()
		return tx.Error
	}

	tx.Commit()

	return nil
}

func (r *Repo) CancelSendDana(id string) error {
	tx := global.DBCON.Begin()

	tx.Exec("UPDATE public.dana_transaction SET status = ? WHERE id = ?;", transactiondomain.DanaStatusCanceled, id)

	if tx.Error != nil {
		tx.Rollback()
		return tx.Error
	}

	tx.Commit()

	return nil
}

func (r *Repo) CheckUserMobilePermission(id string) (bool, error) {
	var count sql.NullInt32

	query := "SELECT COUNT(wu.id) FROM web_user wu " +
		"JOIN role r ON wu.role_id = r.id " +
		"JOIN role_permission rp ON r.id = rp.role_id " +
		"JOIN permission p ON rp.permission_id = p.id " +
		"WHERE wu.id = ? AND p.id = 'mobile'"
	rows, err := global.DBCON.Raw(query, id).Rows()
	if err != nil {
		return false, err
	}
	defer rows.Close()

	for rows.Next() {
		rows.Scan(&count)
	}

	return count.Int32 > 0, nil
}

func (r *Repo) FindUserMobile(id string) ([]transactiondomain.WebUserMobile, error) {
	var ID sql.NullString
	var Name sql.NullString
	query := "SELECT wu.id, wu.name FROM web_user wu " +
		"JOIN role r ON wu.role_id = r.id " +
		"JOIN role_permission rp ON r.id = rp.role_id " +
		"JOIN permission p ON rp.permission_id = p.id " +
		"WHERE wu.id <> ? AND p.id = 'mobile'"
	rows, err := global.DBCON.Raw(query, id).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userMobile []transactiondomain.WebUserMobile
	for rows.Next() {
		rows.Scan(&ID, &Name)
		userMobile = append(userMobile, transactiondomain.WebUserMobile{ID: ID.String, Name: Name.String})
	}

	return userMobile, nil
}

func (r *Repo) CreatePenjualan(userID string, request transactiondomain.TrxCreateRequest) error {
	tx := global.DBCON.Begin()

	ID := stringutil.GenerateUUID()
	date, err := time.Parse("2006-01-02", request.Date)
	if err != nil {
		return err
	}

	tx.Exec("INSERT INTO public.penjualan_tunai(id, date, web_user_id, product_id, quantity, price, created_time) VALUES (?, ?, ?, ?, ?, ?, ?);",
		ID, date, userID, request.ProductID, request.Quantity, request.Price, time.Now())

	if tx.Error != nil {
		tx.Rollback()
		return tx.Error
	}

	tx.Commit()

	return nil
}

func (r *Repo) DeletePenjualan(id string) error {
	tx := global.DBCON.Begin()

	tx.Exec("DELETE FROM public.penjualan_tunai WHERE id = ?;", id)

	if tx.Error != nil {
		tx.Rollback()
		return tx.Error
	}

	tx.Commit()

	return nil
}

func (r *Repo) FindPenjualan(userID string, date time.Time) ([]transactiondomain.TrxInquiryResponse, error) {
	var ID sql.NullString
	var ProductName sql.NullString
	var Quantity sql.NullInt16
	var Price sql.NullFloat64

	log.Println(userID, date)

	query := "SELECT p.id, pr.name, p.quantity, p.price FROM penjualan_tunai p " +
		"JOIN product pr ON p.product_id = pr.id " +
		"WHERE p.web_user_id = ? AND p.date = ? ORDER BY p.created_time DESC"
	rows, err := global.DBCON.Raw(query, userID, date).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var penjualan []transactiondomain.TrxInquiryResponse
	for rows.Next() {
		rows.Scan(&ID, &ProductName, &Quantity, &Price)
		log.Println(ID, ProductName, Quantity, Price)
		penjualan = append(penjualan, transactiondomain.TrxInquiryResponse{
			ID:          ID.String,
			ProductName: ProductName.String,
			Quantity:    Quantity.Int16,
			Price:       Price.Float64,
		})
	}

	return penjualan, nil
}

func (r *Repo) CreateBelanja(userID string, request transactiondomain.TrxCreateRequest) error {
	tx := global.DBCON.Begin()

	ID := stringutil.GenerateUUID()
	date, err := time.Parse("2006-01-02", request.Date)
	if err != nil {
		return err
	}

	tx.Exec("INSERT INTO public.belanja(id, date, web_user_id, product_id, quantity, price, created_time) VALUES (?, ?, ?, ?, ?, ?, ?);",
		ID, date, userID, request.ProductID, request.Quantity, request.Price, time.Now())

	if tx.Error != nil {
		tx.Rollback()
		return tx.Error
	}

	tx.Commit()

	return nil
}

func (r *Repo) DeleteBelanja(id string) error {
	tx := global.DBCON.Begin()

	tx.Exec("DELETE FROM public.belanja WHERE id = ?;", id)

	if tx.Error != nil {
		tx.Rollback()
		return tx.Error
	}

	tx.Commit()

	return nil
}

func (r *Repo) FindBelanja(userID string, date time.Time) ([]transactiondomain.TrxInquiryResponse, error) {
	var ID sql.NullString
	var ProductName sql.NullString
	var Quantity sql.NullInt16
	var Price sql.NullFloat64

	query := "SELECT b.id, pr.name, b.quantity, b.price FROM belanja b " +
		"JOIN product pr ON b.product_id = pr.id " +
		"WHERE b.web_user_id = ? AND b.date = ? ORDER BY b.created_time DESC"
	rows, err := global.DBCON.Raw(query, userID, date).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var belanja []transactiondomain.TrxInquiryResponse
	for rows.Next() {
		rows.Scan(&ID, &ProductName, &Quantity, &Price)
		belanja = append(belanja, transactiondomain.TrxInquiryResponse{
			ID:          ID.String,
			ProductName: ProductName.String,
			Quantity:    Quantity.Int16,
			Price:       Price.Float64,
		})
	}

	return belanja, nil
}

func (r *Repo) CreateOperasional(userID string, request transactiondomain.TrxCreateOperasionalRequest) error {
	tx := global.DBCON.Begin()

	ID := stringutil.GenerateUUID()
	date, err := time.Parse("2006-01-02", request.Date)
	if err != nil {
		return err
	}

	tx.Exec("INSERT INTO public.operasional(id, date, web_user_id, description, quantity, price, created_time) VALUES (?, ?, ?, ?, ?, ?, ?);",
		ID, date, userID, request.Description, request.Quantity, request.Price, time.Now())

	if tx.Error != nil {
		tx.Rollback()
		return tx.Error
	}

	tx.Commit()

	return nil
}

func (r *Repo) FindOperasional(userID string, date time.Time) ([]transactiondomain.TrxInquiryOperasionalResponse, error) {
	var ID sql.NullString
	var Description sql.NullString
	var Quantity sql.NullInt16
	var Price sql.NullFloat64

	query := "SELECT o.id, o.description, o.quantity, o.price FROM operasional o WHERE o.web_user_id = ? AND o.date = ? ORDER BY o.created_time DESC"
	rows, err := global.DBCON.Raw(query, userID, date).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var operasional []transactiondomain.TrxInquiryOperasionalResponse
	for rows.Next() {
		rows.Scan(&ID, &Description, &Quantity, &Price)
		operasional = append(operasional, transactiondomain.TrxInquiryOperasionalResponse{
			ID:          ID.String,
			Description: Description.String,
			Quantity:    Quantity.Int16,
			Price:       Price.Float64,
		})
	}

	return operasional, nil
}

func (r *Repo) DeleteOperasional(id string) error {
	tx := global.DBCON.Begin()

	tx.Exec("DELETE FROM public.operasional WHERE id = ?;", id)

	if tx.Error != nil {
		tx.Rollback()
		return tx.Error
	}

	tx.Commit()

	return nil
}

func (r *Repo) FindSaldo(userID string, date time.Time) (*transactiondomain.SaldoResponse, error) {
	// find saldo awal and dana tambahan
	var saldoAwalVal sql.NullFloat64
	var danaTambahanVal sql.NullFloat64

	query := "SELECT saldo_awal, dana_tambahan FROM dana WHERE web_user_id = ? AND date = ?"
	rows, err := global.DBCON.Raw(query, userID, date).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		rows.Scan(&saldoAwalVal, &danaTambahanVal)
	}

	saldoAwal := 0.0
	if saldoAwalVal.Valid {
		saldoAwal = saldoAwalVal.Float64
	}

	danaTambahan := 0.0
	if danaTambahanVal.Valid {
		danaTambahan = danaTambahanVal.Float64
	}

	// find dana masuk
	query = "SELECT SUM(dt.amount) FROM dana_transaction dt " +
		"WHERE dt.receiver = ? AND dt.date = ? AND dt.status = ?"
	rows2, err := global.DBCON.Raw(query, userID, date, transactiondomain.DanaStatusApproved).Rows()
	if err != nil {
		return nil, err
	}
	defer rows2.Close()

	var danaMasukVal sql.NullFloat64
	for rows2.Next() {
		rows2.Scan(&danaMasukVal)
	}

	danaMasuk := 0.0
	if danaMasukVal.Valid {
		danaMasuk = danaMasukVal.Float64
	}

	// find dana keluar
	query = "SELECT SUM(dt.amount) FROM dana_transaction dt " +
		"WHERE dt.sender = ? AND dt.date = ? AND dt.status = ?"
	rows3, err := global.DBCON.Raw(query, userID, date, transactiondomain.DanaStatusApproved).Rows()
	if err != nil {
		return nil, err
	}
	defer rows3.Close()

	var danaKeluarVal sql.NullFloat64
	for rows3.Next() {
		rows3.Scan(&danaKeluarVal)
	}

	danaKeluar := 0.0
	if danaKeluarVal.Valid {
		danaKeluar = danaKeluarVal.Float64
	}

	// find penjualan
	query = "SELECT SUM(p.quantity * p.price) FROM penjualan_tunai p " +
		"WHERE p.web_user_id = ? AND p.date = ?"
	rows4, err := global.DBCON.Raw(query, userID, date).Rows()
	if err != nil {
		return nil, err
	}
	defer rows4.Close()

	var penjualanVal sql.NullFloat64
	for rows4.Next() {
		rows4.Scan(&penjualanVal)
	}

	penjualan := 0.0
	if penjualanVal.Valid {
		penjualan = penjualanVal.Float64
	}

	danaMasuk = danaMasuk + penjualan

	// find belanja
	query = "SELECT SUM(b.quantity * b.price) FROM belanja b " +
		"WHERE b.web_user_id = ? AND b.date = ?"
	rows5, err := global.DBCON.Raw(query, userID, date).Rows()
	if err != nil {
		return nil, err
	}
	defer rows5.Close()

	var belanjaVal sql.NullFloat64
	for rows5.Next() {
		rows5.Scan(&belanjaVal)
	}

	belanja := 0.0
	if belanjaVal.Valid {
		belanja = belanjaVal.Float64
	}

	// find operasional
	query = "SELECT SUM(o.quantity * o.price) FROM operasional o " +
		"WHERE o.web_user_id = ? AND o.date = ?"
	rows6, err := global.DBCON.Raw(query, userID, date).Rows()
	if err != nil {
		return nil, err
	}
	defer rows6.Close()

	var operasionalVal sql.NullFloat64
	for rows6.Next() {
		rows6.Scan(&operasionalVal)
	}

	operasional := 0.0
	if operasionalVal.Valid {
		operasional = operasionalVal.Float64
	}

	saldoAkhir := saldoAwal + danaTambahan + danaMasuk - belanja - operasional - danaKeluar

	return &transactiondomain.SaldoResponse{
		SaldoAwal:    saldoAwal,
		DanaTambahan: danaTambahan,
		DanaMasuk:    danaMasuk,
		Belanja:      belanja,
		Operasional:  operasional,
		DanaKeluar:   danaKeluar,
		SaldoAkhir:   saldoAkhir,
	}, nil
}

func (r *Repo) FindRekapitulasi(date time.Time) (*transactiondomain.RekapitulasiResponse, error) {
	// find all user mobile
	var ID sql.NullString
	var Name sql.NullString
	query := "SELECT wu.id, wu.name FROM web_user wu " +
		"JOIN role r ON wu.role_id = r.id " +
		"JOIN role_permission rp ON r.id = rp.role_id " +
		"JOIN permission p ON rp.permission_id = p.id " +
		"WHERE p.id = 'mobile'"
	rows, err := global.DBCON.Raw(query).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rekapitulasiList []transactiondomain.RekapitulasiDetail
	var totalBelanja float64
	var totalOperasional float64
	for rows.Next() {
		rows.Scan(&ID, &Name)

		// find belanja
		query = "SELECT SUM(b.quantity * b.price) FROM belanja b " +
			"WHERE b.web_user_id = ? AND b.date = ?"
		rows2, err := global.DBCON.Raw(query, ID.String, date).Rows()
		if err != nil {
			return nil, err
		}
		defer rows2.Close()

		var belanjaVal sql.NullFloat64
		for rows2.Next() {
			rows2.Scan(&belanjaVal)
		}

		belanja := 0.0
		if belanjaVal.Valid {
			belanja = belanjaVal.Float64
		}

		// find operasional
		query = "SELECT SUM(o.quantity * o.price) FROM operasional o " +
			"WHERE o.web_user_id = ? AND o.date = ?"
		rows3, err := global.DBCON.Raw(query, ID.String, date).Rows()
		if err != nil {
			return nil, err
		}
		defer rows3.Close()

		var operasionalVal sql.NullFloat64
		for rows3.Next() {
			rows3.Scan(&operasionalVal)
		}

		operasional := 0.0
		if operasionalVal.Valid {
			operasional = operasionalVal.Float64
		}

		rekapitulasiList = append(rekapitulasiList, transactiondomain.RekapitulasiDetail{
			Name:        Name.String,
			Belanja:     belanja,
			Operasional: operasional,
		})

		totalBelanja += belanja
		totalOperasional += operasional
	}

	return &transactiondomain.RekapitulasiResponse{
		Rekapitulasi:     rekapitulasiList,
		TotalBelanja:     totalBelanja,
		TotalOperasional: totalOperasional,
	}, nil
}
