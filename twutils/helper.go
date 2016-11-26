package twutils

import "github.com/ChimeraCoder/anaconda"

func ProfitIdStr(t *anaconda.Tweet) string {
	if t.RetweetedStatus != nil {
		return ProfitIdStr(t.RetweetedStatus)
	}
	return t.IdStr
}

func ProfitScreenName(t *anaconda.Tweet) string {
	if t.RetweetedStatus != nil {
		return ProfitScreenName(t.RetweetedStatus)
	}
	return t.User.ScreenName
}
