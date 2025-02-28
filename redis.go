package main

import (
	// Go Internal Packages
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"time"

	// External Packages
	"github.com/redis/go-redis/v9"
)

type Client struct {
	client *redis.Client
}

func NewRedis(URI string) (*Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:        URI,
		Password:    "",
		DB:          0,
		DialTimeout: 100 * time.Millisecond,
		ReadTimeout: 100 * time.Millisecond,
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	return &Client{client: client}, nil
}

func (c *Client) Get(ctx context.Context, empID string) (*EmpSalary, error) {
	cmd := c.client.Get(ctx, empID)
	cmdb, err := cmd.Bytes()
	if err != nil {
		return nil, fmt.Errorf("failed to get emp salary due to %v", err)
	}

	b := bytes.NewReader(cmdb)
	var empSalary EmpSalary
	if err := gob.NewDecoder(b).Decode(&empSalary); err != nil {
		return nil, err
	}
	return &empSalary, nil
}

func (c *Client) Set(ctx context.Context, empSalary *EmpSalary) error {
	var b bytes.Buffer
	if err := gob.NewEncoder(&b).Encode(empSalary); err != nil {
		return fmt.Errorf("encoding of struct failed due to: %v", err)
	}

	res := c.client.Set(ctx, empSalary.EmployeeID, b.Bytes(), 24*time.Hour)
	if err := res.Err(); err != nil {
		return fmt.Errorf("failed to set into cache due to: %v", err)
	}
	return nil
}
