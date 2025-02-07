package transactionusecase

import (
	transactiondomain "dromatech/pos-backend/internal/domain/transaction"
	"fmt"
	"time"
)

func (u *Usecase) FindDana(userID string, date time.Time) (*transactiondomain.DanaInquiryResponse, error) {
	// find to database
	dana, err := u.transactionRepo.FindDana(userID, date)
	if err != nil {
		return nil, err
	}

	return dana, nil
}

func (u *Usecase) CreateDana(userID string, request transactiondomain.DanaRequest) error {
	// create dana
	err := u.transactionRepo.CreateDana(userID, request)
	if err != nil {
		return err
	}

	return nil
}

func (u *Usecase) UpdateDana(userID string, request transactiondomain.DanaRequest) error {
	// find dana
	dana, err := u.transactionRepo.FindDanaByID(request.ID)
	if err != nil {
		return err
	}

	// check the userID is the same as the dana.UserID
	if dana.WebUserID != userID {
		return fmt.Errorf("not allowed")
	}

	// update dana
	err = u.transactionRepo.UpdateDana(userID, request)
	if err != nil {
		return err
	}

	return nil
}

func (u *Usecase) SendDana(userID string, request transactiondomain.DanaTransactionRequest) error {
	// check reciever is have permission for mobile
	hasPermission, err := u.transactionRepo.CheckUserMobilePermission(userID)
	if err != nil {
		return err
	}
	if !hasPermission {
		return fmt.Errorf("user not have permission for mobile")
	}

	// send dana
	err = u.transactionRepo.SendDana(userID, request)
	if err != nil {
		return err
	}

	return nil
}

func (u *Usecase) ApproveDana(userID string, id string) error {
	// find dana transaction
	danaTransaction, err := u.transactionRepo.FindDanaTransaction(id)
	if err != nil {
		return err
	}

	// check the receiver is the userID
	if danaTransaction.Receiver != userID {
		return fmt.Errorf("not allowed")
	}

	// check the dana transaction is still pending
	if danaTransaction.Status != transactiondomain.DanaStatusPending {
		return fmt.Errorf("not pending")
	}

	// approve dana
	err = u.transactionRepo.ApproveDana(id)
	if err != nil {
		return err
	}

	return nil
}

func (u *Usecase) RejectDana(userID string, id string) error {
	// find dana transaction
	danaTransaction, err := u.transactionRepo.FindDanaTransaction(id)
	if err != nil {
		return err
	}

	// check the receiver is the userID
	if danaTransaction.Receiver != userID {
		return fmt.Errorf("not allowed")
	}

	// check the dana transaction is still pending
	if danaTransaction.Status != transactiondomain.DanaStatusPending {
		return fmt.Errorf("not pending")
	}

	// reject dana
	err = u.transactionRepo.RejectDana(id)
	if err != nil {
		return err
	}

	return nil
}

func (u *Usecase) CancelSendDana(userID string, id string) error {
	// find dana transaction
	danaTransaction, err := u.transactionRepo.FindDanaTransaction(id)
	if err != nil {
		return err
	}

	// check the sender is the userID
	if danaTransaction.Sender != userID {
		return fmt.Errorf("not allowed")
	}

	// cancel send dana
	err = u.transactionRepo.CancelSendDana(id)
	if err != nil {
		return err
	}

	return nil
}

func (u *Usecase) FindUserMobile(userID string) ([]transactiondomain.WebUserMobile, error) {
	userMobile, err := u.transactionRepo.FindUserMobile(userID)
	if err != nil {
		return nil, err
	}

	return userMobile, nil
}

func (u *Usecase) CreatePenjualan(userID string, request transactiondomain.TrxCreateRequest) error {
	// create penjualan
	err := u.transactionRepo.CreatePenjualan(userID, request)
	if err != nil {
		return err
	}

	return nil
}

func (u *Usecase) DeletePenjualan(userID string, id string) error {
	err := u.transactionRepo.DeletePenjualan(id)
	if err != nil {
		return err
	}

	return nil
}

func (u *Usecase) FindPenjualan(userID string, date time.Time) ([]transactiondomain.TrxInquiryResponse, error) {
	penjualan, err := u.transactionRepo.FindPenjualan(userID, date)
	if err != nil {
		return nil, err
	}

	return penjualan, nil
}

func (u *Usecase) CreateBelanja(userID string, request transactiondomain.TrxCreateRequest) error {
	err := u.transactionRepo.CreateBelanja(userID, request)
	if err != nil {
		return err
	}

	return nil
}

func (u *Usecase) DeleteBelanja(userID string, id string) error {
	err := u.transactionRepo.DeleteBelanja(id)
	if err != nil {
		return err
	}

	return nil
}

func (u *Usecase) FindBelanja(userID string, date time.Time) ([]transactiondomain.TrxInquiryResponse, error) {
	belanja, err := u.transactionRepo.FindBelanja(userID, date)
	if err != nil {
		return nil, err
	}

	return belanja, nil
}

func (u *Usecase) CreateOperasional(userID string, request transactiondomain.TrxCreateOperasionalRequest) error {
	err := u.transactionRepo.CreateOperasional(userID, request)
	if err != nil {
		return err
	}

	return nil
}

func (u *Usecase) DeleteOperasional(userID string, id string) error {
	err := u.transactionRepo.DeleteOperasional(id)
	if err != nil {
		return err
	}

	return nil
}

func (u *Usecase) FindOperasional(userID string, date time.Time) ([]transactiondomain.TrxInquiryOperasionalResponse, error) {
	operasional, err := u.transactionRepo.FindOperasional(userID, date)
	if err != nil {
		return nil, err
	}

	return operasional, nil
}

func (u *Usecase) FindSaldo(userID string, date time.Time) (*transactiondomain.SaldoResponse, error) {
	saldo, err := u.transactionRepo.FindSaldo(userID, date)
	if err != nil {
		return nil, err
	}

	return saldo, nil
}

func (u *Usecase) FindRekapitulasi(date time.Time) (*transactiondomain.RekapitulasiResponse, error) {
	rekapitulasi, err := u.transactionRepo.FindRekapitulasi(date)
	if err != nil {
		return nil, err
	}

	return rekapitulasi, nil
}
