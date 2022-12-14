package lib

const (
	RetSuccess       = 0
	RetNotValidCode  = -10000 // 微信code校验失败
	RetInternalError = -10001 // 内部异常
	RetParamError    = -10002 // 参数错误
	RetTokenNotValid = -10003 // token失效

	RetNotFindPrize = -20002 // 未发现奖励

	RetNotFindBox = -20003 // 未找到盒子
)
