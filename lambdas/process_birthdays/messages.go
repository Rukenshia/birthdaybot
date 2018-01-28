package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"strings"
	"text/template"

	log "github.com/sirupsen/logrus"
)

var templatesSingleToday = []string{
	"*I'm so excited*, today is @{{Username}}'s birthday! Hope you're having an awesome day!!",
	"Sup @{{Username}},\nI've heard it's your birthday today! That's *completely awesome* and I hope you have a great day :tada:",
	"Hey @{{Username}},\n\nHAPPY FREAKING BIRTHDAY! Don't be afraid, you still look as good as you always did :blush::tada:",
	"*Oooooh it's birthday time again*! Happy Birthday, @{{Username}} :tada:",
	"Hey you! Yes, you @{{Username}}. Just wanted to let you know that *everyone* wishes you a happy birthday today!",
	"Drumrolls, _everyone_ :drum: :drum:\n\n\nHAPPY BIRTHDAY to @{{Username}} :tada:",
}

var templatesMultipleToday = []string{
	"Wow, this is more special than usual! Today we have not only @{{.FirstUsername}} celebrating their birthday, but also {{enumerate .OtherUsernames}}!\nHappy birthday to you all :blush: :tada:",
	"Drumrolls, _everyone_ :drum: :drum:\n\nHAPPY BIRTHDAY to {{enumerate .OtherUsernames}}!\nDon't worry, @{{.FirstUsername}}, I didn't forget about you :blush: Happy birthday to you as well :tada:",
}

// GetRandomTodayBirthdayMessage returns a randomised message for the birthdays today mentioning all users
func GetRandomTodayBirthdayMessage(usernames []string) string {
	return GetRandomMultiTodayBirthdayMessage(usernames)
}

// GetRandomSingleTodayBirthdayMessage returns a randomised message for a single birthday personalised for the user
func GetRandomSingleTodayBirthdayMessage(username string) string {
	return strings.Replace(templatesSingleToday[rand.Intn(len(templatesSingleToday))], "{{Username}}", username, 0)
}

// GetRandomMultiTodayBirthdayMessage returns a randomised message for multiple birthdays personalised for all users
func GetRandomMultiTodayBirthdayMessage(usernames []string) string {
	if len(usernames) == 0 {
		return "Somethings wrong with me, help me @jan"
	}

	if len(usernames) == 1 {
		return GetRandomSingleTodayBirthdayMessage(usernames[0])
	}

	tpl, err := template.New("multi").Funcs(template.FuncMap{
		"enumerate": func(names []string) string {
			if len(names) == 1 {
				return fmt.Sprintf("@%s", names[0])
			}

			tempStr := fmt.Sprintf("@%s", names[0])
			for idx, name := range names[1:] {
				if idx == len(names)-2 {
					tempStr = fmt.Sprintf("%s and @%s", tempStr, name)
					continue
				}

				tempStr = fmt.Sprintf("%s, @%s", tempStr, name)
			}

			return tempStr
		},
	}).Parse(templatesMultipleToday[rand.Intn(len(templatesMultipleToday))])

	if err != nil {
		log.Errorf("Could not parse template: %v", err)
		return "Something is really wrong with me, help @jan :worried:"
	}

	var buf bytes.Buffer

	if err := tpl.Execute(&buf, map[string]interface{}{
		"Usernames":      usernames,
		"FirstUsername":  usernames[0],
		"OtherUsernames": usernames[1:],
	}); err != nil {
		return "@jan, you messed something up. Can you please take a look at me? :worried:"
	}

	return buf.String()

}
