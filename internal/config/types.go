package config

import "go.uber.org/zap/zapcore"

type Config struct {
	Init       Init       `mapstructure:"init"`
	Management Management `mapstructure:"management"`
	Auth       Auth       `mapstructure:"auth"`
	Client     Client     `mapstructure:"client"`
	Repository Repository `mapstructure:"repository"`
	OAuth      OAuth      `mapstructure:"oauth"`
	Cluster    Cluster    `mapstructure:"cluster"`
	Email      Email      `mapstructure:"email"`
	Registry   Registry   `mapstructure:"registry"`
}

type Init struct {
	Root InitRoot `mapstructure:"root"`
}

type InitRoot struct {
	Email    string `mapstructure:"email"`
	Password string `mapstructure:"password"`
}

type Management struct {
	Port      int                 `mapstructure:"port"`
	PublicURL string              `mapstructure:"public_url"`
	Telemetry ManagementTelemetry `mapstructure:"telemetry"`
}

type ManagementTelemetry struct {
	Enabled bool `mapstructure:"enabled"`
}

type Auth struct {
	JWT AuthJWT `mapstructure:"jwt"`
}

type AuthJWT struct {
	Secret             string `mapstructure:"secret"`
	CertificateFile    string `mapstructure:"certificate_file"`
	CertificateKeyFile string `mapstructure:"certificate_key_file"`
}

type Client struct {
	Postgres ClientPostgres `mapstructure:"postgres"`
	Mongo    ClientMongo    `mapstructure:"mongo"`
	Minio    ClientMinio    `mapstructure:"minio"`
	Docker   ClientDocker   `mapstructure:"docker"`
	Mailjet  ClientMailjet  `mapstructure:"mailjet"`
	SMTP     ClientSMTP     `mapstructure:"smtp"`
}

type ClientSMTP struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

type ClientPostgres struct {
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Host     string `mapstructure:"host"`
}

type ClientMongo struct {
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Host     string `mapstructure:"host"`
}

type ClientMinio struct {
	AccessKeyID     string `mapstructure:"access_key_id"`
	SecretAccessKey string `mapstructure:"secret_access_key"`
	Host            string `mapstructure:"host"`
	Secure          bool   `mapstructure:"secure"`
}

type ClientDocker struct {
	Host string `mapstructure:"host"`
}

type ClientMailjet struct {
	From      string `mapstructure:"from"`
	APIKey    string `mapstructure:"api_key"`
	SecretKey string `mapstructure:"secret_key"`
}

type Repository struct {
	Storage          RepositoryStore       `mapstructure:"storage"`
	Capsule          RepositoryStore       `mapstructure:"capsule"`
	Database         RepositoryStore       `mapstructure:"database"`
	ServiceAccount   RepositoryStore       `mapstructure:"service_account"`
	Group            RepositoryStore       `mapstructure:"group"`
	Project          RepositoryStore       `mapstructure:"project"`
	Secret           RepositoryStoreSecret `mapstructure:"secret"`
	Session          RepositoryStore       `mapstructure:"session"`
	User             RepositoryStore       `mapstructure:"user"`
	VerificationCode RepositoryStore       `mapstructure:"verification_code"`
}

type RepositoryStore struct {
	Store string `mapstructure:"store"`
}

type RepositoryStoreSecret struct {
	Store   string                       `mapstructure:"store"`
	MongoDB RepositoryStoreSecretMongoDB `mapstructure:"mongodb"`
}

type RepositoryStoreSecretMongoDB struct {
	Key string `mapstructure:"key"`
}

type OAuth struct {
	Google   OAuthClientCredentials `mapstructure:"google"`
	Github   OAuthClientCredentials `mapstructure:"github"`
	Facebook OAuthClientCredentials `mapstructure:"facebook"`
}

type OAuthClientCredentials struct {
	ClientID     string `mapstructure:"client_id"`
	ClientSecret string `mapstructure:"client_secret"`
}

type Cluster struct {
	Type ClusterType `mapstructure:"type"`
}

type ClusterType string

const (
	ClusterTypeDocker     ClusterType = "docker"
	ClusterTypeKubernetes ClusterType = "k8s"
)

type Email struct {
	From string `mapstructure:"from"`
	Type string `mapstructure:"type"`
}

type EmailType string

const (
	EmailTypeNoEmail = ""
	EmailTypeMailjet = "mailjet"
	EmailTypeSMTP    = "smtp"
)

type Registry struct {
	Enabled  bool          `mapstructure:"enabled"`
	Port     int           `mapstructure:"port"`
	LogLevel zapcore.Level `mapstructure:"log_level"`
}