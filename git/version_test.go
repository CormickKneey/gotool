package git

import (
	"fmt"
	"testing"
)

func TestRepo(t *testing.T) {
	repo := Repo()
	fmt.Printf("Repo: %v",repo)
}