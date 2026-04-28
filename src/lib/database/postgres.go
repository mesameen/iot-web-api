package database

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
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

func (p *Postgres) GetTelematicsData(ctx context.Context, req *model.GetTelematicsDataRequest) ([]*model.TelematicsData, error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	// Dynamic where condition
	args := []any{}
	conditions := []string{}
	argIdx := 1
	if req.IMEI != "" {
		conditions = append(conditions, fmt.Sprintf("imei = $%d", argIdx))
		args = append(args, req.IMEI)
		argIdx++
	}
	if req.From > 0 {
		conditions = append(conditions, fmt.Sprintf("device_date_time >= to_timestamp($%d::bigint / 1000.0)", argIdx))
		args = append(args, req.From)
		argIdx++
	}
	if req.To > 0 {
		conditions = append(conditions, fmt.Sprintf("device_date_time <= to_timestamp($%d::bigint / 1000.0)", argIdx))
		args = append(args, req.To)
		argIdx++
	}

	where := ""
	if len(conditions) > 0 {
		where = "WHERE " + strings.Join(conditions, " AND ")
	}

	query := fmt.Sprintf(`select 
							imei,
							device_date_time,
							listener_date_time,
							gps_data,
							sensor_data,
							network_data
							from %s %s ORDER BY id DESC LIMIT 25`,
		config.Config.Postgres.TelematicsDataTable, where)
	rows, err := p.conn.Query(ctxWithTimeout, query, args...)
	if err != nil {
		p.telem.Errorf(ctx, "Failed to retrieve telematics rows. Error: %v", err)
		return nil, err
	}
	defer rows.Close()
	records := make([]*model.TelematicsData, 0)
	for rows.Next() {
		var imei string
		var deviceDateTime time.Time
		var listenerDateTime time.Time
		var gpsDataStr string
		var sensorDataStr string
		var networkDataStr string
		if err := rows.Scan(
			&imei,
			&deviceDateTime,
			&listenerDateTime,
			&gpsDataStr,
			&sensorDataStr,
			&networkDataStr,
		); err != nil {
			p.telem.Errorf(ctx, "Failed to scan telematics row. Error: %v", err)
			continue
		}
		var gpsData model.GpsData
		var sensorData model.SensorData
		var networkData model.NetworkData
		if err := json.Unmarshal([]byte(gpsDataStr), &gpsData); err != nil {
			p.telem.Errorf(ctx, "Failed to unmarshal gps data fetched from db. Error: %v", err)
		}
		if err := json.Unmarshal([]byte(sensorDataStr), &sensorData); err != nil {
			p.telem.Errorf(ctx, "Failed to unmarshal sensor data fetched from db. Error: %v", err)
		}
		if err := json.Unmarshal([]byte(networkDataStr), &networkData); err != nil {
			p.telem.Errorf(ctx, "Failed to unmarshal network data fetched from db. Error: %v", err)
		}
		record := model.TelematicsData{
			Imei:             imei,
			ListenerDatetime: uint64(listenerDateTime.UnixMilli()),
			DeviceDatetime:   uint64(deviceDateTime.UnixMilli()),
			GpsData:          &gpsData,
			SensorData:       &sensorData,
			NetworkData:      &networkData,
		}
		records = append(records, &record)
	}
	return records, nil
}

func (p *Postgres) GetRecentTelematicsData(ctx context.Context, req *model.GetTelematicsDataRequest) ([]*model.TelematicsData, error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	args := []any{}
	conditions := []string{}
	argIdx := 1
	if req.IMEI != "" {
		args = append(args, req.IMEI)
		conditions = append(conditions, fmt.Sprintf("imei=$%d", argIdx))
		argIdx++
	}
	where := ""
	if len(conditions) != 0 {
		where = "where " + strings.Join(conditions, " AND ")
	}
	query := fmt.Sprintf(`select 
							imei,
							device_date_time,
							listener_date_time,
							gps_data,
							sensor_data,
							network_data
							from %s %s ORDER BY updated_at DESC LIMIT 25`,
		config.Config.Postgres.RecentTelematicsDataTable, where)
	p.telem.Infof(ctx, "Query to fetch recent telematics data: %s", strings.Join(strings.Fields(query), " "))
	rows, err := p.conn.Query(ctxWithTimeout, query, args...)
	if err != nil {
		p.telem.Errorf(ctx, "Failed to retrieve recent telematics rows. Error: %v", err)
		return nil, err
	}
	defer rows.Close()
	records := make([]*model.TelematicsData, 0)
	for rows.Next() {
		var imei string
		var deviceDateTime time.Time
		var listenerDateTime time.Time
		var gpsDataStr string
		var sensorDataStr string
		var networkDataStr string
		if err := rows.Scan(
			&imei,
			&deviceDateTime,
			&listenerDateTime,
			&gpsDataStr,
			&sensorDataStr,
			&networkDataStr,
		); err != nil {
			p.telem.Errorf(ctx, "Failed to scan telematics row. Error: %v", err)
			continue
		}
		var gpsData model.GpsData
		var sensorData model.SensorData
		var networkData model.NetworkData
		if err := json.Unmarshal([]byte(gpsDataStr), &gpsData); err != nil {
			p.telem.Errorf(ctx, "Failed to unmarshal gps data fetched from db. Error: %v", err)
		}
		if err := json.Unmarshal([]byte(sensorDataStr), &sensorData); err != nil {
			p.telem.Errorf(ctx, "Failed to unmarshal sensor data fetched from db. Error: %v", err)
		}
		if err := json.Unmarshal([]byte(networkDataStr), &networkData); err != nil {
			p.telem.Errorf(ctx, "Failed to unmarshal network data fetched from db. Error: %v", err)
		}
		record := model.TelematicsData{
			Imei:             imei,
			ListenerDatetime: uint64(listenerDateTime.UnixMilli()),
			DeviceDatetime:   uint64(deviceDateTime.UnixMilli()),
			GpsData:          &gpsData,
			SensorData:       &sensorData,
			NetworkData:      &networkData,
		}
		records = append(records, &record)
	}
	return records, nil
}

func (p *Postgres) GetConnectionEvents(ctx context.Context, req *model.GetConnectionsDataRequest) ([]*model.ConnectionsData, error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	// Dynamic where condition
	args := []any{}
	conditions := []string{}
	argIdx := 1
	if req.IMEI != "" {
		conditions = append(conditions, fmt.Sprintf("imei = $%d", argIdx))
		args = append(args, req.IMEI)
		argIdx++
	}
	if req.From > 0 {
		conditions = append(conditions, fmt.Sprintf("connected_at >= to_timestamp($%d::bigint / 1000.0)", argIdx))
		args = append(args, req.From)
		argIdx++
	}
	if req.To > 0 {
		conditions = append(conditions, fmt.Sprintf("connected_at <= to_timestamp($%d::bigint / 1000.0)", argIdx))
		args = append(args, req.To)
		argIdx++
	}

	where := ""
	if len(conditions) > 0 {
		where = "WHERE " + strings.Join(conditions, " AND ")
	}
	query := fmt.Sprintf(`select 
			imei, 
			connected_at, 
			disconnected_at, 
			duration_ms, 
			reason,
			sent, 
			recv, 
			action 
		from %s %s ORDER BY updated_at DESC LIMIT 25`,
		config.Config.Postgres.ConnectionEventsTable, where)
	p.telem.Infof(ctx, "Query to fetch connection events data: %s", strings.Join(strings.Fields(query), " "))
	rows, err := p.conn.Query(ctxWithTimeout, query, args...)
	if err != nil {
		p.telem.Errorf(ctx, "Failed to retrieve connection events. Error: %v", err)
		return nil, err
	}
	defer rows.Close()
	records := make([]*model.ConnectionsData, 0)
	for rows.Next() {
		var record model.ConnectionsData
		var connectedAt time.Time
		var disconnectedAt *time.Time

		if err := rows.Scan(
			&record.Imei,
			&connectedAt,
			&disconnectedAt,
			&record.DurationMS,
			&record.Reason,
			&record.Sent,
			&record.Recv,
			&record.Action,
		); err != nil {
			p.telem.Errorf(ctx, "Failed to scan connection events row. Error: %v", err)
			continue
		}
		record.ConnectedAt = connectedAt.UnixMilli()
		if disconnectedAt != nil {
			record.DisconnectedAt = disconnectedAt.UnixMilli()
		}
		records = append(records, &record)
	}
	return records, nil
}

func (p *Postgres) GetRecentConnectionEvents(ctx context.Context, req *model.GetConnectionsDataRequest) ([]*model.ConnectionsData, error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	args := []any{}
	conditions := []string{}
	argIdx := 1
	if req.IMEI != "" {
		args = append(args, req.IMEI)
		conditions = append(conditions, fmt.Sprintf("imei=$%d", argIdx))
		argIdx++
	}
	where := ""
	if len(conditions) != 0 {
		where = "where " + strings.Join(conditions, " AND ")
	}
	query := fmt.Sprintf(`select 
			imei, 
			connected_at, 
			disconnected_at, 
			duration_ms, 
			reason,
			sent, 
			recv, 
			action 
		from %s %s ORDER BY updated_at DESC LIMIT 25`,
		config.Config.Postgres.RecentConnectionEventsTable, where)
	p.telem.Infof(ctx, "Query to fetch recent connections data: %s", strings.Join(strings.Fields(query), " "))

	rows, err := p.conn.Query(ctxWithTimeout, query, args...)
	if err != nil {
		p.telem.Errorf(ctx, "Failed to retrieve recent connection events. Error: %v", err)
		return nil, err
	}
	defer rows.Close()
	records := make([]*model.ConnectionsData, 0)
	for rows.Next() {
		var record model.ConnectionsData
		var connectedAt time.Time
		var disconnectedAt *time.Time

		if err := rows.Scan(
			&record.Imei,
			&connectedAt,
			&disconnectedAt,
			&record.DurationMS,
			&record.Reason,
			&record.Sent,
			&record.Recv,
			&record.Action,
		); err != nil {
			p.telem.Errorf(ctx, "Failed to scan recent connection events row. Error: %v", err)
			continue
		}
		record.ConnectedAt = connectedAt.UnixMilli()
		if disconnectedAt != nil {
			record.DisconnectedAt = disconnectedAt.UnixMilli()
		}
		records = append(records, &record)
	}
	return records, nil
}

func (p *Postgres) GetEntities(ctx context.Context) ([]*model.ConnectionsData, error) {
	return nil, nil
}

func (p *Postgres) GetRegisteredDevices(ctx context.Context, req *model.GetRegisteredDevicesRequest) ([]*model.RegisteredDevice, error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	args := []any{}
	conditions := []string{}
	argIdx := 1
	if req.IMEI != "" {
		args = append(args, req.IMEI)
		conditions = append(conditions, fmt.Sprintf("imei=$%d", argIdx))
		argIdx++
	}
	where := ""
	if len(conditions) != 0 {
		where = "where " + strings.Join(conditions, " AND ")
	}
	query := fmt.Sprintf(`select 
			imei, 
			tenant_id, 
			parser_id, 
			status
		from %s %s LIMIT 25`,
		config.Config.Postgres.RegisteredDevicesTable, where)
	p.telem.Infof(ctx, "Query to fetch registered devices: %s", strings.Join(strings.Fields(query), " "))

	rows, err := p.conn.Query(ctxWithTimeout, query, args...)
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
			response_at_ms
		FROM %s ORDER BY id DESC LIMIT 25`,
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
