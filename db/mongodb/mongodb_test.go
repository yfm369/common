package mongodb

import (
	"fmt"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// func Test_DailDB(t *testing.T) {
// 	dburl := "mongodb://127.0.0.1:27017"
// 	t.Log("url: ", dburl)
// 	dbo := NewDbOperate(dburl, 5*time.Second)
// 	dbo.OpenDB(nil)
// 	dbo.CloseDB()
// }

// //go test -v github.com\viphxin\xingo\db\mongo -run ^Test_CommonOperate$
func Test_CommonOperate(t *testing.T) {
	dburl := "mongodb://127.0.0.1:27017"
	t.Log("url: ", dburl)
	dbo := NewDbOperate(dburl, 5*time.Second, 99)
	dbo.OpenDB(func(ms *mongo.Client) {
		IndexTable(ms, "test", "bulkUpdateInsert", "name", []string{"name"}, true)
	})
	defer func() {
		dbo.CloseDB()
	}()

	err := dbo.RemoveAll("test", "bulkUpdateInsert", bson.M{"name": "1"})
	if err != nil {
		//not do anything
	}
	//------------------------------------------------------------------------
	//err = dbo.Insert("test", "bulkUpdateInsert", bson.M{"age": "99", "name": "pass1111", "sex":"男"})
	//if err != nil {
	//	return
	//}
	//
	//err = dbo.Insert("test", "bulkUpdateInsert",bson.M{"age": "100", "name": "文文", "sex":"女"})
	//if err != nil {
	//	return
	//}

	var result interface{}
	err = dbo.DBFindOne("test", "bulkUpdateInsert", bson.M{"name": "文文"}, &result)
	if err != nil {
		return
	}
	t.Log("11111111111111 findOne ", result)
	//-------------------------------------------------------------------------
	docs := make([]interface{}, 0)
	for i := 0; i < 500; i++ {
		docs = append(docs, bson.M{"name": fmt.Sprintf("xingo_%d", i), "sex": "女", "age": 90})
	}
	err = dbo.BulkInsertDoc("test", "bulkUpdateInsert", docs)
	if err != nil {
		return
	}
	//
	//_, err = dbo.DeleteAll("test", bson.M{"pass": "pass1111"})
	//if err != nil {
	//	dbo.CloseDB()
	//	t.Fatal(err)
	//	return
	//}
}

// //go test -v github.com\viphxin\xingo\db\mongo -bench ^Benchmark_CommonOperate$
// func Benchmark_CommonOperate(b *testing.B) {
// 	dburl := "mongodb://127.0.0.1:27017"
// 	b.Log("url: ", dburl)
// 	dbo := NewDbOperate(dburl, 5*time.Second)
// 	dbo.OpenDB(func(ms *mgo.Session) {
// 		ms.DB("").C("test").EnsureIndex(mgo.Index{
// 			Key:    []string{"username"},
// 			Unique: true,
// 		})
// 	})

// 	for i := 0; i < b.N; i++ {
// 		_, err := dbo.DeleteAll("test", bson.M{"pass": "pass1111"})
// 		if err != nil {
// 			//not do anything
// 		}
// 		//------------------------------------------------------------------------
// 		err = dbo.Insert("test", bson.M{"username": "xingo", "pass": "pass1111"})
// 		if err != nil {
// 			dbo.CloseDB()
// 			b.Fatal(err)
// 			return
// 		}

// 		err = dbo.Insert("test", bson.M{"username": "xingo_0", "pass": "pass1111"})
// 		if err != nil {
// 			dbo.CloseDB()
// 			b.Fatal(err)
// 			return
// 		}

// 		err = dbo.DBFindOne("test", bson.M{"username": "xingo"}, func(a bson.M) error {
// 			if a != nil {
// 				b.Log(a)
// 				return nil
// 			} else {
// 				dbo.CloseDB()
// 				b.Fatal("DBFindOne error")
// 				return errors.New("DBFindOne error")
// 			}

// 		})
// 		if err != nil {
// 			dbo.CloseDB()
// 			b.Fatal(err)
// 			return
// 		}
// 		_, err = dbo.DeleteAll("test", bson.M{"pass": "pass1111"})
// 		if err != nil {
// 			dbo.CloseDB()
// 			b.Fatal(err)
// 			return
// 		}
// 		//-------------------------------------------------------------------------
// 		//bulk
// 		docs := make([]bson.M, 0)
// 		for i := 0; i < 500; i++ {
// 			docs = append(docs, bson.M{"username": fmt.Sprintf("xingo_%d", i), "pass": "pass1111"})
// 		}

// 		err = dbo.BulkInsert("test", docs)
// 		if err != nil {
// 			dbo.CloseDB()
// 			b.Fatal(err)
// 			return
// 		}

// 		_, err = dbo.DeleteAll("test", bson.M{"pass": "pass1111"})
// 		if err != nil {
// 			dbo.CloseDB()
// 			b.Fatal(err)
// 			return
// 		}
// 	}
// 	dbo.CloseDB()
// }

// //go test -v github.com\viphxin\xingo\db\mongo -bench ^Benchmark_CommonOperatePP$
// func Benchmark_CommonOperatePP(b *testing.B) {
// 	dburl := "mongodb://127.0.0.1:27017"
// 	b.Log("url: ", dburl)
// 	dbo := NewDbOperate(dburl, 5*time.Second)
// 	dbo.OpenDB(func(ms *mgo.Session) {
// 		ms.DB("").C("test").DropIndex("username")
// 	})
// 	_, err := dbo.DeleteAll("test", bson.M{"pass": "pass1111"})
// 	if err != nil {
// 		dbo.CloseDB()
// 		b.Fatal(err)
// 		return
// 	}
// 	b.RunParallel(func(pb *testing.PB) {
// 		for pb.Next() {
// 			//------------------------------------------------------------------------
// 			err := dbo.Insert("test", bson.M{"username": "xingo", "pass": "pass1111"})
// 			if err != nil {
// 				dbo.CloseDB()
// 				b.Fatal(err)
// 				return
// 			}

// 			err = dbo.DBFindOne("test", bson.M{"username": "xingo"}, func(a bson.M) error {
// 				if a != nil {
// 					b.Log(a)
// 					return nil
// 				} else {
// 					dbo.CloseDB()
// 					b.Fatal("DBFindOne error")
// 					return errors.New("DBFindOne error")
// 				}

// 			})
// 			if err != nil {
// 				dbo.CloseDB()
// 				b.Fatal(err)
// 				return
// 			}
// 			//-------------------------------------------------------------------------
// 			//bulk
// 			docs := make([]bson.M, 0)
// 			for i := 0; i < 500; i++ {
// 				docs = append(docs, bson.M{"username": fmt.Sprintf("xingo_%d", i), "pass": "pass1111"})
// 			}

// 			err = dbo.BulkInsert("test", docs)
// 			if err != nil {
// 				dbo.CloseDB()
// 				b.Fatal(err)
// 				return
// 			}
// 		}
// 	})
// 	dbo.CloseDB()
// }

//func Test_bulkUpInsert(t *testing.T) {
//	dbUrl := "mongodb://127.0.0.1:27017"
//	t.Log("url: ", dbUrl)
//	dbo := NewDbOperate(dbUrl, 5*time.Second, 99)
//	defer dbo.CloseDB()
//
//	dbo.OpenDB(func(ms *mongo.Client) {
//		//ms.DB("").C("test.gfs").EnsureIndexKey("filename")
//	})
//
//	type TestBulk struct {
//		ID   primitive.ObjectID `bson:"_id"`
//		Name string
//		Age  int32
//		Sex  string
//	}
//
//	nums := 0
//	doc := make([]interface{}, 0)
//	for i := 0; i < 10; i++ {
//		tmp := &TestBulk{ID: primitive.NewObjectID(), Name: fmt.Sprintf("%d", i), Age: 188, Sex: "男"}
//		//idb := bson.D{{"_id", tmp.ID}}
//		//upb := bson.D{{"$set", bson.D{{"ID", primitive.NewObjectID()}, {"Name",
//			//fmt.Sprintf("%d", i)}, {"Age", 188}, {"Sex", "男"}}}}
//		idb := bson.M{"_id":tmp.ID}
//		upb := bson.M{"$set":tmp}
//		doc = append(doc, idb, upb)
//		nums++
//		if nums >= 10 {
//			nums = 0
//			err := dbo.BulkUpsert("test", "bulkUpdateInsert", doc)
//			if err != nil {
//				t.Log(err.Error())
//			}
//
//			doc = make([]interface{}, 0)
//		}
//	}
//}
