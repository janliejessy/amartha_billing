DROP SCHEMA IF EXISTS billing ;
CREATE SCHEMA billing;

CREATE TABLE billing.user_account (
    `userid` bigint(20) NOT NULL AUTO_INCREMENT,
    `firstname` varchar(255) NOT NULL,
    `lastname` varchar(255) NOT NULL,
    PRIMARY KEY (`userid`),
    KEY `idx_user` (`userid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE billing.transaction (
    `txnid` bigint(20) NOT NULL AUTO_INCREMENT,
    `loanid` bigint(20) NOT NULL,
    `payment` decimal(20, 2) NOT NULL,
    `paymentdate` timestamp DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (txnid),
    KEY `idx_txn` (`txnid`, `loanid`, `paymentdate`, `payment`)
 ) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE billing.loan (
    `loanid` bigint(20) NOT NULL AUTO_INCREMENT,
    `userid` bigint(20) NOT NULL,
    `amount` decimal(20, 2) NOT NULL,
    `startdate` timestamp NOT NULL,
    `installment` bigint(20) NOT NULL DEFAULT 50,
    `interestrate` decimal(2,2) NOT NULL DEFAULT 0.1,
    `amountpaid` decimal(20, 2) NOT NULL DEFAULT 0,
    PRIMARY KEY (loanid),
    KEY `idx_loan` (`userid`, `loanid`, `amount`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

INSERT INTO billing.user_account VALUES
    (1, 'john', 'doe');

INSERT INTO billing.transaction VALUES 
    (1, 1, 110, NOW()),
    (2, 2, 110, '2024-05-14'),
    (3, 4, 110, DATE_ADD(NOW(), INTERVAL -21 DAY)),
    (4, 4, 110, DATE_ADD(NOW(), INTERVAL -14 DAY)),
    (5, 4, 110, DATE_ADD(NOW(), INTERVAL -7 DAY)),
    (6, 3, 110, DATE_ADD(NOW(), INTERVAL -21 DAY)),
    (7, 2, 505, DATE_ADD(NOW(), INTERVAL -28 DAY)),
    (8, 2, 505, DATE_ADD(NOW(), INTERVAL -21 DAY));

INSERT INTO billing.loan VALUES
    (1, 1, 1000, NOW(), 50, 0.1, 0), -- new loan
    (2, 1, 1000, DATE_ADD(NOW(), INTERVAL -1 MONTH), 2, 0.1, 1100), -- loan repaid
    (3, 1, 1000, DATE_ADD(NOW(), INTERVAL -1 MONTH), 50, 0.1, 30), -- current deliquent
    (4, 1, 1000, DATE_ADD(NOW(), INTERVAL -1 MONTH), 4, 0.1, 825), -- payment with 1 installment left
	(5, 1, 1000, DATE_ADD(NOW(), INTERVAL -5 YEAR), 50, 0.1, 0); -- deliquent long past installment period