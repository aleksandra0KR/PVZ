package postgres

import (
	"context"
	"final/internal/domain"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
)

type ProductRepository struct {
	db                  *sqlx.DB
	receptionRepository *ReceptionRepository
}

func NewProductRepository(db *sqlx.DB, receptionRepository *ReceptionRepository) *ProductRepository {
	return &ProductRepository{db: db, receptionRepository: receptionRepository}
}

func (r *ProductRepository) CreateProduct(ctx context.Context, inputProduct *domain.InputProduct) (*domain.Product, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		log.Error(err)
		return nil, domain.ErrCreateProduct
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Errorf("rollback failed: %v", rollbackErr)
			}
		} else {
			commitErr := tx.Commit()
			if commitErr != nil {
				log.Errorf("commit failed: %v", commitErr)
			}
		}
	}()

	receptionID, err := r.receptionRepository.getLastReceptionIDInProgressForPVZTx(ctx, tx, inputProduct.PvzId)
	if err != nil {
		return nil, domain.ErrCreateProduct
	} else if receptionID == nil {
		return nil, domain.ErrReceptionNotFound
	}

	var product domain.Product
	product.ReceptionID = receptionID
	product.Type = inputProduct.Type

	query := `INSERT INTO products (reception_id, type)
              VALUES ($1, $2)
              RETURNING id, date_time`
	row := tx.QueryRowContext(ctx, query, product.ReceptionID, product.Type)
	err = row.Scan(&product.ID, &product.DateTime)
	if err != nil {
		log.Error(err)
		return nil, domain.ErrCreateProduct
	}
	return &product, nil
}

func (r *ProductRepository) DeleteLastProductForPVZ(ctx context.Context, pvzId *string) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		log.Error(err)
		return domain.ErrDeleteProduct
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Errorf("rollback failed: %v", rollbackErr)
			}
		} else {
			commitErr := tx.Commit()
			if commitErr != nil {
				log.Errorf("commit failed: %v", commitErr)
			}
		}
	}()

	receptionID, err := r.receptionRepository.getLastReceptionIDInProgressForPVZTx(ctx, tx, pvzId)
	if err != nil {
		return domain.ErrDeleteProduct
	} else if receptionID == nil {
		return domain.ErrReceptionNotFound
	}

	productCount, err := r.getAmountOfProductsForReceptionTx(ctx, tx, receptionID)
	if err != nil {
		return domain.ErrDeleteProduct
	} else if productCount == 0 {
		return domain.ErrNoProductsInReception
	}

	query := `
		DELETE FROM products
		WHERE id = (
			SELECT id
			FROM products
			WHERE reception_id = $1
			ORDER BY date_time DESC
			LIMIT 1
		)`
	_, err = tx.ExecContext(ctx, query, *receptionID)
	if err != nil {
		log.Error(err)
		return domain.ErrDeleteProduct
	}
	return nil
}

func (r *ProductRepository) getAmountOfProductsForReceptionTx(ctx context.Context, tx *sqlx.Tx, receptionId *string) (int, error) {
	var productCount int
	query := `SELECT COUNT(*) FROM products WHERE reception_id = $1`
	err := tx.QueryRowContext(ctx, query, *receptionId).Scan(&productCount)
	if err != nil {
		log.Error(err)
		return 0, err
	}
	return productCount, nil
}
