---
subcategory: "PaaS"
layout: "aws"
page_title: "CROC Cloud: aws_paas_service"
description: |-
  Manages a PaaS service.
---

[paas]: https://docs.cloud.croc.ru/en/services/paas/index.html
[technical support]: https://support.croc.ru/app/#/project/CS
[timeouts]: https://www.terraform.io/docs/configuration/blocks/resources/syntax.html#operation-timeouts

# Resource: aws_paas_service

Manages a PaaS service. For details about PaaS, see the [user documentation][paas].

## Example Usage

### Elasticsearch Service

```terraform
resource "aws_vpc" "example" {
  cidr_block = "172.16.0.0/16"

  tags = {
    Name = "tf-vpc"
  }
}

resource "aws_subnet" "example" {
  vpc_id            = aws_vpc.example.id
  cidr_block        = cidrsubnet(aws_vpc.example.cidr_block, 4, 1)
  availability_zone = "ru-msk-vol52"

  tags = {
    Name = "tf-subnet"
  }
}

resource "aws_paas_service" "elasticsearch" {
  name          = "tf-service"
  instance_type = "c5.large"

  root_volume {
    type = "st2"
    size = 32
  }

  data_volume {
    type = "st2"
    size = 32
  }

  delete_interfaces_on_destroy = true
  security_group_ids           = [aws_vpc.example.default_security_group_id]
  subnet_ids                   = [aws_subnet.example.id]

  elasticsearch {
    version = "8.2.2"
    kibana  = true
  }
}
```

### Memcached Service with Enabled Monitoring

~> This example uses the VPC and subnet defined in [elasticsearch service example](#elasticsearch-service).

```terraform
resource "aws_paas_service" "memcached" {
  name          = "tf-service"
  instance_type = "c5.large"

  root_volume {
    type = "st2"
    size = 32
  }

  data_volume {
    type = "st2"
    size = 32
  }

  delete_interfaces_on_destroy = true
  security_group_ids           = [aws_vpc.example.default_security_group_id]
  subnet_ids                   = [aws_subnet.example.id]

  memcached {
    monitoring {
      monitor_by = "fm-cluster-12345678"
      monitoring_labels = {
        key1 = "value1"
        key3 = "value3"
      }
    }
  }
}
```

### PostgreSQL Service with Arbitrator

```terraform
resource "aws_vpc" "example" {
  cidr_block = "172.33.0.0/16"

  tags = {
    Name = "tf-vpc"
  }
}

resource "aws_subnet" "subnet_vol52" {
  vpc_id            = aws_vpc.example.id
  cidr_block        = cidrsubnet(aws_vpc.example.cidr_block, 4, 15)
  availability_zone = "ru-msk-vol52"

  tags = {
    Name = "tf-subnet"
  }
}

resource "aws_subnet" "subnet_vol51" {
  vpc_id            = aws_vpc.example.id
  cidr_block        = cidrsubnet(aws_vpc.example.cidr_block, 4, 14)
  availability_zone = "ru-msk-vol51"

  tags = {
    Name = "tf-subnet"
  }
}

resource "aws_subnet" "subnet_comp1p" {
  vpc_id            = aws_vpc.example.id
  cidr_block        = cidrsubnet(aws_vpc.example.cidr_block, 4, 13)
  availability_zone = "ru-msk-comp1p"

  tags = {
    Name = "tf-subnet"
  }
}

resource "aws_s3_bucket" "example" {
  bucket = "tf-paas-backup"

  # Use the predefined provider configuration to connect to CROC Cloud S3 storage
  # https://docs.cloud.croc.ru/en/api/tools/terraform.html#providers-tf
  provider = aws.noregion
}

resource "aws_paas_service" "pgsql" {
  name = "tf-service"

  arbitrator_required = true
  high_availability   = true

  instance_type = "c5.large"

  root_volume {
    type = "st2"
    size = 32
  }

  data_volume {
    type = "st2"
    size = 32
  }

  delete_interfaces_on_destroy = true
  security_group_ids           = [aws_vpc.example.default_security_group_id]
  subnet_ids                   = [aws_subnet.subnet_vol52.id, aws_subnet.subnet_vol51.id, aws_subnet.subnet_comp1p.id]

  backup_settings {
    enabled            = true
    expiration_days    = 5
    notification_email = "example@mail.com"
    start_time         = "15:10"
    bucket_name        = aws_s3_bucket.example.id
    user_login         = "user@company"
  }

  pgsql {
    version = "10.21"

    autovacuum_analyze_scale_factor = 0.3
    min_wal_size                    = 85 * 1024 * 1024
    max_wal_size                    = 85 * 1024 * 1024
    work_mem                        = 4 * 1024 * 1024
    maintenance_work_mem            = 1024 * 1024
    wal_keep_segments               = 0
    replication_mode                = "synchronous"

    user {
      name     = "user1"
      password = "********"
    }

    database {
      name           = "test_db1"
      owner          = "user1"
      backup_enabled = true
      extensions     = ["bloom", "dict_int"]
      user {
        name = "user1"
      }
    }

    options = {
      logDestination = "csvlog"
    }
  }
}
```

### Redis Service with Enabled Logging

~> This example uses the VPC and subnet defined in [elasticsearch service example](#elasticsearch-service).

```terraform
resource "aws_paas_service" "redis" {
  name          = "tf-service"
  instance_type = "c5.large"

  root_volume {
    type = "st2"
    size = 32
  }

  data_volume {
    type = "st2"
    size = 32
  }

  delete_interfaces_on_destroy = true
  security_group_ids           = [aws_vpc.example.default_security_group_id]
  subnet_ids                   = [aws_subnet.example.id]

  redis {
    class   = "database"
    version = "5.0.14"

    password = "********"

    persistence_rdb = false
    persistence_aof = true

    monitoring = true

    databases     = 1
    timeout       = 50
    tcp_backlog   = 300
    tcp_keepalive = 600

    logging {
      log_to       = "fm-cluster-87654321"
      logging_tags = ["tag1", "tag2", "tag3"]
    }
  }
}
```

## Argument Reference

~> Arguments are not editable (changes force new resource) except for blocks with service parameters and `backup_settings`.

* `arbitrator_required` - (Optional) Indicates whether to create a cluster with an arbitrator. Defaults to `false`.
  The parameter can be set to `true` only if `high_availability` is `true`.
  The parameter is supported only for [Elasticsearch](#elasticsearch-argument-reference) and [PostgreSQL](#postgresql-argument-reference) services.
* `backup_settings` - (Optional) The backup settings for the service. The structure of this block is [described below](#backup_settings).
  The parameter is supported only for [PostgreSQL](#postgresql-argument-reference) service.
* `data_volume` - (Optional) The data volume parameters for the service. The structure of this block is [described below](#data_volume).
  The parameter is required for [Elasticsearch](#elasticsearch-argument-reference), [Memcached](#memcached-argument-reference),
  [PostgreSQL](#postgresql-argument-reference) and [Redis](#redis-argument-reference) services.
* `delete_interfaces_on_destroy` - (Optional) Indicates whether to delete instance network interfaces when the service is destroyed. Defaults to `false`.
* `high_availability` - (Optional) Indicates whether to create a high availability service. Defaults to `false`.
  The parameter is supported only for [Elasticsearch](#elasticsearch-argument-reference),
  [PostgreSQL](#postgresql-argument-reference) and [Redis](#redis-argument-reference) services.
* `instance_type` - (Required) The instance type.
* `name` - (Required) The service name. The value must start and end with a Latin letter or number and
  can only contain lowercase Latin letters, numbers, periods (.) and hyphens (-).
* `network_interface_ids` - (Required if `subnet_ids` is not specified) List of network interface IDs.
* `root_volume` - (Required) The root volume parameters for the service. The structure of this block is [described below](#root_volume).
* `security_group_ids` - (Required) List of security group IDs.
* `ssh_key_name` - (Optional) The name of the SSH key for accessing instances.
* `subnet_ids` - (Required if `network_interface_ids` is not specified) List of subnet IDs.
* `user_data` - (Required if `user_data_content_type` is specified) User data.
* `user_data_content_type` - (Required if `user_data` is specified) The type of `user_data`. Valid values are `cloud-config`, `x-shellscript`.

One of the following blocks with service parameters must be specified:

* `elasticsearch` - Elasticsearch parameters. The structure of this block is [described below](#elasticsearch-argument-reference).
* `memcached` - Memcached parameters. The structure of this block is [described below](#memcached-argument-reference).
* `pgsql` - PostgreSQL parameters. The structure of this block is [described below](#postgresql-argument-reference).
* `redis` - Redis parameters. The structure of this block is [described below](#redis-argument-reference).

### backup_settings

~> All the parameters in the `backup_settings` block are editable.

The `backup_settings` block has the following structure:

* `bucket_name` - (Optional) The name of the bucket in object storage where the service backup is saved.
  The parameter must be set if `enabled` is `true`.
* `enabled` -  (Optional) Indicates whether backup is enabled for the service. Defaults to `false`.
* `expiration_days` - (Optional) The backup retention period in days. Valid values are from 1 to 3650.
* `notification_email` - (Optional) The email address to which a notification that backup was created is sent.
* `start_time` - (Optional) The time when the daily backup process starts. It is set as a string in the HH:MM format Moscow time.
  The parameter must be set if `enabled` is `true`.
* `user_login` - (Optional) The login of a user with write permissions to the bucket in object storage (e.g. `user@company`).
  The parameter must be set if `enabled` is `true`.

### data_volume

The `data_volume` block has the following structure:

* `iops` - (Optional) The number of read/write operations per second for the data volume.
  The parameter must be set if `type` is `io2`.
* `size` - (Optional) The size of the data volume in GiB. Defaults to `32`.
* `type` - (Optional) The type of the data volume. Valid values are `st2`, `gp2`, `io2`. Defaults to `st2`.

### root_volume

The `root_volume` block has the following structure:

* `iops` - (Optional) The number of read/write operations per second for the root volume.
  The parameter must be set if `type` is `io2`.
* `size` - (Optional) The size of the root volume in GiB. Defaults to `32`.
* `type` - (Optional) The type of the root volume. Valid values are `st2`, `gp2`, `io2`. Defaults to `st2`.

## Elasticsearch Argument Reference

In addition to the common arguments for all services [described above](#argument-reference),
the `elasticsearch` block can contain the following arguments:

* `class` - (Optional) The service class. Valid value is `search`. Defaults to `search`.
* `kibana` - (Optional) Indicates whether Kibana deployment is enabled. Defaults to `false`.
* `logging` - (Optional) The logging settings for the service. The structure of this block is [described below](#logging).
* `monitoring` - (Optional) The monitoring settings for the service. The structure of this block is [described below](#monitoring).
* `options` - (Optional) Map containing other Elasticsearch parameters.
  Parameter names must be in camelCase. Values are strings.

~> If the parameter name includes a dot, then it cannot be passed in the `options`.
If you need to use such a parameter, contact [technical support].

* `password` - (Optional) The Elasticsearch user password.
  The value must be 8 to 128 characters long and must not contain `-`, `!`, `:`, `;`, `%`, `'`, `"`,  `` ` `` and `\`.
* `version` - (Required) The version to install.
  Valid values are `7.11.2`, `7.12.1`, `7.13.1`, `7.14.2`, `7.15.2`, `7.16.3`, `7.17.4`, `8.0.1`, `8.1.3`, `8.2.2`.

## Memcached Argument Reference

In addition to the common arguments for all services [described above](#argument-reference),
the `memcached` block can contain the following arguments:

* `class` - (Optional) The service class. Valid value is `cacher`. Defaults to `cacher`.
* `logging` - (Optional) The logging settings for the service. The structure of this block is [described below](#logging).
* `monitoring` - (Optional) The monitoring settings for the service. The structure of this block is [described below](#monitoring).

## PostgreSQL Argument Reference

In addition to the common arguments for all services [described above](#argument-reference),
the `pgsql` block can contain the following arguments:

* `autovacuum` - (Optional) Indicates whether the server must run the autovacuum launcher daemon.
  Valid values are `ON`, `OFF`. Defaults to `ON`.
* `autovacuum_max_workers` - (Optional) The maximum number of autovacuum processes (other than the autovacuum launcher)
  that can be running simultaneously. Valid values are from 1 to 262143. Defaults to `3`.
* `autovacuum_vacuum_cost_delay` - (Optional) The cost delay value in milliseconds used in automatic `VACUUM` operations.
  Valid values are `-1`, from 1 to 100.
* `autovacuum_vacuum_cost_limit` - (Optional) The cost limit value used in automatic `VACUUM` operations.
  Valid values are `-1`, from 1 to 10000. Defaults to `-1`.
* `autovacuum_analyze_scale_factor` - (Optional) The fraction of the table size to add to `autovacuum_analyze_threshold`
  when deciding whether to trigger an `ANALYZE`. Valid values are from 0 to 100. Defaults to `0.1`.
* `autovacuum_vacuum_scale_factor` - (Optional) The fraction of the table size to add to `autovacuum_vacuum_threshold`
  when deciding whether to trigger a `VACUUM`. Valid values are from 0 to 100. Defaults to `0.2`.
* `class` - (Optional) The service class. Valid value is `database`. Defaults to `database`.
* `database` - (Optional) List of PostgreSQL databases with parameters. The maximum number of databases is 1000. The structure of this block is [described below](#postgresql-database).
* `effective_cache_size` - (Optional) The planner’s assumption about the effective size of the disk cache
  that is available to a single query. Valid values are from 1 to 2147483647. Defaults to `524288`.
* `effective_io_concurrency` -  (Optional) The number of concurrent disk I/O operations. Valid values are from 0 to 1000. Defaults to `1`.
* `logging` - (Optional) The logging settings for the service. The structure of this block is [described below](#logging).
* `maintenance_work_mem` -  (Optional) The maximum amount of memory in bytes (multiple of 1 KiB) used by maintenance operations,
  such as `VACUUM`, `CREATE INDEX`, and `ALTER TABLE ADD FOREIGN KEY`.
  Valid values are from 1 MiB to 2 GiB. Defaults to `67108864` (64 MiB).
* `max_connections` - (Optional) The maximum number of simultaneous connections to the database server.
  Valid values are from 1 to 262143. Defaults to `100`.
* `max_wal_size` - (Optional) The maximum size in bytes (multiple of 1 MiB) that WAL can reach at automatic checkpoints.
  Valid values are from 2 to 2147483647 MiB. Defaults to `83886080` (80 MiB).
* `max_parallel_maintenance_workers` - (Optional) The maximum number of parallel workers that a single utility command can start.
  This parameter is relevant only for PostgreSQL versions 11 and higher. Valid values are from 0 to 1024.
* `max_parallel_workers` - (Optional) The maximum number of workers that the system can support for parallel operations.
* `max_parallel_workers_per_gather` - (Optional) The maximum number of workers that a single _Gather_ node can start.
  Valid values are from 0 to 1024. Defaults to `2`.
* `max_worker_processes` - (Optional) The maximum number of background processes that the system can support.
  Valid values are from 0 to 262143. Defaults to `8`.
* `min_wal_size` - (Optional) The minimum size in bytes (multiple of 1 MiB) to shrink the WAL to. As long as WAL disk usage stays below this setting,
  old WAL files are always recycled for future use at a checkpoint, rather than removed.
  Valid values are from 32 to 2147483647 MiB. Defaults to `83886080` (80 MiB).
* `monitoring` - (Optional) The monitoring settings for the service. The structure of this block is [described below](#monitoring).
* `options` - (Optional) Map containing other PostgreSQL parameters.
  Parameter names must be in camelCase. Values are strings.

~> If the parameter name includes a dot, then it cannot be passed in the `options`.
If you need to use such a parameter, contact [technical support].

* `replication_mode` - (Optional) The replication mode in the _Patroni_ cluster.
  The parameter must be set if `high_availability` is `true`. Valid values are `asynchronous`, `synchronous`, `synchronous_strict`.
* `shared_buffers` - (Optional) The amount of memory in 8 KiB pages the database server uses for shared memory buffers.
  Valid values are from 16 to 1073741823. Defaults to `1024`.
* `user` - (Optional) List of PostgreSQL users with parameters. The maximum number of users is 1000.
  The structure of this block is [described below](#postgresql-user).
* `version` - (Required) The version to install. Valid values are `10.21`, `11.16`, `12.11`, `13.7`, `14.4`, `15.2`.
* `wal_buffers` - (Optional) The amount of shared memory in 8 KiB pages used for WAL data not yet written to a volume.
  Valid values are from 8 to 262143.
* `wal_keep_segments` - (Optional) The minimum number of log files segments that must be kept in the _pg_xlog_ directory,
  in case a standby server needs to fetch them for streaming replication.
  This parameter is relevant only for PostgreSQL versions 10, 11, 12. Valid values are from 0 to 2147483647.
* `work_mem` - (Optional) The base maximum amount of memory in bytes (multiple of 1 KiB) to be used by a query operation
  (such as a sort or hash table) before writing to temporary disk files.
  Valid values are from 64 to 2147483647 KiB. Defaults to `4194304` (4 MiB).

### PostgreSQL database

~> All the parameters in the `database` block are editable.

The `database` block has the following structure:

* `backup_enabled` - (Optional) Indicates whether backup is enabled for the database. Defaults to `false`.
* `backup_id` - (Optional) The database backup ID.
* `backup_db_name` - The name of a database from the backup specified in the `backup_id` parameter.
* `encoding` - (Optional) The database encoding. Defaults to `UTF8`.
* `extensions` - (Optional) List of extensions for the database. Valid values are
  `address_standardizer`, `address_standardizer_data_us`, `amcheck`, `autoinc`, `bloom`, `btree_gin`, `btree_gist`,
  `citext`, `cube`, `dblink`, `dict_int`, `dict_xsyn`, `earthdistance`, `fuzzystrmatch`, `hstore`, `intarray`, `isn`,
  `lo`, `ltree`, `moddatetime`, `pg_buffercache`, `pg_trgm`, `pg_visibility `, `pgcrypto`, `pgrowlocks`, `pgstattuple`,
  `postgis`, `postgis_tiger_geocoder`, `postgis_topology`, `postgres_fdw`, `seg`, `tablefunc`, `tcn`, `timescaledb`,
  `tsm_system_rows`, `tsm_system_time`, `unaccent`, `uuid-ossp`, `xml2`.
* `locale` - (Optional) The database locale. Defaults to `ru_RU.UTF-8`.
* `name` - (Required) The database name.
* `owner` - (Required) The name of the user who is the database owner. This must be one of the existing users.
  Such a user cannot be deleted as long as it is the database owner.
* `user` - (Optional) List of PostgreSQL users with parameters. The maximum number of databases is 1000.
  The structure of this block is [described below](#postgresql-database-user).

### PostgreSQL database user

~> All the parameters in the `user` block are editable.

The `user` block has the following structure:

* `name` - (Required) The PostgreSQL user name.

### PostgreSQL user

~> All the parameters in the `user` block are editable.

The `user` block has the following structure:

* `name` - (Required) The PostgreSQL user name.
* `password` - (Required) The PostgreSQL user password.
  The value must be 8 to 128 characters long and must not contain `'`, `"`,  `` ` `` and `\`.

## Redis Argument Reference

In addition to the common arguments for all services [described above](#argument-reference),
the `redis` block can contain the following arguments:

* `class` - (Optional) The service class. Valid values are `cacher`, `database`. Defaults to `cacher`.
* `cluster_type` - (Optional) The clustering option. Valid values are `native`, `sentinel`.
  The parameter must be set if `high_availability` is `true`.
* `databases` - (Optional) The number of databases. Valid values are from 1 to 2147483647. Defaults to `16`.
* `logging` - (Optional) The logging settings for the service. The structure of this block is [described below](#logging).
* `maxmemory_policy` - (Optional) The memory management mode.
  Valid values are `noeviction`, `allkeys-lru`, `allkeys-lfu`, `volatile-lru`, `volatile-lfu`, `allkeys-random`, `volatile-random`, `volatile-ttl`.
  Defaults to `noeviction`.
* `monitoring` - (Optional) The monitoring settings for the service. The structure of this block is [described below](#monitoring).
* `options` - (Optional) Map containing other Redis parameters.
  Parameter names must be in camelCase. Values are strings.

~> If the parameter name includes a dot, then it cannot be passed in the `options`.
If you need to use such a parameter, contact [technical support].

* `password` - (Optional) The Redis user password.
  The value must be 8 to 128 characters long and must not contain `'`, `"`,  `` ` `` and `\`.
* `persistence_aof` - (Optional) Indicates whether AOF storage mode is enabled. Defaults to `false`.
* `persistence_rdb` - (Optional) Indicates whether RDB storage mode is enabled. Defaults to `false`.
* `timeout` - (Optional) The time in seconds for which connection to an inactive client is maintained.
  Valid values are from 0 to 2147483647. Defaults to `0`.
* `tcp_backlog` - (Optional) The size of a connection queue. Valid values are from 1 to 4096. Defaults to `511`.
* `tcp_keepalive` - (Optional) The time in seconds for which the service sends ACKs to detect dead peers
  (clients that cannot be reached even if they look connected). The value must be non-negative. Defaults to `300`.
* `version` - (Required) The version to install. Valid values are `5.0.14`, `6.2.6`, `7.0.11`.

## Common Service Argument Reference

### logging

~> All the parameters in the `logging` block are editable.

The `logging` block has the following structure:

* `log_to` - (Required) The ID of the logging service. It must run in the same VPC as the service.
* `logging_tags` - (Optional) List of tags that are assigned to the log records of the service.
  Each value in the list must be 1 to 256 characters long.

### monitoring

~> All the parameters in the `monitoring` block are editable.

The `monitoring` block has the following structure:

* `monitor_by` - (Required) The ID of the monitoring service. It must run in the same VPC as the service.
* `monitoring_labels` - (Optional) Map containing labels that are assigned to the metrics of the service.
  Keys must be 1 to 64 characters long.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `auto_created_security_group_ids` - List of security group IDs that CROC Cloud created for the service.
* `endpoints` - List of endpoints for connecting to the service.
* `error_code` - The service error code.
* `error_description` - The detailed description of the service error.
* `id` - The ID of the PaaS service.
* `instances` - List of instances that refers to the service. The structure of this block is [described below](#instances).
* `service_class` - The service class. The value matches `class` parameter of the specified block with service parameters.
* `service_type` - The service type. The value matches the name of the specified block with service parameters.
* `status` - The current status of the service.
* `supported_features` - List of service features.
* `total_cpu_count` - Total number of CPU cores in use.
* `total_memory` - Total RAM in use in MiB.

For `backup_settings` the following attribute is also exported:

* `user_id` - The ID of the user whose login is set to `backup_settings.user_login`.

For `*.database` the following attribute is also exported:

* `id` - The ID of the database.

For `*.user` the following attribute is also exported:

* `id` - The ID of the user.

### instances

* `endpoint` - The service endpoint on the instance.
* `index` - The instance index.
* `instance_id` - The ID of the instance.
* `interface_id` - The ID of the instance network interface.
* `name` - The instance name.
* `private_ip` - The private IP address of the instance.
* `role` - The instance role.
* `status` - The current status of the instance.

## Timeouts

`aws_paas_service` provides the following [Timeouts][timeouts] configuration options:

* `create` - (Default `30 minutes`) How long to wait for the service to be created.
* `update` - (Default `60 minutes`) How long to wait for the service to be updated.
* `delete` - (Default `15 minutes`) How long to wait for the service to be deleted.

## Import

PaaS service can be imported using `id`, e.g.,

```
$ terraform import aws_paas_service.example fm-cluster-12345678
```
