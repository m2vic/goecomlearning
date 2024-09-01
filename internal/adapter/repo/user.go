package repo

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"golearning/internal/core/domain"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

type MongoUserRepo struct {
	client *mongo.Client
	col    *mongo.Collection
}

func NewMongoRepo(client *mongo.Client, dbName string, colName string) *MongoUserRepo {
	if client == nil {
		fmt.Println("missing db")
		return nil
	}
	col := client.Database(dbName).Collection(colName)
	return &MongoUserRepo{client: client, col: col}
}
func (r *MongoUserRepo) GetUser(ctx context.Context, userId string) (*domain.User, error) {
	result := domain.User{}
	id, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return nil, fmt.Errorf("fail to convert to primitive:%w", err)
	}
	filter := bson.M{"_id": id}
	err = r.col.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("not found")
		} else {
			return nil, fmt.Errorf("from repo layer:%w", err)
		}
	}
	return &result, nil
}

func (r *MongoUserRepo) CheckUsername(ctx context.Context, username string) (*domain.User, error) {
	result := domain.User{}
	filter := bson.M{"username": username}
	err := r.col.FindOne(context.TODO(), filter).Decode(&result)
	if err == mongo.ErrNoDocuments {
		return nil, err
	}
	return &result, nil
}
func (r *MongoUserRepo) Register(ctx context.Context, username, email string, hash []byte) error {
	data := domain.User{}
	data.UserId = primitive.NewObjectID()
	data.Username = username
	data.Password = string(hash)
	data.Email = email
	if username == "admin" {
		data.Role = "admin"
	} else {
		data.Role = "user"
	}
	data.Cart = []domain.UsersProduct{}
	_, err := r.col.InsertOne(ctx, data)
	if err != nil {
		log.Fatal(err)

	}
	return nil
}

func (r *MongoUserRepo) UpdateUser(ctx context.Context, info domain.User) error {

	filter := bson.M{"_id": info.UserId}
	update := bson.M{"$set": bson.M{"firstname": info.FirstName,
		"lastname": info.LastName, "role": info.Role, "address": info.AddressDetails}}

	_, err := r.col.UpdateOne(ctx, filter, update)

	if err != nil {
		log.Fatal(err)
	}
	return nil
}
func (r *MongoUserRepo) ChangePassword(ctx context.Context, userId, oldPass, newPass string) error {
	id, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		log.Fatal(err)
	}
	result := domain.User{}
	filter := bson.M{"_id": id}
	err = r.col.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return fmt.Errorf("from repo layer:%w", err)
		} else {
			log.Fatal(err)
		}
	}
	err = bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(oldPass))
	if err != nil {
		return fmt.Errorf("password invalid from repo layer:%w", err)
	}
	password, err := bcrypt.GenerateFromPassword([]byte(newPass), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println("afterGen:", password)
	update := bson.M{"$set": bson.M{"password": string(password)}}
	_, err = r.col.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
func (r *MongoUserRepo) CheckEmail(ctx context.Context, email string) (bool, error) {
	result := domain.User{}
	filter := bson.M{"email": email}
	err := r.col.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		} else {
			log.Fatal(err)
		}
	}
	return true, nil
}
func (r *MongoUserRepo) ResetPassword(ctx context.Context, email string) (string, error) {
	filter := bson.M{"email": email}
	// Generate Password
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()"
	charsetLen := len(charset)
	password := make([]byte, 12)
	if _, err := rand.Read(password); err != nil {
		log.Fatal(err)
	}
	for i := range password {
		password[i] = charset[int(password[i])%charsetLen]
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	hashPassword := string(hash)
	update := bson.M{"$set": bson.M{"password": hashPassword}}
	_, err = r.col.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Fatal(err)
	}
	return string(password), nil
}
func (r *MongoUserRepo) CheckRefresh(ctx context.Context, Token string) (*domain.User, error) {
	filter := bson.M{"token": Token}
	result := domain.User{}
	err := r.col.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("checkrefresh from repo layer:%w", err)
		} else {
			log.Fatal(err)
		}
	}
	return &result, nil
}
func (r *MongoUserRepo) GetCart(ctx context.Context, userId string) ([]domain.Cart, error) {
	var result struct {
		Cart []domain.Cart `bson:"cart"`
	}
	id, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return nil, fmt.Errorf("from repo:%w", err)
	}
	filter := bson.M{"_id": id}
	projection := bson.M{"cart": 1}
	err = r.col.FindOne(ctx, filter, options.FindOne().SetProjection(projection)).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("databse:%w", err)
	}
	return result.Cart, nil

}
func (r *MongoUserRepo) AddtoCart(ctx context.Context, userProduct domain.Cart, userId string) error {
	uid, _ := primitive.ObjectIDFromHex(userId)
	// Try to add the product to the cart if it doesn't already exist
	filter := bson.M{"_id": uid, "cart.productid": bson.M{"$ne": userProduct.ProductId}}
	update := bson.M{"$push": bson.M{"cart": userProduct}}
	result, err := r.col.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	// If no document was modified, it means the product is already in the cart
	if result.ModifiedCount == 0 {
		filter = bson.M{"_id": uid, "cart.productid": userProduct.ProductId}
		update = bson.M{"$inc": bson.M{"cart.$.amount": userProduct.Amount, "cart.$.totalprice": userProduct.PricePerPiece}}
		_, err = r.col.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			return err
		}
	}
	return nil
}
func (r *MongoUserRepo) DeleteItemInCart(ctx context.Context, userId string, productId primitive.ObjectID) error {
	id, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		log.Fatal(err)
	}
	filter := bson.M{"_id": id}
	update := bson.M{"$pull": bson.M{"cart": bson.M{"productid": productId}}}
	_, err = r.col.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}
func (r *MongoUserRepo) DeleteItemFromSystem(ctx context.Context, productId primitive.ObjectID) error {
	// query first, if there is product in user cart , update it and if not return nil
	filter := bson.M{"cart.productid": productId}
	count, err := r.col.CountDocuments(ctx, filter)
	if err != nil {
		return err
	}
	if count == 0 {
		return nil
	}
	update := bson.M{"$pull": bson.M{"cart": bson.M{"productid": productId}}}
	_, err = r.col.UpdateMany(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("fail to delete item from user cart:%w", err)
	}
	return nil
}
func (r *MongoUserRepo) EditItemFromSystem(ctx context.Context, product domain.Product) error {
	filter := bson.M{"cart.productid": product.ProductID}
	update := bson.M{"$set": bson.M{"productname": product.ProductName, "priceeach": product.Price, "priceid": product.PriceId, "details": product.Details}}
	_, err := r.col.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("fail to update from user cart:%w", err)
	}
	return nil
}
func (r *MongoUserRepo) IncreaseCartProduct(ctx context.Context, userId string, productId primitive.ObjectID) error {
	id, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": id, "cart.productid": productId}
	update := bson.M{"$inc": bson.M{"cart.$.amount": 1}}
	_, err = r.col.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}
func (r *MongoUserRepo) DecreaseCartProduct(ctx context.Context, userId string, productId primitive.ObjectID) error {
	id, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": id, "cart.productid": productId}
	update := bson.M{"$inc": bson.M{"cart.$.amount": -1}}
	_, err = r.col.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}
func (r *MongoUserRepo) ClearCart(ctx context.Context, userId string) error {
	id, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return fmt.Errorf("fail to parse primitiveId:%w", err)
	}
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"cart": bson.A{}}}
	_, err = r.col.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("from repo:%w", err)
	}
	return nil

}
