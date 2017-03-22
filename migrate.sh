#!/usr/bin/env bash
migrate -url postgres://johndoe:foobar@localhost:5432/clientid_dev\?sslmode=disable -path ./migrations up
