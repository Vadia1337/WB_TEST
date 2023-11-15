package config

import "time"

type Config struct {
	HttpServerPort     string
	JWTSalt            string
	AccessTokenExpire  time.Duration // мин
	RefreshTokenExpire int           // сек
	PostgresConfig
	GameRules
}

type GameRules struct {
	CustomerCapitalMin      int
	CustomerCapitalMax      int
	LoaderMinPortableWeight int
	LoaderMaxPortableWeight int
	LoaderSalaryMin         int
	LoaderSalaryMax         int
	LoaderFatigue           int
	CountJobsToGenerate     int
	CargoWeightMin          int
	CargoWeightMax          int
}

type PostgresConfig struct {
	PostgresqlHost     string
	PostgresqlPort     string
	PostgresqlUser     string
	PostgresqlPassword string
	PostgresqlDbname   string
	PgDriver           string
}

// можно было использовать viper и сделать файл yaml, но руки не дошли
func GetConfig() *Config {
	pgConfig := PostgresConfig{
		PostgresqlHost:     "localhost",
		PostgresqlPort:     "5433",
		PostgresqlUser:     "PgUser",
		PostgresqlPassword: "PgUser",
		PostgresqlDbname:   "game",
		PgDriver:           "pgx",
	}

	Rules := GameRules{
		CustomerCapitalMin:      10000,
		CustomerCapitalMax:      100000,
		LoaderMinPortableWeight: 5,
		LoaderMaxPortableWeight: 30,
		LoaderSalaryMin:         10000,
		LoaderSalaryMax:         30000,
		LoaderFatigue:           20,
		CountJobsToGenerate:     10,
		CargoWeightMin:          10,
		CargoWeightMax:          80,
	}

	return &Config{
		HttpServerPort:     ":8080",
		JWTSalt:            "adwakndaoidia344r3hussigfse8gfse89",
		AccessTokenExpire:  15,
		RefreshTokenExpire: 5184000,
		PostgresConfig:     pgConfig,
		GameRules:          Rules,
	}
}
