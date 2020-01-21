package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"

	redis "github.com/go-redis/redis/v7"
)

func writeKey(ctx context.Context, cfg *Config, client redis.UniversalClient) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		b, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			return fmt.Errorf("Read failed: %s", err.Error())
		}

		if _, err := client.Set(cfg.SetKey, b, 0).Result(); err != nil {
			return fmt.Errorf("SET failed: %s", err.Error())
		}

		break
	}

	return nil
}

func produceBlob(ctx context.Context, cfg *Config, client redis.UniversalClient) error {

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		b, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			return fmt.Errorf("Read failed: %s", err.Error())
		}

		xAddArgs := &redis.XAddArgs{
			Stream: cfg.Stream,
			Values: map[string]interface{}{
				cfg.FieldName: b,
			},
			ID: cfg.Id,
		}

		if _, err := client.XAdd(xAddArgs).Result(); err != nil {
			return fmt.Errorf("XADD failed; values = %s; error = %s", xAddArgs.Values, err.Error())
		}

		break
	}

	return nil
}

func produceCSV(ctx context.Context, cfg *Config, client redis.UniversalClient) error {

	r := csv.NewReader(os.Stdin)
	r.ReuseRecord = true

	// read header
	header, err := readHeader(r)
	if err != nil {
		return err
	}

	xAddArgs := &redis.XAddArgs{
		Stream: cfg.Stream,
		Values: make(map[string]interface{}, len(header)),
	}

	lineNum := 1
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("Read failed: %s", err.Error())
		}

		setValues(xAddArgs.Values, header, record)
		if cfg.Id == "linenum" {
			xAddArgs.ID = "0-" + strconv.Itoa(lineNum)
			lineNum++
		}

		id, err := client.XAdd(xAddArgs).Result()
		if err != nil {
			return fmt.Errorf("XADD failed; values = %s; error = %s", xAddArgs.Values, err.Error())
		} else if !cfg.Silent {
			fmt.Println(id)
		}
	}

	return nil
}

func readHeader(r *csv.Reader) ([]string, error) {

	record, err := r.Read()
	if err != nil {
		return nil, fmt.Errorf("Read failed: %s", err.Error())
	}
	header := make([]string, len(record))
	copy(header, record)

	return header, nil
}

func setValues(m map[string]interface{}, header []string, record []string) {
	for i := 0; i < len(header); i++ {
		m[header[i]] = record[i]
	}
}
