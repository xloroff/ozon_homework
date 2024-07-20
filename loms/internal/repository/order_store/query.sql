-- name: AddOrder :one
insert into "order" ("id", "user", status, created_at, updated_at)
    values (nextval('order_id_manual_seq') + $1, $2, $3, now(), now())
    returning id;

-- name: AddOrderItem :exec
insert into order_item (order_id, sku, count)
    values ($1, $2, $3);

-- name: SetStatus :exec
update "order" 
    set status = $2, updated_at = now()
    where id = $1;

-- name: GetOrder :one
select id, "user" as user, status, created_at, updated_at 
    from "order" 
    where id = $1 limit 1;

-- name: GetOrderItemsByOrderIDs :many
select id, sku, count
    from order_item
    where order_id = $1;