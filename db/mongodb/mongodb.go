package mongodb

import (
	"context"
	"errors"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

var (
	errMongodbSessionNil   = errors.New("DbOperate session nil")
	errMongodbDbFindAll    = errors.New("DBFindAll failed,q is nil")
	LoggerName = "mongodb"
)

type DbOperate struct {
	client *mongo.Client
	timeout time.Duration
	dbUrl   string
	poolSize uint64
	ctx context.Context
}

func NewDbOperate(dbUrl string, timeout time.Duration, pool uint64) *DbOperate {
	return &DbOperate{nil, timeout, dbUrl, pool, nil}
}

func (db *DbOperate) OpenDB(setIndexFunc func(view *mongo.Client)) error {
	var (
		ctx = context.TODO()
	)

	g.Log(LoggerName).Info(ctx, "DbOperate mongodb connect url: ", db.dbUrl)

	ctxMongo, cancel := context.WithTimeout(context.Background(), db.timeout * time.Second)
	defer func() {
		cancel()
	}()

	client, err := mongo.NewClient(options.Client().
		ApplyURI(db.dbUrl).
		SetMaxPoolSize(db.poolSize))
	if err != nil {
		g.Log(LoggerName).Error(ctx, "OpenDB NewClient error ", err.Error())
		return  err
	}
	err = client.Connect(ctxMongo)
	if err != nil {
		g.Log(LoggerName).Error(ctx, "OpenDB Connect error ", err.Error())
		return err
	}
	db.client = client

	//由于阿里云数据库的连接数限制 暂时设置为900
	if setIndexFunc != nil {
		setIndexFunc(client)
	}

	g.Log(LoggerName).Info(ctx,"DbOperate connect ", db.dbUrl, " mongodb...OK")

	return nil
}

func (db *DbOperate) CloseDB() {
	var (
		ctx = context.TODO()
	)

	if db.client != nil {
		err := db.client.Disconnect(db.ctx)
		if err != nil {
			g.Log(LoggerName).Error(context.TODO(), "CloseDB error ", err.Error())
		}
		db.client = nil

		g.Log(LoggerName).Info(ctx,"Disconnect mongodb url: ", db.dbUrl)
	}
}

func (db *DbOperate) Insert(dbname, collection string, doc interface{}) error {
	var (
		ctx = context.TODO()
	)

	if db.client == nil {
		g.Log(LoggerName).Info(ctx,"Insert session = nil dbname = ", dbname,
			"collection is:", collection, " doc:", doc)
		return errMongodbSessionNil
	}

	c := db.client.Database(dbname).Collection(collection)
	_, err := c.InsertOne(context.TODO(), doc)
	if err != nil {
		g.Log(LoggerName).Error(ctx,"Insert error :", dbname, " col:", collection, " doc:",
			doc, " err:", err.Error())
	}

	return err
}

func (db *DbOperate) Cover(dbname, collection string, cond interface{}, change interface{}) error {
	if db.client == nil {
		return errMongodbSessionNil
	}

	c := db.client.Database(dbname).Collection(collection)
	_, err := c.UpdateOne(context.TODO(), cond, change)

	return err
}

func (db *DbOperate) Update1(dbname string, collection string, selector interface{}, update interface{}) error {
	if db.client == nil {
		g.Log(LoggerName).Error(context.TODO(),"Update1 session = nil dbname = ", dbname,
			"collection is:", collection, " selector:", selector, " update:", update)
		return errMongodbSessionNil
	}
	
	col := db.client.Database(dbname).Collection(collection)
	_, err := col.UpdateOne(context.TODO(), selector, update)
	if err != nil && err.Error() != "not found" {
		g.Log(LoggerName).Error(context.TODO(), fmt.Sprintf("Update1 error:%s", err.Error()))
		return err
	}

	return nil
}

func (db *DbOperate) Update(dbname, collection string, cond interface{}, change interface{}) error {
	if db.client == nil {
		g.Log(LoggerName).Error(context.TODO(), "Update session = nil dbname = ", dbname,
			"collection is:", collection, " cond:", cond, " change:", change)
		return errMongodbSessionNil
	}

	c := db.client.Database(dbname).Collection(collection)
	_, err := c.UpdateOne(context.TODO(), cond, bson.M{"$set": change})
	if err != nil {
		g.Log(LoggerName).Error(context.TODO(), "Update error dbname = ", dbname, "collection is:",
			collection, " cond:", cond, " change:", change, " err:", err.Error())
	}

	return err
}

func (db *DbOperate) UpsertID(dbname, collection string, cond interface{}, doc interface{}) error {
	if db.client == nil {
		g.Log(LoggerName).Error(context.TODO(), "UpsertID session = nil dbname = ", dbname,
			"collection is:", collection, " cond is:", cond, " doc:", doc)
		return errMongodbSessionNil
	}

	c := db.client.Database(dbname).Collection(collection)
	_, err := c.UpdateByID(context.TODO(), cond, doc)
	if err != nil {
		g.Log(LoggerName).Error(context.TODO(), "UpsertID failed dbname = ", dbname,
			"collection is:", collection, " cond is:", cond, " doc:", doc, " err:", err.Error())
	}

	return err
}

func (db *DbOperate) UpdateInsert(dbname, collection string, cond interface{}, doc interface{}) error {
	if db.client == nil {
		g.Log(LoggerName).Error(context.TODO(), "UpdateInsert session = nil dbname = ", dbname,
			"collection is:", collection, " cond is:", cond, " doc:", doc)
		return errMongodbSessionNil
	}

	opts := options.Update().SetUpsert(true)
	c := db.client.Database(dbname).Collection(collection)
	_, err := c.UpdateOne(context.TODO(), cond, bson.M{"$set": doc}, opts)
	if err != nil {
		g.Log(LoggerName).Error(context.TODO(), "UpdateInsert failed dbname = ", dbname,
			"collection is:", collection, " cond is:", cond, " doc:", doc, " err:", err.Error())
	}

	return err
}

func (db *DbOperate) UpdateInsert1(dbname, collection string, cond interface{}, doc interface{}) error {
	if db.client == nil {
		g.Log(LoggerName).Error(context.TODO(), "UpdateInsert1 session = nil dbname = ", dbname,
			"collection is:", collection, " cond is:", cond, " doc:", doc)
		return errMongodbSessionNil
	}

	c := db.client.Database(dbname).Collection(collection)
	opts := options.Update().SetUpsert(true)
	_, err := c.UpdateOne(context.TODO(), cond, doc, opts)
	if err != nil {
		g.Log(LoggerName).Error(context.TODO(), "UpdateInsert1 failed dbname = ", dbname,
			"collection is:", collection, " cond is:", cond, " doc:", doc, " err:", err.Error())
	}

	return err
}

func (db *DbOperate) RemoveOne(dbname, collection string, condName string, condValue int64) error {
	if db.client == nil {
		g.Log(LoggerName).Error(context.TODO(), "RemoveOne session = nil dbname = ", dbname,
			"collection is:", collection, " cond:", condName, " cond_value:", condValue)
		return errMongodbSessionNil
	}

	c := db.client.Database(dbname).Collection(collection)
	//opts := options.Delete().SetCollation(&options.Collation{
	//	Locale:    "en_US",
	//	Strength:  1,
	//	CaseLevel: false,
	//})
	_, err := c.DeleteOne(context.TODO(), bson.M{condName: condValue})
	if err != nil {//mgo.ErrNotFound
		g.Log(LoggerName).Error(context.TODO(), "remove failed from dbname :", dbname,
			" collection:", collection, " name:", condName, " value:", condValue)
	}

	return err
}

func (db *DbOperate) RemoveOneByCond(dbname, collection string, cond interface{}) error {
	if db.client == nil {
		g.Log(LoggerName).Error(context.TODO(), "RemoveOneByCond session = nil dbname = ", dbname,
			"collection is:", collection, " cond:", cond)
		return errMongodbSessionNil
	}

	c := db.client.Database(dbname).Collection(collection)
	_, err := c.DeleteOne(context.TODO(), cond)
	if err != nil {
		g.Log(LoggerName).Error(context.TODO(), "remove failed from dbname:", dbname, " collection:",
			collection, " cond :", cond, " err:", err.Error())
	}

	return err
}

func (db *DbOperate) RemoveAll(dbname, collection string, cond interface{}) error {
	if db.client == nil {
		g.Log(LoggerName).Error(context.TODO(), "RemoveAll session = nil dbname = ", dbname,
			"collection is:", collection, " cond:", cond)
		return errMongodbSessionNil
	}

	c := db.client.Database(dbname).Collection(collection)
	change, err := c.DeleteMany(context.TODO(), cond)
	if err != nil {
		g.Log(LoggerName).Error(context.TODO(), "DbOperate.RemoveAll failed from dbname:", dbname,
			" collection:", collection, " cond :", cond, " err:", err.Error())
		return err
	}

	g.Log(LoggerName).Info(context.TODO(),"DbOperate.RemoveAll :", change.DeletedCount)

	return nil
}

func (db *DbOperate) DBFindOne(dbname, collection string, cond interface{}, result interface{}) error {
	if db.client == nil {
		g.Log(LoggerName).Error(context.TODO(), "DBFindOne session = nil dbname = ", dbname,
			"collection is:", collection, " cond:", cond)
		return errMongodbSessionNil
	}

	c := db.client.Database(dbname).Collection(collection)
	err := c.FindOne(context.TODO(), cond).Decode(result)
	if err != nil {
		if err.Error() != "not found" {
			g.Log(LoggerName).Error(context.TODO(), "DBFindOne query failed,return error: ", err.Error(),
				"dbname", dbname, " name: ", collection)
		}
		return err
	}

	return nil
}

func (db *DbOperate) DBFindAll(dbname, collection string, cond interface{}, result interface{}) error {
	if db.client == nil {
		g.Log(LoggerName).Error(context.TODO(), "DBFindAll session = nil dbname = ", dbname,
			"collection is:", collection, " cond:", cond)
		return errMongodbSessionNil
	}

	c := db.client.Database(dbname).Collection(collection)
	cursor, err := c.Find(context.TODO(), cond)
	if err != nil {
		g.Log(LoggerName).Error(context.TODO(), "FindAll error:", err.Error(), " coll:", collection)
		return err
	}
	if err = cursor.All(context.TODO(), result); err != nil {
		g.Log(LoggerName).Error(context.TODO(), "FindAll cursor error:", err.Error(), " coll:", collection)
		return err
	}

	return nil
}

func (db *DbOperate) DBFindSortAll(dbname, collection string, cond interface{},
	result interface{}, sort interface{}) error {
	if db.client == nil {
		g.Log(LoggerName).Error(context.TODO(), "DBFindSortAll session = nil dbname = ", dbname,
			"collection is:", collection, " cond:", cond)
		return errMongodbSessionNil
	}

	c := db.client.Database(dbname).Collection(collection)
	opts := options.Find().SetSort(sort)
	cursor, err := c.Find(context.TODO(), cond, opts)
	if err != nil {
		g.Log(LoggerName).Error(context.TODO(), "DBFindSortAll error:", err.Error(), " coll:", collection)
		return err
	}
	if err = cursor.All(context.TODO(), result); err != nil {
		g.Log(LoggerName).Error(context.TODO(), "DBFindSortAll cursor error:", err.Error(), " coll:", collection)
		return err
	}

	return nil
}

func (db *DbOperate) FindBySortLimit(dbname, collection string, find, result interface{},
	limit int64, sort interface{}) error {
	if db.client == nil {
		g.Log(LoggerName).Error(context.TODO(), "FindBySortLimit session = nil dbname = ", dbname,
			"collection is:", collection, " cond:", find)
		return errMongodbSessionNil
	}

	c := db.client.Database(dbname).Collection(collection)
	opts := options.Find().SetSort(sort).SetLimit(limit)
	cursor, err := c.Find(context.TODO(), find, opts)
	if err != nil {
		g.Log(LoggerName).Error(context.TODO(), "FindBySortLimit error:", err.Error(), " coll:", collection)
		return err
	}
	if err = cursor.All(context.TODO(), result); err != nil {
		g.Log(LoggerName).Error(context.TODO(), "FindBySortLimit cursor error:", err.Error(), " coll:", collection)
		return err
	}

	return nil
}

func (db *DbOperate) FindByProjectSort(dbname string, collection string, find interface{},
	project interface{}, sort interface{}, limit int64, result interface{}) error {
	if db.client == nil {
		g.Log(LoggerName).Error(context.TODO(), "FindByProjectSort session = nil dbname = ", dbname,
			"collection is:", collection)
		return errMongodbSessionNil
	}

	c := db.client.Database(dbname).Collection(collection)
	opts := options.Find()
	if project != nil {
		opts = opts.SetProjection(project)
	}
	if sort != nil {
		opts = opts.SetSort(sort)
	}
	if limit > 0 {
		opts = opts.SetLimit(limit)
	}
	cursor, err := c.Find(context.TODO(), find, opts)
	if err != nil {
		g.Log(LoggerName).Error(context.TODO(), "FindByProjectSort error:", err.Error(), " coll:", collection)
		return err
	}
	if err = cursor.All(context.TODO(), result); err != nil {
		g.Log(LoggerName).Error(context.TODO(), "FindByProjectSort cursor error:", err.Error(), " coll:", collection)
		return err
	}

	return nil
}

func (db *DbOperate) FindCount(dbname string, collection string, find interface{}) (error, int) {
	if db.client == nil {
		g.Log(LoggerName).Error(context.TODO(), "DBFindAll session = nil dbname = ", dbname,
			"collection is:", collection, " find:", find)
		return errMongodbSessionNil, -1
	}

	c := db.client.Database(dbname).Collection(collection)
	cursor, err := c.Find(context.TODO(), find)
	if err != nil {
		g.Log(LoggerName).Error(context.TODO(), "FindByProjectSort error:", err.Error(), " coll:", collection)
		return err, 0
	}

	count := cursor.RemainingBatchLength()
	return nil, count
}

func (db *DbOperate) DBFindAllEx(dbname, collection string, cond interface{}, sort interface{}, resHandler func(cursor *mongo.Cursor) error) error {
	if db.client == nil {
		g.Log(LoggerName).Error(context.TODO(), "DBFindAllEx session = nil dbname = ", dbname,
			"collection is:", collection, " cond:", cond)
		return errMongodbSessionNil
	}

	c := db.client.Database(dbname).Collection(collection)
	opts := options.Find().SetSort(sort)
	cursor, err := c.Find(context.TODO(), cond, opts)
	if err != nil {
		g.Log(LoggerName).Error(context.TODO(), "DBFindAllEx error:", err.Error(), " coll:", collection)
		return err
	}

	g.Log(LoggerName).Info(context.TODO(),"[DbOperate.DBFindAll] name:", collection, "query:", cond)

	if nil == cursor {
		return errMongodbDbFindAll
	}

	if nil != resHandler {
		return resHandler(cursor)
	}

	return nil
}

func (db *DbOperate) FindAndModify(dbname, collection string, cond interface{}, change interface{}, result interface{}) error {
	if db.client == nil {
		g.Log(LoggerName).Error(context.TODO(), "FindAndModify session = nil dbname = ", dbname,
			"collection is:", collection, " cond:", cond)
		return errMongodbSessionNil
	}

	c := db.client.Database(dbname).Collection(collection)
	opts := options.FindOneAndUpdate().SetUpsert(true)
	err := c.FindOneAndUpdate(context.TODO(), cond, change, opts).Decode(result)
	if err != nil {
		g.Log(LoggerName).Error(context.TODO(), "FindAndModify error:", err.Error(), " coll:", collection)
		return err
	}

	return nil
}

//BulkInsertDoc 批量插入文档
func (db *DbOperate) BulkInsertDoc(dbname, collection string, docs []interface{}) error {
	if db.client == nil {
		g.Log(LoggerName).Error(context.TODO(), "BulkInsertDoc session = nil dbname = ", dbname,
			"collection is:", collection, " docs:", docs)
		return errMongodbSessionNil
	}

	c := db.client.Database(dbname).Collection(collection)
	res, err := c.InsertMany(context.TODO(), docs)
	if err != nil {
		g.Log(LoggerName).Error(context.TODO(), "BulkInsertDoc InsertMany error:", err.Error(), " coll:", collection,
			" docs = ", docs)
		return err
	}

	g.Log(LoggerName).Info(context.TODO(),"BulkInsertDoc matched:", len(res.InsertedIDs), " coll:", collection)

	return nil
}

//BulkUpdate 批量更新 参数为interface Bson.D
func (db *DbOperate) BulkUpdate(dbname, collection string, pairs []interface{}) error {
	if db.client == nil {
		g.Log(LoggerName).Error(context.TODO(), "BulkUpdate session = nil dbname = ", dbname,
			"collection is:", collection, " pairs:", pairs)
		return errMongodbSessionNil
	}

	c := db.client.Database(dbname).Collection(collection)

	var models []mongo.WriteModel
	for i := 0; i < len(pairs); i += 2 {
		selector := pairs[i]
		update := pairs[i+1]
		updateOne := mongo.NewUpdateOneModel()
		updateOne.SetFilter(selector).SetUpdate(update)
		models = append(models, updateOne)
	}

	opts := options.BulkWrite().SetOrdered(false)
	res, err := c.BulkWrite(context.TODO(), models, opts)
	if err != nil {
		g.Log(LoggerName).Error(context.TODO(), "BulkUpdate BulkWrite error ", err.Error(), " res : ", res)
		return err
	}

	g.Log(LoggerName).Info(context.TODO(),"BulkUpdate over res = ", res)

	return nil
}

func (db *DbOperate) BulkUpsert(dbname, collection string, pairs []interface{}) error {
	if db.client== nil {
		g.Log(LoggerName).Error(context.TODO(), "BulkUpsert session = nil dbname = ", dbname,
			"collection is:", collection, " pairs:", pairs)
		return errMongodbSessionNil
	}

	c := db.client.Database(dbname).Collection(collection)

	var models []mongo.WriteModel
	for i := 0; i < len(pairs); i += 2 {
		selector := pairs[i]
		update := pairs[i+1]
		updateOne := mongo.NewUpdateOneModel()
		updateOne.SetFilter(selector).SetUpdate(update).SetUpsert(true)
		models = append(models, updateOne)
	}

	opts := options.BulkWrite().SetOrdered(false)
	res, err := c.BulkWrite(context.TODO(), models, opts)
	if err != nil {
		g.Log(LoggerName).Error(context.TODO(), "BulkUpsert BulkWrite error ", err.Error(), " res : ", res)
		return err
	}

	g.Log(LoggerName).Info(context.TODO(),"BulkUpsert over res = ", res)

	return nil
}

func (db *DbOperate) BulkUpdAll(dbname, collection string, pairs []interface{}) error {
	if db.client == nil {
		g.Log(LoggerName).Error(context.TODO(), "BulkUpdAll session = nil dbname = ", dbname,
			"collection is:", collection)
		return errMongodbSessionNil
	}

	c := db.client.Database(dbname).Collection(collection)

	var models []mongo.WriteModel
	for i := 0; i < len(pairs); i += 2 {
		selector := pairs[i]
		update := pairs[i+1]
		updateMany := mongo.NewUpdateManyModel()
		updateMany.SetFilter(selector).SetUpdate(update)
		models = append(models, updateMany)
	}

	opts := options.BulkWrite().SetOrdered(false)
	res, err := c.BulkWrite(context.TODO(), models, opts)
	if err != nil {
		g.Log(LoggerName).Error(context.TODO(), "BulkUpdAll BulkWrite error ", err.Error(), " res : ", res)
		return err
	}

	g.Log(LoggerName).Info(context.TODO(),"BulkUpdAll over res = ", res)

	return nil
}

func (db *DbOperate) GetMaxId(dbname, collection string, field string) (int64, error) {
	var id int64
	var present bool

	fnc := func(mq *mongo.Cursor) error {
		result := make(bson.M)
		err := mq.All(context.TODO(), &result)
		if err != nil {
			g.Log(LoggerName).Error(context.TODO(), "GetMaxId db error ", err.Error())
			return  err
		}
		id, present = result[field].(int64)
		if !present {
			id = 0
		}
		return nil
	}

	opts := options.Find().SetSort(bson.D{{field, -1 }}).SetLimit(1)
	err := db.DBFindAllEx(dbname, collection, nil, opts, fnc)
	if nil != err {
		return 0, nil
	}

	return id, nil
}

//IndexTable 创建索引
func IndexTable(client *mongo.Client, dbname string, collectionName string, indexName string,
	key []string, unique bool) {
	c := client.Database(dbname).Collection(collectionName)
	view := c.Indexes()

	var keys bson.D
	for _,v := range key {
			keys = append(keys, bson.E{Key: v, Value: 1})
	}

	model := mongo.IndexModel{Keys:keys,Options: options.Index().SetName(indexName).SetUnique(unique)}
	str,err := view.CreateOne(context.TODO(), model)
	if err != nil {
		g.Log(LoggerName).Error(context.TODO(), "creat index error ", err.Error())
		return
	}

	g.Log(LoggerName).Info(context.TODO(), "create index succeed ", str)
}
