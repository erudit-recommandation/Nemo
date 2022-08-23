package domain

import "fmt"

func build_erudit_url(a *Article) string {
	return fmt.Sprintf("https://id.erudit.org/iderudit/%v", a.ID)
}
