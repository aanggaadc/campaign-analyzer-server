package pdf

import (
	"bytes"
	"fmt"

	"campaign-analyzer/internal/dto"

	"github.com/jung-kurt/gofpdf"
)

type PDFService struct{}

func NewPDFService() *PDFService {
	return &PDFService{}
}

func (p *PDFService) GenerateAnalysisPDF(data dto.AnalysisExportData) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Title
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(0, 10, "Campaign Analysis Report")

	pdf.Ln(12)

	// Campaign Info
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(0, 8, "Campaign: "+data.CampaignName)
	pdf.Ln(6)
	pdf.Cell(0, 8, "Platform: "+data.Platform)

	pdf.Ln(10)

	pdf.Cell(0, 8, fmt.Sprintf("CTR: %.2f", data.CTR))
	pdf.Ln(6)
	pdf.Cell(0, 8, fmt.Sprintf("CPC: %.2f", data.CPC))
	pdf.Ln(6)
	pdf.Cell(0, 8, fmt.Sprintf("CPA: %.2f", data.CPA))

	pdf.Ln(10)

	// Summary
	pdf.MultiCell(0, 6, "Summary:\n"+data.Summary, "", "", false)

	pdf.Ln(5)

	// Issues
	pdf.MultiCell(0, 6, "Issues:\n"+joinLines(data.Issues), "", "", false)

	pdf.Ln(5)

	// Recommendations
	pdf.MultiCell(0, 6, "Recommendations:\n"+joinLines(data.Recommendations), "", "", false)

	pdf.Ln(5)

	// Priority Actions
	pdf.MultiCell(0, 6, "Priority Actions:\n"+joinLines(data.PriorityActions), "", "", false)

	var buf bytes.Buffer
	err := pdf.Output(&buf)

	return buf.Bytes(), err
}

func joinLines(items []string) string {
	result := ""
	for i, v := range items {
		result += fmt.Sprintf("%d. %s\n", i+1, v)
	}
	return result
}
