package model

type PaymentInfo struct {
	Id         int64  //
	UserId     string // 用户id
	GroupId    string // 团id(与userid联合唯一)
	CreateTime int64  // 创建时间
	PayTime    int64  // 支付时间
	Status     int    // 状态
	OrderId    string // 订单id
	ProveInfo  string // 支付截图
}
