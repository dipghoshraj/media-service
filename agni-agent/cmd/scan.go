package cmd

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/odio4u/agni-tunnels/agni-agent/pkg/bridge"
	"github.com/spf13/cobra"
)

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan for available Seeders",
	Long: `This command scans the local network for available Seeders and their advertised tunnels.
It uses mDNS to discover Seeders and retrieves their tunnel information for potential connections.`,
	Run: runScan,
}

func runScan(_ *cobra.Command, _ []string) {
	bridge.Logger.Info("scanning for Seeders and tunnels")
	result, err := bridge.ScanForSeeders()
	if err != nil {
		bridge.Logger.Error("failed to scan for Seeders", "error", err)
		return
	}
	printSeedersTable(result)
}

// ── string helpers ────────────────────────────────────────────────────────────

// runeWidth returns the rune count of s for terminal column alignment.
func runeWidth(s string) int {
	return utf8.RuneCountInString(s)
}

// padRight pads s with trailing spaces to the given rune width.
func padRight(s string, width int) string {
	if n := runeWidth(s); n < width {
		return s + strings.Repeat(" ", width-n)
	}
	return s
}

// centerPad centers s within a field of the given rune width.
func centerPad(s string, width int) string {
	n := runeWidth(s)
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
	if w := runeWidth(summaryTitle); w > boxWidth {
		boxWidth = w
	}
	for _, line := range summaryLines {
		if w := runeWidth(line) + 2; w > boxWidth {
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
		widths[i] = runeWidth(h)
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
			if w := runeWidth(cell); w > widths[j] {
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
