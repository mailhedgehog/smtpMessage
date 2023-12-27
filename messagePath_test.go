package smtpMessage

import (
	"github.com/mailhedgehog/gounit"
	"testing"
)

func TestPath_HasAddress(t *testing.T) {
	messagePath := &MessagePath{
		Relays:  []string{},
		Mailbox: "foo",
		Domain:  "bar.com",
		Params:  "",
	}
	(*gounit.T)(t).AssertEqualsString("foo@bar.com", messagePath.Address())
}

func TestPath_FromString(t *testing.T) {
	from, err := MessagePathFromString("<baz@foo.com>")
	(*gounit.T)(t).AssertNotError(err)
	(*gounit.T)(t).AssertTrue(len(from.Relays) == 0)
	(*gounit.T)(t).AssertEqualsString("baz", from.Mailbox)
	(*gounit.T)(t).AssertEqualsString("foo.com", from.Domain)

	from2, err := MessagePathFromString("<foo-10@quib.com>")
	(*gounit.T)(t).AssertNotError(err)
	(*gounit.T)(t).AssertTrue(len(from2.Relays) == 0)
	(*gounit.T)(t).AssertEqualsString("foo-10", from2.Mailbox)
	(*gounit.T)(t).AssertEqualsString("quib.com", from2.Domain)
	(*gounit.T)(t).AssertEqualsString("foo-10@quib.com", from2.Address())

	to, err := MessagePathFromString("<@foo,@bar,@baz:quix@quib.com> <foo,bar>")
	(*gounit.T)(t).AssertNotError(err)
	(*gounit.T)(t).AssertTrue(len(to.Relays) == 3)
	(*gounit.T)(t).AssertEqualsString("@bar", to.Relays[1])
	(*gounit.T)(t).AssertEqualsString("quix", to.Mailbox)
	(*gounit.T)(t).AssertEqualsString("quib.com", to.Domain)

	_, err = MessagePathFromString("@foo,@bar,@baz:quix@quib.com")
	(*gounit.T)(t).ExpectError(err)
}

func TestPath_ToString(t *testing.T) {
	mailPath := &MessagePath{
		Relays: []string{
			"@foo",
			"@bar",
			"@baz",
		},
		Mailbox: "quix",
		Domain:  "quib.com",
		Params:  "foo,bar",
	}

	(*gounit.T)(t).AssertEqualsString("<@foo,@bar,@baz:quix@quib.com> <foo,bar>", mailPath.ToString())
}
