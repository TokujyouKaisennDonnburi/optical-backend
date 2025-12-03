package command

import (
	"github.com/TokujouKaisenDonburi/optical-backend/internal/calendar"
	"github.com/google/uuid"
	"golang.org/x/net/context"
)

type MemberCreateInput struct {
	UserId 		uuid.UUID
	CalendarId  uuid.UUID
	Email 		string 
}

func (c *CalendarCommand) CreateMember (ctx context.Context, input MemberCreateInput) (*calendar.Member, error) {
	// TODO: c.memberRepository.Create の引数を修正する
	//   Create(ctx, input.CalendarId, input.Email) の3つの引数が必要

	// TODO: Createの戻り値を受け取って使う（今は _ で捨てている）
	//   member, err := c.memberRepository.Create(...) として受け取る

	// TODO: 受け取ったmemberを返す（今のnewMemberは使われていない）

	var newMember *calendar.Member
	_, err := c.memberRepository.Create(input.Email)
	if err != nil {
		return nil, err
	}
	return newMember, nil
}
