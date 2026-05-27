package cmd

import (
	"fmt"
	"strings"

	"github.com/mattn/go-runewidth"
	"github.com/odio4u/agni-tunnels/agni-agent/pkg/bridge"
	"github.com/spf13/cobra"
)

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan for available Seeders",
	Long: `This command queries the configured registry for available Seeders.
It lists the Seeders returned by the registry so you can identify potential connections.`,
	Run: runScan,
}

func runScan(_ *cobra.Command, _ []string) {
	result, err := bridge.ScanForSeeders()
	if err != nil {
		bridge.Logger.Error("failed to scan for Seeders", "error", err)
		return
	}
	printSeedersTable(result)
}

// ── string helpers ────────────────────────────────────────────────────────────

// isRegionalIndicator reports whether r is a Unicode Regional Indicator Symbol
// (U+1F1E6–U+1F1FF), used in pairs to encode country flag emoji.
func isRegionalIndicator(r rune) bool {
	return r >= 0x1F1E6 && r <= 0x1F1FF
}

// displayWidth returns the terminal display width of s for column alignment.
// go-runewidth counts each Regional Indicator pair (flag emoji like 🇮🇳) as
// width 1, but most terminals render the combined flag as 2 columns. We
// add 1 correction per detected pair so padding stays aligned.
func displayWidth(s string) int {
	w := runewidth.StringWidth(s)
	runes := []rune(s)
	for i := 0; i+1 < len(runes); i++ {
		if isRegionalIndicator(runes[i]) && isRegionalIndicator(runes[i+1]) {
			w++
			i++ // skip the second indicator — it's already part of this pair
		}
	}
	return w
}

// padRight pads s with trailing spaces until its display width reaches width.
func padRight(s string, width int) string {
	if n := displayWidth(s); n < width {
		return s + strings.Repeat(" ", width-n)
	}
	return s
}

// centerPad centers s within a field of the given display width.
func centerPad(s string, width int) string {
	n := displayWidth(s)
	if n >= width {
		return s
	}
	pad := width - n
	return strings.Repeat(" ", pad/2) + s + strings.Repeat(" ", pad-pad/2)
}

// ── display formatters ────────────────────────────────────────────────────────

// statusLabel formats a seeder status string into a human-readable indicator.
func statusLabel(status string) string {
	if strings.EqualFold(status, "online") {
		return "● ONLINE"
	}
	return "○ OFFLINE"
}

// regionLabel formats a SeederRegion into a single display string.
func regionLabel(r bridge.SeederRegion) string {
	if r.Flag != "" {
		return r.Flag + " " + r.Country
	}
	return r.Country
}

// ── table renderer ────────────────────────────────────────────────────────────

func printSeedersTable(result bridge.SeederInfoResponse) {
	online := 0
	for _, s := range result.Seeders {
		if strings.EqualFold(s.Status, "online") {
			online++
		}
	}

	// summary header box
	summaryTitle := "AGNISTACK SEEDERS"
	summaryLines := []string{
		"Network      : " + result.Network,
		fmt.Sprintf("Online Nodes : %d / %d", online, len(result.Seeders)),
		"Discovery    : Decentralized",
	}
	boxWidth := 46
	if w := displayWidth(summaryTitle); w > boxWidth {
		boxWidth = w
	}
	for _, line := range summaryLines {
		if w := displayWidth(line) + 2; w > boxWidth {
			boxWidth = w
		}
	}

	fmt.Printf("┌%s┐\n", strings.Repeat("─", boxWidth))
	fmt.Printf("│%s│\n", centerPad(summaryTitle, boxWidth))
	fmt.Printf("├%s┤\n", strings.Repeat("─", boxWidth))
	for _, line := range summaryLines {
		fmt.Printf("│ %s │\n", padRight(line, boxWidth-2))
	}
	fmt.Printf("└%s┘\n\n", strings.Repeat("─", boxWidth))

	if len(result.Seeders) == 0 {
		fmt.Println("  No seeders found.")
		return
	}

	// build rows and compute column widths from content
	headers := []string{"#", "IP", "STATUS", "REGION", "MAINTAINER"}
	widths := make([]int, len(headers))
	for i, h := range headers {
		widths[i] = displayWidth(h)
	}

	rows := make([][]string, len(result.Seeders))
	for i, s := range result.Seeders {
		rows[i] = []string{
			fmt.Sprintf("%d", i+1),
			s.IP,
			statusLabel(s.Status),
			regionLabel(s.Region),
			s.Maintainer,
		}
		for j, cell := range rows[i] {
			if w := displayWidth(cell); w > widths[j] {
				widths[j] = w
			}
		}
	}

	// horizontal rule: left/mid/right are box-drawing junction characters
	hline := func(left, mid, right string) {
		fmt.Print(left)
		for i, w := range widths {
			fmt.Print(strings.Repeat("─", w+2))
			if i < len(widths)-1 {
				fmt.Print(mid)
			}
		}
		fmt.Println(right)
	}

	printRow := func(cells []string) {
		fmt.Print("│")
		for i, cell := range cells {
			fmt.Printf(" %s │", padRight(cell, widths[i]))
		}
		fmt.Println()
	}

	hline("┌", "┬", "┐")
	printRow(headers)
	hline("├", "┼", "┤")
	for _, row := range rows {
		printRow(row)
	}
	hline("└", "┴", "┘")
}

func init() {
	rootCmd.AddCommand(scanCmd)
}
