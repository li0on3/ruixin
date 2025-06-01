package services

import (
	"backend/internal/models"
	"backend/internal/utils"
)

// MaskCardForResponse masks card codes in API responses
func MaskCardForResponse(card *models.Card) {
	if card != nil {
		card.CardCode = utils.MaskCardCode(card.CardCode)
	}
}

// MaskCardsForResponse masks card codes in a slice of cards
func MaskCardsForResponse(cards []*models.Card) {
	for _, card := range cards {
		MaskCardForResponse(card)
	}
}

// MaskCardLogForResponse masks card codes in usage logs
func MaskCardLogForResponse(log *models.CardUsageLog) {
	if log != nil {
		log.CardCode = utils.MaskCardCode(log.CardCode)
		if log.Card != nil {
			MaskCardForResponse(log.Card)
		}
	}
}

// MaskCardLogsForResponse masks card codes in a slice of logs
func MaskCardLogsForResponse(logs []*models.CardUsageLog) {
	for _, log := range logs {
		MaskCardLogForResponse(log)
	}
}