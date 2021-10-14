package repository

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os/exec"
	"strings"
)

func ExecCmd(c string) error {
	cmd := exec.Command("sh", "-c", c)
	stderr, _ := cmd.StderrPipe()
	stdout, _ := cmd.StdoutPipe()
	multi := io.MultiReader(stdout, stderr)
	scanner := bufio.NewScanner(multi)
	cmd.Start()
	for scanner.Scan() {
		m := scanner.Text()
		fmt.Println(m)
	}
	return cmd.Wait()
}

func makeBucketUrl(c Configuration, e ExportConfig, fileNameOrPattern string) string {
	return fmt.Sprintf("gs://%s/%s/%s/%s", c.GcsBucket, c.GcsPrefix, e.Table, fileNameOrPattern)
}

func ExportSQLToCSV(instance, dbName string, c Configuration, e ExportConfig, filename string, startId int64, endId int64) error {
	bucketUrl := makeBucketUrl(c, e, filename)

	fields := strings.Join(e.Fields, ",")
	query := fmt.Sprintf("SELECT id,%s FROM %s WHERE %d < id AND id <= %d",
		fields, e.Table, startId, endId,
	)

	cmd := fmt.Sprintf(`
        gcloud sql export csv %s %s \
            --database=%s \
            --query="%s"
    `, instance, bucketUrl, dbName, query)
	log.Println(cmd)
	// return ExecCmd(cmd)
	return nil
}

func BQTransferJobCreate(c Configuration, e ExportConfig) error {
	bucketPattern := makeBucketUrl(c, e, "*")
	cmd := fmt.Sprintf(`
    bq mk --transfer_config \
          --target_dataset=%s \
          --display_name="%s" \
          --params='{
              "data_path_template":"%s",
              "destination_table_name_template":"%s",
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
    `, c.BqDataset, c.GcsPrefix+"-"+e.Table, bucketPattern, e.Table)
	log.Println(cmd)
	// return ExecCmd(cmd)
	return nil
}

func BQTableCreateFromSchema(c Configuration, e ExportConfig) error {
    return nil
}
