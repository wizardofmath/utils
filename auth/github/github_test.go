package github

import "testing"

func TestKeysProfile(t *testing.T) {
	user := "hibooboo2"
	keys := getKeysProfile(user)
	if len(keys) == 0 {
		t.Fatal("Failed to get keys for", user)
	}
}

func TestKeysAPI(t *testing.T) {
	user := "hibooboo2"
	keys := getKeysApi(user)
	if len(keys) == 0 {
		t.Fatal("Failed to get keys for", user)
	}
}

func TestKeysBothSame(t *testing.T) {
	user := "hibooboo2"
	keys := getKeysApi(user)
	if len(keys) == 0 {
		t.Fatal("Failed to get keys for", user)
	}
	keys2 := getKeysProfile(user)
	if len(keys) == 0 {
		t.Fatal("Failed to get keys for", user)
	}

	if len(keys) != len(keys2) {
		t.Fatalf("keys api: %d keys profile: %d\n", len(keys), len(keys2))
	}
	for i := range keys {
		if keys[i] != keys2[i] {
			t.Error("key pub does not match key api")
		}
	}
}

func TestCanSSHGithub(t *testing.T) {
	_, err := GetUser()
	if err != nil {
		t.Fatal(err)
	}
}
