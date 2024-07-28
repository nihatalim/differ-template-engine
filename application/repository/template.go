package repository

import (
	"context"
	"differ-template-engine/application/customerror"
	"differ-template-engine/application/domain"
	"differ-template-engine/log"
	"differ-template-engine/pkg/config"
	"fmt"
	"github.com/jackc/pgx/v4"
	"time"
)

type TemplateRepository interface {
	SaveTemplate(userId string, template domain.Template) (*domain.Template, error)
	DeleteTemplate(userId string, templateId int64) error
	GetUserTemplate(userId string) ([]domain.Template, error)
	GetTemplateByUserIdAndTemplateId(userId string, templateId int64) (*domain.Template, error)
}

type templateRepository struct {
	cfg    *config.ApplicationConfig
	conn   *pgx.Conn
	logger log.Logger
}

var InsertQuery = "insert into request_templates (user_id, name, content) values ($1, $2, $3) returning id, created_date"
var DeleteQuery = "delete from request_templates where user_id = $1 and id = $2"
var SelectAllTemplatesQuery = "select id, user_id, name, content, created_date from request_templates where user_id = $1"
var SelectTemplateByIdQuery = "select id, user_id, name, content, created_date from request_templates where user_id = $1 and id = $2"

func NewTemplateRepository(cfg *config.ApplicationConfig, conn *pgx.Conn, logger log.Logger) TemplateRepository {
	return &templateRepository{cfg: cfg, conn: conn, logger: logger}
}

func (t *templateRepository) SaveTemplate(userId string, template domain.Template) (*domain.Template, error) {
	var id int64
	var createdDate time.Time
	if err := t.conn.QueryRow(context.Background(), InsertQuery, userId, template.Name, template.Content).Scan(&id, &createdDate); err != nil {
		return nil, err
	}

	template.Id = id
	template.UserId = userId
	template.CreatedDate = createdDate

	return &template, nil
}

func (t *templateRepository) DeleteTemplate(userId string, templateId int64) error {
	command, err := t.conn.Exec(context.Background(), DeleteQuery, userId, templateId)
	if err != nil {
		return err
	}

	if command.RowsAffected() == 0 {
		return customerror.ErrTemplateIsNotFound
	}

	return nil
}

func (t *templateRepository) GetUserTemplate(userId string) ([]domain.Template, error) {
	rows, err := t.conn.Query(context.Background(), SelectAllTemplatesQuery, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	templates := make([]domain.Template, 0)

	for rows.Next() {
		var template domain.Template
		if err := rows.Scan(&template.Id, &template.UserId, &template.Name, &template.Content, &template.CreatedDate); err != nil {
			// TODO: burada düzgün maplenemeyenler ignore edilecek, belki loglanabilir
			continue
		}

		templates = append(templates, template)
	}

	return templates, nil
}

func (t *templateRepository) GetTemplateByUserIdAndTemplateId(userId string, templateId int64) (*domain.Template, error) {
	rows, err := t.conn.Query(context.Background(), SelectTemplateByIdQuery, userId, templateId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var template domain.Template

	for rows.Next() {
		if err := rows.Scan(&template.Id, &template.UserId, &template.Name, &template.Content, &template.CreatedDate); err != nil {
			continue
		}
	}

	if template.Id != templateId {
		return nil, customerror.ErrTemplateIsNotFound
	}

	return &template, nil
}

func Initialize(cfg *config.ApplicationConfig, logger log.Logger) (*pgx.Conn, error) {
	// PostgreSQL bağlantı bilgilerinizi buraya girin
	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", cfg.Postgresql.Username, cfg.Postgresql.Password, cfg.Postgresql.Host, cfg.Postgresql.Port, cfg.Postgresql.Database)

	// Bağlantıyı oluşturun
	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		logger.Errorf("Unable to connect to database: %v\n", err)
	}

	// Tablo oluşturma request_templates
	_, err = conn.Exec(context.Background(), `
        CREATE TABLE IF NOT EXISTS request_templates (
            id SERIAL PRIMARY KEY,
            user_id VARCHAR NOT NULL,
            name VARCHAR NOT NULL,
            content jsonb NOT NULL,
            created_date TIMESTAMPTZ DEFAULT NOW()
        )
    `)
	if err != nil {
		logger.Errorf("Unable to create table request_templates: %v\n", err)
		return nil, err
	}

	_, err = conn.Exec(context.Background(), "CREATE INDEX IF NOT EXISTS idx_request_templates_user_id ON request_templates(user_id)")
	if err != nil {
		logger.Errorf("Unable to create index for user_id on request_templates table: %v\n", err)
		return nil, err
	}

	// Tablo oluşturma execution result
	_, err = conn.Exec(context.Background(), `
        CREATE TABLE IF NOT EXISTS execution_results (
            operation_id VARCHAR PRIMARY KEY,
            user_id VARCHAR NOT NULL,
            request jsonb NOT NULL,
            diff BOOLEAN DEFAULT true,
            created_date TIMESTAMPTZ DEFAULT NOW()
        )
    `)
	if err != nil {
		logger.Errorf("Unable to create table execution_results: %v\n", err)
		return nil, err
	}

	_, err = conn.Exec(context.Background(), "CREATE INDEX IF NOT EXISTS idx_execution_results_user_id ON execution_results(user_id)")
	if err != nil {
		logger.Errorf("Unable to create index for user_id on execution_results table: %v\n", err)
		return nil, err
	}

	return conn, nil
}
