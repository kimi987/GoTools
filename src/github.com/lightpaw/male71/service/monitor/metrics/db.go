package metrics

import (
	"context"
	"database/sql"
	"github.com/lightpaw/male7/service/db/isql"
	"github.com/prometheus/client_golang/prometheus"
	"sync"
)

func NewMetricsDb(db *sql.DB) *MetricsDB {

	errorCounter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: Namespace,
		Subsystem: "db",
		Name:      "error_count",
		Help:      "database access error count.",
	},
		[]string{"sql"})

	queryTimeCost := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace:  Namespace,
		Subsystem:  "db",
		Name:       "time_cost",
		Help:       "database access time cost.",
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	},
		[]string{"sql"})

	return &MetricsDB{
		db:            db,
		errorCounter:  errorCounter,
		queryTimeCost: queryTimeCost,
		queryMap:      make(map[string]*query),
	}
}

var _ isql.DB = (*MetricsDB)(nil)

type MetricsDB struct {
	db *sql.DB

	errorCounter  *prometheus.CounterVec
	queryTimeCost *prometheus.SummaryVec

	queryMap       map[string]*query
	queryMapLocker sync.RWMutex
}

func (db *MetricsDB) Close() error {
	return db.db.Close()
}

func (db *MetricsDB) Collectors() []prometheus.Collector {
	return []prometheus.Collector{db.errorCounter, db.queryTimeCost}
}

// query

// 最多监控100个sql，如果真的超出，加大这个值
// 这个主要是用来防止bug出现，比如存在一些自己拼接的sql语句，导致每个语句都创建一个Metrics
const max_query = 100

var (
	emptyCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "empty_counter",
	})

	emptySummary = prometheus.NewSummary(prometheus.SummaryOpts{
		Name:       "empty_summary",
		Objectives: map[float64]float64{},
	})

	emptyQuery = &query{
		query:         "empty_query",
		errorCounter:  emptyCounter,
		queryTimeCost: emptySummary,
	}
)

// 单个查询语句
type query struct {
	query string

	errorCounter  prometheus.Counter
	queryTimeCost prometheus.Observer
}

func (db *MetricsDB) getOrCreateQuery(s string) *query {
	q, createNew := db.getCacheQuery(s)
	if q != nil {
		return q
	}

	if createNew {
		db.queryMapLocker.Lock()
		defer db.queryMapLocker.Unlock()

		q = &query{
			query:         s,
			errorCounter:  db.errorCounter.WithLabelValues(s),
			queryTimeCost: db.queryTimeCost.WithLabelValues(s),
		}

		db.queryMap[s] = q

		return q
	}

	return emptyQuery // 使用一个默认的代替，不上报
}

func (db *MetricsDB) getCacheQuery(s string) (q *query, createNew bool) {

	db.queryMapLocker.RLock()
	defer db.queryMapLocker.RUnlock()

	q = db.queryMap[s]
	if q != nil {
		return
	}

	// 如果当前数量已经超过阈值，那么不再创建新的
	createNew = len(db.queryMap) < max_query
	return
}

// db

func (db *MetricsDB) Query(query string, args ...interface{}) (isql.Rows, error) {
	return db.QueryContext(context.Background(), query, args...)
}

func (db *MetricsDB) QueryContext(ctx context.Context, query string, args ...interface{}) (isql.Rows, error) {
	q := db.getOrCreateQuery(query)

	timer := prometheus.NewTimer(q.queryTimeCost)
	defer timer.ObserveDuration()

	rows, err := db.db.QueryContext(ctx, query, args...)
	if err != nil {
		if err != sql.ErrNoRows {
			q.errorCounter.Inc()
		}

		return nil, err
	}

	return db.NewRows(q, rows), nil
}

func (db *MetricsDB) QueryRow(query string, args ...interface{}) isql.Row {
	return db.QueryRowContext(context.Background(), query, args...)
}

func (db *MetricsDB) QueryRowContext(ctx context.Context, query string, args ...interface{}) isql.Row {
	q := db.getOrCreateQuery(query)

	timer := prometheus.NewTimer(q.queryTimeCost)
	defer timer.ObserveDuration()

	return db.NewRow(q, db.db.QueryRowContext(ctx, query, args...))
}

func (db *MetricsDB) Exec(query string, args ...interface{}) (sql.Result, error) {
	return db.ExecContext(context.Background(), query, args...)
}

func (db *MetricsDB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	q := db.getOrCreateQuery(query)

	timer := prometheus.NewTimer(q.queryTimeCost)
	defer timer.ObserveDuration()

	result, err := db.db.ExecContext(ctx, query, args...)
	if err != nil {
		q.errorCounter.Inc()
		return nil, err
	}

	return result, nil
}

// Row
func (db *MetricsDB) NewRow(q *query, r *sql.Row) *row {
	return &row{
		Row: r,
		q:   q,
	}
}

var _ isql.Row = (*row)(nil)

type row struct {
	*sql.Row

	q *query
}

func (r *row) Scan(dest ...interface{}) error {
	err := r.Row.Scan(dest...)
	if err != nil {
		if err != sql.ErrNoRows {
			r.q.errorCounter.Inc()
		}
		return err
	}

	return nil
}

// Rows
func (db *MetricsDB) NewRows(q *query, r *sql.Rows) *rows {
	return &rows{
		Rows: r,
		q:    q,
	}
}

var _ isql.Rows = (*rows)(nil)

type rows struct {
	*sql.Rows

	q *query
}

func (r *rows) Scan(dest ...interface{}) error {
	err := r.Rows.Scan(dest...)
	if err != nil {
		r.q.errorCounter.Inc()
		return err
	}

	return nil
}
