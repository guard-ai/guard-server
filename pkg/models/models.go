package models

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/geojson"
)

type Point struct {
	*geom.Point
}

func NewPoint() Point {
	return Point{
		Point: geom.NewPoint(geom.XY).SetSRID(4326),
	}
}

func (p *Point) AsGeoJSON() string {
	return fmt.Sprintf(`{ "type": "Point", "coordinates": [%f, %f], "crs": {"type": "name", "properties": {"name": "EPSG:%d"}} }`, p.X(), p.Y(), p.SRID())
}

func (p *Point) Scan(src interface{}) error {
	if src == nil {
		return nil // handle null values
	}

	point, ok := src.(string)
	if !ok {
		return fmt.Errorf("expected string, got %T", src)
	}

	// handle null values
	if point == "" {
		return nil
	}

	var t geom.T
	geojson.Unmarshal([]byte(point), &t)
	p.Point = geom.NewPoint(geom.XY).SetSRID(4326)
	p.Point.SetCoords(t.FlatCoords())
	return nil
}

func (p *Point) UnmarshalJSON(b []byte) error {
	values := strings.Split(strings.Trim(string(b), `"`), ",")
	x, err := strconv.ParseFloat(values[0], 64)
	if err != nil {
		return err
	}
	y, err := strconv.ParseFloat(values[1], 64)
	if err != nil {
		return err
	}
	if p.Point == nil {
		p.Point = geom.NewPoint(geom.XY).SetSRID(4326)
	}

	p.Point.SetCoords(geom.Coord{x, y})
	return nil
}

func (p *Point) MarshalJSON() ([]byte, error) {
	if p.Point == nil {
		return []byte(`""`), nil
	}

	return []byte(fmt.Sprintf(`"%f,%f"`, p.X(), p.Y())), nil
}

type Log struct {
	Id        uuid.UUID          `json:"id"`
	Region    string             `json:"region"`
	Utterance string             `json:"utterance"`
	CreatedAt pgtype.Timestamptz `json:"created_at"`
}

type Event struct {
	Id          uuid.UUID          `json:"id"`
	Level       string             `json:"level"`
	Location    Point              `json:"location"`
	Category    string             `json:"category"`
	LogId       uuid.UUID          `json:"log_id"`
	CreatedAt   pgtype.Timestamptz `json:"created_at"`
	Description string             `json:"description"`
}

type User struct {
	Id        uuid.UUID          `json:"id"`
	LastPing  pgtype.Timestamptz `json:"last_ping"`
	CreatedAt pgtype.Timestamptz `json:"created_at"`
	Location  Point              `json:"location"`
	PushToken string             `json:"push_token"`
}
