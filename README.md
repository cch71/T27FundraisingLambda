# T27FundraisingListOrderLambda

## Database SQL Schema

```
CREATE TABLE mulch_orders (order_id UUID PRIMARY KEY DEFAULT gen_random_uuid(), order_owner_id STRING, cash_amount_collected DECIMAL(13, 4), check_amount_collected DECIMAL(13, 4), check_numbers STRING, amount_from_donations DECIMAL(13, 4), amount_from_purchases DECIMAL(13, 4), will_collect_money_later BOOL, total_amount_collected DECIMAL(13,4), special_instructions STRING, is_verified BOOL, last_modified_time TIMESTAMP, purchases JSONB, delivery_id INT, customer_addr1 STRING, customer_addr2 STRING, customer_neighborhood STRING, known_addr_id UUID, customer_email STRING, customer_phone STRING, customer_name STRING);
```
```
CREATE TABLE mulch_spreaders (order_id UUID PRIMARY KEY, spreaders JSONB);
```
```
CREATE TABLE mulch_delivery_timecards (uid STRING, delivery_id INT, last_modified_time TIMESTAMP, time_in TIME, time_out TIME, time_total TIME, PRIMARY KEY (uid, delivery_id, time_in));
```
```
CREATE TABLE fundraiser_config (kind STRING PRIMARY KEY, description STRING, last_modified_time TIMESTAMP, is_locked BOOL, products JSONB, mulch_delivery_configs JSONB, finalization_data JSONB);
```
```
CREATE TABLE neighborhoods (name STRING PRIMARY KEY, zipcode INTEGER, city STRING, dist_pt STRING, is_visible BOOL, last_modified_time TIMESTAMP, meta JSONB);
```
```
CREATE TABLE users (id STRING, group_id STRING, name STRING, created_time TIMESTAMP, last_modified_time TIMESTAMP, has_auth_creds BOOL);
```
```
CREATE TABLE allocation_summary (uid STRING PRIMARY KEY, bags_sold INT, bags_spread DECIMAL(13,4), delivery_minutes DECIMAL(13,4), total_donations DECIMAL(13,4), allocation_from_bags_sold DECIMAL(13,4), allocation_from_bags_spread DECIMAL(13,4), allocation_from_delivery DECIMAL(13,4), allocation_total DECIMAL(13,4));
```
```
CREATE TABLE known_addrs (id UUID, addr STRING, zipcode INTEGER, city STRING, lat STRING, lng STRING, last_modified_time TIMESTAMP, PRIMARY KEY (addr, zipcode, city);
```
