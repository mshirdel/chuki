package config

var builtinConfig = `
logging:
  level: "debug"
database:
  host: "localhost"
  port: 3306
  user: "chuki"
  password: "chuki"
  dbname: "chuki"
  charset: "utf8mb4"
  collation: "utf8mb4_unicode_ci"
  parse_time: true
  location: "Asia/Tehran"
  max_life_time: 5m
  max_idle_time: 0s
  max_open_connections: 10
  max_idle_connections: 5
  skip_initialize_with_version: false
  logger:
    slow_threshold: 200ms
    level: info
    colorful: true
    ignore_record_not_found_error: true
`
