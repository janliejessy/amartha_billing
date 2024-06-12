package main

import (
	"errors"
	"flag"
	"fmt"

	"amartha/database"
	"amartha/service"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db := database.ConnectDB()
	defer db.Close()

	svc := &service.BillingService{
		Db: db,
	}

	loanid, operation := parseFlag()

	switch operation {
	case 1:
		deliquent, err := svc.IsDeliquent(loanid)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Printf("\n\n** is deliquent: %t **\n\n", deliquent)
	case 2:
		outstanding, err := svc.GetOutstanding(loanid)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Printf("\n\n** get outstanding: %f **\n\n", outstanding)
	case 3:
		err := svc.MakePayment(loanid)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Printf("\n\n** payment completed **\n\n")
	}
}

func parseFlag() (int64, int64) {
	flagLoanId := flag.Int64("l", 1, "loan id")
	flagOperation := flag.Int64("op", 1, "1=IsDeliquent, 2=GetOutstanding, 3=MakePayment")

	flag.Parse()

	if flagLoanId == nil {
		panic(errors.New("loan id required"))
	} else if flagOperation == nil {
		panic(errors.New("operation required"))
	}

	return *flagLoanId, *flagOperation
}
