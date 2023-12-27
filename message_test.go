package smtpMessage

import (
	"github.com/mailhedgehog/gounit"
	"io"
	"testing"
)

var emailString = `HELO:[127.0.0.1]
FROM:<sender@example.test>
TO:<@foo,@bar.com:your-email@here.test> <parma1,param2>
TO:<second@here.test>

From: Sender <sender@example.test>
To: ReceiverName <your-email@here.test>
X-Priority: 1 (Highest)
Subject: Mail test 2023-03-25 22:16:52
Message-ID: <3ac479f00c5d9ea7519ade0784ed1060@example.test>
MIME-Version: 1.0
Date: Sat, 25 Mar 2023 22:16:52 +0200
Content-Type: multipart/alternative; boundary=UhEiB9Sb

--UhEiB9Sb
Content-Type: text/plain; charset=utf-8
Content-Transfer-Encoding: quoted-printable

This is an important message!
--UhEiB9Sb
Content-Type: text/html; charset=utf-8
Content-Transfer-Encoding: quoted-printable

<strong>This is an test message!</strong>
--UhEiB9Sb--`

func TestNewMessageIDIsUuid(t *testing.T) {
	id := NewMessageID()
	(*gounit.T)(t).AssertLengthString(36, string(id))
}

func TestFromString(t *testing.T) {
	smtpMessage := FromString(emailString, NewMessageID())
	(*gounit.T)(t).AssertEqualsString("[127.0.0.1]", smtpMessage.Helo)
	(*gounit.T)(t).AssertEqualsString("sender@example.test", smtpMessage.From.Address())
	(*gounit.T)(t).AssertEqualsString("your-email@here.test", smtpMessage.To[0].Address())
	(*gounit.T)(t).AssertEqualsString("second@here.test", smtpMessage.To[1].Address())
	(*gounit.T)(t).AssertLengthString(586, smtpMessage.GetOrigin())
	(*gounit.T)(t).AssertEqualsString("Mail test 2023-03-25 22:16:52", smtpMessage.GetEmail().Subject)
	(*gounit.T)(t).AssertEqualsString("Sender", smtpMessage.GetEmail().From[0].Name)
	(*gounit.T)(t).AssertEqualsString("sender@example.test", smtpMessage.GetEmail().From[0].Address)
}

func TestToReader(t *testing.T) {
	smtpMessage := FromString(emailString, NewMessageID())

	bytes, err := io.ReadAll(smtpMessage.ToReader())
	(*gounit.T)(t).AssertNotError(err)

	(*gounit.T)(t).AssertLengthString(755, string(bytes))
}

func TestSetOrigin(t *testing.T) {
	msg := &SmtpMessage{}
	msg.SetOrigin("")

	(*gounit.T)(t).AssertNil(msg.From)
	(*gounit.T)(t).AssertEqualsInt(0, len(msg.To))

	msg.SetOrigin(`From: Rares <quix@quib.com>
Date: Thu, 2 May 2019 11:25:35 +0300
Subject: Re: kern/54143 (virtualbox)
To: foo@quib.com
Content-Type: multipart/mixed; boundary="0000000000007e2bb40587e36196"

--0000000000007e2bb40587e36196
Content-Type: text/html; charset="UTF-8"

<div dir="ltr"><div>html text part</div><div><br></div><div><br><br></div></div>

--0000000000007e2bb40587e36196--
`)

	(*gounit.T)(t).AssertEqualsString("quix@quib.com", msg.From.Address())
	(*gounit.T)(t).AssertEqualsInt(1, len(msg.To))
	(*gounit.T)(t).AssertEqualsString("foo@quib.com", msg.To[0].Address())
}
