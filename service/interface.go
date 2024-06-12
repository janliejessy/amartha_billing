package service

import "github.com/jmoiron/sqlx"

type BillingService struct {
	Db *sqlx.DB
}

type BillingInterface interface {
	IsDeliquent(loanid int64) (bool, error)
	GetOutstanding(loanid int64) (float64, error)
	MakePayment(loanid int64) error
}
