-- name: GetProductByID :one
SELECT * FROM product WHERE id = ?1;

-- name: ListProductsWithSkusAsc :many
WITH filtered_products AS (
    SELECT p.*
    FROM product p
    WHERE p.status = 'ACTIVE'
      AND p.category_id = ?1
      AND ?1 IS NOT NULL

    UNION ALL

    SELECT p.*
    FROM product p
    WHERE p.status = 'ACTIVE'
      AND ?1 IS NULL
),
ranked_skus AS (
    SELECT
        s.id,
        s.product_id,
        s.price_minor,
        ROW_NUMBER() OVER (
            PARTITION BY s.product_id
            ORDER BY s.price_minor ASC, s.id ASC
        ) AS rn
    FROM sku s
    JOIN filtered_products p
        ON p.id = s.product_id
    WHERE s.active = 1
)
SELECT
    p.*,
    s.id AS sku_id,
    s.price_minor AS min_price,
    COUNT(*) OVER () as total_count
FROM filtered_products p
JOIN ranked_skus s
    ON s.product_id = p.id
   AND s.rn = 1
ORDER BY s.price_minor ASC, p.id ASC
LIMIT ?2 OFFSET ?3;

-- name: ListProductMediaByProductID :many
SELECT * FROM product_media
WHERE product_id = ?1