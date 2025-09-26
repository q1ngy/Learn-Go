package service

import (
	"context"
	"errors"
	"fmt"
	"math/rand/v2"

	"github.com/q1ngy/Learn-Go/webook/internal/repository"
	"github.com/q1ngy/Learn-Go/webook/internal/service/sms"
)

var ErrCodeSendTooMany = repository.ErrCodeVerifyTooMany

type CodeService interface {
	Send(ctx context.Context, biz, phone string) error
	Verify(ctx context.Context, biz, phone, code string) (bool, error)
}

type CodeServiceImpl struct {
	repo repository.CodeRepository
	sms  sms.Service
}

func NewCodeService(repo repository.CodeRepository, sms sms.Service) CodeService {
	return &CodeServiceImpl{
		repo: repo,
		sms:  sms,
	}
}

func (s *CodeServiceImpl) Send(ctx context.Context, biz, phone string) error {
	code := s.generate()
	err := s.repo.Set(ctx, biz, phone, code)
	if err != nil {
		return err
	}
	const tplId = "x"
	return s.sms.Send(ctx, tplId, []string{code}, phone)
}

func (s *CodeServiceImpl) Verify(ctx context.Context, biz, phone, code string) (bool, error) {
	ok, err := s.repo.Verify(ctx, biz, phone, code)
	if errors.Is(err, ErrCodeSendTooMany) {
		// 屏蔽了底层错误，让攻击者不知道具体问题
		return false, nil
	}
	return ok, err
}

func (s *CodeServiceImpl) generate() string {
	code := rand.IntN(1000000)
	return fmt.Sprintf("%06d", code)
}
