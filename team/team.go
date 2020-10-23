package team

const (
	TEAM_CALC_BEGIN = 2 //团队收益开始结算时间
	TEAM_CALC_END   = 3 //团队收益结束结算时间
)

type TeamHistory struct {
	DiZ    int32 `bson:"diz"`    //弟子人数
	TZ     int32 `bson:"tz"`     //堂主人数
	DuoZ   int32 `bson:"duoz"`   //舵主人数
	ZM     int32 `bson:"zm"`     //掌门人数
	TD     int32 `bson:"td"`     //团队人数
	Income int32 `bson:"income"` //总收益
	Date   int32 `bson:"date"`   //日期
}

//昨日团队贡献
type LstICMDetail struct {
	UID    string  `bson:"uid"`    //贡献者唯一id
	Devote float64 `bson:"devote"` //贡献的金币
	Job    int8    `bson:"job"`    //贡献时的职位
	Golds  int32   `bson:"golds"`  //总金币数
}

//团队收益相关的记录信息
type DBTeamIncome struct {
	UID     string         `bson:"uid"`    //玩家唯一id
	Income  int32          `bson:"income"` //总收益
	HisT    []TeamHistory  `bson:"hist"`   //团队收益历史记录
	YDay    int16          `bson:"yday"`   //更新的标记(一年中的天数)
	TmpIcm  float64        `bson:"-"`      //临时的收益
	LstIcm  int32          `bson:"lsticm"` //昨日团队收益
	LstDvt  []LstICMDetail `bson:"lstdvt"` //昨日贡献明细
	HisIcm  int64          `bson:"hisicm"` //历史总收益
	TmpDay  int16          `bson:"-"`      //临时的标识
	TmpHour int16          `bson:"-"`      //临时的标识
}
