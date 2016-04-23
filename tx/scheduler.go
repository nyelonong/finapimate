package tx

import (
	"gopkg.in/robfig/cron.v2"
	"log"
)

type Scheduler struct{
	format string // cron expression format.
	cron *cron.Cron
	txm *TxModule
}

func NewScheduler(f string, txm *TxModule) *Scheduler {
	s := &Scheduler{
		format: f,
		txm: txm,
	}
	s.start()

	return s

}

func (s *Scheduler) start() {
	s.cron = cron.New()
	s.cron.AddFunc(s.format, s.notify)
	s.cron.Start()
	return
}

// notify needs 
func (s *Scheduler) notify() {
	// get all user_id to notify from tx database.
	q := `
		SELECT borrower_id
		FROM fm_tx
		WHERE deadline >= CURRENT_DATE-'1 Day'::Interval
	`

	var borrowers []int64
	err := s.txm.DBConn.Select(&borrowers, q)
	if err != nil {
		log.Println(err)
		return
	}

	// TODO: send notification.
	return
}