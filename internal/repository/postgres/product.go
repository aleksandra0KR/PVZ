package postgres

import (
	"final/internal/domain"
	"fmt"
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

func (r *ProductRepository) CreateProduct(inputProduct *domain.InputProduct) (*domain.Product, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Printf("rollback failed: %v", rollbackErr)
			}
		} else {
			tx.Commit()
		}
	}()

	receptionID, err := r.receptionRepository.getLastReceptionIDInProgressForPVZTx(tx, inputProduct.PvzId)
	if err != nil {
		return nil, err
	} else if receptionID == nil {
		return nil, fmt.Errorf("no reception in progress for PVZ with ID: %s", *inputProduct.PvzId)
	}
	var product domain.Product
	product.ReceptionID = receptionID
	product.Type = inputProduct.Type

	query := `INSERT INTO products (reception_id, type)
              VALUES ($1, $2)
              RETURNING id, date_time`
	row := tx.QueryRow(query, product.ReceptionID, product.Type)
	err = row.Scan(&product.ID, &product.DateTime)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepository) DeleteLastProductForPVZ(pvzId *string) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Printf("rollback failed: %v", rollbackErr)
			}
		} else {
			tx.Commit()
		}
	}()
	receptionID, err := r.receptionRepository.getLastReceptionIDInProgressForPVZTx(tx, pvzId)
	if err != nil {
		return err
	} else if receptionID == nil {
		return fmt.Errorf("no reception in progress for PVZ with ID: %s", *pvzId)
	}

	productCount, err := r.getAmountOfProductsForReceptionTx(tx, receptionID)
	if productCount == 0 {
		return fmt.Errorf("no products to delete in the reception with ID: %s", *receptionID)
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
	_, err = tx.Exec(query, *receptionID)
	if err != nil {
		return err
	}
	return nil
}

func (r *ProductRepository) getAmountOfProductsForReceptionTx(tx *sqlx.Tx, receptionId *string) (int, error) {
	var productCount int
	query := `SELECT COUNT(*) FROM products WHERE reception_id = $1`
	err := tx.QueryRow(query, *receptionId).Scan(&productCount)
	if err != nil {
		return 0, err
	}
	return productCount, nil
}
