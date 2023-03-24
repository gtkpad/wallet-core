package gateway

import "github.com/gtkpad/wallet-core/internal/entity"


type AccountGateway interface {
	Save(account *entity.Account) error
	FindById(id string) (*entity.Account, error)
	UpdateBalance(account *entity.Account) error
}