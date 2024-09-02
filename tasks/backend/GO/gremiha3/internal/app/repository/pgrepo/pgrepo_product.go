package pgrepo

import (
	"context"
	"fmt"

	"github.com/KozlovNikolai/test-task/internal/app/domain"
	"github.com/KozlovNikolai/test-task/internal/app/repository/models"
	"github.com/KozlovNikolai/test-task/internal/pkg/pg"
)

type ProductRepo struct {
	db *pg.DB
}

func NewProductRepo(db *pg.DB) *ProductRepo {
	return &ProductRepo{
		db: db,
	}
}

// CreateProduct implements services.IProductRepository.
func (r *ProductRepo) CreateProduct(ctx context.Context, product domain.Product) (domain.Product, error) {
	dbProduct := domainToProduct(product)

	var insertedProduct models.Product

	// Начинаем транзакцию
	tx, err := r.db.WR.Begin(ctx)
	if err != nil {
		return domain.Product{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Вставка данных о аре и получение ID
	err = tx.QueryRow(ctx, `
			INSERT INTO products (name, provider_id, price, stock)
			VALUES ($1, $2, $3, $4)
			RETURNING id,name,provider_id,price,stock`,
		dbProduct.Name, dbProduct.ProviderID, dbProduct.Price, dbProduct.Stock).
		Scan(
			&insertedProduct.ID,
			&insertedProduct.Name,
			&insertedProduct.ProviderID,
			&insertedProduct.Price,
			&insertedProduct.Stock)
	if err != nil {
		return domain.Product{}, fmt.Errorf("failed to insert product: %w", err)
	}
	// Фиксация транзакции
	if err := tx.Commit(ctx); err != nil {
		return domain.Product{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	domainProduct, err := productToDomain(insertedProduct)
	if err != nil {
		return domain.Product{}, fmt.Errorf("failed to create domain product: %w", err)
	}

	return domainProduct, nil
}

// DeleteProduct implements service.IProductRepository.
func (r *ProductRepo) DeleteProduct(ctx context.Context, id int) error {
	if id == 0 {
		return fmt.Errorf("%w: id", domain.ErrRequired)
	}
	// Начинаем транзакцию
	tx, err := r.db.WR.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)
	// Проверяем, что товар не связан ни с одним заказом.
	var count int
	err = tx.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM items
		WHERE product_id = (SELECT id FROM products WHERE id = $1)`, id).Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to request the products products: %w", err)
	}
	if count > 0 {
		return fmt.Errorf("error, there are product-related products.: %w", err)
	}
	// Удаляем поставщика
	_, err = tx.Exec(ctx, `
		DELETE FROM products
		WHERE product_id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete product with id %d: %w", id, err)
	}
	// Фиксация транзакции
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

// GetProducts implements service.IProductRepository.
func (r *ProductRepo) GetProducts(ctx context.Context, limit, offset int) ([]domain.Product, error) {

	query := `
		SELECT id, name, provider_id, price, stock
		FROM products
		ORDER BY id
		LIMIT $1 OFFSET $2
	`
	// Запрос
	rows, err := r.db.RO.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()
	// Заполняем массив поставщиков
	var products []models.Product
	for rows.Next() {
		var product models.Product
		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.ProviderID,
			&product.Price,
			&product.Stock)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		products = append(products, product)
	}

	// Проверка на ошибки, возникшие при итерации по строкам
	if rows.Err() != nil {
		return nil, fmt.Errorf("error occurred during row iteration: %w", rows.Err())
	}
	// мапим модель в домен
	domainProducts := make([]domain.Product, len(products))
	for i, product := range products {
		domainProduct, err := productToDomain(product)
		if err != nil {
			return nil, fmt.Errorf("failed to create domain product: %w", err)
		}

		domainProducts[i] = domainProduct
	}
	return domainProducts, nil
}

// GetProduct implements service.IProductRepository.
func (r *ProductRepo) GetProduct(ctx context.Context, id int) (domain.Product, error) {
	if id == 0 {
		return domain.Product{}, fmt.Errorf("%w: id", domain.ErrRequired)
	}

	// SQL-запрос на получение данных Товара по ID
	query := `
SELECT id, name, provider_id, price, stock
FROM products
WHERE id = $1
`
	var product models.Product
	// Выполняем запрос и сканируем результат в структуру Product
	err := r.db.RO.QueryRow(ctx, query, id).Scan(
		&product.ID,
		&product.Name,
		&product.ProviderID,
		&product.Price,
		&product.Stock)
	if err != nil {
		return domain.Product{}, fmt.Errorf("failed to get product by id: %w", err)
	}

	domainProduct, err := productToDomain(product)
	if err != nil {
		return domain.Product{}, fmt.Errorf("failed to create domain product: %w", err)
	}

	return domainProduct, nil
}

// UpdateProduct implements service.IProductRepository.
func (r *ProductRepo) UpdateProduct(ctx context.Context, product domain.Product) (domain.Product, error) {
	dbProduct := domainToProduct(product)
	//dbProduct.UpdatedAt = time.Now()
	// Начинаем транзакцию
	tx, err := r.db.WR.Begin(ctx)
	if err != nil {
		return domain.Product{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)
	// SQL-запрос на обновление данных Поставщика
	query := `
		UPDATE products
		SET name = $1, provider_id = $2, price = $3, stock = $4
		WHERE id = $5
		RETURNING id, name, provider_id, price, stock
	`
	var updatedProduct models.Product

	// Выполняем запрос и сканируем обновленный результат в структуру Product
	err = tx.QueryRow(ctx, query,
		dbProduct.Name,
		dbProduct.ProviderID,
		dbProduct.Price,
		dbProduct.Stock,
		dbProduct.ID).
		Scan(
			&updatedProduct.ID,
			&updatedProduct.Name,
			&updatedProduct.ProviderID,
			&updatedProduct.Price,
			&updatedProduct.Stock)
	if err != nil {
		return domain.Product{}, fmt.Errorf("failed to update product: %w", err)
	}
	// Фиксация транзакции
	if err := tx.Commit(ctx); err != nil {
		return domain.Product{}, fmt.Errorf("failed to commit transaction: %w", err)
	}
	domainProduct, err := productToDomain(updatedProduct)
	if err != nil {
		return domain.Product{}, fmt.Errorf("failed to create domain product: %w", err)
	}

	return domainProduct, nil
}
