package service

import (
	"strings"

	"github.com/tidwall/gjson"
)

// normalizeAnthropicBillingModel recognizes provider, context, snapshot and
// thinking aliases without mistaking a different numeric version for an old one.
func normalizeAnthropicBillingModel(model string) string {
	id := strings.ToLower(strings.TrimSpace(model))
	if slash := strings.LastIndexByte(id, '/'); slash >= 0 {
		id = id[slash+1:]
	}
	id = strings.TrimPrefix(id, "anthropic.")
	if colon := strings.IndexByte(id, ':'); colon >= 0 {
		id = id[:colon]
	}
	id = strings.TrimSuffix(id, "[1m]")
	id = strings.TrimSuffix(id, "-thinking")
	id = strings.TrimSuffix(id, "-latest")
	if len(id) > 9 && id[len(id)-9] == '-' {
		suffix := id[len(id)-8:]
		if strings.Trim(suffix, "0123456789") == "" {
			id = id[:len(id)-9]
		}
	}
	return strings.ReplaceAll(id, ".", "-")
}

// JWS supplier standard cards before the independent group/user multiplier.
// Opus 5.5, Sonnet 5 and Haiku 4.5 input/output were replayed against usage
// records on 2026-10-11. Fable and older versions use the supplier's public card:
// https://web.aijws.com/models
// Cache TTL ratios for models without matching usage remain provisional.
// The public cards say "all" context lengths; do not inherit an older ladder.
// Use Fable's public 12/60 price directly, with no additional retail premium.
// Explicit group/channel cards still take precedence without a hidden premium.
func anthropicJWSBillingCard(model string) *ModelPricing {
	var input, output, read float64
	switch normalizeAnthropicBillingModel(model) {
	case "claude-opus-5-5":
		input, output, read = 4e-6, 20e-6, .2e-6
	case "claude-fable-5":
		input, output, read = 12e-6, 60e-6, 1.2e-6
	case "claude-fable-5-1":
		input, output, read = 12e-6, 60e-6, .3e-6
	case "claude-sonnet-5":
		input, output, read = 2e-6, 10e-6, .2e-6
	case "claude-sonnet-5-5":
		// The deployed model directory agrees with Sonnet 5; no JWS usage
		// sample for this exact version was available in the last 30 days.
		input, output, read = 2e-6, 10e-6, .2e-6
	case "claude-haiku-4-5":
		input, output, read = 1e-6, 5e-6, .1e-6
	case "claude-opus-5", "claude-opus-4-8":
		input, output, read = 5e-6, 25e-6, .5e-6
	case "claude-sonnet-4-6":
		input, output, read = 3e-6, 15e-6, .3e-6
	default:
		return nil
	}
	p := &ModelPricing{
		InputPricePerToken: input, OutputPricePerToken: output, CacheReadPricePerToken: read,
		CacheCreationPricePerToken: input * 1.25, CacheCreation5mPrice: input * 1.25,
		CacheCreation1hPrice: input * 2, SupportsCacheBreakdown: true,
	}
	if normalizeAnthropicBillingModel(model) == "claude-opus-5-5" {
		// Retain the existing Fast policy; JWS samples only verified Standard.
		// A whole-card multiplier also scales the separate 5m/1h cache prices.
		fastMultiplier := 2.0
		p.FastMultiplier = &fastMultiplier
	}
	return p
}

// Preserve the existing provider-reference column separately from supplier
// pricing. Fable's reference excludes JWS's 20% markup; request billing uses
// the supplier card directly. Unverified versions are not usage evidence.
func anthropicReferenceBillingCard(model string) *ModelPricing {
	p := anthropicJWSBillingCard(model)
	if p == nil {
		return nil
	}
	switch normalizeAnthropicBillingModel(model) {
	case "claude-fable-5", "claude-fable-5-1":
		p.InputPricePerToken /= 1.2
		p.OutputPricePerToken /= 1.2
		p.CacheReadPricePerToken /= 1.2
		p.CacheCreationPricePerToken /= 1.2
		p.CacheCreation5mPrice /= 1.2
		p.CacheCreation1hPrice /= 1.2
	}
	return p
}

// WithRequestBillingEffort carries the current request's effort into pre-holds.
// Final settlement uses the gateway's effective effort after any mapping.
func WithRequestBillingEffort(input CostInput, body []byte) CostInput {
	for _, path := range []string{"output_config.effort", "reasoning.effort", "reasoning_effort"} {
		if raw := gjson.GetBytes(body, path); raw.Exists() {
			input.ReasoningEffort = strings.ToLower(strings.TrimSpace(raw.String()))
			break
		}
	}
	return input
}
