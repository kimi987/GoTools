package entity

import "time"

func (d *HeroGenXuanyuan) TryResetDaily(resetTime time.Time) bool {

	resetTimeUnix := resetTime.Unix()
	if d.lastResetTime < resetTimeUnix {
		d.lastResetTime = resetTimeUnix
		d.Reset()
		return true
	}
	return false
}
