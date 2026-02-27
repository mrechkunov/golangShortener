package repository

import (
	"bufio"
	"encoding/json"
	"io"
	"os"

	"github.com/mrechkunov/golangShortener.git/internal/model"
)

type Producer struct {
	file *os.File
	// добавляем Writer в Producer
	writer *bufio.Writer
}

func NewProducer(filename string) (*Producer, error) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	return &Producer{
		file: file,
		// создаём новый Writer
		writer: bufio.NewWriter(file),
	}, nil
}

func (p *Producer) WriteEvent(event *model.Event) error {
	data, err := json.Marshal(&event)
	if err != nil {
		return err
	}

	// записываем событие в буфер
	if _, err := p.writer.Write(data); err != nil {
		return err
	}

	// добавляем перенос строки
	if err := p.writer.WriteByte('\n'); err != nil {
		return err
	}

	// записываем буфер в файл
	return p.writer.Flush()
}

func (p *Producer) WriteEvents(events *[]model.Event) error {
	data, err := json.MarshalIndent(events, "", "")
	if err != nil {
		return err
	}

	// записываем событие в буфер
	if _, err := p.writer.Write(data); err != nil {
		return err
	}

	// записываем буфер в файл
	return p.writer.Flush()
}

func (p *Producer) Close() {
	p.file.Close()
}

type Consumer struct {
	file *os.File
	// добавляем reader в Consumer
	reader *bufio.Reader
}

func NewConsumer(filename string) (*Consumer, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		file: file,
		// создаём новый Reader
		reader: bufio.NewReader(file),
	}, nil
}

func (c *Consumer) ReadEvent() (*model.Event, error) {
	// читаем данные до символа переноса строки
	data, err := c.reader.ReadBytes('\n')
	if err != nil {
		return nil, err
	}

	// преобразуем данные из JSON-представления в структуру
	event := model.Event{}
	err = json.Unmarshal(data, &event)
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func (c *Consumer) ReadEvents() (*[]model.Event, error) {
	// читаем данные до символа переноса строки
	data, err := io.ReadAll(c.reader)
	if err != nil {
		return nil, err
	}

	// преобразуем данные из JSON-представления в структуру
	var events []model.Event
	err = json.Unmarshal(data, &events)
	if err != nil {
		return nil, err
	}
	return &events, nil
}

func (c *Consumer) Close() {
	c.file.Close()
}
