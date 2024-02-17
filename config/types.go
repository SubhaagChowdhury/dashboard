package config

/* Confuguration Structure */
type Configurations struct {
	// app details
	APPLICATION_DETAILS struct {
		APPLICATION_NAME string `json:"APPLICATION_NAME"`
		HOST_IP          string `validate:"required" json:"HOST_IP"`
		HOST_PORT        int    `validate:"required" json:"HOST_PORT"`
	} `validate:"required" json:"APPLICATION_DETAILS"`

	// logger details
	LOGGER_CONFIG struct {
		FILE_PATH string `validate:"required" json:"FILE_PATH"`
		FILE_NAME string `validate:"required" json:"FILE_NAME"`
		LEVEL     string `validate:"required" json:"LEVEL"`
	} `validate:"required" json:"LOGGER_CONFIG_DETAILS"`

	// mysql details
	MYSQL_CONFIG_DETAILS struct {
		DATABASE struct {
			HOST     string `validate:"required" json:"HOST"`
			PORT     int    `validate:"required" json:"PORT"`
			USERNAME string `validate:"required" json:"USERNAME"`
			PASSWORD string `validate:"required" json:"PASSWORD"`
			DATABASE string `validate:"required" json:"DATABASE"`

			SSL_FLAG         bool   `json:"SSL_FLAG"`
			TLS_ONLY_CA_CERT bool   `json:"TLS_ONLY_CA_CERT"`
			CA_CERT          string `validate:"required" json:"CA_CERT"`
			CLIENT_CERT      string `validate:"required" json:"CLIENT_CERT"`
			CLIENT_KEY       string `validate:"required" json:"CLIENT_KEY"`
		} `json:"FUND_MANAGER"`

		BATCH_SIZE               int  `validate:"required" json:"BATCH_SIZE"`
		TLS_INSECURE_SKIP_VERIFY bool `json:"TLS_INSECURE_SKIP_VERIFY"`

		IDLE_CONNS   int `json:"IDLE_CONNS"`
		OPEN_CONNS   int `json:"OPEN_CONNS"`
		MAX_LIFETIME int `json:"MAX_LIFETIME"`
	} `validate:"required" json:"MYSQL_CONFIG_DETAILS"`

	// rabbit mq details
	RABBITMQ_CONFIG_DETAILS struct {
		HOST           string `validate:"required" json:"HOST"`
		PORT           int    `validate:"required" json:"PORT"`
		EXCHANGE       string `validate:"required" json:"EXCHANGE"`
		VHOST          string `validate:"required" json:"VHOST"`
		DURABLE        bool   `validate:"required" json:"DURABLE"`
		USERNAME       string `validate:"required" json:"USERNAME"`
		PASSWORD       string `validate:"required" json:"PASSWORD"`
		ADD_TIMEOUT    int    `validate:"required" json:"ADD_TIMEOUT"`
		PREFETCH_COUNT int    `validate:"required" json:"PREFETCH_COUNT"`

		PUBLISH_QUEUE  string `validate:"required" json:"PUBLISH_QUEUE"`
		CONSUMER_QUEUE string `validate:"required" json:"CONSUMER_QUEUE"`

		CA_CERT                  string `validate:"required" json:"CA_CERT"`
		CLIENT_CERT              string `validate:"required" json:"CLIENT_CERT"`
		CLIENT_KEY               string `validate:"required" json:"CLIENT_KEY"`
		SSL_FLAG                 bool   `json:"SSL_FLAG"`
		TLS_INSECURE_SKIP_VERIFY bool   `json:"TLS_INSECURE_SKIP_VERIFY"`
		TLS_ONLY_CA_CERT         bool   `json:"TLS_ONLY_CA_CERT"`
	} `validate:"required" json:"RABBITMQ_CONFIG_DETAILS"`

	RETRY_COUNT int `validate:"required" json:"RETRY_COUNT"`
	DELAY_TIME  int `validate:"required" json:"DELAY_TIME"`

	RATE_LIMIT float64 `validate:"required" json:"RATE_LIMIT"`

	// flag objects
	DEV_ENV bool `json:"dev_env"`
}
