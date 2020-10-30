package dal

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/skeris/flat-grabber/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

const (
	HouseColl = "houses"
	CRColl    = "call_requests"
	LastPage  = "config"
)

type DAL struct {
	db     *mongo.Database
	logger *zap.Logger
	reqInd int64
}

func New(db *mongo.Database, log *zap.Logger) *DAL {
	cnt, _ := db.Collection(CRColl).CountDocuments(context.Background(), bson.M{})
	return &DAL{
		db:     db,
		logger: log,
		reqInd: cnt,
	}
}

func (d *DAL) PutHouse(ctx context.Context, house *model.Flat) error {

	err := d.db.Collection(HouseColl).FindOneAndReplace(ctx, bson.M{"_id": house.ID}, house, options.FindOneAndReplace().SetUpsert(true)).Decode(&model.Flat{})
	if err == mongo.ErrNoDocuments {
		return nil
	}
	return err
}

type config struct {
	ID string `bson:"_id"`
	LP int64  `bson:"lp"`
}

func (d *DAL) GetLastPage(ctx context.Context) (int64, error) {
	var conf config

	if err := d.db.Collection(LastPage).FindOne(ctx, bson.M{"_id": "lasstPage"}).Decode(&conf); err != nil {
		return 0, err
	}

	return conf.LP, nil
}

func (d *DAL) SaveLastPage(ctx context.Context, lp int64) error {
	var conf config

	if err := d.db.Collection(LastPage).FindOneAndUpdate(ctx, bson.M{
		"_id": "lasstPage",
	}, bson.M{
		"$set": bson.M{
			"lp": lp,
		},
	}, options.FindOneAndUpdate().SetUpsert(true)).Decode(&conf); err != nil {
		return err
	}
	return nil
}

func (d *DAL) PutRequest(ctx context.Context, data, street string, square int64) (int64, error) {
	d.reqInd++
	_, err := d.db.Collection(CRColl).InsertOne(ctx, &model.CallReq{
		ID:           d.reqInd,
		Data:         data,
		CreationTime: time.Now().Unix(),
		Street:       street,
		Square:       square,
	})
	return d.reqInd, err
}

func (d *DAL) UpdateRecall(ctx context.Context, id, price0, price1 int64) (*model.CallReq, error) {
	var res model.CallReq
	if err := d.db.Collection(CRColl).FindOneAndUpdate(ctx, bson.M{
		"_id": id,
	}, bson.M{
		"$set": bson.M{
			"Response1": price0,
			"Response2": price1,
		},
	}).Decode(&res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (d *DAL) GetRequest(ctx context.Context, id int64) (*model.CallReq, error) {
	var res model.CallReq

	if err := d.db.Collection(CRColl).FindOne(ctx, bson.M{"_id": id}).Decode(&res); err != nil {
		fmt.Println("errrrr", err)
		return nil, err
	}

	return &res, nil

}

func (d *DAL) GetAverage(ctx context.Context, square int64, city, addr string) (int64, error) {
	type resp struct {
		ID  int32   `bson:"_id"`
		Avg float64 `bson:"average"`
	}

	var r []resp

	// valuableField := ""
	// if square == 0 {
	// 	valuableField = "$Price"
	// } else {
	// 	valuableField = "$PerSquare"
	// }

	//fmt.Println()

	splitted := strings.Split(addr, " ")

	filtered := []string{}

	for _, word := range splitted {
		if len(word) > 4 {
			filtered = append(filtered, word)
		}
	}

	matchQ := bson.M{
		"$project": bson.M{
			"_id":       "$_id",
			"PerSquare": "$PerSquare",
			"CntSrch": bson.M{
				"$size": bson.M{
					"$setIntersection": bson.A{"$SplittedAddress", splitted},
				},
			},
		},
	}

	cur, err := d.db.Collection(HouseColl).Aggregate(ctx, []bson.M{
		matchQ,
		bson.M{
			"$group": bson.M{
				"_id": "$CntSrch",
				"average": bson.M{
					"$avg": "$PerSquare",
				},
			},
		},
		bson.M{
			"$sort": bson.M{
				"_id": -1,
			},
		},
	})

	if err != nil {
		return 0, err
	}

	if err = cur.All(ctx, &r); err != nil {
		return 0, err
	}

	if r[0].Avg > 0 {

		return int64(r[0].Avg), nil
	}

	return 0, nil
}

func (d *DAL) GetCategories(ctx context.Context) ([]string, error) {
	type ret struct {
		data string `bson:"_id"`
	}

	var result []ret

	cur, err := d.db.Collection(HouseColl).Aggregate(ctx, mongo.Pipeline{
		bson.D{{
			"$project", bson.D{{
				"_id", bson.D{{
					"$concat", bson.A{
						"$CategoryName", " ", "$PricePostfix", " ", bson.D{{
							"$arrayElemAt", bson.A{
								bson.D{{
									"$split", bson.A{"$Title", " "},
								}}, 0,
							},
						}},
					},
				}},
			}},
		}},
		bson.D{{
			"$group", bson.D{{
				"_id", "$_id",
			}},
		}},
	})
	if err != nil {
		return nil, err
	}

	if err := cur.All(ctx, &result); err != nil {
		return nil, err
	}

	var res []string

	for _, item := range result {
		res = append(res, item.data)
	}

	return res, nil
}

// func Dto2dao_House(in *model.ArticleItem) *model.Flat {
// 	var (
// 		square      float64
// 		squarePoint string
// 		err         error
// 	)

// 	splittedTitle := strings.Split(in.Title, " ")

// 	//fmt.Println("splitted ", splittedTitle)

// 	for i, word := range splittedTitle {
// 		if strings.Contains(word, "м²") {
// 			squarePoint = "м²"
// 		} else if strings.Contains(word, "сот") {
// 			squarePoint = "сот"
// 		} else {
// 			continue
// 		}
// 		square, err = strconv.ParseFloat(splittedTitle[i-1], 64)
// 		if err != nil {
// 			panic(err)
// 		}
// 	}

// 	return &model.Flat{
// 		ID:            in.ID,
// 		Price:         in.PriceDetailed.Value,
// 		LocationID:    in.Location.ID,
// 		LocationName:  in.Location.Name,
// 		CategoryID:    in.Category.ID,
// 		CategoryName:  in.Category.Name,
// 		PointName:     squarePoint,
// 		Square:        int64(square),
// 		PerSquare:     int64(in.PriceDetailed.Value / int64(square)),
// 		Title:         in.Title,
// 		SortTimeStamp: in.SortTimestamp,
// 		PricePostfix:  in.PriceDetailed.Postfix,
// 		AddressDetailed: model.AddressDetailed{
// 			Metro:    in.AddressDetailed.MetroName,
// 			District: in.AddressDetailed.DistrictName,
// 			Road: in.AddressDetailed.RoadName + " " + in.GeoForItems.FormattedAddress + " " + func() string {
// 				if len(in.GeoForItems.GeoReferences) > 0 {
// 					return in.GeoForItems.GeoReferences[0].Content
// 				}
// 				return ""
// 			}(),
// 			Name: in.AddressDetailed.LocationName,
// 		},
// 	}
// }
