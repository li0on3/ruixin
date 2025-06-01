package services

import (
	"backend/internal/models"
	"backend/internal/utils"
	"errors"
	"strings"
)

type SecurityService struct{}

func NewSecurityService() *SecurityService {
	return &SecurityService{}
}

// MaskCard masks sensitive card information
func (s *SecurityService) MaskCard(card *models.Card) *models.Card {
	if card == nil {
		return nil
	}
	
	// Create a copy to avoid modifying the original
	maskedCard := *card
	maskedCard.CardCode = utils.MaskCardCode(card.CardCode)
	
	return &maskedCard
}

// MaskCards masks a slice of cards
func (s *SecurityService) MaskCards(cards []*models.Card) []*models.Card {
	maskedCards := make([]*models.Card, len(cards))
	for i, card := range cards {
		maskedCards[i] = s.MaskCard(card)
	}
	return maskedCards
}

// MaskCardInLog masks card code in a log entry
func (s *SecurityService) MaskCardInLog(log *models.CardUsageLog) *models.CardUsageLog {
	if log == nil {
		return nil
	}
	
	maskedLog := *log
	maskedLog.CardCode = utils.MaskCardCode(log.CardCode)
	if maskedLog.Card != nil {
		maskedLog.Card = s.MaskCard(log.Card)
	}
	
	return &maskedLog
}

// SanitizeOrderResponse removes sensitive information from order response
func (s *SecurityService) SanitizeOrderResponse(order *models.Order) map[string]interface{} {
	return map[string]interface{}{
		"order_no":      order.OrderNo,
		"status":        order.Status,
		"total_amount":  order.TotalAmount,
		"store_name":    order.StoreName,
		"store_address": order.StoreAddress,
		"created_at":    order.CreatedAt,
		"updated_at":    order.UpdatedAt,
		// Explicitly exclude: card_code, card_id, cost_amount, profit_amount
	}
}

// ValidateAndSanitizeInput validates and sanitizes user input
func (s *SecurityService) ValidateAndSanitizeInput(input string, inputType string) (string, error) {
	// Remove dangerous characters
	sanitized := utils.SanitizeInput(input)
	
	// Validate based on type
	switch inputType {
	case "card_code":
		if !utils.IsValidCardCode(sanitized) {
			return "", models.ErrInvalidCardCode
		}
	case "phone":
		// Validate phone number format
		if !s.isValidPhoneNumber(sanitized) {
			return "", models.ErrInvalidPhoneNumber
		}
	case "order_no":
		// Validate order number format
		if !s.isValidOrderNo(sanitized) {
			return "", models.ErrInvalidOrderNo
		}
	}
	
	return sanitized, nil
}

func (s *SecurityService) isValidPhoneNumber(phone string) bool {
	// Chinese phone number validation
	if len(phone) != 11 {
		return false
	}
	
	// Must start with 1
	if phone[0] != '1' {
		return false
	}
	
	// Check if all digits
	for _, ch := range phone {
		if ch < '0' || ch > '9' {
			return false
		}
	}
	
	return true
}

func (s *SecurityService) isValidOrderNo(orderNo string) bool {
	// Order number should be alphanumeric with possible hyphens
	if len(orderNo) < 10 || len(orderNo) > 50 {
		return false
	}
	
	for _, ch := range orderNo {
		if !((ch >= 'A' && ch <= 'Z') || (ch >= 'a' && ch <= 'z') || 
			(ch >= '0' && ch <= '9') || ch == '-' || ch == '_') {
			return false
		}
	}
	
	return true
}

// RemoveCardCodesFromError removes any card codes from error messages
func (s *SecurityService) RemoveCardCodesFromError(err error) error {
	if err == nil {
		return nil
	}
	
	errMsg := err.Error()
	// Replace any string that looks like a card code with masked version
	// This is a simple implementation - in production, use more sophisticated pattern matching
	words := strings.Fields(errMsg)
	for i, word := range words {
		if utils.IsValidCardCode(word) {
			words[i] = utils.MaskCardCode(word)
		}
	}
	
	return errors.New(strings.Join(words, " "))
}