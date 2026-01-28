package notice

import (
	"context"
	"fmt"
	"strings"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/calendar/repository"
	"github.com/TokujouKaisenDonburi/optical-backend/internal/notice/creator"
	"github.com/TokujouKaisenDonburi/optical-backend/internal/option"
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

// カレンダーカラー変更時にメンバーへ通知を送信
func (s *CalendarNoticeService) NotifyCalendarColorUpdated(
	ctx context.Context,
	calendarId uuid.UUID,
	calendarName string,
	oldColor string,
	newColor string,
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
			"カラーが変更されました",
			fmt.Sprintf("%sさんにより「%s」から「%s」に変更されました", user.Name, oldColor, newColor),
			creator.WithCalendarID(calendarId),
		)
	}

	return nil
}

// カレンダーオプション変更時にメンバーへ通知を送信
func (s *CalendarNoticeService) NotifyCalendarOptionsUpdated(
	ctx context.Context,
	calendarId uuid.UUID,
	calendarName string,
	oldOptions []option.Option,
	newOptions []option.Option,
	updatedByUserId uuid.UUID,
) error {
	// 追加・削除されたオプションを検出
	oldOptionIds := make(map[int32]string)
	for _, opt := range oldOptions {
		oldOptionIds[opt.Id] = opt.Name
	}

	newOptionIds := make(map[int32]string)
	for _, opt := range newOptions {
		newOptionIds[opt.Id] = opt.Name
	}

	var addedNames []string
	var deletedNames []string

	// 削除されたオプション（古いにあって新しいにない）
	for id, name := range oldOptionIds {
		if _, exists := newOptionIds[id]; !exists {
			deletedNames = append(deletedNames, name)
		}
	}

	// 追加されたオプション（新しいにあって古いにない）
	for id, name := range newOptionIds {
		if _, exists := oldOptionIds[id]; !exists {
			addedNames = append(addedNames, name)
		}
	}

	// 変更がなければ何もしない
	if len(addedNames) == 0 && len(deletedNames) == 0 {
		return nil
	}

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

	// タイトルとコンテンツを決定
	var title string
	var content string

	if len(addedNames) > 0 && len(deletedNames) > 0 {
		// パターン3: 追加&削除
		title = "オプションが変更されました"
		content = fmt.Sprintf("%sさんにより「%s」のオプションが変更されました\n追加されたオプション: %s\n削除されたオプション: %s",
			user.Name, calendarName, strings.Join(addedNames, ", "), strings.Join(deletedNames, ", "))
	} else if len(addedNames) > 0 {
		// パターン1: 追加のみ
		title = "オプションが追加されました"
		content = fmt.Sprintf("%sさんにより「%s」にオプションが追加されました\n追加されたオプション: %s",
			user.Name, calendarName, strings.Join(addedNames, ", "))
	} else {
		// パターン2: 削除のみ
		title = "オプションが削除されました"
		content = fmt.Sprintf("%sさんにより「%s」からオプションが削除されました\n削除されたオプション: %s",
			user.Name, calendarName, strings.Join(deletedNames, ", "))
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
			title,
			content,
			creator.WithCalendarID(calendarId),
		)
	}

	return nil
}
