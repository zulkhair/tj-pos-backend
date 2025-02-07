package transactionhandler

import (
	transactiondomain "dromatech/pos-backend/internal/domain/transaction"
	restutil "dromatech/pos-backend/internal/util/rest"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *Handler) FindDana(c *gin.Context) {
	// get user id
	userID := restutil.GetSession(c).UserID
	dateString := c.Query("date")
	if dateString == "" {
		restutil.SendResponseFail(c, "Harap isi tanggal")
		return
	}

	date, err := time.Parse("2006-01-02", dateString)
	if err != nil {
		restutil.SendResponseFail(c, "Tanggal tidak valid")
		return
	}

	dana, err := h.transactionUsecase.FindDana(userID, date)
	if err != nil {
		restutil.SendResponseFail(c, err.Error())
		return
	}

	restutil.SendResponseOk(c, "", dana)

}

func (h *Handler) CreateDana(c *gin.Context) {
	userID := restutil.GetSession(c).UserID

	// get request body
	var request transactiondomain.DanaRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		restutil.SendResponseFail(c, "Gagal mengambil data")
		return
	}

	// create dana
	err := h.transactionUsecase.CreateDana(userID, request)
	if err != nil {
		restutil.SendResponseFail(c, "Gagal membuat dana")
		return
	}

	restutil.SendResponseOk(c, "Dana berhasil dibuat", nil)
}

func (h *Handler) UpdateDana(c *gin.Context) {
	userID := restutil.GetSession(c).UserID

	// get request body
	var request transactiondomain.DanaRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		restutil.SendResponseFail(c, "Gagal mengambil data")
		return
	}

	err := h.transactionUsecase.UpdateDana(userID, request)
	if err != nil {
		if err.Error() == "not allowed" {
			restutil.SendResponseFail(c, "Anda tidak diperbolehkan mengubah data ini")
			return
		}
		restutil.SendResponseFail(c, "Gagal mengubah data")
		return
	}

	restutil.SendResponseOk(c, "Data berhasil diubah", nil)
}

func (h *Handler) SendDana(c *gin.Context) {
	userID := restutil.GetSession(c).UserID

	// get request body
	var request transactiondomain.DanaTransactionRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		restutil.SendResponseFail(c, "Gagal mengambil data")
		return
	}

	err := h.transactionUsecase.SendDana(userID, request)
	if err != nil {
		restutil.SendResponseFail(c, "Gagal mengirim dana")
		return
	}

	restutil.SendResponseOk(c, "Dana berhasil dikirim", nil)
}

func (h *Handler) ApproveDana(c *gin.Context) {
	userID := restutil.GetSession(c).UserID
	// get id from request body
	var request map[string]string
	if err := c.ShouldBindJSON(&request); err != nil {
		restutil.SendResponseFail(c, "Gagal mengambil data")
		return
	}

	err := h.transactionUsecase.ApproveDana(userID, request["id"])
	if err != nil {
		if err.Error() == "not allowed" {
			restutil.SendResponseFail(c, "Anda tidak diperbolehkan mengubah dana ini")
			return
		}
		if err.Error() == "not pending" {
			restutil.SendResponseFail(c, "Pengiriman dana sudah tidak dapat diapprove")
			return
		}
		restutil.SendResponseFail(c, "Gagal mengubah dana")
		return
	}

	restutil.SendResponseOk(c, "Dana berhasil diapprove", nil)
}

func (h *Handler) RejectDana(c *gin.Context) {
	userID := restutil.GetSession(c).UserID

	// get id from request body
	var request map[string]string
	if err := c.ShouldBindJSON(&request); err != nil {
		restutil.SendResponseFail(c, "Gagal mengambil data")
		return
	}

	err := h.transactionUsecase.RejectDana(userID, request["id"])
	if err != nil {
		if err.Error() == "not allowed" {
			restutil.SendResponseFail(c, "Anda tidak diperbolehkan menolak pengiriman dana ini")
			return
		}
		if err.Error() == "not pending" {
			restutil.SendResponseFail(c, "Pengiriman dana sudah tidak dapat ditolak")
			return
		}
		restutil.SendResponseFail(c, "Gagal menolak pengiriman dana")
		return
	}

	restutil.SendResponseOk(c, "Pengiriman dana berhasil ditolak", nil)
}

func (h *Handler) CancelSendDana(c *gin.Context) {
	userID := restutil.GetSession(c).UserID

	// get id from request body
	var request map[string]string
	if err := c.ShouldBindJSON(&request); err != nil {
		restutil.SendResponseFail(c, "Gagal mengambil data")
		return
	}

	err := h.transactionUsecase.CancelSendDana(userID, request["id"])
	if err != nil {
		restutil.SendResponseFail(c, "Gagal membatalkan pengiriman dana")
		return
	}

	restutil.SendResponseOk(c, "Pengiriman dana berhasil dibatalkan", nil)
}

func (h *Handler) FindUserMobile(c *gin.Context) {
	userID := restutil.GetSession(c).UserID

	userMobile, err := h.transactionUsecase.FindUserMobile(userID)
	if err != nil {
		restutil.SendResponseFail(c, "Gagal mengambil data")
		return
	}

	restutil.SendResponseOk(c, "Data berhasil diambil", userMobile)
}

func (h *Handler) CreatePenjualan(c *gin.Context) {
	userID := restutil.GetSession(c).UserID

	// get request body
	var request transactiondomain.TrxCreateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		restutil.SendResponseFail(c, "Gagal mengambil data")
		return
	}

	err := h.transactionUsecase.CreatePenjualan(userID, request)
	if err != nil {
		restutil.SendResponseFail(c, "Gagal membuat penjualan")
		return
	}

	restutil.SendResponseOk(c, "Penjualan berhasil dibuat", nil)
}

func (h *Handler) DeletePenjualan(c *gin.Context) {
	userID := restutil.GetSession(c).UserID

	// get id from request body
	var request map[string]string
	if err := c.ShouldBindJSON(&request); err != nil {
		restutil.SendResponseFail(c, "Gagal mengambil data")
		return
	}

	err := h.transactionUsecase.DeletePenjualan(userID, request["id"])
	if err != nil {
		restutil.SendResponseFail(c, "Gagal menghapus penjualan")
		return
	}

	restutil.SendResponseOk(c, "Penjualan berhasil dihapus", nil)
}

func (h *Handler) FindPenjualan(c *gin.Context) {
	userID := restutil.GetSession(c).UserID
	dateString := c.Query("date")
	if dateString == "" {
		restutil.SendResponseFail(c, "Harap isi tanggal")
		return
	}

	date, err := time.Parse("2006-01-02", dateString)
	if err != nil {
		restutil.SendResponseFail(c, "Tanggal tidak valid")
		return
	}

	penjualan, err := h.transactionUsecase.FindPenjualan(userID, date)
	if err != nil {
		restutil.SendResponseFail(c, "Gagal mengambil data")
		return
	}

	restutil.SendResponseOk(c, "Data berhasil diambil", penjualan)
}

func (h *Handler) CreateBelanja(c *gin.Context) {
	userID := restutil.GetSession(c).UserID

	// get request body
	var request transactiondomain.TrxCreateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		restutil.SendResponseFail(c, "Gagal mengambil data")
		return
	}

	err := h.transactionUsecase.CreateBelanja(userID, request)
	if err != nil {
		restutil.SendResponseFail(c, "Gagal membuat belanja")
		return
	}

	restutil.SendResponseOk(c, "Belanja berhasil dibuat", nil)
}

func (h *Handler) DeleteBelanja(c *gin.Context) {
	userID := restutil.GetSession(c).UserID

	// get id from request body
	var request map[string]string
	if err := c.ShouldBindJSON(&request); err != nil {
		restutil.SendResponseFail(c, "Gagal mengambil data")
		return
	}

	err := h.transactionUsecase.DeleteBelanja(userID, request["id"])
	if err != nil {
		restutil.SendResponseFail(c, "Gagal menghapus belanja")
		return
	}

	restutil.SendResponseOk(c, "Belanja berhasil dihapus", nil)
}

func (h *Handler) FindBelanja(c *gin.Context) {
	userID := restutil.GetSession(c).UserID
	dateString := c.Query("date")
	if dateString == "" {
		restutil.SendResponseFail(c, "Harap isi tanggal")
		return
	}

	date, err := time.Parse("2006-01-02", dateString)
	if err != nil {
		restutil.SendResponseFail(c, "Tanggal tidak valid")
		return
	}

	belanja, err := h.transactionUsecase.FindBelanja(userID, date)
	if err != nil {
		restutil.SendResponseFail(c, "Gagal mengambil data")
		return
	}

	restutil.SendResponseOk(c, "Data berhasil diambil", belanja)
}

func (h *Handler) CreateOperasional(c *gin.Context) {
	userID := restutil.GetSession(c).UserID

	// get request body
	var request transactiondomain.TrxCreateOperasionalRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		restutil.SendResponseFail(c, "Gagal mengambil data")
		return
	}

	err := h.transactionUsecase.CreateOperasional(userID, request)
	if err != nil {
		restutil.SendResponseFail(c, "Gagal membuat operasional")
		return
	}

	restutil.SendResponseOk(c, "Operasional berhasil dibuat", nil)
}

func (h *Handler) DeleteOperasional(c *gin.Context) {
	userID := restutil.GetSession(c).UserID

	// get id from request body
	var request map[string]string
	if err := c.ShouldBindJSON(&request); err != nil {
		restutil.SendResponseFail(c, "Gagal mengambil data")
		return
	}

	err := h.transactionUsecase.DeleteOperasional(userID, request["id"])
	if err != nil {
		restutil.SendResponseFail(c, "Gagal menghapus operasional")
		return
	}

	restutil.SendResponseOk(c, "Operasional berhasil dihapus", nil)
}

func (h *Handler) FindOperasional(c *gin.Context) {
	userID := restutil.GetSession(c).UserID
	dateString := c.Query("date")
	if dateString == "" {
		restutil.SendResponseFail(c, "Harap isi tanggal")
		return
	}

	date, err := time.Parse("2006-01-02", dateString)
	if err != nil {
		restutil.SendResponseFail(c, "Tanggal tidak valid")
		return
	}

	operasional, err := h.transactionUsecase.FindOperasional(userID, date)
	if err != nil {
		restutil.SendResponseFail(c, "Gagal mengambil data")
		return
	}

	restutil.SendResponseOk(c, "Data berhasil diambil", operasional)
}

func (h *Handler) FindDescriptionOperasional(c *gin.Context) {
	descriptions := []string{
		"Parkir Mobil",
		"Parkir Motor",
		"Masuk Pasar",
		"Keluar Pasar",
		"Pak Ogah",
		"Makanan",
		"Minuman",
		"Tambal Ban",
		"Isi Angin Ban",
		"Keresek",
		"Muat",
		"Bensin",
		"Solar",
		"Biaya Admin",
		"Peralatan",
		"Alat Tulis Kantor",
		"Cuci Kendaraan",
		"Pajak Kendaraan",
		"Top Up E-Money",
		"Biaya Pengiriman",
		"Biaya Transportasi",
		"Biaya Perawatan Kendaraan",
	}

	restutil.SendResponseOk(c, "Data berhasil diambil", descriptions)
}

func (h *Handler) FindSaldo(c *gin.Context) {
	userID := restutil.GetSession(c).UserID
	dateString := c.Query("date")
	if dateString == "" {
		restutil.SendResponseFail(c, "Harap isi tanggal")
		return
	}

	date, err := time.Parse("2006-01-02", dateString)
	if err != nil {
		restutil.SendResponseFail(c, "Tanggal tidak valid")
		return
	}

	saldo, err := h.transactionUsecase.FindSaldo(userID, date)
	if err != nil {
		restutil.SendResponseFail(c, "Gagal mengambil data")
		return
	}

	restutil.SendResponseOk(c, "Data berhasil diambil", saldo)
}

func (h *Handler) FindRekapitulasi(c *gin.Context) {
	dateString := c.Query("date")
	if dateString == "" {
		restutil.SendResponseFail(c, "Harap isi tanggal")
		return
	}

	date, err := time.Parse("2006-01-02", dateString)
	if err != nil {
		restutil.SendResponseFail(c, "Tanggal tidak valid")
		return
	}

	rekapitulasi, err := h.transactionUsecase.FindRekapitulasi(date)
	if err != nil {
		restutil.SendResponseFail(c, "Gagal mengambil data")
		return
	}

	restutil.SendResponseOk(c, "Data berhasil diambil", rekapitulasi)
}
