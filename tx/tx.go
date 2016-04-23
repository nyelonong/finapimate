package tx

import (
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/nyelonong/finapimate/user"
)

const (
	STATUS_REQUEST  = 1
	STATUS_APPROVED = 2
	STATUS_DECLINE  = 3
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
		// usr.Status = status

		if err := trx.Update(tx); err != nil {
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
