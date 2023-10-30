package daos

import (
	"database/sql"
	"errors"
	"github.com/mahendraintelops/home-automation-solution-v2/device-service/pkg/rest/server/daos/clients/sqls"
	"github.com/mahendraintelops/home-automation-solution-v2/device-service/pkg/rest/server/models"
	log "github.com/sirupsen/logrus"
)

type DeviceDao struct {
	sqlClient *sqls.SQLiteClient
}

func migrateDevices(r *sqls.SQLiteClient) error {
	query := `
	CREATE TABLE IF NOT EXISTS devices(
		Id INTEGER PRIMARY KEY AUTOINCREMENT,
        
		Name TEXT NOT NULL,
		Volume INTEGER NOT NULL,
        CONSTRAINT id_unique_key UNIQUE (Id)
	)
	`
	_, err1 := r.DB.Exec(query)
	return err1
}

func NewDeviceDao() (*DeviceDao, error) {
	sqlClient, err := sqls.InitSqliteDB()
	if err != nil {
		return nil, err
	}
	err = migrateDevices(sqlClient)
	if err != nil {
		return nil, err
	}
	return &DeviceDao{
		sqlClient,
	}, nil
}

func (deviceDao *DeviceDao) CreateDevice(m *models.Device) (*models.Device, error) {
	insertQuery := "INSERT INTO devices(Name, Volume)values(?, ?)"
	res, err := deviceDao.sqlClient.DB.Exec(insertQuery, m.Name, m.Volume)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	m.Id = id

	log.Debugf("device created")
	return m, nil
}

func (deviceDao *DeviceDao) ListDevices() ([]*models.Device, error) {
	selectQuery := "SELECT * FROM devices"
	rows, err := deviceDao.sqlClient.DB.Query(selectQuery)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)
	var devices []*models.Device
	for rows.Next() {
		m := models.Device{}
		if err = rows.Scan(&m.Id, &m.Name, &m.Volume); err != nil {
			return nil, err
		}
		devices = append(devices, &m)
	}
	if devices == nil {
		devices = []*models.Device{}
	}

	log.Debugf("device listed")
	return devices, nil
}

func (deviceDao *DeviceDao) GetDevice(id int64) (*models.Device, error) {
	selectQuery := "SELECT * FROM devices WHERE Id = ?"
	row := deviceDao.sqlClient.DB.QueryRow(selectQuery, id)
	m := models.Device{}
	if err := row.Scan(&m.Id, &m.Name, &m.Volume); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sqls.ErrNotExists
		}
		return nil, err
	}

	log.Debugf("device retrieved")
	return &m, nil
}
