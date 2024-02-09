package repository

import (
	"context"
)

func (r *Repository) InsertUser(ctx context.Context, input InsertUser) (int, error) {
	var id int
	err := r.Db.QueryRowContext(
		ctx,
		"SELECT * FROM insert_user( $1, $2, $3 )",
		input.PhoneNumber, input.Name, input.Password,
	).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *Repository) GetUserByPhoneNumber(ctx context.Context, phoneNumber string) (User, error) {
	var user User
	err := r.Db.QueryRowContext(
		ctx,
		"SELECT * FROM get_user_by_phone_number( $1 )",
		phoneNumber,
	).Scan(&user.ID, &user.Name, &user.PhoneNumber, &user.Password)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (r *Repository) InsertUserLoginHistory(ctx context.Context, input InsertUserLoginHistory) error {
	if _, err := r.Db.ExecContext(
		ctx,
		"SELECT * FROM insert_user_login_history( $1, $2, $3, $4)",
		input.ID, input.PhoneNumber, input.Password, input.Name,
	); err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetUserByID(ctx context.Context, id int) (User, error) {
	var user User
	err := r.Db.QueryRowContext(
		ctx,
		"SELECT * FROM get_user_by_id( $1 )",
		id,
	).Scan(&user.ID, &user.Name, &user.PhoneNumber, &user.Password)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (r *Repository) UpdateUser(ctx context.Context, input UpdateUser) (User, error) {
	var user User
	err := r.Db.QueryRowContext(
		ctx,
		"SELECT * FROM update_user( $1, $2, $3 )",
		input.ID, input.Name, input.PhoneNumber,
	).Scan(&user.ID, &user.Name, &user.PhoneNumber, &user.Password)
	if err != nil {
		return user, err
	}
	return user, nil
}
