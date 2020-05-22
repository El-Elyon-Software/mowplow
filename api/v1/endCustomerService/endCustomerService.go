package endCustomerService

import (
	"net/http"
	"strconv"
	"strings"
	"../dal"
	e "../errors"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/render"
	"github.com/pkg/errors"
	"time"
)

type EndCustomerService struct {
	ID           				int64  	`json:"id"`
	EndCustomerID				int64	`json:"endCustomerID"`
	ServiceProviderServiceID	int64	`json:"serviceProviderServiceID"`
	Description  				string 	`json:"description"`
	EstimatedJobLength			float64	`json:"estimatedJobLength"`
	ContractStartDate			string 	`json:"contractStartDate"`
	ContractEndDate				string 	`json:"contractEndDate"`
	DateAdded    				string 	`json:"dateAdded"`
	DateModified 				string 	`json:"dateModified"`
	dal          				dal.DB
}

type GeneralResponse struct {
	MSG string `json:"msg"`
	ID  int64  `json:"id"`
}

func New() *EndCustomerService {
	db := dal.DB{
		DBType:     "",
		DBName:     "",
		DBUser:     "",
		DBPassword: ""}
	db.NewDB()
	ecs := EndCustomerService{dal: db}
	return &ecs
}

func (ecs *EndCustomerService) Routes() *chi.Mux {
	router := chi.NewRouter()
	router.Post("/", ecs.create)
	router.Get("/{ID}", ecs.read)
	router.Get("/", ecs.readFilter)
	router.Put("/{ID}", ecs.update)
	router.Delete("/{ID}", ecs.delete)
	router.Patch("/{ID}", ecs.update)
	return router
}

func (ecs *EndCustomerService) create(rw http.ResponseWriter, r *http.Request) {
	if ecs.bindData(rw, r) != nil {
		return
	}

	err := ecs.dal.OpenDB()
	if err != nil {
		e.HandleError(rw, r, err)
		return
	}
	defer ecs.dal.CloseDB()

	stmt := `INSERT INTO end_customer_service (
				end_customer_id
				,service_provider_service_id
				,description
				,estimated_job_length
				,contract_start_date
				,contract_end_date) 
			VALUES (?, ?, ?, ?, ?, ?);`
	
	res, err := ecs.dal.DB.Exec(stmt, ecs.EndCustomerID, ecs.ServiceProviderServiceID, 
		ecs.Description, ecs.EstimatedJobLength, ecs.ContractStartDate, ecs.ContractEndDate)

	if err != nil {
		e.HandleError(rw, r, err)
		return
	}

	id, err := res.LastInsertId()
	if err != nil {
		e.HandleError(rw, r, err)
		return
	}

	if id > 0 {
		ecs.ID = id
	}
	render.JSON(rw, r, ecs)
}

func (ecs *EndCustomerService) read(rw http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "ID"), 10, 64)
	if err != nil {
		e.HandleError(rw, r, err)
		return
	}

	err = ecs.dal.OpenDB()
	if err != nil {
		e.HandleError(rw, r, err)
		return
	}
	defer ecs.dal.CloseDB()

	stmt := `SELECT 
				end_customer_service_id
				,end_customer_id
				,service_provider_service_id
				,description
				,estimated_job_length
				,contract_start_date
				,contract_end_date
				,date_added
				,date_modified
			FROM 
				end_customer_service 
			WHERE 
				end_customer_service_id=?
				AND date_deleted IS NULL;`

	err = ecs.dal.DB.QueryRow(stmt, id).Scan(
		&ecs.ID, &ecs.EndCustomerID, &ecs.ServiceProviderServiceID, 
		&ecs.Description, &ecs.EstimatedJobLength, &ecs.ContractStartDate,
		&ecs.ContractEndDate, &ecs.DateAdded, &ecs.DateModified)

	if err != nil {
		e.HandleError(rw, r, err)
		return
	}

	render.JSON(rw, r, ecs)
}

func (ecs *EndCustomerService) readFilter(rw http.ResponseWriter, r *http.Request) {
	wc, vals := dal.ParseQueryStringParams(r.URL.RawQuery)

	stmt := `SELECT 
				end_customer_service_id
				,end_customer_id
				,service_provider_service_id
				,description
				,estimated_job_length
				,contract_start_date
				,contract_end_date
				,date_added
				,date_modified
			FROM 
				end_customer_service  
			WHERE 
				` + strings.Join(wc, " ") + `;`

	err := ecs.dal.OpenDB()
	if err != nil {
		e.HandleError(rw, r, err)
		return
	}
	defer ecs.dal.CloseDB()

	rows, err := ecs.dal.DB.Query(stmt, vals...)
	if err != nil {
		e.HandleError(rw, r, err)
		return
	}
	defer rows.Close()

	var ecss []*EndCustomerService
	for rows.Next() {
		var e EndCustomerService
		rows.Scan(&e.ID, &e.EndCustomerID, &e.ServiceProviderServiceID, 
			&e.Description, &e.EstimatedJobLength, &e.ContractStartDate,
			&e.ContractEndDate, &e.DateAdded, &e.DateModified)
		ecss = append(ecss, &e)
	}

	render.JSON(rw, r, ecss)
}

func (ecs *EndCustomerService) update(rw http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "ID"), 10, 64)
	if err != nil {
		e.HandleError(rw, r, err)
		return
	}
	ecs.ID = id

	if ecs.bindData(rw, r) != nil {
		return
	}

	err = ecs.dal.OpenDB()
	if err != nil {
		e.HandleError(rw, r, err)
		return
	}
	defer ecs.dal.CloseDB()

	stmt := `UPDATE 
				end_customer_service 
			SET
				end_customer_id=?
				,service_provider_service_id=?
				,description=?
				,estimated_job_length=?
				,contract_start_date=?
				,contract_end_date=?
				,date_modified=NOW()
			WHERE
				end_customer_service_id=?;`

	_, err = ecs.dal.DB.Exec(stmt, ecs.EndCustomerID, ecs.ServiceProviderServiceID, 
		ecs.Description, ecs.EstimatedJobLength, ecs.ContractStartDate,
		ecs.ContractEndDate, ecs.ID)

	if err != nil {
		e.HandleError(rw, r, err)
		return
	}

	ecs.read(rw, r)
}

func (ecs *EndCustomerService) delete(rw http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "ID"), 10, 64)
	if err != nil {
		e.HandleError(rw, r, err)
		return
	}

	err = ecs.dal.OpenDB()
	if err != nil {
		e.HandleError(rw, r, err)
		return
	}
	defer ecs.dal.CloseDB()

	stmt := `UPDATE 
				service_provider_service 
			SET
				date_deleted=NOW()
			WHERE
				service_id=?;`

	_, err = ecs.dal.DB.Exec(stmt, id)

	if err != nil {
		e.HandleError(rw, r, err)
		return
	}

	gr := &GeneralResponse{
		MSG: "The end customer was deleted",
		ID:  id,
	}
	render.JSON(rw, r, gr)
}

func (ecs *EndCustomerService) bindData(rw http.ResponseWriter, r *http.Request) error {
	if err := render.Bind(r, ecs); err != nil {
		e.HandleError(rw, r, err)
		return err
	}

	return nil
}

func (ecs *EndCustomerService) Bind(r *http.Request) error {
	if ecs.EndCustomerID <= 0 {
		return errors.New("endCustomerID is required.")
	}

	if ecs.ServiceProviderServiceID <= 0 {
		return errors.New("serviceProviderServiceID is required.")
	}

	if ecs.Description == "" || len(ecs.Description) < 1 {
		return errors.New("description is required and must be at least one characters.")
	}

	if ecs.EstimatedJobLength <= 0.0 {
		return errors.New("estimatedJobLength is required.")
	}
	
	var err error
	_, err = time.Parse("2006-01-02", ecs.ContractStartDate)
	if err != nil {
		return errors.New("contractStartDate must be a valid data in the YYYY-MM-DD format.")
	}

	_, err = time.Parse("2006-01-02", ecs.ContractEndDate)
	if err != nil {
		return errors.New("contractEndDate must be a valid data in the YYYY-MM-DD format.")
	}

	return nil
}