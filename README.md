# T27 Fundraising Lambda

This is a bundle of functionality that is used by an AWS and CLI
application to expose a GraphQL interface for the fundraiser db.

This is written in a modular way to allow for easy porting to
Azure/GCP/Other if the time comes.

## Database SQL Schema

```SQL
CREATE TABLE mulch_orders (
    order_id UUID PRIMARY KEY DEFAULT gen_random_uuid(), order_owner_id STRING,
    cash_amount_collected DECIMAL(13, 4), check_amount_collected DECIMAL(13, 4), check_numbers STRING,
    amount_from_donations DECIMAL(13, 4), amount_from_purchases DECIMAL(13, 4),
    will_collect_money_later BOOL, total_amount_collected DECIMAL(13,4), special_instructions STRING,
    is_verified BOOL, last_modified_time TIMESTAMP, purchases JSONB, delivery_id INT,
    customer_addr1 STRING, customer_addr2 STRING, customer_zipcode INT, customer_city STRING,
    customer_neighborhood STRING, known_addr_id UUID, customer_email STRING,
    customer_phone STRING, customer_name STRING);
```

```SQL
CREATE TABLE mulch_spreaders (order_id UUID PRIMARY KEY, spreaders JSONB);
```

```SQL

CREATE TABLE mulch_delivery_timecards (uid STRING, delivery_id INT, last_modified_time TIMESTAMP, time_in TIME, time_out TIME, time_total TIME, PRIMARY KEY (uid, delivery_id, time_in));

```

```SQL
CREATE TABLE fundraiser_config (kind STRING PRIMARY KEY, description STRING, last_modified_time TIMESTAMP, is_locked BOOL, products JSONB, mulch_delivery_configs JSONB, finalization_data JSONB);
```

```SQL
CREATE TABLE neighborhoods (name STRING PRIMARY KEY, zipcode INTEGER, city STRING, dist_pt STRING, is_visible BOOL, last_modified_time TIMESTAMP, meta JSONB);
```

```SQL
CREATE TABLE users (id STRING, group_id STRING, first_name STRING, last_name STRING, created_time TIMESTAMP, last_modified_time TIMESTAMP, has_auth_creds BOOL);
```

```SQL
CREATE TABLE allocation_summary (uid STRING PRIMARY KEY, bags_sold INT, bags_spread DECIMAL(13,4), delivery_minutes DECIMAL(13,4), total_donations DECIMAL(13,4), allocation_from_bags_sold DECIMAL(13,4), allocation_from_bags_spread DECIMAL(13,4), allocation_from_delivery DECIMAL(13,4), allocation_total DECIMAL(13,4));
```

```SQL
CREATE TABLE known_addrs (id UUID PRIMARY KEY DEFAULT gen_random_uuid(), addr STRING, zipcode INTEGER, city STRING, lat STRING, lng STRING, last_modified_time TIMESTAMP, created_time TIMESTAMP);
```
