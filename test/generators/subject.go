package generators

import (
	"time"

	"github.com/google/uuid"

	"github.com/DanilaKorobkov/defi-monitoring/internal/domain"
)

const (
	minCheckIntervalMinutes = 1
	maxCheckIntervalMinutes = 60

	minTelegramUserID = 100
	maxTelegramUserID = 100_000
)

type SubjectGenerator struct {
	buffer domain.Subject
}

func NewSubjectGenerator() *SubjectGenerator {
	return &SubjectGenerator{}
}

func (gen *SubjectGenerator) Base(subject domain.Subject) *SubjectGenerator {
	gen.buffer = subject
	return gen
}

func (gen *SubjectGenerator) Result() domain.Subject {
	return gen.buffer
}

func (gen *SubjectGenerator) Slim() *SubjectGenerator {
	return gen.
		WithTelegramUserID().
		WithWallets().
		WithCheckInterval()
}

func (gen *SubjectGenerator) WithCheckInterval(interval ...time.Duration) *SubjectGenerator {
	set(&gen.buffer.CheckInterval, generateCheckInterval, interval...)
	return gen
}

func (gen *SubjectGenerator) WithTelegramUserID(id ...int64) *SubjectGenerator {
	set(&gen.buffer.TelegramUserID, generateTelegramUserID, id...)
	return gen
}

func (gen *SubjectGenerator) WithWallets(wallets ...[]string) *SubjectGenerator {
	set(&gen.buffer.Wallets, generateWallets, wallets...)
	return gen
}

func generateTelegramUserID() int64 {
	return int64(RandomInt(minTelegramUserID, maxTelegramUserID))
}

func generateWallets() []string {
	return GeneratePlenty(RandomCount(), generateWallet)
}

func generateWallet() string {
	return uuid.NewString()
}

func generateCheckInterval() time.Duration {
	minutes := RandomInt(minCheckIntervalMinutes, maxCheckIntervalMinutes)
	return time.Duration(minutes) * time.Minute
}
