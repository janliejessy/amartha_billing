package service

import (
	"database/sql"
	"errors"
	"fmt"
	"math"
	"time"

	"amartha/database/gen/billing/model"
)

// IsDeliquent checks the given loanid if the payment has not been made for more than two weeks.
func (b *BillingService) IsDeliquent(loanid int64) (bool, error) {
	loanDetails, err := b.getLoanDetails(loanid)
	if err != nil {
		return false, fmt.Errorf("IsDeliquent get error: %w", err)
	}

	deliquentDeadline := time.Now().AddDate(0, 0, -14).Format("2006-01-02")
	loanStartdate := loanDetails.Startdate.Format("2006-01-02")
	installmentEndDate := loanDetails.Startdate.AddDate(0, 0, int(loanDetails.Installment)*7).Format("2006-01-02")

	// note: check for edge-case scenarios
	if loanStartdate > deliquentDeadline {
		// this loan is less than 14 days old
		return false, nil
	} else if loanDetails.Amountpaid >= getTotalLoan(*loanDetails) {
		// this loan has been fully repaid
		return false, nil
	} else if installmentEndDate < deliquentDeadline {
		// this loan has past its expiry date for more than 2 weeks and has not been fully paid
		return true, nil
	}

	expectedAmountpaid := getCurrentInstallmentWeek(*loanDetails) * getWeeklyInstallmentPayment(*loanDetails)
	deliquentAmountThreshold := expectedAmountpaid - (2 * getWeeklyInstallmentPayment(*loanDetails))

	if deliquentAmountThreshold > loanDetails.Amountpaid {
		return true, nil
	}

	return false, nil
}

// GetOutstanding returns the current outstanding on a loan, or 0 if the loan has been fully paid.
func (b *BillingService) GetOutstanding(loanid int64) (float64, error) {
	loanDetails, err := b.getLoanDetails(loanid)
	if err != nil {
		return 0, fmt.Errorf("GetOutstanding get error: %w", err)
	}

	totalLoan := getTotalLoan(*loanDetails)

	if loanDetails.Amountpaid >= totalLoan {
		// this loan has been fully repaid
		return 0, nil
	}

	return (totalLoan - loanDetails.Amountpaid), nil
}

// MakePayment makes the required payment for that particular installment week.
// Note: this is under the assumption made from the case study whereby; "borrower can
// only pay the exact amount of payable that week".
func (b *BillingService) MakePayment(loanid int64) error {
	loanDetails, err := b.getLoanDetails(loanid)
	if err != nil {
		return fmt.Errorf("MakePayment error: %w", err)
	}

	if loanDetails.Amountpaid >= getTotalLoan(*loanDetails) {
		// this loan has been fully repaid, do nothing
		return nil
	}

	expectedAmountpaid := getCurrentInstallmentWeek(*loanDetails) * getWeeklyInstallmentPayment(*loanDetails)
	minPaymentNeeded := expectedAmountpaid - loanDetails.Amountpaid

	if minPaymentNeeded == 0 {
		// payment has been made for the week, do nothing
		return nil
	}

	// inserts a new transaction
	_, err = b.Db.Exec(`
		INSERT INTO billing.transaction (loanid, payment)
		VALUES (?, ?)`, loanid, minPaymentNeeded)
	if err != nil {
		return fmt.Errorf("MakePayment insert transaction: %w", err)
	}

	// update the loan details
	_, err = b.Db.Exec(`
		UPDATE billing.loan
		SET amountpaid = ?
		WHERE loanid = ?`, expectedAmountpaid, loanid)
	if err != nil {
		return fmt.Errorf("MakePayment update loan: %w", err)
	}

	return nil
}

// getTotalLoan returns the capital amount plus its interest rate.
func getTotalLoan(loanDetails model.Loan) float64 {
	return loanDetails.Amount + (loanDetails.Amount * (loanDetails.Interestrate))
}

// getWeeklyInstallmentPayment returns the amount the user needs to pay for the week,
// assuming the amount is divided equally over the whole installment period.
func getWeeklyInstallmentPayment(loanDetails model.Loan) float64 {
	return getTotalLoan(loanDetails) / float64(loanDetails.Installment)
}

// getCurrentInstallmentWeek returns the current installment week the loan is at.
// If it has past the alloted period, then return the maximum installment weel.
func getCurrentInstallmentWeek(loanDetails model.Loan) float64 {
	timelapse := time.Since(loanDetails.Startdate).Hours()
	curWeek := math.Ceil(timelapse / 24 / 7)

	if curWeek > float64(loanDetails.Installment) {
		return float64(loanDetails.Installment)
	}

	fmt.Printf("\ntimelapse: %f\n", curWeek)
	return curWeek
}

// getLoanDetails returns model.Loan struct for the given loanid.
func (b *BillingService) getLoanDetails(loanid int64) (*model.Loan, error) {
	loanDetails := &model.Loan{}

	err := b.Db.Get(loanDetails, `
		SELECT *
		FROM billing.loan
		WHERE loanid = ?
		`, loanid)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("loan not found")
		}
		return nil, fmt.Errorf("sql getLoanDetails error: %w", err)
	}

	return loanDetails, nil
}
