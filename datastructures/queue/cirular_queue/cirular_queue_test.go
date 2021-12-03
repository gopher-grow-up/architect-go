package cirular_queue

import (
	"reflect"
	"testing"
)

func TestCircularQueue_Dequeue(t *testing.T) {
	type fields struct {
		Head  int
		Tail  int
		Items []string
		N     int
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{name: "enqueue and dequeue", fields: fields{
			Head:  0,
			Tail:  3,
			Items: []string{"a", "b", "c"},
			N:     3,
		}, want: "a"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := CircularQueue{
				Head:  tt.fields.Head,
				Tail:  tt.fields.Tail,
				Items: tt.fields.Items,
				N:     tt.fields.N,
			}
			if got := q.Dequeue(); got != tt.want {
				t.Errorf("Dequeue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCircularQueue_Enqueue(t *testing.T) {
	type fields struct {
		Head  int
		Tail  int
		Items []string
		N     int
	}
	type args struct {
		item string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{name: "enqueue and dequeue", args: args{item: "d"}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := NewCircularQueue(5)
			if err := q.Enqueue(tt.args.item); (err != nil) != tt.wantErr {
				t.Errorf("Enqueue() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCircularQueue_IsEmpty(t *testing.T) {
	type fields struct {
		Head  int
		Tail  int
		Items []string
		N     int
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := CircularQueue{
				Head:  tt.fields.Head,
				Tail:  tt.fields.Tail,
				Items: tt.fields.Items,
				N:     tt.fields.N,
			}
			if got := q.IsEmpty(); got != tt.want {
				t.Errorf("IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCircularQueue_IsFull(t *testing.T) {
	type fields struct {
		Head  int
		Tail  int
		Items []string
		N     int
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := CircularQueue{
				Head:  tt.fields.Head,
				Tail:  tt.fields.Tail,
				Items: tt.fields.Items,
				N:     tt.fields.N,
			}
			if got := q.IsFull(); got != tt.want {
				t.Errorf("IsFull() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewCircularQueue(t *testing.T) {
	type args struct {
		capacity int
	}
	tests := []struct {
		name string
		args args
		want CircularQueue
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCircularQueue(tt.args.capacity); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCircularQueue() = %v, want %v", got, tt.want)
			}
		})
	}
}
