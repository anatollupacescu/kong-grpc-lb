package product

import "time"

func transform(in []CSVRow, date time.Time) []productDTO {
	type countingProduct struct {
		product *CSVRow
		count   int
	}

	var (
		unique = make(map[string]countingProduct)
		order  []string
	)

	for i := range in {
		p := in[i]
		if cp, found := unique[p.Name]; found {
			cp.count++
			cp.product = &p
			unique[p.Name] = cp
			continue
		}

		unique[p.Name] = countingProduct{product: &p}
		order = append(order, p.Name)
	}

	out := make([]productDTO, 0, len(order))
	for _, name := range order {
		v := unique[name]
		record := productDTO{
			LastUpdated: date,
			UpdateCount: v.count,
			Name:        v.product.Name,
			Price:       v.product.Price,
		}

		out = append(out, record)
	}

	return out
}
