package config

var (
	DefaultMongoUri                  = "mongodb://0.0.0.0:27017/"
	DefaultMongoDatabase             = "recommender"
	DefaultMongoIndexesFile          = "indicies.json"
	DefaultMongoTermFile             = "terms.json"
	DefaultMongoSimilarityCollection = "similarity"
)

type MongoConfig interface {
	Uri() string
	Database() string
}

type mongoConfig struct {
	uri      string
	database string
}

func (mc *mongoConfig) Uri() string {
	return mc.uri
}

func (mc *mongoConfig) Database() string {
	return mc.database
}

func NewMongoConfig() *mongoConfig {
	mongoConfig := &mongoConfig{
		uri:      LookupEnvVariableOrDefault("MONGO_URI", DefaultMongoUri),
		database: LookupEnvVariableOrDefault("MONGO_DB", DefaultMongoDatabase),
	}

	return mongoConfig
}
