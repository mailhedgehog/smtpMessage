package smtpMessage

import (
    "bytes"
    "fmt"
    "github.com/google/uuid"
    "github.com/mailhedgehog/email"
    "github.com/mailhedgehog/logger"
    "io"
    "strings"
)

var configuredLogger *logger.Logger

func logManager() *logger.Logger {
    if configuredLogger == nil {
        configuredLogger = logger.CreateLogger("smtpMessage")
    }
    return configuredLogger
}

// SmtpMessage represents parsed SMTP message what allows easily get and manipulate data
type SmtpMessage struct {
    ID     MessageID
    Helo   string
    From   *MessagePath
    To     []*MessagePath
    email  *email.Email
    origin string
}

// MessageID represents the ID of an SMTP message
type MessageID string

// NewMessageID generates a new mail identificatior
func NewMessageID() MessageID {
    return MessageID(uuid.New().String())
}

// SetOrigin add to object new Origin Data and automatically parse it to email structure
func (message *SmtpMessage) SetOrigin(origin string) error {
    var err error

    message.origin = origin
    message.email, err = email.Parse(strings.NewReader(message.origin))

    if message.ID == "" {
        message.ID = NewMessageID()
    }
    if message.From == nil {
        message.From, _ = MessagePathFromString(fmt.Sprintf("<%s>", message.email.From[0].Address))
    }
    if len(message.To) == 0 {
        for _, to := range message.email.To {
            if to != nil {
                toConverted, _ := MessagePathFromString(fmt.Sprintf("<%s>", to.Address))
                if toConverted != nil {
                    message.To = append(message.To, toConverted)
                }
            }
        }
    }

    return err
}

// GetOrigin data value (string)
func (message *SmtpMessage) GetOrigin() string {
    return message.origin
}

// GetEmail data object (parsed origin value)
func (message *SmtpMessage) GetEmail() *email.Email {
    if message.email == nil {
        message.SetOrigin(message.origin)
    }
    return message.email
}

// ToReader returns an io.Reader containing the raw message data
func (message *SmtpMessage) ToReader() io.Reader {
    var bufferReader = new(bytes.Buffer)

    if message != nil {
        bufferReader.WriteString("ID:" + string(message.ID) + "\r\n")
        bufferReader.WriteString("HELO:" + message.Helo + "\r\n")
        if message.From != nil {
            bufferReader.WriteString("FROM:" + message.From.ToString() + "\r\n")
        }
        for _, to := range message.To {
            if to != nil {
                bufferReader.WriteString("TO:" + to.ToString() + "\r\n")
            }
        }
        bufferReader.WriteString("\r\n")
        bufferReader.WriteString(message.origin)
    }

    return bufferReader
}

// FromString returns a SmtpMessage from raw message bytes (as output by SmtpMessage.ToReader())
func FromString(messageString string, messageId MessageID) *SmtpMessage {
    var messagePath *MessagePath
    msg := &SmtpMessage{
        ID: messageId,
    }
    var headerDone bool
    var origin string
    for _, l := range strings.Split(messageString, "\n") {
        if !headerDone {

            if strings.HasPrefix(l, "ID:") {
                if msg.ID == "" {
                    l = strings.TrimPrefix(l, "ID:")
                    l = strings.Trim(l, " \n\r")
                    msg.ID = MessageID(l)
                }
                continue
            }

            if strings.HasPrefix(l, "HELO:") {
                l = strings.TrimPrefix(l, "HELO:")
                l = strings.Trim(l, " \n\r")
                msg.Helo = l
                continue
            }

            if strings.HasPrefix(l, "FROM:") {
                l = strings.TrimPrefix(l, "FROM:")
                l = strings.Trim(l, " \n\r")
                messagePath, _ = MessagePathFromString(l)
                if messagePath != nil {
                    msg.From = messagePath
                }
                continue
            }

            if strings.HasPrefix(l, "TO:") {
                l = strings.TrimPrefix(l, "TO:")
                l = strings.Trim(l, " \n\r")
                messagePath, _ := MessagePathFromString(l)
                if messagePath != nil {
                    msg.To = append(msg.To, messagePath)
                }
                continue
            }

            if strings.TrimSpace(l) == "" {
                headerDone = true
                continue
            }
        }

        origin += l + "\n"
    }

    err := msg.SetOrigin(origin)
    if err != nil {
        logManager().Error(err.Error())
    }

    return msg
}
