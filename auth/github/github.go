package github

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"

	"github.com/hibooboo2/utils/http"
)

func getKeysProfile(user string) []string {
	resp, err := http.DefaultClient.Get(fmt.Sprintf("https://github.com/%s.keys", user))
	if err != nil {
		return nil
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil
	}
	keys := strings.Split(string(data), "\n")
	keysKeep := []string{}
	for i := range keys {
		if keys[i] != "" {
			keysKeep = append(keysKeep, keys[i])
		}
	}
	return keysKeep
}

func getKeysApi(user string) []string {
	var sshKeys []struct {
		ID  int64  `json:"id"`
		Key string `json:"key"`
	}

	err := http.DefaultClient.GetAsObj(fmt.Sprintf("https://api.github.com/users/%s/keys", user), &sshKeys)
	if err != nil {
		return nil
	}
	keys := []string{}
	for _, k := range sshKeys {
		keys = append(keys, k.Key)
	}
	return keys
}

func GetUser() (string, error) {
	cmd := exec.Command("ssh", "git@github.com")
	data, _ := cmd.CombinedOutput()

	val := string(data)
	val = strings.ToLower(val)
	if !strings.Contains(val, "successfully authenticated") {
		return "", fmt.Errorf("failed to authenticate with github")
	}
	val = strings.Split(val, "hi ")[1]
	val = strings.Split(val, "!")[0]
	val = strings.TrimSpace(val)
	return val, nil
}
