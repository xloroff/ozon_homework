/* При вставке сообщения в outbox гарантируем что если есть заблокированные записи c entity_id, статус lock и время блокировки будет выставлено таким же */
-- name: AddOutbox :one
with locked_entity as (
    select status, locked_to
    from "outbox"
    where entity_id = sqlc.arg(entity_id)::text
      and status = 'locked'
    order by locked_to desc
    LIMIT 1
)
insert into "outbox" (entity_id, "payload", "metadata", "status", locked_to)
    values (
               sqlc.arg(entity_id)::text,
               sqlc.arg(payload)::text,
               sqlc.arg(metadata)::jsonb,
               coalesce((SELECT status FROM locked_entity), 'new'),
               coalesce((SELECT locked_to FROM locked_entity), now())
            )
    returning id;

/* При чтении сообщений из outbox, не берем заблокированные записи время блокировки которых не истекло, а так же берем  записи с ошибкой отправки для ретрая */
-- name: Outbox :many
with updated_rows as(
    select *
    from "outbox"
    where status = sqlc.arg(status_new)::outbox_status_type
        or ((status = sqlc.arg(status_locked)::outbox_status_type or  status = sqlc.arg(status_failed)::outbox_status_type)
        and locked_to < now() - (sqlc.arg(locked_to)::text || ' seconds')::interval)
    order by created_at asc
        for update nowait
) update "outbox"
    set status = sqlc.arg(status_locked)::outbox_status_type,
        locked_to = now() + (sqlc.arg(locked_to) ||'seconds')::interval
from updated_rows
    where "outbox".id = updated_rows.id
returning updated_rows.*;

-- name: SetStatusOutbox :exec
update "outbox"
    set status = sqlc.arg(new_status)::outbox_status_type,
        updated_at = now()
    where id = sqlc.arg(id_msg)::varchar;