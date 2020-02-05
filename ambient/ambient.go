package ambient

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

type xData struct {
	index int
	value float32
}

type Data struct {
	values  [9]xData
	created utcTime
}

func NewData() *Data {
	return &Data{
		values:  [9]xData{},
		created: utcTime{},
	}
}

func (d *Data) Set(i int, v float32) error {
	if i < 1 || i > 8 {
		return errors.New("max index is over")
	}
	d.values[i].index = i
	d.values[i].value = v
	return nil
}

func (d *Data) SetTime(t time.Time) error {
	d.created = utcTime{t}
	return nil
}

func (d *Data) MarshalJSON() ([]byte, error) {
	s := []string{}
	for _, v := range d.values {
		if v.index != 0 {
			s = append(s, fmt.Sprintf("\"d%d\":%g", v.index, v.value))
		}
	}
	if (d.created != utcTime{time.Time{}}) {
		j, err := d.created.MarshalJSON()
		if err != nil {
			return []byte{}, err
		}
		s = append(s, fmt.Sprintf("\"created\":%s", string(j)))
	}
	ss := fmt.Sprintf("%s", strings.Join(s, ","))
	return []byte(fmt.Sprintf("{%s}", ss)), nil
}

type sData struct {
	WriteKey string  `json:"writeKey"`
	Data     []*Data `json:"data"`
}

func SendData(c *Client, d *Data) error {
	return SendBulkData(c, []*Data{d})
}

func SendBulkData(c *Client, ds []*Data) error {
	s := sData{
		WriteKey: c.config.WriteKey,
		Data:     ds,
	}
	b, err := json.Marshal(s)
	if err != nil {
		return err
	}
	// fmt.Println(string(b))
	if err := c.CreateData(context.Background(), b); err != nil {
		return err
	}
	return nil
}
