package tx

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/nyelonong/finapimate/oauth"
	"github.com/nyelonong/finapimate/user"
	"github.com/nyelonong/finapimate/utils"
)

const (
	STATUS_REQUEST  = 1
	STATUS_APPROVED = 2
	STATUS_DECLINE  = 3
	STATUS_PAID     = 4
)

type TxModule struct {
	DBConn     *sqlx.DB
	UserModule *user.UserModule
}

func NewTxModule(db *sqlx.DB, um *user.UserModule) *TxModule {
	return &TxModule{
		DBConn:     db,
		UserModule: um,
	}
}

type Transaction struct {
	ID            int64     `json:"tx_id,omitempty"       	db:"tx_id"`
	LenderID      int64     `json:"lender_id"            	db:"lender_id"`
	BorrowerID    int64     `json:"borrower_id"      		db:"borrower_id"`
	Amount        float64   `json:"amount"                 	db:"amount"`
	Deadline      int64     `json:"deadline"`
	DeadlineValid time.Time `json:"-"          				db:"deadline"`
	Status        int       `json:"status,omitempty"  		db:"status"`
	Notes         string    `json:"notes,omitempty"      	db:"notes"`
	CreateTime    time.Time `json:"create_time,omitempty"	db:"create_time"`
	UserProfile   user.User `json:"user_profile"`
}

type EwalletPayment struct {
	CompanyCode   string
	PrimaryID     string
	TransactionID string
	ReferenceID   string
	RequestDate   string
	Amount        string
	CurrencyCode  string
}

type EwalletPaymentResponse struct {
	CompanyCode     string
	TransactionID   string
	ReferenceID     string
	PaymentID       string
	TransactionDate string
}

type EwalletTopUp struct {
	CompanyCode    string
	CustomerNumber string
	TransactionID  string
	RequestDate    string
	Amount         string
	CurrencyCode   string
}

type EwalletTopUpResponse struct {
	CompanyCode     string
	TransactionID   string
	TopUpID         string
	TransactionDate string
}

type TopUp struct {
	UserID int64   `json:"user_id"`
	Amount float64 `json:"amount"`
}

func (tm *TxModule) RequestBorrow(trxs []Transaction) error {
	tx, err := tm.DBConn.Beginx()
	if err != nil {
		log.Println(err)
		return err
	}

	for _, trx := range trxs {
		trx.DeadlineValid = time.Unix(trx.Deadline, 0)
		if err := trx.Insert(tx); err != nil {
			log.Println(err)
			if err := tx.Rollback(); err != nil {
				log.Println(err)
			}
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (trx *Transaction) Insert(tx *sqlx.Tx) error {
	sqlQuery := `
        INSERT INTO fm_tx (
            lender_id,
            borrower_id,
            amount,
            deadline,
            status,
            create_time
        ) VALUES (
            :lender_id,
            :borrower_id,
            :amount,
            :deadline,
            :status,
            CURRENT_TIMESTAMP
        )
    `
	_, err := tx.NamedExec(sqlQuery, trx)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (tm *TxModule) ListBorrow(trx Transaction) ([]Transaction, error) {
	data := make([]Transaction, 0)

	query := `
        SELECT
			lender_id,
			borrower_id,
			amount,
			deadline,
			status,
			COALESCE(notes, '-') AS notes,
			create_time
        FROM fm_tx
        WHERE borrower_id = $1
    `

	rows, err := tm.DBConn.Queryx(query, trx.BorrowerID)
	if err != nil {
		log.Println(err)
		return data, err
	}

	for rows.Next() {
		var trx Transaction
		if err := rows.StructScan(&trx); err != nil {
			log.Println(err)
		} else {
			usr := user.User{
				ID: trx.LenderID,
			}
			if err := usr.Get(tm.UserModule); err != nil {
				log.Println(err)
			} else {
				trx.UserProfile = usr
				data = append(data, trx)
			}
		}
	}

	return data, nil
}

func (tm *TxModule) ListLend(trx Transaction) ([]Transaction, error) {
	data := make([]Transaction, 0)

	query := `
        SELECT
			lender_id,
			borrower_id,
			amount,
			deadline,
			status,
			COALESCE(notes, '-') AS notes,
			create_time
        FROM fm_tx
        WHERE lender_id = $1
    `

	rows, err := tm.DBConn.Queryx(query, trx.LenderID)
	if err != nil {
		log.Println(err)
		return data, err
	}

	for rows.Next() {
		var trx Transaction
		if err := rows.StructScan(&trx); err != nil {
			log.Println(err)
		} else {
			usr := user.User{
				ID: trx.BorrowerID,
			}
			if err := usr.Get(tm.UserModule); err != nil {
				log.Println(err)
			} else {
				trx.UserProfile = usr
				data = append(data, trx)
			}
		}
	}

	return data, nil
}

func (tm *TxModule) NotifBorrow(trx Transaction) ([]Transaction, error) {
	data := make([]Transaction, 0)

	query := `
        SELECT
			lender_id,
			borrower_id,
			amount,
			deadline,
			status,
			COALESCE(notes, '-') AS notes,
			create_time
        FROM fm_tx
        WHERE lender_id = $1
		AND status = $2
    `

	rows, err := tm.DBConn.Queryx(query, trx.LenderID, STATUS_REQUEST)
	if err != nil {
		log.Println(err)
		return data, err
	}

	for rows.Next() {
		var trx Transaction
		if err := rows.StructScan(&trx); err != nil {
			log.Println(err)
		} else {
			usr := user.User{
				ID: trx.BorrowerID,
			}
			if err := usr.Get(tm.UserModule); err != nil {
				log.Println(err)
			} else {
				trx.UserProfile = usr
				data = append(data, trx)
			}
		}
	}

	return data, nil
}

func (tm *TxModule) ChangeStatusTx(trxs []Transaction) error {
	tx, err := tm.DBConn.Beginx()
	if err != nil {
		log.Println(err)
		return err
	}

	for _, trx := range trxs {
		if err := trx.Update(tx); err != nil {
			log.Println(err)
			if err := tx.Rollback(); err != nil {
				log.Println(err)
			}
			return err
		}

		switch trx.Status {
		case STATUS_APPROVED:
			lender := user.User{
				ID: trx.LenderID,
			}
			if err := lender.Get(tm.UserModule); err != nil {
				log.Println(err)
				return err
			}

			// Debit
			debit := EwalletPayment{
				CompanyCode:   utils.COMPANY_CODE,
				PrimaryID:     lender.Email,
				TransactionID: fmt.Sprintf("%d", trx.ID),
				ReferenceID:   fmt.Sprintf("%d-%d", trx.ID, trx.LenderID),
				RequestDate:   time.Now().Format(time.RFC3339),
				Amount:        fmt.Sprintf("%.2f", trx.Amount),
				CurrencyCode:  "IDR",
			}

			if _, err := debit.Payment(); err != nil {
				log.Println(err)
				if err := tx.Rollback(); err != nil {
					log.Println(err)
				}
				return err
			}

			// Credit
			credit := EwalletTopUp{
				CompanyCode:    utils.COMPANY_CODE,
				CustomerNumber: fmt.Sprintf("%d", trx.BorrowerID),
				TransactionID:  fmt.Sprintf("%d", trx.ID),
				RequestDate:    time.Now().Format(time.RFC3339),
				Amount:         fmt.Sprintf("%.2f", trx.Amount),
				CurrencyCode:   "IDR",
			}

			if _, err := credit.TopUp(); err != nil {
				log.Println(err)
				if err := tx.Rollback(); err != nil {
					log.Println(err)
				}
				return err
			}
		case STATUS_PAID:
			borrower := user.User{
				ID: trx.LenderID,
			}
			if err := borrower.Get(tm.UserModule); err != nil {
				log.Println(err)
				return err
			}

			// Debit
			debit := EwalletPayment{
				CompanyCode:   utils.COMPANY_CODE,
				PrimaryID:     borrower.Email,
				TransactionID: fmt.Sprintf("%d", trx.ID),
				ReferenceID:   fmt.Sprintf("%d-%d", trx.ID, trx.LenderID),
				RequestDate:   time.Now().Format(time.RFC3339),
				Amount:        fmt.Sprintf("%.2f", trx.Amount),
				CurrencyCode:  "IDR",
			}

			if _, err := debit.Payment(); err != nil {
				log.Println(err)
				if err := tx.Rollback(); err != nil {
					log.Println(err)
				}
				return err
			}

			// Credit
			credit := EwalletTopUp{
				CompanyCode:    utils.COMPANY_CODE,
				CustomerNumber: fmt.Sprintf("%d", trx.LenderID),
				TransactionID:  fmt.Sprintf("%d", trx.ID),
				RequestDate:    time.Now().Format(time.RFC3339),
				Amount:         fmt.Sprintf("%.2f", trx.Amount),
				CurrencyCode:   "IDR",
			}

			if _, err := credit.TopUp(); err != nil {
				log.Println(err)
				if err := tx.Rollback(); err != nil {
					log.Println(err)
				}
				return err
			}
		}
	}

	if err := tx.Commit(); err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (trx *Transaction) Update(tx *sqlx.Tx) error {
	sqlQuery := `
        UPDATE
            fm_tx
        SET
            status     	= :status,
            update_time	= CURRENT_TIMESTAMP
        WHERE tx_id = :tx_id
    `
	_, err := tx.NamedExec(sqlQuery, trx)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (ep *EwalletPayment) Payment() (*EwalletPaymentResponse, error) {
	encoded, err := json.Marshal(ep)
	if err != nil {
		return nil, err
		log.Println(err)
	}

	now := oauth.GetTime()
	method := "POST"
	path := "/ewallet/payments"

	// get access token.
	accessToken, err := oauth.GetAccessToken()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	headers := make(map[string]string)
	headers["Authorization"] = "Bearer " + accessToken
	headers["Origin"] = "tokopedia.com"
	headers["X-BCA-Key"] = utils.API_KEY
	headers["X-BCA-Timestamp"] = now
	headers["X-BCA-Signature"] = utils.GetSignature(method, path, accessToken, string(encoded), now)

	agent := utils.NewHTTPRequest()
	agent.Url = utils.API_URL
	agent.Path = path
	agent.Method = method
	agent.IsJson = true
	agent.Json = ep
	agent.Headers = headers

	body, err := agent.DoReq()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var resp EwalletPaymentResponse
	if err := json.Unmarshal(*body, &resp); err != nil {
		log.Println(err)
		var errResp utils.Error
		_ = json.Unmarshal(*body, &errResp)
		log.Println(errResp)
		return nil, err
	}

	return &resp, nil
}

func (et *EwalletTopUp) TopUp() (*EwalletTopUpResponse, error) {
	encoded, err := json.Marshal(et)
	if err != nil {
		return nil, err
		log.Println(err)
	}

	now := oauth.GetTime()
	method := "POST"
	path := "/ewallet/topup"

	// get access token.
	accessToken, err := oauth.GetAccessToken()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	headers := make(map[string]string)
	headers["Authorization"] = "Bearer " + accessToken
	headers["Origin"] = "tokopedia.com"
	headers["X-BCA-Key"] = utils.API_KEY
	headers["X-BCA-Timestamp"] = now
	headers["X-BCA-Signature"] = utils.GetSignature(method, path, accessToken, string(encoded), now)

	agent := utils.NewHTTPRequest()
	agent.Url = utils.API_URL
	agent.Path = path
	agent.Method = method
	agent.IsJson = true
	agent.Json = et
	agent.Headers = headers

	body, err := agent.DoReq()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var resp EwalletTopUpResponse
	if err := json.Unmarshal(*body, &resp); err != nil {
		log.Println(err)
		var errResp utils.Error
		_ = json.Unmarshal(*body, &errResp)
		log.Println(errResp)
		return nil, err
	}

	return &resp, nil
}

func (top TopUp) UserTopUp(tm *TxModule) error {
	user := EwalletTopUp{
		CompanyCode:    utils.COMPANY_CODE,
		CustomerNumber: fmt.Sprintf("%d", top.UserID),
		TransactionID:  fmt.Sprintf("%d", time.Now().Unix()),
		RequestDate:    time.Now().Format(time.RFC3339),
		Amount:         fmt.Sprintf("%.2f", top.Amount),
		CurrencyCode:   "IDR",
	}

	if _, err := user.TopUp(); err != nil {
		log.Println(err)
		return err
	}

	return nil
}
