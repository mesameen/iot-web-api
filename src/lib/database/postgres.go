package database

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/mesameen/iot-web-api/src/config"
	"github.com/mesameen/iot-web-api/src/model"
	"github.com/mesameen/iot-web-api/src/telemetryservice"
)

type Postgres struct {
	telem *telemetryservice.Service
	conn  *pgx.Conn
}

func connectPostgres(ctx context.Context, telem *telemetryservice.Service) (*Postgres, error) {
	conn, err := pgx.Connect(ctx, config.Config.Postgres.Address)
	if err != nil {
		telem.Errorf(ctx, "Failed to connect to postgres. Error: %v", err)
		return nil, err
	}
	if err := conn.Ping(ctx); err != nil {
		telem.Errorf(ctx, "Failed to ping to postgres. Error: %v", err)
		return nil, err
	}
	return &Postgres{
		telem,
		conn,
	}, nil
}

func (p *Postgres) GetTelematicsData(ctx context.Context) ([]*model.TelematicsData, error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	query := fmt.Sprintf("select payload from %s ORDER BY id DESC LIMIT 25", config.Config.Postgres.TelematicsDataTable)
	rows, err := p.conn.Query(ctxWithTimeout, query)
	if err != nil {
		p.telem.Errorf(ctx, "Failed to retrieve telematics rows. Error: %v", err)
		return nil, err
	}
	defer rows.Close()
	records := make([]*model.TelematicsData, 0)
	for rows.Next() {
		var payload []byte
		if err := rows.Scan(&payload); err != nil {
			p.telem.Errorf(ctx, "Failed to scan telematics row. Error: %v", err)
			continue
		}
		var record model.TelematicsData
		if err := json.Unmarshal(payload, &record); err != nil {
			p.telem.Errorf(ctx, "Failed to unmarshal telematics row. Error: %v", err)
			continue
		}
		records = append(records, &record)
	}
	return records, nil
}

func (p *Postgres) Close(ctx context.Context) error {
	return p.conn.Close(ctx)
}
