package server

import (
	"context"
	"math/rand"
	"strings"

	"github.com/117503445/eventize/src/be/ent"
	"github.com/117503445/eventize/src/be/internal/common"
	"github.com/117503445/eventize/src/be/internal/rpc"
	"github.com/rs/zerolog/log"
	"github.com/twitchtv/twirp"
	"google.golang.org/protobuf/types/known/emptypb"

	"database/sql"

	_ "github.com/lib/pq"
)

// Server implements the Haberdasher service
type Server struct {
}

func (s *Server) MakeHat(ctx context.Context, size *rpc.Size) (hat *rpc.Hat, err error) {
	if size.Inches <= 0 {
		return nil, twirp.InvalidArgumentError("inches", "I can't make a hat that small!")
	}
	return &rpc.Hat{
		Inches: size.Inches,
		Color:  []string{"white", "black", "brown", "red", "blue"}[rand.Intn(5)],
		Name:   []string{"bowler", "baseball cap", "top hat", "derby"}[rand.Intn(4)],
	}, nil
}

func (s *Server) GetBuildInfo(context.Context, *emptypb.Empty) (meta *rpc.BuildInfo, err error) {
	return &rpc.BuildInfo{
		Version: common.GetBuildInfo().Version,
	}, nil
}

type DBManager struct {
	client *ent.Client
	db     *sql.DB
}

func NewDBManager() *DBManager {
	client, err := ent.Open("postgres", "host=postgres port=5432 user=postgres password=12345678 dbname=dev sslmode=disable")
	if err != nil {
		log.Fatal().Err(err).Msg("failed opening ent connection to postgres")
	}
	db, err := sql.Open("postgres", "host=postgres port=5432 user=postgres password=12345678 sslmode=disable")
	if err != nil {
		log.Fatal().Err(err).Msg("failed opening sql connection to postgres")
	}
	return &DBManager{
		client: client,
		db:     db,
	}
}

func (m *DBManager) Close() {
	if err := m.client.Close(); err != nil {
		log.Fatal().Err(err).Msg("failed closing connection to postgres")
	}
}

// CreateNewDB 创建新的数据库
func (m *DBManager) CreateNewDB() {
	isExist := false

	// 查看 dev 数据库是否存在
	_, err := m.db.Query("SELECT 1 FROM pg_database WHERE datname='dev'")
	if err != nil {
		if strings.Contains(err.Error(), "does not exist") {
			isExist = false
		} else {
			log.Fatal().Err(err).Msg("failed to check if dev database exists")
		}
	} else {
		isExist = true
	}

	if !isExist {
		_, err := m.db.Exec("CREATE DATABASE dev")
		if err != nil {
			log.Fatal().Err(err).Msg("failed to create dev database")
		}
	}
}

func (m *DBManager) CreateOrMigrationSchema() {
	if err := m.client.Schema.Create(context.Background()); err != nil {
		log.Fatal().Err(err).Msg("failed creating schema resources")
	}
}
