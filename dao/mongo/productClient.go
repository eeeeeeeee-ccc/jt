package mongo

import (
	"context"
	"fmt"
	jkinterface "github.com/eeeeeeeee-ccc/jt/interface"
	Err "github.com/eeeeeeeee-ccc/jt/model/client_err"
	Kv "github.com/eeeeeeeee-ccc/jt/model/kv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

type MogClient struct {
	Client    *mongo.Client
}


func New(ExtMap map[string]string) jkinterface.ProductClientInterface{
	var mog MogClient
	mongodbConnectInfo:=ExtMap["mongodb_connect_info"]
	client, err := mongo.NewClient(options.Client().ApplyURI(mongodbConnectInfo))
	fmt.Println(err)
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	mog.Client = client
	return &mog
}



func (m *MogClient)PutCollection(project, setName string,group *Kv.CollectionGroup,extMap map[string]string)error{
	var errs Err.Error
	if len(group.Collections) == 0 {
		// empty log group
		errs=Err.Error{
			HttpCode: 0,
			Code:     1,
			Msg:      "mongodb 批次内容为空>>",
		}
		return errs
	}
	coll:=m.Client.Database(project).Collection(setName)
	subArr:=[]interface{}{}
	for _,item:=range group.Collections{
		 bD:=bson.D{}
		for _,akv:=range item.Content{
			e:=bson.E{
				Key:   *akv.Key,
				Value: akv.Value,
			}
			bD=append(bD,e)
		}
		subArr=append(subArr,bD)
	}
	_, err := coll.InsertMany(context.TODO(), subArr)
	if err!=nil{
		errs=Err.Error{
			HttpCode: 0,
			Code:     1,
			Msg:      "mongodb 提交错误>>"+err.Error(),
		}
	}
	return errs
}

