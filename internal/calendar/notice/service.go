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

// カレンダー招待時に登録済みユーザーへ通知を送信
func (s *CalendarNoticeService) NotifyCalendarInvited(
	ctx context.Context,
	calendarId uuid.UUID,
	calendarName string,
	invitedByUserId uuid.UUID,
	emails []string,
) error {
	// 招待実行者の情報を取得
	inviter, err := s.userRepository.FindById(ctx, invitedByUserId)
	if err != nil {
		return err
	}

	// メールアドレスから登録済みユーザーを取得
	users, err := s.userRepository.FindByEmails(ctx, emails)
	if err != nil {
		return err
	}

	// 登録済みユーザーに通知を送信
	for _, user := range users {
		_ = s.noticeCreator.CreateNotice(
			ctx,
			user.Id,
			"カレンダーに招待されました",
			fmt.Sprintf("%sさんから「%s」に招待されました。\n招待メールを確認してください。", inviter.Name, calendarName),
			creator.WithCalendarID(calendarId),
		)
	}

	return nil
}

// カレンダー名変更時にメンバーへ通知を送信
func (s *CalendarNoticeService) NotifyCalendarNameUpdated(
	ctx context.Context,
	calendarId uuid.UUID,
	oldName string,
	newName string,
	updatedByUserId uuid.UUID,
) error {
	// 更新実行者の情報を取得
	user, err := s.userRepository.FindById(ctx, updatedByUserId)
	if err != nil {
		return err
	}

	// メンバー一覧を取得
	members, err := s.memberRepository.FindMembers(ctx, calendarId)
	if err != nil {
		return err
	}

	// 参加済みメンバーに通知を送信（更新実行者を除く）
	for _, member := range members {
		if member.UserId == updatedByUserId {
			continue
		}
		if member.JoinedAt.IsZero() {
			continue
		}
		_ = s.noticeCreator.CreateNotice(
			ctx,
			member.UserId,
			"カレンダー名が変更されました",
			fmt.Sprintf("%sさんにより「%s」から「%s」に変更されました", user.Name, oldName, newName),
			creator.WithCalendarID(calendarId),
		)
	}

	return nil
}
