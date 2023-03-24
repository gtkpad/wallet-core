package create_transaction

import (
	"context"
	"testing"

	"github.com/gtkpad/wallet-core/internal/entity"
	"github.com/gtkpad/wallet-core/internal/event"
	"github.com/gtkpad/wallet-core/internal/usecase/mocks"
	"github.com/gtkpad/wallet-core/pkg/events"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)


type TransactionGatewayMock struct {
	mock.Mock
}

func (m *TransactionGatewayMock) Create(transaction *entity.Transaction) error {
	args := m.Called(transaction)
	return args.Error(0)
}

type AccountGatewayMock struct {
	mock.Mock
}

func (m *AccountGatewayMock) Save(account *entity.Account) error {
	args := m.Called(account)
	return args.Error(0)
}

func (m *AccountGatewayMock) FindById(id string) (*entity.Account, error) {
	args := m.Called(id)
	return args.Get(0).(*entity.Account), args.Error(1)
}


func TestCreateTransactionUseCase_Execute(t *testing.T) {
	clientFrom, _ := entity.NewClient("John Doe", "j@j.com")
	clientTo, _ := entity.NewClient("John Doe 2", "j2@j.com")

	mockUow := mocks.UowMock{}
	mockUow.On("Do", mock.Anything, mock.Anything).Return(nil)

	accountFrom := entity.NewAccount(clientFrom)
	accountTo := entity.NewAccount(clientTo)

	accountFrom.Credit(1000)

	inputDTO := CreateTransactionInputDTO{
		AccountIDFrom: accountFrom.ID,
		AccountIDTo: accountTo.ID,
		Amount: 100,
	}


	dispatcher := events.NewEventDispatcher()
	event := event.NewTransactionCreated()
	ctx := context.Background()

	uc := NewCreateTransactionUseCase(&mockUow, dispatcher, event)

	output, err := uc.Execute(ctx, inputDTO)

	assert.Nil(t, err)
	assert.NotNil(t, output.ID)
	mockUow.AssertExpectations(t)
	mockUow.AssertNumberOfCalls(t, "Do", 1)
}