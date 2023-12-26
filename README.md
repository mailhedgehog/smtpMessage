# MailHedgehog smtp message structure adaptor

Converts smtp message form string value to object and back.

## Usage

```go
msg := smtpMessage.FromString(emailString, "")

msg.GetEmail().Subject
smtpMessage.GetEmail().From[0].Name
smtpMessage.GetEmail().From[0].Address
```

## Development

```shell
go mod tidy
go mod verify
go mod vendor
go test --cover
```

## Credits

- [![Think Studio](https://yaroslawww.github.io/images/sponsors/packages/logo-think-studio.png)](https://think.studio/)
