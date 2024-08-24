#!/bin/bash

mongosh -u abc -p abc<<EOF
use test;
db.users.insertOne({
  username: "admin",
  password: "admin",
});
EOF
