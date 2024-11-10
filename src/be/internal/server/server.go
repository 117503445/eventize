package server

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/117503445/eventize/src/be/ent"
	"github.com/117503445/eventize/src/be/internal/common"
	"github.com/117503445/eventize/src/be/internal/rpc"
	"github.com/rs/zerolog/log"
	"github.com/twitchtv/twirp"
	"google.golang.org/protobuf/types/known/emptypb"

	"database/sql"

	"google.golang.org/protobuf/types/known/anypb"

	_ "github.com/lib/pq"
)

// Server implements the Haberdasher service
type Server struct {
	dbManager *DBManager
}

func NewServer() *Server {
	dbManager := NewDBManager()
	dbManager.CreateNewDB("dev", false)
	dbManager.CreateOrMigrationSchema()
	return &Server{
		dbManager: dbManager,
	}
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

func convertMap(input map[string]*anypb.Any) (map[string]interface{}, error) {
	output := make(map[string]interface{})

	for key, anyValue := range input {
		if anyValue == nil {
			continue
		}

		// 解析 *anypb.Any
		value, err := anyValue.UnmarshalNew()
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal Any value for key %s: %w", key, err)
		}

		// 将解析后的值转换为 interface{}
		output[key] = value.ProtoReflect().Interface()
	}

	return output, nil
}

func (s *Server) CreateEvent(ctx context.Context, request *rpc.CreateEventRequest) (response *rpc.CreateEventResponse, err error) {
	data, err := convertMap(request.Data)
	if err != nil {
		return nil, twirp.InvalidArgumentError("data", err.Error())
	}

	event, err := s.dbManager.client.Event.Create().SetData(data).SetType(request.Type).Save(ctx)
	if err != nil {
		return nil, twirp.InternalErrorWith(fmt.Errorf("failed to create event: %w", err))
	}

	log.Info().Int("id", event.ID).Msg("event created")

	return &rpc.CreateEventResponse{
		Id: strconv.Itoa(event.ID),
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

	// Ping the database to check if the connection is working
	err = common.NewRetry().Execute(func() error {
		if err := db.Ping(); err != nil {
			log.Warn().Err(err).Msg("failed to ping database")
			return err
		}
		return nil
	})

	if err != nil {
		log.Fatal().Msg("failed to ping database")
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
func (m *DBManager) CreateNewDB(dbName string, force bool) {
	isExist := false

	// 查看 dev 数据库是否存在
	rows, err := m.db.Query("SELECT 1 FROM pg_database WHERE datname=$1", dbName)
	if err != nil {
		if strings.Contains(err.Error(), "does not exist") {
			isExist = false
		} else {
			log.Fatal().Err(err).Msg("failed to check if dev database exists")
		}
	} else {
		for rows.Next() {
			isExist = true
			break
		}
	}

	if force && isExist {
		_, err := m.db.Exec(fmt.Sprintf("DROP DATABASE %s", identifier(dbName)))
		if err != nil {
			log.Fatal().Err(err).Msg("failed to drop dev database")
		}
		isExist = false
	}

	if !isExist {
		_, err := m.db.Exec(fmt.Sprintf("CREATE DATABASE %s", identifier(dbName)))
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

// TODO: 更好的处理标识符
// sql.Identifier 是一个辅助函数，用于正确地引用标识符（如表名、列名等）
func identifier(s string) string {
	return fmt.Sprintf(`"%s"`, strings.ReplaceAll(s, `"`, `""`))
}
