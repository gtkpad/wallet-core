package create_transaction

import (
	"context"

	"github.com/gtkpad/wallet-core/internal/entity"
	"github.com/gtkpad/wallet-core/internal/gateway"
	"github.com/gtkpad/wallet-core/pkg/events"
	"github.com/gtkpad/wallet-core/pkg/uow"
)

type CreateTransactionInputDTO struct {
	AccountIDFrom string 		`json:"account_id_from"`
	AccountIDTo 	string		`json:"account_id_to"`
	Amount 				float64 	`json:"amount"`
}

type CreateTransactionOutputDTO struct {
	ID 			string					`json:"id"`
	AccountIDFrom string 		`json:"account_id_from"`
	AccountIDTo 	string		`json:"account_id_to"`
	Amount 				float64 	`json:"amount"`
}

type CreateTransactionUseCase struct {
	Uow uow.UowInterface
	EventDispatcher events.EventDispatcherInterface
	TransactionCreated events.EventInterface
}

func NewCreateTransactionUseCase(
	uow uow.UowInterface,
	eventDispatcher events.EventDispatcherInterface,
	transactionCreated events.EventInterface,
) *CreateTransactionUseCase {
	return &CreateTransactionUseCase{
		Uow: uow,
		EventDispatcher: eventDispatcher,
		TransactionCreated: transactionCreated,
	}
}

func (useCase *CreateTransactionUseCase) Execute(ctx context.Context, input CreateTransactionInputDTO) (*CreateTransactionOutputDTO, error) {
	output := &CreateTransactionOutputDTO{}

	err := useCase.Uow.Do(ctx, func(_ *uow.Uow) error {
		accountRepository := useCase.GetAccountRepository(ctx)
		transactionRepository := useCase.GetTransactionRepository(ctx)
		accountFrom, err := accountRepository.FindById(input.AccountIDFrom)

		if err != nil {
			return err
		}
	
		accountTo, err := accountRepository.FindById(input.AccountIDTo)
		if err != nil {
			return err
		}
	
		transaction, err := entity.NewTransaction(accountFrom, accountTo, input.Amount)
		if err != nil {
			return err
		}

		err = accountRepository.UpdateBalance(accountFrom)
		if err != nil {
			return err
		}

		err = accountRepository.UpdateBalance(accountTo)
		if err != nil {
			return err
		}
	
		err = transactionRepository.Create(transaction)
		if err != nil {
			return err
		}
	
		output.ID = transaction.ID
		output.AccountIDFrom = transaction.AccountFrom.ID
		output.AccountIDTo = transaction.AccountTo.ID
		output.Amount = transaction.Amount

		return nil
	})

	if err != nil {
		return nil, err
	}

	useCase.TransactionCreated.SetPayload(output)
	useCase.EventDispatcher.Dispatch(useCase.TransactionCreated)

	return output, nil
}

func (useCase *CreateTransactionUseCase) GetAccountRepository (ctx context.Context) gateway.AccountGateway {
	repo, err := useCase.Uow.GetRepository(ctx, "AccountDB")

	if err != nil {
		panic(err)
	}

	return repo.(gateway.AccountGateway)
}

func (useCase *CreateTransactionUseCase) GetTransactionRepository (ctx context.Context) gateway.TransactionGateway {
	repo, err := useCase.Uow.GetRepository(ctx, "TransactionDB")

	if err != nil {
		panic(err)
	}

	return repo.(gateway.TransactionGateway)
}

