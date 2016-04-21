package user

import (
	"log"
	"time"

	"github.com/jmoiron/sqlx"
)

type UserModule struct {
	DBConn *sqlx.DB
}

func NewUserModule(db *sqlx.DB) *UserModule {
	return &UserModule{
		DBConn: db,
	}
}

type User struct {
	ID              int64     `json:"user_id,omitempty"      db:"user_id"`
	Email           string    `json:"email"        db:"email"`
	Name            string    `json:"name"         db:"name"`
	Password        string    `json:"password"     db:"password"`
	Gender          int       `json:"gender"       db:"gender"`
	BirthDate       time.Time `json:"birth_date"   db:"birth_date"`
	NIK             string    `json:"nik"          db:"nik"`
	NIKValid        int       `json:"nik_valid,omitempty"    db:"nik_valid"`
	MSISDN          string    `json:"msidn"        db:"msisdn"`
	ThresholdAmount float64   `json:"th_amount"    db:"th_amount"`
	CreateTime      time.Time `json:"create_time,omitempty"  db:"create_time"`
	Photo           string    `json:"photo,omitempty"        db:"photo"`
}

func (um *UserModule) RegisterUser(user User) error {
	tx, err := um.DBConn.Beginx()
	if err != nil {
		log.Println(err)
		return err
	}

	user.NIKValid = 0
	user.ThresholdAmount = 100000

	if err := user.Insert(tx); err != nil {
		log.Println(err)
		if err := tx.Rollback(); err != nil {
			log.Println(err)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (user *User) Insert(tx *sqlx.Tx) error {
	sqlQuery := `
        INSERT INTO fm_user (
            email,
            name,
            password,
            gender,
            birth_date,
            nik,
            nik_valid,
            msisdn,
            th_amount,
            create_time
        ) VALUES (
            :email,
            :name,
            :password,
            :gender,
            :birth_date,
            :nik,
            :nik_valid,
            :msisdn,
            :th_amount,
            CURRENT_TIMESTAMP
        )
    `
	_, err := tx.NamedExec(sqlQuery, user)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
