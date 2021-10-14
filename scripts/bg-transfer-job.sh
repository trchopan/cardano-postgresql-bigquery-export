# Manually run this to create a BQ schedule job
# Run every day at 00:00 => Update every 24 hour

TABLE="tx"

bq mk --transfer_config \
  --target_dataset=cardano_mainnet \
  --display_name="cardano-$TABLE" \
  --params='{
  "data_path_template":"gs://cardano-bucket/mainnet-export/tx/*",
  "destination_table_name_template":"tx",
  "file_format":"CSV",
  "max_bad_records":"0",
  "ignore_unknown_values":"true",
  "field_delimiter":",",
  "skip_leading_rows":"0",
  "write_disposition": "MIRROR",
  "allow_quoted_newlines":"true",
  "allow_jagged_rows":"true",
  "delete_source_files":"false"
  }' \
  --data_source=google_cloud_storage \
  --schedule="every day 00:00"
