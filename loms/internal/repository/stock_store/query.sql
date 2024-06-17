-- name: GetAvailableForReserve :one
select sku, total_count, reserved 
    from stock
    where sku = $1 limit 1;

-- name: AddReserve :exec
update stock
    set reserved = reserved + $2  
    where sku = $1;

-- name: DelItemFromReserve :exec
update stock
    set total_count = total_count - $2,
        reserved = reserved - $2  
    where sku = $1;

-- name: CancelReserve :exec
update stock
    set reserved = reserved - $2  
    where sku = $1;
