package prompt

import (
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/nicarl/somafm/radioChannels"
)

func SelectChannel(radioCh []radioChannels.RadioCh) (radioChannels.RadioCh, error) {
	templates := &promptui.SelectTemplates{
		Label:    "Select channel",
		Active:   "{{ .Title | cyan }}",
		Inactive: "{{ .Title }}",
		Selected: "{{ .Title | cyan }}",
		Details: `-----------------------
{{ .Description }}
DJ: {{ .Dj }}
Genre: {{ .Genre }}`,
	}

	searcher := func(input string, index int) bool {
		selectedCh := radioCh[index]
		name := strings.Replace(strings.ToLower(selectedCh.Title), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)

		return strings.Contains(name, input)
	}

	prompt := promptui.Select{
		Label:     "Select channel",
		Items:     radioCh,
		Templates: templates,
		Size:      8,
		Searcher:  searcher,
	}

	index, _, err := prompt.Run()
	return radioCh[index], err
}
