package repository

import (
	"challengephp/lib"
	"challengephp/src/db"
	"challengephp/src/types"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"strings"
)

type Repository struct {
	ctx  context.Context
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool, ctx context.Context) Repository {
	return Repository{
		ctx:  ctx,
		pool: pool,
	}
}

func (repo Repository) FindUser(id int64) (types.User, bool, lib.Error) {
	var user types.User
	err := repo.pool.QueryRow(repo.ctx, "SELECT id FROM "+db.TableUsers+" WHERE id = $1", id).Scan(
		&user.Id,
	)
	if err == nil {
		return user, true, nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return types.User{}, false, nil
	}
	return types.User{}, false, lib.Err(err)
}

func (repo Repository) FindEventType(name string) (types.EventType, bool, lib.Error) {
	var eventType types.EventType
	err := repo.pool.QueryRow(repo.ctx, "SELECT id, name FROM "+db.TableEventTypes+" WHERE name = $1", name).Scan(
		&eventType.Id,
		&eventType.Name,
	)
	if err == nil {
		return eventType, true, nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return types.EventType{}, false, nil
	}
	return types.EventType{}, false, lib.Err(err)
}

func (repo Repository) FindOrCreateUser(id int64) (types.User, lib.Error) {
	var user types.User
	err := repo.pool.QueryRow(repo.ctx, "SELECT id FROM "+db.TableUsers+" WHERE id = $1", id).Scan(
		&user.Id,
	)
	if err == nil {
		return user, nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		err = repo.pool.QueryRow(repo.ctx, `
            INSERT INTO `+db.TableUsers+` (id, name, created_at)
            VALUES ($1, $2, NOW())
            RETURNING id, name
        `, id, fmt.Sprintf("user%d", id)).Scan(&user.Id, &user.Name)
		if err != nil {
			return types.User{}, lib.Err(err)
		}
		return user, nil
	}
	return types.User{}, lib.Err(err)
}

func (repo Repository) FindOrCreateEvent(name string) (types.EventType, lib.Error) {
	var eventType types.EventType
	err := repo.pool.QueryRow(repo.ctx, "SELECT id, name FROM "+db.TableEventTypes+" WHERE name = $1", name).Scan(
		&eventType.Id,
		&eventType.Name,
	)
	if err == nil {
		return eventType, nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		err = repo.pool.QueryRow(repo.ctx, `
			INSERT INTO `+db.TableEventTypes+` (name)
			VALUES ($1)
			RETURNING id, name
		`, name).Scan(&eventType.Id, &eventType.Name)
		if err != nil {
			return types.EventType{}, lib.Err(err)
		}
		return eventType, nil
	}
	return types.EventType{}, lib.Err(err)
}

func (repo Repository) CreateEvent(event types.Event) (int64, lib.Error) {
	var eventID int64
	metadataBytes, err := json.Marshal(event.Metadata)
	if err != nil {
		return 0, lib.Err(err)
	}
	query := `
        INSERT INTO ` + db.TableEvents + ` (timestamp, metadata, user_id, type_id)
        VALUES ($1, $2, $3, $4)
        RETURNING id`
	err = repo.pool.QueryRow(repo.ctx, query,
		event.Timestamp,
		metadataBytes,
		event.UserId,
		event.TypeId,
	).Scan(&eventID)
	if err != nil {
		return 0, lib.Err(err)
	}
	return eventID, nil
}

func (repo Repository) EventsTotal() (int64, lib.Error) {
	var total int64
	query := "SELECT COUNT(*) FROM " + db.TableEvents
	err := repo.pool.QueryRow(repo.ctx, query).Scan(&total)
	if err != nil {
		return 0, lib.Err(err)
	}
	return total, nil
}

func (repo Repository) List(page, limit int64) ([]types.Event2, lib.Error) {
	offset := (page - 1) * limit
	query := `
        SELECT e.id, e.timestamp, e.metadata, e.type_id, e.user_id, et.name AS type
		FROM ` + db.TableEvents + ` e
			JOIN ` + db.TableEventTypes + ` et ON e.type_id = et.id
		ORDER BY e.timestamp DESC
		LIMIT $1 OFFSET $2`
	rows, err := repo.pool.Query(repo.ctx, query, limit, offset)
	if err != nil {
		return nil, lib.Err(err)
	}
	defer rows.Close()
	var events []types.Event2
	for rows.Next() {
		var event types.Event2
		var metadataMap map[string]interface{}
		err := rows.Scan(
			&event.Id,
			&event.Timestamp,
			&metadataMap,
			&event.TypeId,
			&event.UserId,
			&event.TypeName,
		)
		if err != nil {
			return nil, lib.Err(err)
		}

		event.Metadata = metadataMap
		events = append(events, event)
	}
	if err := rows.Err(); err != nil {
		return nil, lib.Err(err)
	}
	return events, nil
}

func (repo Repository) ListForUser(userId int64) ([]types.Event2, lib.Error) {
	query := `
		SELECT e.id, e.timestamp, e.metadata, e.type_id, e.user_id, et.name AS type
		FROM ` + db.TableEvents + ` e
			JOIN ` + db.TableEventTypes + ` et ON e.type_id = et.id
		WHERE e.user_id = $1
		ORDER BY e.timestamp DESC
		LIMIT 1000
	`
	rows, err := repo.pool.Query(repo.ctx, query, userId)
	if err != nil {
		return nil, lib.Err(err)
	}
	defer rows.Close()
	var events []types.Event2
	for rows.Next() {
		var event types.Event2
		var metadataMap map[string]interface{}
		err := rows.Scan(
			&event.Id,
			&event.Timestamp,
			&metadataMap,
			&event.TypeId,
			&event.UserId,
			&event.TypeName,
		)
		if err != nil {
			return nil, lib.Err(err)
		}
		event.Metadata = metadataMap
		events = append(events, event)
	}
	if err := rows.Err(); err != nil {
		return nil, lib.Err(err)
	}
	return events, nil
}

func (repo Repository) GetStatTotal(typeId int64, from, to string) (int64, lib.Error) {
	whereConditions, whereParams := getWhere(typeId, from, to)
	var total int64
	query := "SELECT COUNT(*) FROM " + db.TableEvents + " " + whereConditions
	err := repo.pool.QueryRow(repo.ctx, query, whereParams...).Scan(&total)
	if err != nil {
		return 0, lib.Err(err)
	}
	return total, nil
}

func (repo Repository) GetStatUniqueUsers(typeId int64, from, to string) (int64, lib.Error) {
	whereConditions, whereParams := getWhere(typeId, from, to)
	query := `
		select count(distinct user_id)
		from ` + db.TableEvents + `
		` + whereConditions
	total := int64(0)
	err := repo.pool.QueryRow(repo.ctx, query, whereParams...).Scan(&total)
	if err != nil {
		return 0, lib.Err(err)
	}
	return total, nil
}

func (repo Repository) GetStatsPages(typeId int64, from, to string) (map[string]int64, lib.Error) {
	whereConditions, whereParams := getWhere(typeId, from, to)
	query := `
		SELECT
			metadata->>'page' as page,
			COUNT(*) as page_count
		FROM ` + db.TableEvents + `
		` + whereConditions + `
		group by page
		order by page_count desc
	`
	rows, err := repo.pool.Query(repo.ctx, query, whereParams...)
	if err != nil {
		return nil, lib.Err(err)
	}
	defer rows.Close()
	pages := make(map[string]int64)
	for rows.Next() {
		var key string
		var val int64
		err := rows.Scan(
			&key,
			&val,
		)
		if err != nil {
			return nil, lib.Err(err)
		}
		pages[key] = val
	}
	if err := rows.Err(); err != nil {
		return nil, lib.Err(err)
	}
	return pages, nil
}

func getWhere(typeId int64, from, to string) (string, []any) {
	whereConditions := make([]string, 0, 3)
	whereParams := make([]any, 0, 3)

	if typeId > 0 {
		whereConditions = append(whereConditions, fmt.Sprintf("type_id = $%d", len(whereConditions)+1))
		whereParams = append(whereParams, typeId)
	}
	if from != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("$%d <= timestamp", len(whereConditions)+1))
		whereParams = append(whereParams, from)
	}
	if to != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("timestamp <= $%d", len(whereConditions)+1))
		whereParams = append(whereParams, to)
	}
	where := ""
	if len(whereConditions) > 0 {
		where = " WHERE " + strings.Join(whereConditions, " AND ")
	}
	return where, whereParams
}
