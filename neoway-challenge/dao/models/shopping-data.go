package models

import (
	"time"
)

type (
	ShoppingData struct {
		ID                 uint   `gorm:"primaryKey"`
		CPF                string `gorm:"index" validate:"regexp=^\d{3}\.\d{3}\.\d{3}\-\d{2}$"`
		Private            int
		Incompleto         int
		DataUltimaCompra   *time.Time `gorm:"type:date"`
		TicketMedio        int
		TicketUltimaCompra int
		LojaMaisFrequente  *string `validate:"regexp=^\d{2}\.\d{3}\.\d{3}\/\d{4}\-\d{2}$"`
		LojaUltimaCompra   *string `validate:"regexp=^\d{2}\.\d{3}\.\d{3}\/\d{4}\-\d{2}$"`
	}
)
