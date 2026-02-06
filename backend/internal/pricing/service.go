package pricing

import (
	"context"
	"fmt"
	"math"

	"github.com/gablelbm/gable/internal/customer"
	"github.com/google/uuid"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CalculatePrice(ctx context.Context, cust *customer.Customer, productID uuid.UUID, basePrice float64) (CalculatedPrice, error) {
	// 1. Check Contract
	contract, err := s.repo.GetContract(ctx, cust.ID, productID)
	if err != nil {
		return CalculatedPrice{}, err
	}
	if contract != nil {
		discountPct := 0.0
		if basePrice > 0 {
			discountPct = (basePrice - contract.ContractPrice) / basePrice * 100
		}
		return CalculatedPrice{
			ProductID:     productID,
			OriginalPrice: basePrice,
			FinalPrice:    contract.ContractPrice,
			DiscountPct:   math.Round(discountPct*100) / 100,
			Source:        SourceContract,
			Details:       "Specific Contract Price",
		}, nil
	}

	// 2. Check Price Level (Tier)
	// Priority: Custom Price Level > Tier Enum > Retail

	multiplier := 1.0
	details := ""
	source := SourceRetail

	if cust.PriceLevel != nil {
		multiplier = cust.PriceLevel.Multiplier
		details = fmt.Sprintf("%s (Level)", cust.PriceLevel.Name)
		source = SourceTier
	} else if cust.Tier != "" && cust.Tier != customer.TierRetail {
		switch cust.Tier {
		case customer.TierSilver:
			multiplier = 0.90
			details = "Silver Tier (10%)"
		case customer.TierGold:
			multiplier = 0.85
			details = "Gold Tier (15%)"
		case customer.TierPlatinum:
			multiplier = 0.80
			details = "Platinum Tier (20%)"
		}
		if multiplier < 1.0 {
			source = SourceTier
		}
	}

	if source == SourceTier {
		final := basePrice * multiplier
		discountPct := (1 - multiplier) * 100
		return CalculatedPrice{
			ProductID:     productID,
			OriginalPrice: basePrice,
			FinalPrice:    final,
			DiscountPct:   math.Round(discountPct*100) / 100,
			Source:        SourceTier,
			Details:       details,
		}, nil
	}

	// 3. Retail
	return CalculatedPrice{
		ProductID:     productID,
		OriginalPrice: basePrice,
		FinalPrice:    basePrice,
		DiscountPct:   0,
		Source:        SourceRetail,
		Details:       "Base Retail Price",
	}, nil
}
