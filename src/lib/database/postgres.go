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

func (p *Postgres) GetConnectionSnapshotsData(ctx context.Context) ([]*model.ConnectionsData, error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	query := fmt.Sprintf(`select 
			imei, 
			connected_at_ms, 
			disconnected_at_ms, 
			duration, 
			reason, 
			sent, 
			recv, 
			action 
		from %s ORDER BY id DESC LIMIT 25`,
		config.Config.Postgres.ConnectionSnapshotTable)
	rows, err := p.conn.Query(ctxWithTimeout, query)
	if err != nil {
		p.telem.Errorf(ctx, "Failed to retrieve telematics rows. Error: %v", err)
		return nil, err
	}
	defer rows.Close()
	records := make([]*model.ConnectionsData, 0)
	for rows.Next() {
		var record model.ConnectionsData
		if err := rows.Scan(
			&record.Imei,
			&record.ConnectedAt,
			&record.DisconnectedAt,
			&record.Duration,
			&record.Reason,
			&record.Sent,
			&record.Recv,
			&record.Action,
		); err != nil {
			p.telem.Errorf(ctx, "Failed to scan telematics row. Error: %v", err)
			continue
		}
		records = append(records, &record)
	}
	return records, nil
}

func (p *Postgres) GetEntities(ctx context.Context) ([]*model.ConnectionsData, error) {
	return nil, nil
}

func (p *Postgres) GetRegistereddevices(ctx context.Context) ([]*model.RegisteredDevice, error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	query := fmt.Sprintf(`select 
			imei, 
			tenant_group_id, 
			tenant_id, 
			parser_id, 
			status, 
		from %s ORDER BY id DESC LIMIT 25`,
		config.Config.Postgres.RegisteredDevicesTable)
	rows, err := p.conn.Query(ctxWithTimeout, query)
	if err != nil {
		p.telem.Errorf(ctx, "Failed to retrieve registered devices. Error: %v", err)
		return nil, err
	}
	defer rows.Close()
	records := make([]*model.RegisteredDevice, 0)
	for rows.Next() {
		var record model.RegisteredDevice
		if err := rows.Scan(
			&record.Imei,
			&record.TenantGroupID,
			&record.TenantID,
			&record.ParserID,
			&record.Status,
		); err != nil {
			p.telem.Errorf(ctx, "Failed to scan registered devices. Error: %v", err)
			continue
		}
		records = append(records, &record)
	}
	return records, nil
}

func (p *Postgres) GetCommands(ctx context.Context) ([]*model.Command, error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	query := fmt.Sprintf(`select 
			id,
			imei,
			tenant_group_id, 
			tenant_id, 
			data,
			response,
			is_response_required,
			max_retries,
			retries_count,
			expires_at_ms,
			sent_to_device,
			sent_at_ms,
			response_at_ms,
		from %s ORDER BY id DESC LIMIT 25`,
		config.Config.Postgres.CommandsTable)
	rows, err := p.conn.Query(ctxWithTimeout, query)
	if err != nil {
		p.telem.Errorf(ctx, "Failed to retrieve commands. Error: %v", err)
		return nil, err
	}
	defer rows.Close()
	records := make([]*model.Command, 0)
	for rows.Next() {
		var record model.Command
		if err := rows.Scan(
			&record.ID,
			&record.Imei,
			&record.TenantGroupID,
			&record.TenantID,
			&record.Data,
			&record.Response,
			&record.IsResponseRequired,
			&record.MaxRetries,
			&record.RetriesCount,
			&record.ExpiresAtMs,
			&record.SentToDevice,
			&record.SentAtMs,
			&record.ResponseAtMs,
		); err != nil {
			p.telem.Errorf(ctx, "Failed to scan commands. Error: %v", err)
			continue
		}
		records = append(records, &record)
	}
	return records, nil
}

func (p *Postgres) Close(ctx context.Context) error {
	return p.conn.Close(ctx)
}
