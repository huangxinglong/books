package mongo

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"strings"
	"zcm_tools/email"
)

// mongodb
//var (
//    MgoSession  *mgo.Session
//    MgoDbName = "zcm"
//    MgoUri = "mongodb://root:zcm2016@10.253.0.218:27017" // 生产阿里云
//)

// 群发收件人
var ToUsers = []string{
	"zxh@zcmlc.com",
	//"dongju@zcmlc.com",
	//"lxy@zcmlc.com",
	//"yiwei@zcmlc.com",
}

// 创建连接池
// 手动创建会话记得关闭连接 defer mgoSession.Close()
// @Author 朱学煌 zxh@zcmlc.com
// @Date 2017-04-27 15:30
// @Param mgoSession mgo.Session mgo会话对象
// @Param mgoUri string mgo连接地址
// @Return mgoSession mgo.Session mgo会话对象
func NewMgoSession(mgoSession *mgo.Session, mgoUri string) *mgo.Session {
	defer func() {
		if r := recover(); r != nil {
			email.SendEmail("NewMgoSessionErr", "Mongo重新连接 panic", strings.Join(ToUsers, ";"))
			var ok bool
			err, ok := r.(error)
			if !ok {
				email.SendEmail("NewMgoSessionErr", fmt.Sprintf("Mongo重新连接 %v \n error: %s", r, err.Error()), strings.Join(ToUsers, ";"))
			}
		}
	}()

	if mgoSession == nil {
		var err error
		mgoSession, err = mgo.Dial(mgoUri)
		if err != nil {
			email.SendEmail("NewMgoSessionErr", fmt.Sprintf("Mongo重新连接%s", err.Error()), strings.Join(ToUsers, ";"))
			mgoSession = NewMgoSession(mgoSession, mgoUri)
		}

		// 会话为单调行为 参数资料 http://m.blog.csdn.net/article/details?id=52314977
		//mgoSession.SetMode(mgo.Monotonic, true)
		mgoSession.SetMode(mgo.Eventual, true)

		// 设置连接池100 （默认最大连接池为4096）
		mgoSession.SetPoolLimit(200)
	}

	// Ping()
	pingRs := mgoSession.Ping()
	if pingRs != nil {
		email.SendEmail("MgoPingErr", fmt.Sprintf("Mongo重新连接%s", pingRs.Error()), strings.Join(ToUsers, ";"))
		mgoSession = NewMgoSession(mgoSession, mgoUri)
	}
	return mgoSession
}
