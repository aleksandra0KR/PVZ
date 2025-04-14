package postgres

import (
	"database/sql"
	"final/internal/domain"
	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
	"time"
)

type PVZRepository struct {
	db *sqlx.DB
}

func NewPvzRepository(db *sqlx.DB) *PVZRepository {
	return &PVZRepository{db: db}
}

func (r *PVZRepository) CreatePVZ(pvz *domain.PVZ) (*domain.PVZ, error) {
	query := `INSERT INTO pvz (id, registration_date, city)
              VALUES (COALESCE($1, gen_random_uuid()), COALESCE($2, now()), $3)
              RETURNING id, registration_date`

	row := r.db.QueryRow(query, pvz.ID, pvz.RegistrationDate, pvz.City)
	if err := row.Scan(&pvz.ID, &pvz.RegistrationDate); err != nil {
		log.Error(err)
		return nil, domain.ErrCreatePVZ
	}
	return pvz, nil
}

func (r *PVZRepository) GetAllPVZ() ([]*domain.PVZ, error) {
	query := `SELECT * FROM pvz`

	rows, err := r.db.Query(query)
	if err != nil {
		log.Error(err)
		return nil, domain.ErrGetPVZ
	}
	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			log.Error(err)
		}
	}(rows)

	result := make([]*domain.PVZ, 0)
	for rows.Next() {
		var pvz domain.PVZ
		err = rows.Scan(&pvz.ID, &pvz.RegistrationDate, &pvz.City)
		if err != nil {
			log.Error(err)
			return nil, domain.ErrGetPVZ
		}
		result = append(result, &pvz)

	}
	return result, nil
}

func (r *PVZRepository) GetPVZInfo(startDate, endDate *time.Time, offset, limit int) ([]domain.PVZWithReceptions, error) {
	sq := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	query := sq.
		Select(
			"pvz.id as pvz_id", "pvz.registration_date", "pvz.city",
			"receptions.id as reception_id", "receptions.date_time", "receptions.status",
			"products.id as product_id", "products.date_time as product_datetime", "products.type", "products.reception_id as product_reception_id",
		).
		From("pvz").
		LeftJoin("receptions ON pvz.id = receptions.pvz_id").
		LeftJoin("products ON receptions.id = products.reception_id")

	if startDate != nil {
		query = query.Where(squirrel.GtOrEq{"receptions.date_time": *startDate})
	}
	if endDate != nil {
		query = query.Where(squirrel.LtOrEq{"receptions.date_time": *endDate})
	}

	query = query.OrderBy("pvz.registration_date DESC").
		Limit(uint64(limit)).
		Offset(uint64(offset))

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		log.Error(err)
		return nil, domain.ErrGetPVZ
	}

	rows, err := r.db.Query(sqlQuery, args...)
	if err != nil {
		log.Error(err)
		return nil, domain.ErrGetPVZ
	}
	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			log.Error(err)
		}
	}(rows)

	pvzMap := make(map[string]*domain.PVZWithReceptions)
	for rows.Next() {
		var (
			pvzID, receptionID, productID sql.NullString
			pvzDate                       time.Time
			city                          string
			receptionTime, productTime    sql.NullTime
			receptionStatus, productType  sql.NullString
		)

		err = rows.Scan(
			&pvzID, &pvzDate, &city,
			&receptionID, &receptionTime, &receptionStatus,
			&productID, &productTime, &productType, new(sql.NullString),
		)
		if err != nil {
			log.Error(err)
			return nil, domain.ErrGetPVZ
		}

		id := pvzID.String
		if _, ok := pvzMap[id]; !ok {
			pvzMap[id] = &domain.PVZWithReceptions{
				PVZ: domain.PVZ{
					ID:               &id,
					RegistrationDate: &pvzDate,
					City:             &city,
				},
				Receptions: []domain.ReceptionWithProducts{},
			}
		}

		if receptionID.Valid {
			receptionUUID := receptionID.String
			receptionExists := false

			for i := range pvzMap[id].Receptions {
				if (*pvzMap[id].Receptions[i].Reception.ID) == receptionUUID {
					receptionExists = true
					if productID.Valid {
						pvzMap[id].Receptions[i].Products = append(pvzMap[id].Receptions[i].Products, domain.Product{
							ID:          &productID.String,
							DateTime:    &productTime.Time,
							Type:        &productType.String,
							ReceptionID: &receptionUUID,
						})
					}
					break
				}
			}

			if !receptionExists {
				newReception := domain.ReceptionWithProducts{
					Reception: domain.Reception{
						ID:       &receptionUUID,
						DateTime: &receptionTime.Time,
						Status:   &receptionStatus.String,
						PvzId:    &id,
					},
					Products: []domain.Product{},
				}
				if productID.Valid {
					newReception.Products = append(newReception.Products, domain.Product{
						ID:          &productID.String,
						DateTime:    &productTime.Time,
						Type:        &productType.String,
						ReceptionID: &receptionUUID,
					})
				}
				pvzMap[id].Receptions = append(pvzMap[id].Receptions, newReception)
			}
		}
	}

	result := make([]domain.PVZWithReceptions, 0, len(pvzMap))
	for _, v := range pvzMap {
		result = append(result, *v)
	}

	return result, nil
}
