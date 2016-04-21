package user

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
)

const (
	GENDER_MALE   int = 1
	GENDER_FEMALE int = 2
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

func (um *UserModule) UserRegister(user User) error {
	if !user.ValidateNIK() {
		return fmt.Errorf("NIK is not valid.")
	}

	user.NIKValid = 0
	user.ThresholdAmount = 100000

	tx, err := um.DBConn.Beginx()
	if err != nil {
		log.Println(err)
		return err
	}

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

func (user *User) ValidateNIK() bool {
	nik := user.NIK
	gender := user.Gender
	year, month, day := user.BirthDate.Date()

	// Default length of NIK is 16 digits
	if len(nik) != 16 {
		return false
	}

	// Year
	if nik[10:12] != strconv.Itoa(year)[2:4] {
		return false
	}

	// Month
	if nik[8:10] != strconv.Itoa(int(month)) {
		return false
	}

	// Date
	bornDay := nik[6:8]
	if gender == GENDER_MALE {
		if bornDay != strconv.Itoa(day) {
			return false
		}
	} else {
		if bornDay != strconv.Itoa(day-40) {
			return false
		}
	}

	return true
}

func (user *User) UserLogin(um *UserModule) error {
	query := `
        SELECT
			email,
            name,
            gender,
            birth_date,
            nik,
            nik_valid,
            msisdn,
            th_amount,
            create_time
        FROM fm_user
        WHERE email = $1
        AND password = $2
    `
	if err := um.DBConn.Get(user, query, user.Email, user.Password); err != nil {
		log.Println(err)
		return err
	}

	return nil
}
