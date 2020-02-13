package core

import (
	"database/sql"
	"errors"
	"fmt"
)

var ErrInvalidPass = errors.New("invalid password")

type QueryError struct { // alt + enter
	Query string
	Err   error
}

type DbError struct {
	Err error
}

type DbTxError struct {
	Err         error
	RollbackErr error
}

type Atm struct {
	Id int64
	Name string
	Adress string
}
type Account struct {
	Id int64
	Name string
	Balance int64
	Client_id int64
}

type Service struct {
	Id int64
	Name string
	Balance int64
}



func (receiver *QueryError) Unwrap() error {
	return receiver.Err
}

func (receiver *QueryError) Error() string {
	return fmt.Sprintf("can't execute query %s: %s", loginSQL, receiver.Err.Error())
}



func queryError(query string, err error) *QueryError {
	return &QueryError{Query: query, Err: err}
}

func (receiver *DbError) Error() string {
	return fmt.Sprintf("can't handle db operation: %v", receiver.Err.Error())
}

func (receiver *DbError) Unwrap() error {
	return receiver.Err
}

func dbError(err error) *DbError {
	return &DbError{Err: err}
}


func Init(db *sql.DB) (err error) {
	ddls := []string{ managersDDL,atmPlace, clientDDL,servicesDDL,accountsDDL}
	for _, ddl := range ddls {
		_, err = db.Exec(ddl)
		if err != nil {
			return err
		}
	}

	initialData := []string{managersInitialData}
	for _, datum := range initialData {
		_, err = db.Exec(datum)
		if err != nil {
			return err
		}
	}

	return nil
}

// Log IN STUFFS
var clientId int64

func Login(login, password string, db *sql.DB) (bool, error) {
	var dbLogin, dbPassword string

	err := db.QueryRow(
		loginSQL,
		login).Scan(&dbLogin, &dbPassword)

	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}

		return false, queryError(loginSQL, err)
	}

	if dbPassword != password {
		return false, ErrInvalidPass
	}

	return true, nil
}

func LoginClient(login, password string, db *sql.DB) (bool, error) {
	var dbLogin, dbPassword string

	err := db.QueryRow(
		loginCliSQL,
		login).Scan(&clientId,&dbLogin, &dbPassword)

	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}

		return false, queryError(loginCliSQL, err)
	}


	if dbPassword != password {
		return false, ErrInvalidPass
	}

	return true, nil
}

//GETALL STUFFS

func GetAllAccounts(db *sql.DB) (accounts []Account, err error) {
	rows, err := db.Query(getListAccountsSQL,clientId)
	if err != nil {
		return nil, queryError(getListAccountsSQL, err)
	}
	defer func() {
		if innerErr := rows.Close(); innerErr != nil {
			accounts, err = nil, dbError(innerErr)
		}
	}()

	for rows.Next() {
		account := Account{}
		err = rows.Scan(&account.Id,&account.Name, &account.Balance, &account.Client_id)
		if err != nil {
			return nil, dbError(err)
		}
		accounts = append(accounts, account)
	}
	if rows.Err() != nil {
		return nil, dbError(rows.Err())
	}

	return accounts, nil
}

func GetAllAtms(db *sql.DB) (atms []Atm, err error) {
	rows, err := db.Query(getAllAtmsSQL)
	if err != nil {
		return nil, queryError(getAllAtmsSQL, err)
	}
	defer func() {
		if innerErr := rows.Close(); innerErr != nil {
			atms, err = nil, dbError(innerErr)
		}
	}()

	for rows.Next() {
		atm := Atm{}
		err = rows.Scan(&atm.Id, &atm.Name, &atm.Adress)
		if err != nil {
			return nil, dbError(err)
		}
		atms = append(atms, atm)
	}
	if rows.Err() != nil {
		return nil, dbError(rows.Err())
	}

	return atms, nil
}

func GetAllServices(db *sql.DB) (services []Service, err error) {
	rows, err := db.Query(getAllServicesSQL)
	if err != nil {
		return nil, queryError(getAllServicesSQL, err)
	}
	defer func() {
		if innerErr := rows.Close(); innerErr != nil {
			services, err = nil, dbError(innerErr)
		}
	}()

	for rows.Next() {
		service := Service{}
		err = rows.Scan(&service.Id, &service.Name, &service.Balance)
		if err != nil {
			return nil, dbError(err)
		}
		services = append(services, service)
	}
	if rows.Err() != nil {
		return nil, dbError(rows.Err())
	}

	return services, nil
}

//ADD STUFFS

func AddAtm( atmName string,atmAdress string, db *sql.DB) (err error) {

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}
		err = tx.Commit()
	}()


	_, err = tx.Exec(
		insertAtmSQL,

		sql.Named("name", atmName),
		sql.Named("adress", atmAdress),
	)
	if err != nil {
		return err
	}

	return nil
}

func AddClients( clientName string,clientLog string,clientPassword string,clientPhoneNum int64, db *sql.DB) (err error) {

	tx, err := db.Begin()
	if err != nil {
		return err
	}


	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	_, err = tx.Exec(
		insertClientSQL,
		sql.Named("name", clientName),
		sql.Named("login", clientLog),
		sql.Named("password", clientPassword),
		sql.Named("phoneNum", clientPhoneNum),
	)
	if err != nil {
		return err
	}

	return nil
}

func AddService( serviceName string,serviceBalance int, db *sql.DB) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}
		err = tx.Commit()
	}()
	_, err = tx.Exec(
		insertServiceSQL,

		sql.Named("name", serviceName),
		sql.Named("balance", serviceBalance),
	)
	if err != nil {
		return err
	}
	return nil
}

func AddAccount( accountName string, accountBalance int64, accountUser_id int64, db *sql.DB) (err error) {
	
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}
		err = tx.Commit()
	}()
	
	
	_, err = tx.Exec(
		insertAccountSQL,

		sql.Named("name", accountName),
		sql.Named("balance", accountBalance),
		sql.Named("client_id", accountUser_id),
	)
	
	
	if err != nil {
		return err
	}

	return nil
}

//TRANSFER

func TransferMoneyToService(currency int,serviceId int , db *sql.DB) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	var curren int
	err =tx.QueryRow(selectBalanceSQL,clientId,currency).Scan(&curren)
	
	curren=curren-currency

	_, err = tx.Exec(
		updateAccountsSQL,curren,clientId,
	)
	
	
	var curren2 int
	err =tx.QueryRow(selectBalanceServiceSQL,serviceId).Scan(&curren2)
	
	curren2=curren2+currency

	_, err = tx.Exec(
		updateServiceSQL,curren2,serviceId,
	)

	return nil
}

func TransferMoneyWithPhoneNumber(currency int,phoneNumber int , db *sql.DB) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	var curren int
	var clientIdForTransfer int
	err =tx.QueryRow(selectIdNumberPhoneSQL,phoneNumber).Scan(&curren, &clientIdForTransfer)
	if err != nil {
		return err
	}
	
	curren=curren-currency

	_, err = tx.Exec(
		updateBalanceFirstAccountsSQL,curren,clientId,
	)


	var curren2 int
	err =tx.QueryRow( selectNumberPhoneAccountsSQL,clientIdForTransfer).Scan(&curren2)

	curren2=curren2+currency

	_, err = tx.Exec(
		updateBalanceSecondAccountSQL,curren2,clientIdForTransfer,
	)

	return nil
}

func TransferMoneyWithAccountId(currency int,accountsId int , db *sql.DB) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	var curren int
	err =tx.QueryRow(selectBalanceAccountIdSQL,clientId,currency).Scan(&curren)

	curren=curren-currency

	_, err = tx.Exec(
		updateBalanceAccountsIdSQL,curren,clientId,
	)


	var curren2 int
	err =tx.QueryRow(selectBalanceAccountsSQL,accountsId).Scan(&curren2)

	curren2=curren2+currency

	_, err = tx.Exec(
		updateBalanceAccountsSQL,curren2,accountsId,
	)

	return nil
}
