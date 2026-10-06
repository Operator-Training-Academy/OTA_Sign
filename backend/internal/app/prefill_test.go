package app

import "testing"

func TestDocuSealPrefillValuesNameRankOrdering(t *testing.T) {
	values := docusealPrefillValues(User{
		FullName:      "Example, Jane Q.",
		FirstName:     "Jane",
		LastName:      "Example",
		MiddleInitial: "Q",
		Rank:          "SGT",
		PayGrade:      "E-5",
		DoDID:         "1234567890",
		Email:         "jane.example@example.mil",
		UIC:           "WABC12",
	})

	nameFirstFields := []string{"Name & Rank", "Name / Rank", "Name Rank", "Name and Rank"}
	for _, field := range nameFirstFields {
		if values[field] != "Example, Jane Q. SGT" {
			t.Fatalf("expected %q to be name then rank, got %q", field, values[field])
		}
	}

	rankFirstFields := []string{"Rank & Name", "Rank / Name", "Rank Name", "Rank and Name"}
	for _, field := range rankFirstFields {
		if values[field] != "SGT Example, Jane Q." {
			t.Fatalf("expected %q to be rank then name, got %q", field, values[field])
		}
	}

	nameAliasFields := []string{"Operator Name", "Commander Name", "Commander's Name"}
	for _, field := range nameAliasFields {
		if values[field] != "Example, Jane Q." {
			t.Fatalf("expected %q to use formatted name, got %q", field, values[field])
		}
	}
}

func TestUserFromKeycloakClaimsRequiresOTAIdentity(t *testing.T) {
	cfg := Config{
		KeycloakClientID:       "ota-sign",
		KeycloakDoDIDClaim:     "dod_id",
		KeycloakUICClaim:       "uic",
		KeycloakRankClaim:      "rank",
		KeycloakArmyEmailClaim: "army_email",
	}
	claims := map[string]any{
		"name":        "Jane Example",
		"given_name":  "Jane",
		"family_name": "Example",
		"email":       "jane.example@example.mil",
		"dod_id":      "1234567890",
		"uic":         "WABC12",
		"rank":        "SGT",
		"army_email":  "jane.example@army.mil",
		"resource_access": map[string]any{
			"ota-sign": map[string]any{
				"roles": []any{"viewown", "viewunit", "unexpected-role", "viewunit"},
			},
		},
	}

	user, err := userFromKeycloakClaims(cfg, "https://keycloak.example.mil/realms/ota", "subject-1", claims)
	if err != nil {
		t.Fatalf("expected valid Keycloak claims, got %v", err)
	}
	if user.DoDID != "1234567890" || user.UIC != "WABC12" || user.ArmyEmail != "jane.example@army.mil" || !hasAny(user.Capabilities, "viewown", "viewunit") || len(user.Capabilities) != 2 {
		t.Fatalf("unexpected mapped Keycloak user: %#v", user)
	}

	delete(claims, "dod_id")
	if _, err := userFromKeycloakClaims(cfg, "https://keycloak.example.mil/realms/ota", "subject-1", claims); err == nil {
		t.Fatal("expected missing DoD ID claim to be rejected")
	}
}

func TestUserFromKeycloakClaimsRequiresViewOwnRole(t *testing.T) {
	cfg := Config{KeycloakClientID: "ota-sign", KeycloakDoDIDClaim: "dod_id", KeycloakUICClaim: "uic"}
	claims := map[string]any{
		"name": "Jane Example", "email": "jane.example@example.mil", "dod_id": "1234567890", "uic": "WABC12",
		"resource_access": map[string]any{
			"ota-sign": map[string]any{"roles": []any{"viewunit"}},
		},
	}

	if _, err := userFromKeycloakClaims(cfg, "https://keycloak.example.mil/realms/ota", "subject-1", claims); err == nil {
		t.Fatal("expected Keycloak user without viewown to be rejected")
	}
}
