# Event Contract

TipDrop uses RabbitMQ for domain events and async jobs.

## Envelope

```json
{
  "id": "event_uuid",
  "version": 1,
  "type": "tip.slip_uploaded",
  "occurred_at": "2026-05-04T00:00:00Z",
  "idempotency_key": "tip:123:slip_uploaded",
  "payload": {}
}
```

## Events

- `tip.created`
- `tip.slip_uploaded`
- `tip.confirmed`
- `tip.disputed`
- `tip.expired`
- `notification.send`
- `leaderboard.refresh`
- `payment_profile.verified`

## Retry Rules

Consumers should be idempotent. Transient failures may be retried by RabbitMQ dead-letter/retry queues. Permanent payload errors should be logged and dead-lettered.
