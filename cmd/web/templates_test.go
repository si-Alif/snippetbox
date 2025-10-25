package main

import (
	"testing"
	"time"
	"snippetbox._alif__.net/internal/assert"
)

//Your unit tests are contained in a normal Go function with the signature func(*testing.T).

func TestHumanDate(t *testing.T){

	/*

	For executing a single testcase this is the recommended way

	tm:= time.Date(2025 , 10 , 26 , 2 , 59 , 59 , 0 , time.UTC)

	hd := human_date(tm)

	if hd != "26 Oct 2025 at 02:59" {
		t.Errorf("want %q got %q" , "26 Oct 2025 at 02:59" , hd)
	}

	*/

	// Since we want to execute multiple testcases , we can use a table driven test
	tests := []struct{
		name string
		tm time.Time
		want string
	}{
		{
		name : "UTC",
		tm : time.Date(2025 , 10 , 26 , 2 , 59 , 59 , 0 , time.UTC),
		want : "26 Oct 2025 at 02:59",
		},
		{
			name : "Empty",
			tm : time.Time{},
			want : "",
		},
		{
			name : "CET",
			tm : time.Date(2025 , 10 , 26 , 2 , 59 , 59 , 0 , time.FixedZone("CET" , 1 * 60 * 60)),
			want : "26 Oct 2025 at 02:59",
		},
	}

	for _ , tt := range tests{

		// t.Run() executes a single testcase with the signature <test_name> , func(*testing.T)

		t.Run(tt.name , func(t *testing.T) {
			hd := human_date(tt.tm)

			assert.Equal(t , hd , tt.want)
		})
	}

}