package auth

import "testing"

func TestPassword(t *testing.T) {
	t.Parallel()

	// table testing
	cases := []struct {
		Name            string
		Password        string
		ComparePassword string
		Equal           bool
	}{
		{"Success", "superman", "superman", true},
		{"Compare with empty password", "superman", "", false},
		{"Compare with invalid password", "superman", "batman", false},
	}

	for _, tC := range cases {
		t.Run(tC.Name, func(t *testing.T) {
			hashed := HashPassword(tC.Password)
			if hashed == "" {
				t.Error("expected non-empty string from hashPassword; got empty")
			}
			if compareHashAndPassword(hashed, tC.ComparePassword) != tC.Equal {
				if tC.Equal {
					t.Error("expected hashed and password are equal; got not equal")
				} else {
					t.Error("expected hashed and password are equal; got equal")
				}
			}
		})
	}
}
