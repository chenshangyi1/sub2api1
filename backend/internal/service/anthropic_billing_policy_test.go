//go:build unit

package service

import (
	"context"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/stretchr/testify/require"
)

func TestAnthropicJWSCardsAndAliases(t *testing.T) {
	// A stale catalog must not silently substitute an older version's price or ladder.
	stale := &LiteLLMModelPricing{InputCostPerToken: 5e-6, OutputCostPerToken: 25e-6, LongContextInputTokenThreshold: 200000, LongContextInputCostMultiplier: 2, LongContextOutputCostMultiplier: 1.5}
	bs := NewBillingService(&config.Config{}, &PricingService{pricingData: map[string]*LiteLLMModelPricing{
		"claude-opus-5": stale, "claude-opus-5-5": stale, "claude-fable-5-1": stale,
	}})
	for _, tc := range []struct {
		model                            string
		input, output, write, hour, read float64
	}{
		{"claude-opus-5-5", 4e-6, 20e-6, 5e-6, 8e-6, .2e-6},
		{"anthropic/Claude-Opus-5.5-20260901", 4e-6, 20e-6, 5e-6, 8e-6, .2e-6},
		{"claude-fable-5", 12e-6, 60e-6, 15e-6, 24e-6, 1.2e-6},
		{"claude-fable-5[1m]", 12e-6, 60e-6, 15e-6, 24e-6, 1.2e-6},
		{"claude-fable-5-1", 12e-6, 60e-6, 15e-6, 24e-6, .3e-6},
		{"anthropic.claude-fable-5.1:0", 12e-6, 60e-6, 15e-6, 24e-6, .3e-6},
		{"claude-sonnet-5-5", 2e-6, 10e-6, 2.5e-6, 4e-6, .2e-6},
		{"claude-opus-5", 5e-6, 25e-6, 6.25e-6, 10e-6, .5e-6},
		{"claude-opus-4.8", 5e-6, 25e-6, 6.25e-6, 10e-6, .5e-6},
		{"claude-sonnet-4.6", 3e-6, 15e-6, 3.75e-6, 6e-6, .3e-6},
	} {
		p, err := bs.GetModelPricing(tc.model)
		require.NoError(t, err, tc.model)
		require.True(t, bs.HasIdentifiedTokenPricing(tc.model), tc.model)
		require.InDelta(t, tc.input, p.InputPricePerToken, 1e-15, tc.model)
		require.InDelta(t, tc.output, p.OutputPricePerToken, 1e-15, tc.model)
		require.InDelta(t, tc.write, p.CacheCreationPricePerToken, 1e-15, tc.model)
		require.InDelta(t, tc.hour, p.CacheCreation1hPrice, 1e-15, tc.model)
		require.InDelta(t, tc.read, p.CacheReadPricePerToken, 1e-15, tc.model)
		require.True(t, p.SupportsCacheBreakdown)
		require.Zero(t, p.LongContextInputThreshold)
	}
	require.Equal(t, 5e-6, stale.InputCostPerToken)
}

func TestFableGatewayAllTokenCostsAndEffectiveEffort(t *testing.T) {
	bs := NewBillingService(&config.Config{}, nil)
	for _, useResolver := range []bool{false, true} {
		gateway := &GatewayService{billingService: bs}
		if useResolver {
			gateway.resolver = NewModelPricingResolver(nil, bs)
		}
		group := &Group{ID: 1, Platform: PlatformAnthropic, RateMultiplier: 1.1, LongContextPricingEnabled: true}
		for _, tc := range []struct {
			effort                                    string
			input, output, write, read, total, actual float64
		}{
			{"", 12, 60, 39, .3, 111.3, 122.43},
			{"xhigh", 12, 60, 39, .3, 111.3, 122.43},
			{"max", 12, 60, 39, .3, 111.3, 122.43},
		} {
			requested := "max" // Effective xhigh must not be charged as requested max.
			cost := gateway.calculateTokenCost(context.Background(), &ForwardResult{
				ReasoningEffort: &tc.effort, RequestedReasoningEffort: &requested,
				Usage: ClaudeUsage{InputTokens: 1000000, OutputTokens: 1000000, CacheReadInputTokens: 1000000, CacheCreationInputTokens: 2000000, CacheCreation5mTokens: 1000000, CacheCreation1hTokens: 1000000},
			}, &APIKey{Group: group}, "claude-fable-5-1", 1.1, time.Time{}, &recordUsageOpts{})
			require.InDelta(t, tc.input, cost.InputCost, 1e-10)
			require.InDelta(t, tc.output, cost.OutputCost, 1e-10)
			require.InDelta(t, tc.write, cost.CacheCreationCost, 1e-10)
			require.InDelta(t, tc.read, cost.CacheReadCost, 1e-10)
			require.InDelta(t, tc.total, cost.TotalCost, 1e-10)
			require.InDelta(t, tc.actual, cost.ActualCost, 1e-10)
			require.False(t, cost.LongContextBillingApplied)
			require.Equal(t, 1.1, group.RateMultiplier)
		}
	}
}

func TestFableJWSDoesNotAddPremiumToCustomCard(t *testing.T) {
	bs := NewBillingService(&config.Config{}, nil)
	gateway := &GatewayService{billingService: bs, resolver: NewModelPricingResolver(nil, bs)}
	for _, tc := range []struct {
		effort string
		want   float64
	}{{"max", 2.2}, {"high", 2.2}, {"max", 2.2}} {
		t.Run(tc.effort, func(t *testing.T) {
			t.Parallel()
			price := 2e-6
			group := &Group{ID: 1, RateMultiplier: 1.1, ModelPricing: []ChannelModelPricing{{Models: []string{"claude-fable-5"}, BillingMode: BillingModeToken, InputPrice: &price}}}
			cost := gateway.calculateTokenCost(context.Background(), &ForwardResult{ReasoningEffort: &tc.effort, Usage: ClaudeUsage{InputTokens: 1000000}}, &APIKey{Group: group}, "claude-fable-5", 1.1, time.Time{}, &recordUsageOpts{})
			require.InDelta(t, tc.want, cost.ActualCost, 1e-10)
			p, err := bs.GetModelPricing("claude-fable-5")
			require.NoError(t, err)
			require.InDelta(t, 12e-6, p.InputPricePerToken, 1e-15)
			require.Equal(t, 2e-6, price)
		})
	}
}

func TestOpus55FlatPricingAndFablePolicyDoesNotAffectOtherModels(t *testing.T) {
	bs := NewBillingService(&config.Config{}, nil)
	gateway := &GatewayService{billingService: bs, resolver: NewModelPricingResolver(nil, bs)}
	effort := "max"
	for _, n := range []int{199999, 200000, 200001, 900000} {
		cost := gateway.calculateTokenCost(context.Background(), &ForwardResult{ReasoningEffort: &effort, Usage: ClaudeUsage{InputTokens: n, OutputTokens: 1000000}}, &APIKey{Group: &Group{ID: 1, LongContextPricingEnabled: true}}, "claude-opus-5.5", 1.1, time.Time{}, &recordUsageOpts{})
		require.InDelta(t, float64(n)*4e-6, cost.InputCost, 1e-10)
		require.InDelta(t, 20, cost.OutputCost, 1e-10)
		require.InDelta(t, (float64(n)*4e-6+20)*1.1, cost.ActualCost, 1e-10)
		require.False(t, cost.LongContextBillingApplied)
	}
	cost := gateway.calculateTokenCost(context.Background(), &ForwardResult{ReasoningEffort: &effort, Usage: ClaudeUsage{InputTokens: 1000000}}, &APIKey{Group: &Group{ID: 1}}, "claude-sonnet-4", 1.1, time.Time{}, &recordUsageOpts{})
	require.InDelta(t, 3.3, cost.ActualCost, 1e-10)
}

func TestFableAdaptiveHoldPricesMaxLikeSettlement(t *testing.T) {
	bs := NewBillingService(&config.Config{}, nil)
	body := []byte(`{"model":"claude-fable-5-1","messages":[{"role":"user","content":"hi"}],"max_tokens":128,"output_config":{"effort":"max"}}`)
	held, err := EstimateAdaptiveHoldBaseFromPreauthorization(bs, CostInput{Model: "claude-fable-5-1", RateMultiplier: 1.1}, body)
	require.NoError(t, err)
	// The estimator's floor is 500 input tokens; 128 output tokens.
	value, _ := held.Float64()
	require.InDelta(t, .015048, value, 1e-6)
}

func TestFableRequestEffortAcrossProtocolsAndFreeGroup(t *testing.T) {
	bs := NewBillingService(&config.Config{}, nil)
	for _, body := range []string{
		`{"output_config":{"effort":" MAX "}}`,
		`{"reasoning":{"effort":"max"}}`,
		`{"reasoning_effort":"max"}`,
	} {
		input := WithRequestBillingEffort(CostInput{Model: "claude-fable-5-1", RateMultiplier: 1.1, Tokens: UsageTokens{InputTokens: 1000000}}, []byte(body))
		cost, err := bs.CalculateCostUnified(input)
		require.NoError(t, err)
		require.InDelta(t, 13.2, cost.ActualCost, 1e-10)
		input.RateMultiplier = 0
		free, err := bs.CalculateCostUnified(input)
		require.NoError(t, err)
		require.Zero(t, free.ActualCost)
	}
}

func TestFablePreauthorizationOutputTopupsUseMaxPrice(t *testing.T) {
	bs := NewBillingService(&config.Config{}, nil)
	preauth := &BalancePreauthorizationService{costCalculator: bs}
	estimate, err := preauth.estimateHold(context.Background(), BalancePreauthorizationRequest{
		EstimatedInputTokens: 500, InitialOutputWindowTokens: 128,
		CostInput: CostInput{Model: "claude-fable-5-1", RateMultiplier: 1.1, ReasoningEffort: "max"},
	})
	require.NoError(t, err)
	require.InDelta(t, .015048, estimate.HoldAmount, 1e-7)
	require.InDelta(t, .000066, estimate.OutputUnitPrice, 1e-12)
}

func TestFableOpenAICompatibleSettlement(t *testing.T) {
	bs := NewBillingService(&config.Config{}, nil)
	for _, useResolver := range []bool{false, true} {
		gateway := &OpenAIGatewayService{billingService: bs}
		if useResolver {
			gateway.resolver = NewModelPricingResolver(nil, bs)
		}
		cost, err := gateway.calculateOpenAIRecordUsageTokenCost(context.Background(), &APIKey{Group: &Group{ID: 1}}, "anthropic/claude-fable-5-1", 1.1, time.Time{}, UsageTokens{InputTokens: 1000000, CacheReadTokens: 1000000}, "", "max", nil)
		require.NoError(t, err)
		require.InDelta(t, 13.53, cost.ActualCost, 1e-10)
	}
}

func TestOpus55FastScalesBothCacheTTLs(t *testing.T) {
	bs := NewBillingService(&config.Config{}, nil)
	cost, err := bs.CalculateTokenCostForRequest(TokenCostRequest{Model: "claude-opus-5-5", RateMultiplier: 1.1, ServiceTier: "fast", Tokens: UsageTokens{InputTokens: 1000000, OutputTokens: 1000000, CacheReadTokens: 1000000, CacheCreationTokens: 2000000, CacheCreation5mTokens: 1000000, CacheCreation1hTokens: 1000000}})
	require.NoError(t, err)
	require.InDelta(t, 8, cost.InputCost, 1e-10)
	require.InDelta(t, 40, cost.OutputCost, 1e-10)
	require.InDelta(t, 26, cost.CacheCreationCost, 1e-10)
	require.InDelta(t, .4, cost.CacheReadCost, 1e-10)
	require.InDelta(t, 81.84, cost.ActualCost, 1e-10)
}

func TestFablePlazaRetailAndOfficialPrices(t *testing.T) {
	bs := NewBillingService(&config.Config{}, nil)
	resolver := NewModelPricingResolver(nil, bs)
	retail, err := bs.ResolveContextPricingSchedule(context.Background(), resolver, ContextPricingScheduleInput{Model: "claude-fable-5-1", Group: &Group{ID: 1, LongContextPricingEnabled: true}})
	require.NoError(t, err)
	require.Len(t, retail.Tiers, 1)
	require.InDelta(t, 12e-6, *retail.Tiers[0].Input, 1e-15)
	require.InDelta(t, .3e-6, *retail.Tiers[0].CacheRead, 1e-15)
	official, err := bs.ResolveContextPricingSchedule(context.Background(), resolver, ContextPricingScheduleInput{Model: "claude-fable-5-1"})
	require.NoError(t, err)
	require.Len(t, official.Tiers, 1)
	require.InDelta(t, 10e-6, *official.Tiers[0].Input, 1e-15)
	require.InDelta(t, .25e-6, *official.Tiers[0].CacheRead, 1e-15)
}

func TestFableJWSWithChannelIntervalAndTimePrice(t *testing.T) {
	bs := NewBillingService(&config.Config{}, nil)
	resolver := NewModelPricingResolver(nil, bs)
	highInput := 2e-6
	resolved := &ResolvedPricing{Mode: BillingModeToken, Source: PricingSourceChannel, BasePricing: anthropicJWSBillingCard("claude-fable-5"), Intervals: []PricingInterval{{MinTokens: 200000, InputPrice: &highInput}}, longContextPricingEnabled: true, channelPricing: &ChannelModelPricing{TimePricing: &ChannelTimePricing{Timezone: "UTC", Periods: []ChannelTimePricingPeriod{{StartTime: "00:00", EndTime: "23:59", Multiplier: 1.5}}}}}
	cost, err := bs.CalculateCostUnified(CostInput{Ctx: context.Background(), Model: "claude-fable-5", Tokens: UsageTokens{InputTokens: 300000}, ReasoningEffort: "max", RateMultiplier: 1.1, Resolver: resolver, Resolved: resolved, PricingAt: time.Date(2026, 10, 11, 12, 0, 0, 0, time.UTC)})
	require.NoError(t, err)
	require.InDelta(t, .9, cost.TotalCost, 1e-10)
	require.InDelta(t, .99, cost.ActualCost, 1e-10)
}

// Exact supplier fixtures catch wrong TTL pricing, duplicate cached input and
// repeated application of the independent group multiplier.
func TestAnthropicJWSActualUsageReplay(t *testing.T) {
	bs := NewBillingService(&config.Config{}, nil)
	for _, tc := range []struct {
		name, model                                     string
		tokens                                          UsageTokens
		rate, input, output, write, read, total, actual float64
	}{
		{"opus_5m_20261011_205130", "claude-opus-5-5", UsageTokens{InputTokens: 28, OutputTokens: 216, CacheCreationTokens: 2451, CacheCreation5mTokens: 2451, CacheReadTokens: 130469}, 1.1, .000112, .004320, .012255, .0260938, .0427808, .04705888},
		{"opus_1h_20261011_203712", "claude-opus-5-5", UsageTokens{InputTokens: 28, OutputTokens: 35790, CacheCreationTokens: 1044, CacheCreation1hTokens: 1044, CacheReadTokens: 136353}, 1.1, .000112, .7158, .008352, .0272706, .7515346, .82668806},
		{"sonnet_5_20261001_010200", "claude-sonnet-5", UsageTokens{InputTokens: 94, OutputTokens: 9, CacheCreationTokens: 215, CacheCreation5mTokens: 215, CacheReadTokens: 88327}, .85, .000188, .00009, .0005375, .0176654, .0184809, .015708765},
		{"haiku_45_20260930_225644", "claude-haiku-4-5", UsageTokens{InputTokens: 843, OutputTokens: 40}, .85, .000843, .0002, 0, 0, .001043, .00088655},
	} {
		for _, useResolver := range []bool{false, true} {
			t.Run(tc.name+map[bool]string{false: "/legacy", true: "/resolver"}[useResolver], func(t *testing.T) {
				gateway := &GatewayService{billingService: bs}
				if useResolver {
					gateway.resolver = NewModelPricingResolver(nil, bs)
				}
				group := &Group{ID: 1, Platform: PlatformAnthropic, RateMultiplier: tc.rate, LongContextPricingEnabled: true}
				u := tc.tokens
				cost := gateway.calculateTokenCost(context.Background(), &ForwardResult{Usage: ClaudeUsage{InputTokens: u.InputTokens, OutputTokens: u.OutputTokens, CacheCreationInputTokens: u.CacheCreationTokens, CacheCreation5mTokens: u.CacheCreation5mTokens, CacheCreation1hTokens: u.CacheCreation1hTokens, CacheReadInputTokens: u.CacheReadTokens}}, &APIKey{Group: group}, tc.model, tc.rate, time.Time{}, &recordUsageOpts{})
				require.InDelta(t, tc.input, cost.InputCost, 1e-12)
				require.InDelta(t, tc.output, cost.OutputCost, 1e-12)
				require.InDelta(t, tc.write, cost.CacheCreationCost, 1e-12)
				require.InDelta(t, tc.read, cost.CacheReadCost, 1e-12)
				require.InDelta(t, tc.total, cost.TotalCost, 1e-12)
				require.InDelta(t, tc.actual, cost.ActualCost, 1e-12)
				require.Equal(t, tc.rate, group.RateMultiplier)
			})
		}
	}
}

func TestFableJWSPlazaKeepsReferenceSeparate(t *testing.T) {
	bs := NewBillingService(&config.Config{}, nil)
	plaza := &ModelPlazaService{billingService: bs, resolver: NewModelPricingResolver(nil, bs)}
	reference := plaza.lookupOfficialPricing(context.Background(), "claude-fable-5-1", map[string]*PlazaOfficialPricing{})
	require.NotNil(t, reference)
	require.InDelta(t, 10e-6, *reference.InputPrice, 1e-15)
	require.InDelta(t, .25e-6, *reference.CacheReadPrice, 1e-15)
}
