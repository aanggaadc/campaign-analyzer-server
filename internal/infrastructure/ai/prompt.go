package ai

import (
	"campaign-analyzer/internal/usecase"
	"fmt"
)

func BuildAnalyzePrompt(input usecase.AnalyzeInput) string {
	return fmt.Sprintf(`
You are a digital marketing analyst.

Analyze the following campaign data.

Campaign:
- Name: %s
- Platform: %s
- Impressions: %d
- Clicks: %d
- Conversions: %d
- Cost: %.2f
- CTR: %.4f
- CPC: %.2f
- CPA: %.2f

Return STRICTLY valid JSON.

Rules:
- Do NOT return null values.
- If no data, return empty array [] instead.
- All fields MUST be present.
- "priority_actions" MUST contain at least 1 actionable item.

{
  "summary": "string",
  "issues": ["string"],
  "recommendations": ["string"],
  "priority_actions": ["string"]
}
`,
		input.Name,
		input.Platform,
		input.Impressions,
		input.Clicks,
		input.Conversions,
		input.Cost,
		input.CTR,
		input.CPC,
		input.CPA,
	)
}
