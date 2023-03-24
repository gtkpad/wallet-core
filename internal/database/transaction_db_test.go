package database

import (
	"database/sql"
	"testing"

	"github.com/gtkpad/wallet-core/internal/entity"
	"github.com/stretchr/testify/suite"
)

type TransactionDBTestSuite struct {
	suite.Suite
	db *sql.DB
	clientFrom *entity.Client
	clientTo *entity.Client
	accountFrom *entity.Account
	accountTo *entity.Account
	TransactionDB *TransactionDB
}

func (s *TransactionDBTestSuite) SetupSuite() {
	db, err := sql.Open("sqlite3", ":memory:")
	s.Nil(err)
	s.db = db

	_, err = db.Exec("CREATE TABLE clients (id varchar(255), name varchar(255), email varchar(255), created_at date)")
	s.Nil(err)

	_, err = db.Exec("CREATE TABLE accounts (id varchar(255), client_id varchar(255), balance int, created_at date)")
	s.Nil(err)

	_, err = db.Exec("CREATE TABLE transactions (id varchar(255), account_id_from varchar(255), account_id_to varchar(255), amount int, created_at date)")
	s.Nil(err)

	clientFrom, err := entity.NewClient("John Doe", "j@j.com")
	s.Nil(err)
	s.clientFrom = clientFrom

	clientTo, err := entity.NewClient("Jane Doe", "jane@j.com")
	s.Nil(err)
	s.clientTo = clientTo

	accountFrom := entity.NewAccount(clientFrom)
	s.NotNil(accountFrom)
	accountFrom.Credit(1000)
	s.accountFrom = accountFrom

	accountTo := entity.NewAccount(clientTo)
	s.NotNil(accountTo)
	accountTo.Credit(1000)
	s.accountTo = accountTo

	s.TransactionDB = NewTransactionDB(db)
}

func (s *TransactionDBTestSuite) TearDownSuite() {
	defer s.db.Close()
	s.db.Exec("DROP TABLE transactions")
	s.db.Exec("DROP TABLE accounts")
	s.db.Exec("DROP TABLE clients")
}

func TestTransactionDBTestSuite(t *testing.T) {
	suite.Run(t, new(TransactionDBTestSuite))
}

func (s *TransactionDBTestSuite) TestCreate() {
	transaction, err := entity.NewTransaction(s.accountFrom, s.accountTo, 100)
	s.Nil(err)

	err = s.TransactionDB.Create(transaction)
	s.Nil(err)
}