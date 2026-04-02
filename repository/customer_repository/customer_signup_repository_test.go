package customerrepository

import (
	"context"
	"e-commerce/model"
	"regexp"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func SetupMockDB(t *testing.T) (*gorm.DB , sqlmock.Sqlmock){

	sqlDB , mock , err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock %v",err)
	}

	gormDB , err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,

	}),&gorm.Config{})

	if err != nil{
		t.Fatalf("Failed to open gorm DB %v",err)
	}

	t.Cleanup(func() {
		sqlDB.Close()
	})

	return gormDB , mock

}

func Test_CreateCustomerAccount(t *testing.T) {
	testCases := []struct {
		description	string
		mockSetup func(mock sqlmock.Sqlmock)
		expectedError bool
	}{
		{
			description: "Should return no error while creating the customer's account when all details are valid and no error in inserting the record to the database",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(`INSERT INTO "customers"`).WillReturnRows(sqlmock.NewRows([]string{"CustomerId"}).AddRow(1))
				mock.ExpectCommit()
			},
			expectedError: false,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.description, func(t *testing.T) {
			db , mock := SetupMockDB(t)
			tC.mockSetup(mock)

			repo := NewCustomerSignupRepository(db)
			err := repo.CreateCustomerAccount(&model.Customer{
				Name : "rahul",
				Email: "rahul@gmail.com",
				PhoneNumber: "+917665718067",
				PasswordHash: "hash",
			})

			if tC.expectedError{
				assert.Error(t,err)
			} else {
				assert.NoError(t,err)
			}
		})
	}
}

func Test_FindByEmail(t *testing.T) {
	testCases := []struct {
		description	string
		email string
		mockSetup func(mock sqlmock.Sqlmock)
		expectedErr bool
	}{
		{
			description: "Should return no error when the email is valid and email already exists",
			email : "rahul@gmail.com",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"customer_id","name","email","phone_number","password_hash"}).AddRow(1,"rahul","rahul@gmail.com","+917665718067","hash")

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "customers" WHERE email = $1`)).WithArgs("rahul@gmail.com",1).WillReturnRows(rows)
			},
			expectedErr: false,
		},
		{
			description: "Should return Record Not Found error when the email is valid and email does not exists already",
			email : "rahul@gmail.com",
			mockSetup: func(mock sqlmock.Sqlmock) {

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "customers" WHERE email = $1`)).WithArgs("rahul@gmail.com",1).WillReturnError(gorm.ErrRecordNotFound)
			},
			expectedErr: true,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.description, func(t *testing.T) {
			db , mock := SetupMockDB(t)
			tC.mockSetup(mock)
			repo := NewCustomerSignupRepository(db)

			customer , err := repo.FindByEmail(context.Background(),tC.email)
			if tC.expectedErr {
				assert.Error(t,err)
				assert.Nil(t,customer)
			} else {
				assert.NoError(t,err)
				assert.NotNil(t,customer)
			}
		})
	}
}

func Test_FindByPhoneNumber(t *testing.T) {
	testCases := []struct {
		description	string
		phoneNumber string
		mockSetup func(mock sqlmock.Sqlmock)
		expectedErr bool
	}{
		{
			description: "Should return no error when the phone number is valid and phone number already exists in the database",
			phoneNumber: "+917665718067",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"customer_id","name","email","phone_number","password_hash"}).AddRow(1,"rahul","rahul@gmail.com","+917665718067","hash")

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "customers" WHERE phone_number = $1`)).WithArgs("+917665718067",1).WillReturnRows(rows)
			},
			expectedErr: false,
		},
		{
			description: "Should return record not found error when the phone number is valid and phone number does not exists in the database",
			phoneNumber: "+917665718067",
			mockSetup: func(mock sqlmock.Sqlmock) {

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "customers" WHERE phone_number = $1`)).WithArgs("+917665718067",1).WillReturnError(gorm.ErrRecordNotFound)
			},
			expectedErr: true,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.description, func(t *testing.T) {
			db , mock := SetupMockDB(t)
			tC.mockSetup(mock)
			repo := NewCustomerSignupRepository(db)

			customer , err := repo.FindByPhoneNumber(context.Background(),tC.phoneNumber)
			if tC.expectedErr {
				assert.Error(t,err)
				assert.Nil(t,customer)
			} else {
				assert.NoError(t,err)
				assert.NotNil(t,customer)
			}

		})
	}
}