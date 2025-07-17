package repositories

import (
	"context"
	"time"

	"center/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoMetricRepository реализует интерфейс MetricRepository для MongoDB
type MongoMetricRepository struct {
	db *mongo.Database
}

func NewMongoMetricRepository(db *mongo.Database) *MongoMetricRepository {
	return &MongoMetricRepository{db: db}
}

func (r *MongoMetricRepository) SaveSystemMetrics(ctx context.Context, metrics *models.SystemMetrics) error {
	collection := r.db.Collection("system_metrics")
	_, err := collection.InsertOne(ctx, metrics)
	return err
}

func (r *MongoMetricRepository) SaveProcessMetrics(ctx context.Context, metrics *models.ProcessMetrics) error {
	collection := r.db.Collection("process_metrics")
	_, err := collection.InsertOne(ctx, metrics)
	return err
}

func (r *MongoMetricRepository) SaveContainerMetrics(ctx context.Context, metrics *models.ContainerMetrics) error {
	collection := r.db.Collection("container_metrics")
	_, err := collection.InsertOne(ctx, metrics)
	return err
}

func (r *MongoMetricRepository) SaveNetworkMetrics(ctx context.Context, metrics *models.NetworkMetrics) error {
	collection := r.db.Collection("network_metrics")
	_, err := collection.InsertOne(ctx, metrics)
	return err
}

func (r *MongoMetricRepository) GetLastSystemMetrics(ctx context.Context, hostID int) (*models.SystemMetrics, error) {
	collection := r.db.Collection("system_metrics")
	filter := bson.M{"host_id": hostID}
	opts := options.FindOne().SetSort(bson.D{{Key: "timestamp", Value: -1}})

	var metrics models.SystemMetrics
	err := collection.FindOne(ctx, filter, opts).Decode(&metrics)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &metrics, err
}

func (r *MongoMetricRepository) GetSystemMetricsInRange(ctx context.Context, hostID int, from, to time.Time) ([]models.SystemMetrics, error) {
	collection := r.db.Collection("system_metrics")
	filter := bson.M{
		"host_id": hostID,
		"timestamp": bson.M{
			"$gte": from,
			"$lte": to,
		},
	}
	opts := options.Find().SetSort(bson.D{{Key: "timestamp", Value: 1}})

	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var metrics []models.SystemMetrics
	if err := cursor.All(ctx, &metrics); err != nil {
		return nil, err
	}

	return metrics, nil
}

func (r *MongoMetricRepository) GetProcessMetricsInRange(ctx context.Context, hostID int, from, to time.Time) ([]models.ProcessMetrics, error) {
	collection := r.db.Collection("process_metrics")
	filter := bson.M{
		"host_id": hostID,
		"timestamp": bson.M{
			"$gte": from,
			"$lte": to,
		},
	}
	opts := options.Find().SetSort(bson.D{{Key: "timestamp", Value: 1}})

	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var metrics []models.ProcessMetrics
	if err := cursor.All(ctx, &metrics); err != nil {
		return nil, err
	}

	return metrics, nil
}

func (r *MongoMetricRepository) GetContainerMetricsInRange(ctx context.Context, hostID int, from, to time.Time) ([]models.ContainerMetrics, error) {
	collection := r.db.Collection("container_metrics")
	filter := bson.M{
		"host_id": hostID,
		"timestamp": bson.M{
			"$gte": from,
			"$lte": to,
		},
	}
	opts := options.Find().SetSort(bson.D{{Key: "timestamp", Value: 1}})

	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var metrics []models.ContainerMetrics
	if err := cursor.All(ctx, &metrics); err != nil {
		return nil, err
	}

	return metrics, nil
}

func (r *MongoMetricRepository) GetNetworkMetricsInRange(ctx context.Context, hostID int, from, to time.Time) ([]models.NetworkMetrics, error) {
	collection := r.db.Collection("network_metrics")
	filter := bson.M{
		"host_id": hostID,
		"timestamp": bson.M{
			"$gte": from,
			"$lte": to,
		},
	}
	opts := options.Find().SetSort(bson.D{{Key: "timestamp", Value: 1}})

	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var metrics []models.NetworkMetrics
	if err := cursor.All(ctx, &metrics); err != nil {
		return nil, err
	}

	return metrics, nil
}

func (r *MongoMetricRepository) CleanupOldMetrics(ctx context.Context, collectionName string, threshold time.Time) error {
	collection := r.db.Collection(collectionName)
	filter := bson.M{"timestamp": bson.M{"$lt": threshold}}
	_, err := collection.DeleteMany(ctx, filter)
	return err
}

func (r *MongoMetricRepository) SetupTTLIndex(ctx context.Context, collectionName string, ttlSeconds int32) error {
	collection := r.db.Collection(collectionName)
	model := mongo.IndexModel{
		Keys:    bson.M{"timestamp": 1},
		Options: options.Index().SetExpireAfterSeconds(ttlSeconds),
	}

	_, err := collection.Indexes().CreateOne(ctx, model)
	return err
}

func (r *MongoMetricRepository) Ping(ctx context.Context) error {
	return r.db.Client().Ping(ctx, nil)
}
