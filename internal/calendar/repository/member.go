package repository

type MemberRepository interface {
	create(mail string)(error)
}
