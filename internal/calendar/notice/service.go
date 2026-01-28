package notice

import (
	"context"
	"fmt"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/calendar/repository"
	"github.com/TokujouKaisenDonburi/optical-backend/internal/notice/creator"
	userRepo "github.com/TokujouKaisenDonburi/optical-backend/internal/user/repository"
	"github.com/google/uuid"
)

type CalendarNoticeService struct {
	memberRepository repository.MemberRepository
	userRepository   userRepo.UserRepository
	noticeCreator    *creator.NoticeCreator
}

func NewCalendarNoticeService(
	memberRepository repository.MemberRepository,
	userRepository userRepo.UserRepository,
	noticeCreator *creator.NoticeCreator,
) *CalendarNoticeService {
	if memberRepository == nil {
		panic("memberRepository is nil")
	}
	if userRepository == nil {
		panic("userRepository is nil")
	}
	if noticeCreator == nil {
		panic("noticeCreator is nil")
	}
	return &CalendarNoticeService{
		memberRepository: memberRepository,
		userRepository:   userRepository,
		noticeCreator:    noticeCreator,
	}
}

// カレンダー削除時にメンバーへ通知を送信
func (s *CalendarNoticeService) NotifyCalendarDeleted(
	ctx context.Context,
	calendarId uuid.UUID,
	calendarName string,
	deletedByUserId uuid.UUID,
) error {
	// 削除実行者の情報を取得
	user, err := s.userRepository.FindById(ctx, deletedByUserId)
	if err != nil {
		return err
	}

	// メンバー一覧を取得
	members, err := s.memberRepository.FindMembers(ctx, calendarId)
	if err != nil {
		return err
	}

	// 参加済みメンバーに通知を送信（削除実行者を除く）
	for _, member := range members {
		if member.UserId == deletedByUserId {
			continue
		}
		// 参加済みメンバーのみに通知（JoinedAtがnullでない）
		if member.JoinedAt.IsZero() {
			continue
		}
		_ = s.noticeCreator.CreateNotice(
			ctx,
			member.UserId,
			"カレンダーが削除されました",
			fmt.Sprintf("%sさんにより「%s」が削除されました", user.Name, calendarName),
			creator.WithCalendarID(calendarId),
		)
	}

	return nil
}
