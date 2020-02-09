package remote

import (
	prompt "github.com/c-bata/go-prompt"
)

type Host string

func (h Host) Suggest() prompt.Suggest {
	return prompt.Suggest{
		Text: string(h),
	}
}
