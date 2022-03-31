package global

type ILoop interface {
	Loop()
}

type ITimeLoop interface {
	TimeLoop(nowTime int64)
}

type Looper struct {
}

func NewLooper() *Looper {
	return &Looper{}
}

func (l *Looper) Loop() {

}

func (l *Looper) TimeLoop(nowTime int64) {

}
