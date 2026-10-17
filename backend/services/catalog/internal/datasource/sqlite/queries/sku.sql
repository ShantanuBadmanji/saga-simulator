-- name: GetSkuByID :one
SELECT * FROM sku WHERE id = ?1;

-- name: GetSkuBySkuCode :one
SELECT * FROM sku WHERE sku_code = ?1;

-- name: ListSkusByProductID :many
SELECT * FROM sku WHERE product_id = ?1;
