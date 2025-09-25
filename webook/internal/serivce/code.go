package serivce

import (
	"context"
	"errors"
	"fmt"
	"math/rand/v2"

	"github.com/q1ngy/Learn-Go/webook/internal/repository"
	"github.com/q1ngy/Learn-Go/webook/internal/serivce/sms"
)

var ErrCodeSendTooMany = repository.ErrCodeVerifyTooMany

type CodeService struct {
	repo repository.CodeRepository
	sms  sms.Service
}

func (s *CodeService) NewCodeService(repo repository.CodeRepository, sms sms.Service) CodeService {
	return CodeService{
		repo: repo,
		sms:  sms,
	}
}

func (s *CodeService) Send(ctx context.Context, biz, phone string) error {
	code := s.generate()
	err := s.repo.Set(ctx, biz, phone, code)
	if err != nil {
		return err
	}
	const tplId = "x"
	return s.sms.Send(ctx, tplId, []string{code}, phone)
}

func (s *CodeService) Verify(ctx context.Context, biz, phone, code string) (bool, error) {
	ok, err := s.repo.Verify(ctx, biz, phone, code)
	if errors.Is(err, ErrCodeSendTooMany) {
		// 屏蔽了底层错误，让攻击者不知道具体问题
		return false, nil
	}
	return ok, err
}

func (s *CodeService) generate() string {
	code := rand.IntN(1000000)
	return fmt.Sprintf("%06d", code)
}
