package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCheckOrderNumber(t *testing.T) {
	type args struct {
		number string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "check valid number", args: args{number: "4561261212345467"}, want: true},
		{name: "check invalid number", args: args{number: "4561261212345464"}, want: false},
		{name: "check valid number with non digits", args: args{number: "4561 2612 1234 5467 asd"}, want: true},
		{name: "check invalid number with non digits", args: args{number: "4561 2612 1234 5464"}, want: false},
		{name: "check empty number", args: args{number: ""}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckOrderNumber(tt.args.number); got != tt.want {
				t.Errorf("CheckOrderNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPasswordCheck(t *testing.T) {
	t.Run("checks generated password", func(t *testing.T) {
		password := "123456ASD"
		password2 := "123456asd"
		h, err := PasswordHash(password)
		assert.NoError(t, err)
		assert.True(t, PasswordCheck(password, h))
		assert.False(t, PasswordCheck(password2, h))
	})
}

func TestPasswordHash(t *testing.T) {
	t.Run("generates different hashes for same password", func(t *testing.T) {
		h1, err := PasswordHash("123456")
		assert.NoError(t, err)
		h2, err := PasswordHash("123456")
		assert.NoError(t, err)
		assert.NotEqual(t, h1, h2)
	})
}
