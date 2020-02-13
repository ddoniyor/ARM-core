package core

const managersDDL = `
CREATE TABLE IF NOT EXISTS managers
(
    id      INTEGER PRIMARY KEY AUTOINCREMENT,
    name    TEXT    NOT NULL,
    login   TEXT    NOT NULL UNIQUE,
    password TEXT NOT NULL,
    salary  INTEGER NOT NULL CHECK ( salary > 0 )
);`

const clientDDL = `CREATE TABLE IF NOT EXISTS clients
(
    id     INTEGER PRIMARY KEY AUTOINCREMENT,
	name    TEXT NOT NULL ,
    login   TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
	phoneNum	INTEGER NOT NULL
);`


const atmPlace = `CREATE TABLE IF NOT EXISTS atm
(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL ,
	adress  TEXT NOT NULL
);`


const servicesDDL = `
CREATE TABLE IF NOT EXISTS services
(
   id      INTEGER PRIMARY KEY AUTOINCREMENT,
   name    TEXT    NOT NULL,
   balance INTEGER 
);`

const accountsDDL = `
CREATE TABLE IF NOT EXISTS accounts
(
   id      INTEGER PRIMARY KEY AUTOINCREMENT,
   name    TEXT    NOT NULL,
   balance INTEGER NOT NULL,
   client_id INTEGER  REFERENCES clients(id)
);`


const managersInitialData = `INSERT INTO managers
VALUES (1, 'Vasya', 'vasya', 'secret', 100000),
       (2, 'Petya', 'petya', 'secret', 90000),
       (3, 'Vanya', 'vanya', 'secret', 80000),
       (4, 'Masha', 'masha', 'secret', 80000),
       (5, 'Dasha', 'dasha', 'secret', 60000),
       (6, 'Sasha', 'sasha', 'secret', 40000)
       ON CONFLICT DO NOTHING;`




const loginSQL = `SELECT login, password FROM managers WHERE login = ?`

const loginCliSQL = `SELECT id, login, password FROM clients WHERE login = ?`




const getAllAtmsSQL = `SELECT id, name, adress FROM atm;`

const getAllServicesSQL = `SELECT id, name, balance FROM services;`


const getListAccountsSQL = `SELECT  id, name, balance, client_id FROM accounts WHERE client_id = ?`
 

const  insertClientSQL  =`INSERT INTO clients(name, login, password, phoneNum)VALUES( :name, :login, :password, :phoneNum);`

const insertAtmSQL = `INSERT INTO atm( name , adress)VALUES( :name, :adress);`

const insertServiceSQL = `INSERT INTO services( name , balance) VALUES( :name, :balance);`

const insertAccountSQL = `INSERT INTO accounts( name , balance, client_id) VALUES ( :name, :balance, :client_id);`

//Transfer service

const selectBalanceSQL  = `SELECT balance FROM accounts WHERE client_id = ? AND balance >= ?`

const selectBalanceServiceSQL = `SELECT balance FROM services WHERE id = ? `

const updateAccountsSQL  = `UPDATE accounts SET balance = :balance WHERE client_id = :client_id`

const updateServiceSQL  = `UPDATE services SET balance = ? WHERE id = ?`

const selectIdServiceSQL = `SELECT id FROM services WHERE id = ?`

//Transfer to client

const selectIdNumberPhoneSQL = ` select a.balance, a.id from clients c
                                join accounts a ON a.client_id = c.id where c.phoneNum = ?`

const updateBalanceFirstAccountsSQL  = `UPDATE accounts SET balance = :balance WHERE client_id = :client_id`

const selectNumberPhoneAccountsSQL = `SELECT balance FROM accounts WHERE id = ? `

const updateBalanceSecondAccountSQL  = `UPDATE accounts SET balance = ? WHERE id = ?`

//Transfer to client account
const selectBalanceAccountIdSQL  = `SELECT balance FROM accounts WHERE client_id = ? AND balance >= ?`

const updateBalanceAccountsIdSQL  = `UPDATE accounts SET balance = :balance WHERE client_id = :client_id`

const selectBalanceAccountsSQL = `SELECT balance FROM accounts WHERE id = ? `

const updateBalanceAccountsSQL  = `UPDATE accounts SET balance = ? WHERE id = ?`