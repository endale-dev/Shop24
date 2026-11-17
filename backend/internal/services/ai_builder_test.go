package services

import (
	"strings"
	"testing"

	"shop24/internal/domain/models"
)

func TestChatRemembersEarlierCustomization(t *testing.T) {
	prior := priorTurns(&models.DesignNote{Prompt: "Use outline buttons", Summary: "Applied outline buttons."})
	if len(prior) != 2 || prior[0].Role != "user" || prior[1].Role != "designer" {
		t.Fatalf("prior = %+v", prior)
	}
	turns := appendTurns(prior, "Show the products as a list", "Products are in a list.")
	if len(turns) != 4 {
		t.Fatalf("turns = %d", len(turns))
	}
	history := historyFor(turns)
	if len(history) != 4 || history[2].Text != "Show the products as a list" {
		t.Fatalf("history = %+v", history)
	}
	long := appendTurns(nil, strings.Repeat("a", 500), "ok")
	if len([]rune(long[0].Text)) != 400 {
		t.Fatalf("clip = %d", len([]rune(long[0].Text)))
	}
}
