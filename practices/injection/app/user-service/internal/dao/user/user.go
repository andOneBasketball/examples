package user

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/util/gconv"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Dao 用户数据管理对象。
type Dao interface {
	Create(ctx context.Context, in CreateInput) (string, error)
	Update(ctx context.Context, id primitive.ObjectID, in UpdateInput) error
	Delete(ctx context.Context, id []primitive.ObjectID) error
	GetOne(ctx context.Context, id primitive.ObjectID) (*DataItem, error)
	GetList(ctx context.Context, in GetListInput) ([]*DataItem, error)
}

type implUser struct {
	database   *mongo.Database
	collection *mongo.Collection
	fields     collectionFieldNames
}

type collectionFieldNames struct {
	Id        string
	Name      string
	CreatedAt string
	UpdatedAt string
}

func New(db *mongo.Database) Dao {
	return &implUser{
		database:   db,
		collection: db.Collection("user"),
		fields: collectionFieldNames{
			Id:        "_id",
			Name:      "name",
			CreatedAt: "created_at",
			UpdatedAt: "updated_at",
		},
	}
}

type DataItem struct {
	Id        primitive.ObjectID `bson:"_id,omitempty"`
	Name      string             `bson:"name,omitempty"`
	CreatedAt int64              `bson:"created_at,omitempty"`
	UpdatedAt int64              `bson:"updated_at,omitempty"`
}

type CreateInput struct {
	Name string `bson:"name,omitempty"`
}

// Create 新增用户。
func (r *implUser) Create(ctx context.Context, in CreateInput) (string, error) {
	var dataItem = DataItem{
		CreatedAt: time.Now().UnixMilli(),
		UpdatedAt: time.Now().UnixMilli(),
	}
	if err := gconv.Scan(in, &dataItem); err != nil {
		return "", err
	}
	result, err := r.collection.InsertOne(ctx, dataItem)
	if err != nil {
		return "", errors.Wrap(err, "insert user data failed")
	}
	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

type UpdateInput struct {
	Name string `bson:"name,omitempty"`
}

// Update 修改用户。
func (r *implUser) Update(ctx context.Context, id primitive.ObjectID, in UpdateInput) error {
	var dataItem = DataItem{
		UpdatedAt: time.Now().UnixMilli(),
	}
	if err := gconv.Scan(in, &dataItem); err != nil {
		return err
	}
	filter := bson.D{
		{r.fields.Id, id},
	}
	update := bson.M{
		"$set": dataItem,
	}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return errors.Wrap(err, "update user data failed")
	}
	return nil
}

// Delete 删除用户（硬删除）。
func (r *implUser) Delete(ctx context.Context, ids []primitive.ObjectID) error {
	if len(ids) == 0 {
		return nil
	}
	filter := bson.D{
		{r.fields.Id, bson.M{"$in": ids}},
	}
	_, err := r.collection.DeleteMany(ctx, filter)
	if err != nil {
		return errors.Wrap(err, "delete user data failed")
	}
	return nil
}

// GetOne 查询用户详情。
func (r *implUser) GetOne(ctx context.Context, id primitive.ObjectID) (*DataItem, error) {
	var (
		dataItem *DataItem
		filter   = bson.D{
			{r.fields.Id, id},
		}
	)
	err := r.collection.FindOne(ctx, filter).Decode(&dataItem)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, errors.Wrap(err, "find one user data failed")
	}
	return dataItem, nil
}

type GetListInput struct {
	Ids []primitive.ObjectID // 通过id查询
}

// GetList 查询用户列表。
func (r *implUser) GetList(ctx context.Context, in GetListInput) ([]*DataItem, error) {
	var (
		filter = bson.D{}
		opts   = options.Find().SetSort(bson.M{r.fields.Id: 1})
	)
	if len(in.Ids) > 0 {
		filter = append(filter, bson.E{Key: r.fields.Id, Value: bson.M{"$in": in.Ids}})
	}
	cur, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, errors.Wrap(err, `search GetList failed`)
	}
	defer cur.Close(ctx)
	// 查询数据到实体对象中
	var dataItems = make([]*DataItem, 0)
	if err = cur.All(ctx, &dataItems); err != nil {
		return nil, errors.Wrap(err, `mongodb scan result failed`)
	}
	return dataItems, nil
}
