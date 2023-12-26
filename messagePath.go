package smtpMessage

import (
	"fmt"
	"golang.org/x/exp/slices"
	"regexp"
	"strings"
)

// MessagePath represents an SMTP "MAIL FROM" or "RCPT TO" string content converted to object.
type MessagePath struct {
	Relays  []string
	Mailbox string
	Domain  string
	Params  string
}

// Address returns full address form path eg: "mailbox@domain".
func (path *MessagePath) Address() string {
	return path.Mailbox + "@" + path.Domain
}

// ToString convert path object to standard rfc5321 path string.
func (path *MessagePath) ToString() string {
	line := "<"
	if len(path.Relays) > 0 {
		line += strings.Join(path.Relays, ",")
		line += ":"
	}
	line += path.Address()
	line += ">"
	if len(path.Params) > 0 {
		line += fmt.Sprintf(" <%s>", path.Params)
	}

	return line
}

// MessagePathFromString parses a forward-path or reverse-path into its parts
// According rfc5321, message path have format "<[RELAYS:]BOX@DOMAIN>[ <EMAIL_PARAMS>]"
// so all of these examples are possible
// "<quix@quib.com> <foo,bar>"
// "<@foo,@bar,@baz:quix@quib.com>"
// "<@foo,@bar,@baz:quix@quib.com> <foo,bar>"
func MessagePathFromString(path string) (*MessagePath, error) {
	mailPath := &MessagePath{}

	path = strings.Trim(path, " ")

	logManager().Debug(fmt.Sprintf("Converting path: \"%s\"", path))

	parts := slices.DeleteFunc(strings.Split(path, " "), func(e string) bool {
		return !strings.HasPrefix(e, "<") || !strings.HasSuffix(e, ">")
	})

	if len(parts) < 1 || len(parts) > 2 {
		return nil, fmt.Errorf("incorrect path line \"%s\"", path)
	}

	matchEmail := regexp.MustCompile(`(?m)^\s*<([^>]+)>\s*?$`).FindStringSubmatch(parts[0])
	if len(matchEmail) >= 2 {
		var relays []string
		userEmail := matchEmail[1]
		if strings.Contains(matchEmail[1], ":") {
			x := strings.SplitN(matchEmail[1], ":", 2)
			r, e := x[0], x[1]
			userEmail = e
			relays = strings.Split(r, ",")
		}
		mailbox, domain := "", ""
		if strings.Contains(userEmail, "@") {
			x := strings.SplitN(userEmail, "@", 2)
			mailbox, domain = x[0], x[1]
		} else {
			mailbox = userEmail
		}

		mailPath.Relays = relays
		mailPath.Mailbox = mailbox
		mailPath.Domain = domain
	}

	if len(parts) == 2 {
		matchParams := regexp.MustCompile(`(?m)^\s*<([^>]+)>\s*?$`).FindStringSubmatch(parts[1])
		if len(matchEmail) == 2 {
			mailPath.Params = matchParams[1]
		}
	}

	if len(mailPath.Domain) > 0 {
		return mailPath, nil
	}

	return nil, fmt.Errorf("incorrect path line \"%s\"", path)
}
