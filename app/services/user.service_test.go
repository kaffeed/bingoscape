package services

import (
	"reflect"
	"testing"

	"github.com/kaffeed/bingoscape/app/db"
)

func TestNewUserServices(t *testing.T) {
	type args struct {
		store *db.Queries
	}
	tests := []struct {
		name string
		args args
		want *UserService
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewUserServices(tt.args.store); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUserServices() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserService_GetAllUsers(t *testing.T) {
	type fields struct {
		UserStore *db.Queries
	}
	tests := []struct {
		name    string
		fields  fields
		want    []db.Login
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			us := &UserService{
				UserStore: tt.fields.UserStore,
			}
			got, err := us.GetAllUsers()
			if (err != nil) != tt.wantErr {
				t.Errorf("UserService.GetAllUsers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserService.GetAllUsers() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserService_CreateUser(t *testing.T) {
	type fields struct {
		UserStore *db.Queries
	}
	type args struct {
		params db.CreateLoginParams
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			us := &UserService{
				UserStore: tt.fields.UserStore,
			}
			if err := us.CreateUser(tt.args.params); (err != nil) != tt.wantErr {
				t.Errorf("UserService.CreateUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserService_UpdatePassword(t *testing.T) {
	type fields struct {
		UserStore *db.Queries
	}
	type args struct {
		uid int32
		p   string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    db.Login
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			us := &UserService{
				UserStore: tt.fields.UserStore,
			}
			got, err := us.UpdatePassword(tt.args.uid, tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserService.UpdatePassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserService.UpdatePassword() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserService_DeleteUser(t *testing.T) {
	type fields struct {
		UserStore *db.Queries
	}
	type args struct {
		uid int32
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			us := &UserService{
				UserStore: tt.fields.UserStore,
			}
			if err := us.DeleteUser(tt.args.uid); (err != nil) != tt.wantErr {
				t.Errorf("UserService.DeleteUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserService_CheckUsername(t *testing.T) {
	type fields struct {
		UserStore *db.Queries
	}
	type args struct {
		username string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    db.Login
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			us := &UserService{
				UserStore: tt.fields.UserStore,
			}
			got, err := us.CheckUsername(tt.args.username)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserService.CheckUsername() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserService.CheckUsername() = %v, want %v", got, tt.want)
			}
		})
	}
}
