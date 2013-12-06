package wmi

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"io"

	ole "github.com/mjibson/go-ole"
	"github.com/mjibson/go-ole/oleutil"
)

type drv struct{}

func (d *drv) Open(name string) (driver.Conn, error) {
	return Open(name)
}

func init() {
	sql.Register("wmi", &drv{})
}

type conn struct {
	service *ole.IDispatch
}

func (c *conn) Close() error {
	c.service.Release()
	return nil
}

var ErrUnsupported = errors.New("wmi: unsupported operation")

func (c *conn) Prepare(query string) (driver.Stmt, error) {
	fmt.Println("prepare", query)
	return nil, ErrUnsupported
}

func (c *conn) Begin() (driver.Tx, error) {
	fmt.Println("BEGIN")
	return nil, ErrUnsupported
}

func Open(name string) (driver.Conn, error) {
	unknown, err := oleutil.CreateObject("WbemScripting.SWbemLocator")
	if err != nil {
		return nil, err
	}
	defer unknown.Release()

	wmi, err := unknown.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		return nil, err
	}
	defer wmi.Release()

	// service is a SWbemServices
	serviceRaw, err := oleutil.CallMethod(wmi, "ConnectServer")
	if err != nil {
		return nil, err
	}
	service := serviceRaw.ToIDispatch()

	cn := &conn{service: service}
	return cn, nil
}

func (c *conn) Query(query string, args []driver.Value) (driver.Rows, error) {
	fmt.Println("QUERY", query, args)
	resultRaw, err := oleutil.CallMethod(c.service, "ExecQuery", query)
	if err != nil {
		return nil, err
	}
	r := rows{result: resultRaw.ToIDispatch()}
	count, err := oleInt64(r.result, "Count")
	if err != nil {
		return nil, err
	}
	r.count = count
	return &r, nil
}

type rows struct {
	result  *ole.IDispatch
	count   int64
	current int64
}

func (r *rows) Close() error {
	r.result.Release()
	return nil
}

func (r *rows) Next(dest []driver.Value) error {
	if r.current >= r.count {
		return io.EOF
	}
	return ErrUnsupported
}

func (r *rows) Columns() []string {
	itemRaw, err := oleutil.CallMethod(r.result, "ItemIndex", 0)
	if err != nil {
		l.Println(err)
		return nil
	}
	item := itemRaw.ToIDispatch()
	defer item.Release()

	propsRaw, err := oleutil.GetProperty(item, "Properties_")
	if err != nil {
		l.Println(err)
		return nil
	}
	props := propsRaw.ToIDispatch()
	defer props.Release()

	//x, err := props.GetTypeInfo()
	//fmt.Println(x, err)
	//return nil

	count, err := oleInt64(props, "Count")
	if err != nil {
		l.Println(err)
		return nil
	}
	fmt.Println("count props", count)
	cols := make([]string, count)
	for i := int64(0); i < count; i++ {
		propRaw, err := oleutil.CallMethod(props, "ItemIndex", 0)
		if err != nil {
			l.Println("prop raw", err)
			return nil
		}
		prop := propRaw.ToIDispatch()
		defer prop.Release()
		_ = prop
	}
	return cols
}
