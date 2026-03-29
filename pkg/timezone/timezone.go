package timezone

import "time"

var Dushanbe *time.Location

func Init() error {
	var err error
	Dushanbe, err = time.LoadLocation("Asia/Dushanbe")
	return err
}
