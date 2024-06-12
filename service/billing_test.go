package service

import (
	"testing"

	"amartha/database"

	"github.com/stretchr/testify/assert"
)

func TestIsDeliquent(t *testing.T) {
	svc := BillingService{
		Db: database.ConnectDB(),
	}
	defer svc.Db.Close()

	testcases := []struct {
		name      string
		loanid    int64
		expError  string
		expResult bool
	}{
		{
			name:      "deliquent",
			loanid:    3,
			expError:  "",
			expResult: true,
		},
		{
			name:      "not a deliquent",
			loanid:    4,
			expError:  "",
			expResult: false,
		},
		{
			name:      "new loan, skip deliquent check",
			loanid:    1,
			expError:  "",
			expResult: false,
		},
		{
			name:      "loan has been repaid, skip deliquent check",
			loanid:    2,
			expError:  "",
			expResult: false,
		},
		{
			name:      "loan not found",
			loanid:    777,
			expError:  "loan not found",
			expResult: false,
		},
	}

	for _, tc := range testcases {
		res, err := svc.IsDeliquent(tc.loanid)

		if tc.expError != "" {
			assert.ErrorContains(t, err, tc.expError)
		} else {
			assert.Nil(t, err)
		}
		assert.Equal(t, tc.expResult, res)
	}
}

func TestGetOutstanding(t *testing.T) {
	svc := BillingService{
		Db: database.ConnectDB(),
	}
	defer svc.Db.Close()

	testcases := []struct {
		name      string
		loanid    int64
		expError  string
		expResult float64
	}{
		{
			name:      "loan has not been paid, return full amount",
			loanid:    1,
			expError:  "",
			expResult: 1100,
		},
		{
			name:      "loan has been partially paid, return partial amount",
			loanid:    3,
			expError:  "",
			expResult: 1070,
		},
		{
			name:      "loan has been fully paid, return 0",
			loanid:    2,
			expError:  "",
			expResult: 0,
		},
		{
			name:      "loan not found",
			loanid:    777,
			expError:  "loan not found",
			expResult: 0,
		},
	}

	for _, tc := range testcases {
		res, err := svc.GetOutstanding(tc.loanid)

		if tc.expError != "" {
			assert.ErrorContains(t, err, tc.expError)
		} else {
			assert.Nil(t, err)
		}
		assert.Equal(t, tc.expResult, res)
	}
}

func TestMakePayment(t *testing.T) {
	svc := BillingService{
		Db: database.ConnectDB(),
	}
	defer svc.Db.Close()

	testcases := []struct {
		name               string
		loanid             int64
		expError           string
		expTotalPaidAmount float64
	}{
		{
			name:               "one week payment made",
			loanid:             1,
			expError:           "",
			expTotalPaidAmount: 22,
		},
		{
			name:               "multiple weeks payment made",
			loanid:             3,
			expError:           "",
			expTotalPaidAmount: 110,
		},
		{
			name:               "loan repaid, no payment needed",
			loanid:             2,
			expError:           "",
			expTotalPaidAmount: 1100,
		},
		{
			name:               "loan not found",
			loanid:             777,
			expError:           "loan not found",
			expTotalPaidAmount: 0,
		},
	}

	for _, tc := range testcases {
		err := svc.MakePayment(tc.loanid)

		if tc.expError != "" {
			assert.ErrorContains(t, err, tc.expError)
		} else {
			assert.Nil(t, err)

			res, err := svc.getLoanDetails(tc.loanid)
			assert.Nil(t, err)
			assert.Equal(t, tc.expTotalPaidAmount, res.Amountpaid)
		}
	}
}
