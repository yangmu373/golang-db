package db

import (
	"log"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/astaxie/beego"
)

type Mongo struct {
}

func NewMongo() *Mongo {
	return &Mongo{}
}

var (
	mgosess *mgo.Session
)

func (this *Mongo) conn() (*mgo.Session, error) {
	user := beego.AppConfig.String("mongouser")
	pass := beego.AppConfig.String("mongopass")
	host1 := beego.AppConfig.String("mongohost1")
	port1 := beego.AppConfig.String("mongoport1")
	host2 := beego.AppConfig.String("mongohost2")
	port2 := beego.AppConfig.String("mongoport2")
	db := beego.AppConfig.String("mongodb")

	var err error
	if mgosess == nil {
		mgosess, err = mgo.Dial("mongodb://" + user + ":" + pass + "@" + host1 + ":" + port1 + "," + host2 + ":" + port2 + "/" + db + "")
		if err != nil {
			return mgosess, err
		}
		//mgosess.SetMode(mgo.Monotonic, true)
	}
	return mgosess.Clone(), err
}

func (this *Mongo) Insert(collectionName string, data interface{}) error {
	session, err := this.conn()
	defer session.Close()
	if err != nil {
		log.Println(err)
		return err
	}
	c := session.DB("").C(collectionName)
	err = c.Insert(data)
	return err
}

func (this *Mongo) FindOne(collectionName string, cond bson.M, data interface{}) error {
	session, err := this.conn()
	defer session.Close()
	if err != nil {
		log.Println(err)
		return err
	}
	c := session.DB("").C(collectionName)
	err = c.Find(cond).One(data)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (this *Mongo) Find(collectionName string, data interface{}, find bson.M, skip int, limit int, sort ...string) error {
	session, err := this.conn()
	defer session.Close()
	if err != nil {
		log.Println(err)
		return err
	}
	c := session.DB("").C(collectionName)
	err = c.Find(find).Sort(sort...).Skip(skip).Limit(limit).All(data)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (this *Mongo) Count(collectionName string, find bson.M) (int, error) {
	session, err := this.conn()
	defer session.Close()
	if err != nil {
		log.Println(err)
		return 0, err
	}
	c := session.DB("").C(collectionName)
	n, err := c.Find(find).Count()
	return n, err
}

func (this *Mongo) Remove(collectionName string, selector interface{}) error {
	session, err := this.conn()
	defer session.Close()
	if err != nil {
		log.Println(err)
		return err
	}
	c := session.DB("").C(collectionName)
	err = c.Remove(selector)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (this *Mongo) Update(collectionName string, selector interface{}, update interface{}) error {
	session, err := this.conn()
	defer session.Close()
	if err != nil {
		log.Println(err)
		return err
	}
	c := session.DB("").C(collectionName)
	err = c.Update(selector, update)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (this *Mongo) AddIndexTTL(collectionName string, fieldName string, ttl time.Duration) error {
	session, err := this.conn()
	defer session.Close()
	if err != nil {
		log.Println(err)
		return err
	}
	c := session.DB("").C(collectionName)
	var index = mgo.Index{
		Key:         []string{fieldName},
		ExpireAfter: time.Second * ttl,
	}
	err = c.EnsureIndex(index)
	return err
}

func (this *Mongo) AddIndexUnique(collectionName string, fieldName string) error {
	session, err := this.conn()
	defer session.Close()
	if err != nil {
		log.Println(err)
		return err
	}
	c := session.DB("").C(collectionName)
	var index = mgo.Index{
		Key:    []string{fieldName},
		Unique: true,
	}
	err = c.EnsureIndex(index)
	return err
}

func (this *Mongo) AddIndexUnionUnique(collectionName string, fieldNames []string) error {
	session, err := this.conn()
	defer session.Close()
	if err != nil {
		log.Println(err)
		return err
	}
	c := session.DB("").C(collectionName)
	var index = mgo.Index{
		Key:    fieldNames,
		Unique: true,
	}
	err = c.EnsureIndex(index)
	return err
}

func (this *Mongo) AddIndex(collectionName string, fieldNames []string) error {
	session, err := this.conn()
	defer session.Close()
	if err != nil {
		log.Println(err)
		return err
	}
	c := session.DB("").C(collectionName)
	var index = mgo.Index{
		Key: fieldNames,
	}
	err = c.EnsureIndex(index)
	return err
}

func (this *Mongo) Create(collectionName string) error {
	session, err := this.conn()
	defer session.Close()
	if err != nil {
		log.Println(err)
		return err
	}
	c := session.DB("").C(collectionName)
	var info = &mgo.CollectionInfo{}
	err = c.Create(info)
	return err
}
