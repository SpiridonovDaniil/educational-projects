package mongo

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"serverdb/models"
)

type Repo struct{
	session *mgo.Session
	client *mgo.Collection
}

func NewRepo(address, db, collection string) *Repo {
	session, e := mgo.Dial(address)
	if e != nil {
		log.Fatalln(e)
	}

	userCollection := session.DB(db).C(collection)

	return &Repo{
		session: session,
		client: userCollection,
	}
}

func(r *Repo) Close() {
	r.session.Close()
}

func(r *Repo) FindUser(userID string)(models.User, error) {
	query := bson.M{
		"name": bson.M{
			"$eq": userID,
		},
	}
	user := models.User{}
	err := r.client.Find(query).One(&user)
	if err != nil {
		return user ,err
	}

	return user, nil
}

func (r *Repo) Insert(user models.User) error {
	err := r.client.Insert(user)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repo) Update(user, order, object string, objectChange interface{}) error {
	err := r.client.Update(bson.M{"name": user}, bson.M{order:bson.M{object:objectChange}})
	if err != nil {
		return err
	}

	return nil
}

func (r *Repo) RemoveAll(userID string) error {
	_, err := r.client.RemoveAll(bson.M{"name": userID})
	if err != nil {
		return err
	}

	return nil
}