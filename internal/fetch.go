package product

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

func fetchPriceFile(url *url.URL) (io.Reader, error) {
	resp, err := http.Get(url.String())
	if err != nil {
		return nil, fmt.Errorf("http requets: %v", err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, fmt.Errorf("read response: %v", err)
	}

	return bytes.NewReader(body), nil
}

type csvData [][]string

func parseCSV(in io.Reader) (csvData, error) { //io.Reader is easily mockable for testing
	reader := csv.NewReader(in)
	reader.Comma = ';'

	data, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	return data, nil
}

var ErrBadCSVData = errors.New("expected two values per row")

func mapToProducts(rows csvData) (products []CSVRow, err error) {
	for i := 1; i < len(rows); i++ { //skip header
		row := rows[i]

		if len(row) != 2 {
			return nil, ErrBadCSVData
		}

		products = append(products, CSVRow{
			Name:  row[0],
			Price: row[1],
		})
	}

	return
}
