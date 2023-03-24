package main

import (
	"context"
	"database/sql"
	"fmt"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gtkpad/wallet-core/internal/database"
	"github.com/gtkpad/wallet-core/internal/event"
	"github.com/gtkpad/wallet-core/internal/event/handler"
	"github.com/gtkpad/wallet-core/internal/usecase/create_account"
	"github.com/gtkpad/wallet-core/internal/usecase/create_client"
	"github.com/gtkpad/wallet-core/internal/usecase/create_transaction"
	"github.com/gtkpad/wallet-core/internal/web"
	"github.com/gtkpad/wallet-core/internal/web/webserver"
	"github.com/gtkpad/wallet-core/pkg/events"
	"github.com/gtkpad/wallet-core/pkg/kafka"
	"github.com/gtkpad/wallet-core/pkg/uow"
)

func main() {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", "root", "root", "mysql", "3306", "wallet"))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	configMap := &ckafka.ConfigMap{
		"bootstrap.servers": "kafka:29092",
		"group.id":          "wallet",
	}

	kafkaProducer := kafka.NewKafkaProducer(configMap)

	eventDispatcher := events.NewEventDispatcher()
	eventDispatcher.Register("TransactionCreated", handler.NewTransactionCreatedKafkaHandler(kafkaProducer))
	transactionCreatedEvent := event.NewTransactionCreated()
	// eventDispatcher.Register("TransactionCreated", handler)

	clientDb := database.NewClientDB(db)
	accountDb := database.NewAccountDB(db)

	ctx := context.Background()
	uow := uow.NewUow(ctx, db)

	uow.Register("AccountDB", func (tx *sql.Tx) interface {} {
		return database.NewAccountDB(db)
	})

	uow.Register("TransactionDB", func (tx *sql.Tx) interface {} {
		return database.NewTransactionDB(db)
	})

	createClientUseCase := create_client.NewCreateClientUsecase(clientDb)
	createAccountUseCase := create_account.NewCreateAccountUsecase(accountDb, clientDb)
	createTransactionUseCase := create_transaction.NewCreateTransactionUseCase(uow, eventDispatcher, transactionCreatedEvent)

	webServer := webserver.NewWebServer(":8080")

	clientHandler := web.NewWebClientHandler(*createClientUseCase)
	accountHandler := web.NewWebAccountHandler(*createAccountUseCase)
	transactionHandler := web.NewWebTransactionHandler(*createTransactionUseCase)

	webServer.AddHandler("/clients", clientHandler.CreateClient)
	webServer.AddHandler("/accounts", accountHandler.CreateAccount)
	webServer.AddHandler("/transactions", transactionHandler.CreateTransaction)

	fmt.Println("Starting server on port 8080")
	webServer.Start()
}