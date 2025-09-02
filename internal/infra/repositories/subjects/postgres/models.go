package postgres

import (
	"fmt"
	"time"

	"github.com/lib/pq"

	jsoniter "github.com/json-iterator/go"

	"github.com/DanilaKorobkov/defi-monitoring/internal/domain"
)

type subjectModel struct {
	TelegramUserID int64  `db:"telegram_user_id"`
	Payload        string `db:"payload"`
}

type subjectPayloadModel struct {
	TelegramUserID int64          `db:"telegram_user_id"`
	Wallets        pq.StringArray `db:"wallets"`
	CheckInterval  time.Duration  `db:"check_interval"`
}

func newSubjectModel(subject domain.Subject) (subjectModel, error) {
	payloadModel := subjectPayloadModel{
		TelegramUserID: subject.TelegramUserID,
		Wallets:        subject.Wallets,
		CheckInterval:  subject.CheckInterval,
	}

	dump, err := jsoniter.MarshalToString(payloadModel)
	if err != nil {
		return subjectModel{}, fmt.Errorf("jsoniter.MarshalToString: %w", err)
	}

	model := subjectModel{
		TelegramUserID: subject.TelegramUserID,
		Payload:        dump,
	}

	return model, nil
}

func (model subjectModel) mustToSubject() domain.Subject {
	payloadModel := subjectPayloadModel{}

	err := jsoniter.UnmarshalFromString(model.Payload, &payloadModel)
	if err != nil {
		message := "mustToSubject: jsoniter.UnmarshalFromString: " + err.Error()
		panic(message)
	}

	return domain.Subject{
		TelegramUserID: payloadModel.TelegramUserID,
		Wallets:        payloadModel.Wallets,
		CheckInterval:  payloadModel.CheckInterval,
	}
}
