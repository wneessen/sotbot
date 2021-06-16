package response

import "strings"

func Icon(i string) string {
	i = removeSortString(i)
	iconMap := map[string]string{
		"Gold":        "\U0001F7E1",
		"Doubloon":    "🔵",
		"AncientCoin": "💰",
		"Kraken":      "🐙",
		"Megalodon":   "🦈",
		"Chest":       "🎁",
		"Ship":        "⛵",
		"Vomit":       "🤮",
	}

	if iconMap[i] != "" {
		return iconMap[i]
	}
	return "❌"
}

func IconKey(n string) string {
	n = removeSortString(n)
	keyMap := map[string]string{
		"Gold":        "Gold",
		"Doubloon":    "Doubloons",
		"AncientCoin": "Ancient Coins",
		"Kraken":      "Kraken",
		"Megalodon":   "Megalodon",
		"Chest":       "Chests",
		"Ship":        "Ships",
		"Vomit":       "Vomitted",
	}

	if keyMap[n] != "" {
		return keyMap[n]
	}
	return ""
}

func IconValue(n string) string {
	n = removeSortString(n)
	valMap := map[string]string{
		"Gold":        "Gold",
		"Doubloon":    "Doubloon(s)",
		"AncientCoin": "Coin(s)",
		"Kraken":      "defeated",
		"Megalodon":   "encounter(s)",
		"Chest":       "handed in",
		"Ship":        "sunk",
		"Vomit":       "times",
	}

	if valMap[n] != "" {
		return valMap[n]
	}
	return ""
}

func IsBalanceValue(v string) bool {
	v = removeSortString(v)
	switch v {
	case "Gold":
		return true
	case "AncientCoin":
		return true
	case "Doubloon":
		return true
	default:
		return false
	}
}

func BalanceIcon(k string, b int64) string {
	if IsBalanceValue(k) {
		if b > 0 {
			return "📈 "
		}
		if b < 0 {
			return "📉 "
		}
	}
	return ""
}

func removeSortString(s string) string {
	splitString := strings.SplitN(s, "_", 2)
	return splitString[1]
}
