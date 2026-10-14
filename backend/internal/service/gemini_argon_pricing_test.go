package service

import (
	"context"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/stretchr/testify/require"
)

func TestGeminiArgonGatewayPricingAndOverrides(t *testing.T) {
	bs := NewBillingService(&config.Config{}, &PricingService{pricingData: map[string]*LiteLLMModelPricing{
		"gemini-4-argon": {InputCostPerToken: 99e-6, OutputCostPerToken: 99e-6},
	}})
	for _, model := range []string{
		"gemini-4-argon", "gemini-4-argon-a", "gemini-4-argon-instant",
		"models/gemini-4-argon", "models/gemini-4-argon-a", "models/gemini-4-argon-instant",
		" MODELS/GEMINI-4-ARGON ",
	} {
		t.Run(model, func(t *testing.T) {
			require.True(t, bs.HasIdentifiedTokenPricing(model))
			for _, withResolver := range []bool{false, true} {
				gateway := &GatewayService{billingService: bs}
				if withResolver {
					gateway.resolver = NewModelPricingResolver(nil, bs)
				}
				for _, multiplier := range []float64{0, 0.9, 1.3} {
					cost := gateway.calculateTokenCost(context.Background(), &ForwardResult{
						Usage: ClaudeUsage{InputTokens: 23, OutputTokens: 137, CacheReadInputTokens: 4352},
					}, &APIKey{Group: &Group{ID: 16}}, model, multiplier, time.Time{}, &recordUsageOpts{})
					require.InDelta(t, 0.000046, cost.InputCost, 1e-12)
					require.InDelta(t, 0.001370, cost.OutputCost, 1e-12)
					require.InDelta(t, 0.0004352, cost.CacheReadCost, 1e-12)
					require.InDelta(t, 0.0018512, cost.TotalCost, 1e-12)
					require.InDelta(t, 0.0018512*multiplier, cost.ActualCost, 1e-12)
				}
			}
			for _, tc := range []struct {
				name                 string
				input, output, cache *float64
				want                 float64
			}{
				{name: "partial", input: argonTestPrice(0), want: 10.1},
				{name: "free", input: argonTestPrice(0), output: argonTestPrice(0), cache: argonTestPrice(0), want: 0},
				{name: "custom", input: argonTestPrice(1e-6), output: argonTestPrice(2e-6), cache: argonTestPrice(.3e-6), want: 3.3},
			} {
				t.Run(tc.name, func(t *testing.T) {
					group := &Group{ID: 16, ModelPricing: []ChannelModelPricing{{
						Models: []string{model}, BillingMode: BillingModeToken,
						InputPrice: tc.input, OutputPrice: tc.output, CacheReadPrice: tc.cache,
					}}}
					cost, err := bs.CalculateCostUnified(CostInput{
						Ctx: context.Background(), Model: model, Group: group,
						Resolver: NewModelPricingResolver(nil, bs), RateMultiplier: 1.3,
						Tokens: UsageTokens{InputTokens: 1_000_000, OutputTokens: 1_000_000, CacheReadTokens: 1_000_000},
					})
					require.NoError(t, err)
					require.InDelta(t, tc.want, cost.TotalCost, 1e-12)
					require.InDelta(t, tc.want*1.3, cost.ActualCost, 1e-12)
				})
			}
		})
	}
	for _, model := range []string{"gemini-4-argon-fake", "gemini-4-argon-instant-fake"} {
		require.False(t, bs.HasIdentifiedTokenPricing(model), model)
	}
}

func argonTestPrice(value float64) *float64 { return &value }
