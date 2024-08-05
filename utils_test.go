package GoUtils

import (
	"fmt"
	"testing"
)

func TestTimeDiff(t *testing.T) {
	fmt.Println(TimeDiff("2020-01-01 10:00:00", "2020-01-01 11:00:00"))
}

func TestTimeStampMs(t *testing.T) {
	fmt.Println(TimeStampMs())
}

func TestTimeStampS(t *testing.T) {
	fmt.Println(TimeStampS())
}

func TestNowFmtStr(t *testing.T) {
	fmt.Println(NowFmtStr())
}

func TestTimeStamp2FmtStr(t *testing.T) {
	fmt.Println(TimeStamp2FmtStr(1713424488))
}

func TestFmtStr2TimeStamp(t *testing.T) {
	fmt.Println(FmtStr2TimeStamp("2020-01-01 10:00:00"))
}

func TestCreateGwsClients(t *testing.T) {
	fmt.Println(Utc2Local("2024-04-29 05:29:25"))
}