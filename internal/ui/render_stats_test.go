package ui

import (
	"strings"
	"testing"
	"time"

	"strok/internal/domain"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

// TestStatsBarWidthIsStable verifies the bar's rendered width does not change
// as values grow or shrink, so nothing shifts in the HUD while typing.
func TestStatsBarWidthIsStable(t *testing.T) {
	lipgloss.SetColorProfile(termenv.ANSI256)
	th := DefaultTheme()

	snapshots := []domain.Stats{
		{},
		{WPM: 5, Accuracy: 1, Typed: 3, Elapsed: 2 * time.Second},
		{WPM: 123, Accuracy: 0.83, Errors: 12, Typed: 47, Elapsed: 9*time.Minute + 59*time.Second},
		{WPM: 5000, Accuracy: 1, Typed: 1, Elapsed: time.Second}, // first-keystroke spike, clamped
	}

	want := lipgloss.Width(renderStats(th, snapshots[0]))
	for _, s := range snapshots[1:] {
		if got := lipgloss.Width(renderStats(th, s)); got != want {
			t.Errorf("stats bar width = %d for %+v, want stable %d", got, s, want)
		}
	}
}

// TestStatsThresholdColoring verifies WPM/ACC turn green (78) exactly when they
// clear the advance thresholds.
func TestStatsThresholdColoring(t *testing.T) {
	lipgloss.SetColorProfile(termenv.ANSI256)
	th := DefaultTheme()

	below := renderStats(th, domain.Stats{WPM: domain.AdvanceWPM - 1, Accuracy: domain.AdvanceAccuracy - 0.01})
	if strings.Contains(below, "38;5;78") {
		t.Error("below-threshold WPM/ACC should not be green")
	}

	above := renderStats(th, domain.Stats{WPM: domain.AdvanceWPM, Accuracy: domain.AdvanceAccuracy})
	if strings.Count(above, "38;5;78") != 2 {
		t.Errorf("at-threshold WPM and ACC should both be green; got %q", above)
	}
}
