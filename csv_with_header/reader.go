package csv_with_header

import (
  "encoding/csv"
  "io"
  "os"
)

type Reader struct {
  raw     *csv.Reader
  Columns []string
}

func NewReader(r io.Reader) *Reader {
  raw_reader := csv.NewReader(r)

  columns, err := raw_reader.Read()
  if err != nil {
    panic(err)
  }

  return &Reader{raw_reader, columns}
}

func (reader *Reader) Read() (map[string]string, error) {
  raw_record, err := reader.raw.Read()
  if err != nil {
    return nil, err
  }

  record := make(map[string]string, len(reader.Columns))
  for i, column := range reader.Columns {
    record[column] = raw_record[i]
  }

  return record, nil
}

func ReadInto(filepath string, into chan<- map[string]string) {
  file, err := os.Open(filepath)
  defer file.Close()

  if err != nil {
    panic(err)
  }

  reader := NewReader(file)

  for {
    record, err := reader.Read()
    if err == io.EOF {
      break
    }
    if err != nil {
      panic(err)
    }

    into <- record
  }
  close(into)
}

func ReadAll(filepath string) []map[string]string {
  record_chan := make(chan map[string]string)
  go ReadInto(filepath, record_chan)
  // read chan to end... maybe not the most efficient channel use case
  records := make([]map[string]string, 0)
  for record := range record_chan {
    records = append(records, record)
  }
  return records
}
