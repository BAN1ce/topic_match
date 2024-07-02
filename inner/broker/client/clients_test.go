package client

import (
	"net"
	"reflect"
	"testing"
)

func TestNewClients(t *testing.T) {
	tests := []struct {
		name string
		want *Clients
	}{
		{
			name: "TestNewClients",
			want: &Clients{
				clients: map[string]*Client{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewClients(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewClients() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClients_CreateClient(t *testing.T) {
	type fields struct {
		clients map[UID]*Client
	}
	type args struct {
		client *Client
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "createNewClient",
			fields: fields{
				clients: map[UID]*Client{},
			},
			args: args{
				client: &Client{},
			},
		},
		{
			name: "createExistClient",
			fields: fields{
				clients: map[UID]*Client{
					"1": {conn: new(net.TCPConn)},
				},
			},
			args: args{
				client: &Client{UID: "1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Clients{
				clients: tt.fields.clients,
			}
			c.CreateClient(tt.args.client)
		})
	}
}

func TestClients_DestroyClient(t *testing.T) {

	var (
		server, conn = net.Pipe()
		client       = &Client{conn: conn}
	)
	defer func() {
		client.Close()
		server.Close()
	}()

	type fields struct {
		clients map[UID]*Client
	}
	type args struct {
		uid string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Client
		want1  bool
	}{
		{
			name: "destroyExistClient",
			fields: fields{
				clients: map[UID]*Client{
					"1": client,
				},
			},
			args: args{
				"1",
			},
			want:  client,
			want1: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Clients{
				clients: tt.fields.clients,
			}
			got, got1 := c.DestroyClient(tt.args.uid)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DestroyClient() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("DestroyClient() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestClients_ReadClient(t *testing.T) {
	var (
		client = &Client{conn: new(net.TCPConn)}
	)

	type fields struct {
		clients map[UID]*Client
	}
	type args struct {
		uid string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Client
		want1  bool
	}{
		{
			name: "readExistClient",
			fields: fields{
				clients: map[UID]*Client{
					"1": client,
				},
			},
			args: args{
				uid: "1",
			},
			want:  client,
			want1: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Clients{
				clients: tt.fields.clients,
			}
			got, got1 := c.ReadClient(tt.args.uid)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadClient() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("ReadClient() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
