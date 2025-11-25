package repository

import (
	"database/sql"
	"uas-pbe-praksem5/app/model"
)

type UserRepo struct {
	DB *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo { return &UserRepo{DB: db} }

func (r *UserRepo) Create(req model.CreateUserRequest, passwordHash string) error {
	_, err := r.DB.Exec(`
		INSERT INTO users (username,email,password_hash,full_name,role_id,created_at,updated_at)
		VALUES ($1,$2,$3,$4,$5,NOW(),NOW())`,
		req.Username, req.Email, passwordHash, req.FullName, req.RoleID)
	return err
}

func (r *UserRepo) GetByUsernameOrEmail(q string) (*model.User, string, error) {
	row := r.DB.QueryRow(`
SELECT u.id, u.username, u.email, u.full_name, u.role_id, r.name, u.is_active, u.created_at, u.updated_at, u.password_hash
FROM users u LEFT JOIN roles r ON u.role_id = r.id
WHERE u.username = $1 OR u.email = $1`, q)

	var u model.User
	var pwd string
	err := row.Scan(&u.ID, &u.Username, &u.Email, &u.FullName, &u.RoleID, &u.RoleName, &u.IsActive, &u.CreatedAt, &u.UpdatedAt, &pwd)
	if err == sql.ErrNoRows {
		return nil, "", nil
	}
	if err != nil {
		return nil, "", err
	}
	return &u, pwd, nil
}

func (r *UserRepo) GetByID(id string) (*model.User, error) {
	row := r.DB.QueryRow(`
SELECT u.id, u.username, u.email, u.full_name, u.role_id, r.name, u.is_active, u.created_at, u.updated_at
FROM users u LEFT JOIN roles r ON u.role_id = r.id
WHERE u.id = $1`, id)

	var u model.User
	err := row.Scan(&u.ID, &u.Username, &u.Email, &u.FullName, &u.RoleID, &u.RoleName, &u.IsActive, &u.CreatedAt, &u.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepo) GetAll() ([]model.User, error) {
	rows, err := r.DB.Query(`
SELECT u.id,u.username,u.email,u.full_name,u.role_id,r.name,u.is_active,u.created_at,u.updated_at
FROM users u LEFT JOIN roles r ON u.role_id = r.id ORDER BY u.created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []model.User
	for rows.Next() {
		var u model.User
		rows.Scan(&u.ID, &u.Username, &u.Email, &u.FullName, &u.RoleID, &u.RoleName, &u.IsActive, &u.CreatedAt, &u.UpdatedAt)
		res = append(res, u)
	}
	return res, nil
}

func (r *UserRepo) Update(id string, req model.CreateUserRequest) error {
	_, err := r.DB.Exec(`UPDATE users SET username=$1, email=$2, full_name=$3, role_id=$4, updated_at=NOW() WHERE id=$5`,
		req.Username, req.Email, req.FullName, req.RoleID, id)
	return err
}

func (r *UserRepo) Delete(id string) error {
	_, err := r.DB.Exec(`DELETE FROM users WHERE id=$1`, id)
	return err
}

func (r *UserRepo) UpdateRole(id string, roleID string) error {
	_, err := r.DB.Exec(`UPDATE users SET role_id=$1, updated_at=NOW() WHERE id=$2`, roleID, id)
	return err
}
