package tools

import "github.com/wneessen/sotbot/api"

func CompareSotStats(a, b api.UserStats) bool {
	changed := false
	if &a == &b {
		return true
	}

	if a.KrakenDefeated != b.KrakenDefeated {
		changed = true
	}
	if a.MegalodonEncounters != b.MegalodonEncounters {
		changed = true
	}
	if a.ChestsHandedIn != b.ChestsHandedIn {
		changed = true
	}
	if a.ShipsSunk != b.ShipsSunk {
		changed = true
	}
	if a.VomitedTotal != b.VomitedTotal {
		changed = true
	}

	return changed
}
