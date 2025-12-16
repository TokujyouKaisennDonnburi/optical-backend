package user

import "net/url"

type Avatar struct {
	Url   string
	Valid bool
}

func NewAvatar(url string) *Avatar {
	img := &Avatar{}
	img.SetUrl(url)
	return img
}

// URLを設定する
func (img *Avatar) SetUrl(imageUrl string) {
	_, err := url.Parse(imageUrl)
	if imageUrl == "" || err != nil {
		img.Url = ""
		img.Valid = false
	} else {
		img.Url = imageUrl
		img.Valid = true
	}
}

